package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"strings"
)

type assistantResponseEvent struct {
	Content   string  `json:"content"`
	Input     *string `json:"input,omitempty"`
	Name      string  `json:"name"`
	ToolUseId string  `json:"toolUseId"`
	Stop      bool    `json:"stop"`
}

type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type parserState struct {
	currentBlockType  string
	currentBlockIndex int
	currentToolUseID  string
	currentToolName   string
	nextBlockIndex    int
}

// maxHeaderBytes caps an upstream eventstream header length before any
// allocation. Upstream is trusted-ish, but a proxy/MITM in front can inject a
// corrupt frame; without this a bogus headerLen makes make([]byte, headerLen)
// allocate up to ~4GB. Mirrors 9router EVENTSTREAM_MAX_HEADERS_BYTES.
const maxHeaderBytes = 64 * 1024

func ParseEvents(resp []byte) []SSEEvent {

	events := []SSEEvent{}
	state := parserState{currentBlockIndex: -1}

	r := bytes.NewReader(resp)
	for {
		if r.Len() < 12 {
			break
		}

		var totalLen, headerLen uint32
		if err := binary.Read(r, binary.BigEndian, &totalLen); err != nil {
			break
		}
		if err := binary.Read(r, binary.BigEndian, &headerLen); err != nil {
			break
		}

		// Validate frame bounds BEFORE allocating anything. totalLen includes the
		// 12-byte prelude and 4-byte message CRC, so a valid frame is >= 16 bytes
		// and headerLen must leave room for prelude + CRC (else payloadLen goes
		// negative and make([]byte, payloadLen) panics).
		if totalLen < 16 {
			log.Println("Frame length invalid (totalLen < 16)")
			break
		}
		if headerLen > maxHeaderBytes {
			log.Printf("Frame header too large: %d > %d", headerLen, maxHeaderBytes)
			break
		}
		if headerLen > totalLen-16 {
			log.Printf("Frame header exceeds total length: header=%d total=%d", headerLen, totalLen)
			break
		}
		if int(totalLen) > r.Len()+8 {
			log.Println("Frame length invalid (beyond stream)")
			break
		}

		// Skip header
		header := make([]byte, headerLen)
		if _, err := io.ReadFull(r, header); err != nil {
			break
		}

		payloadLen := int(totalLen) - int(headerLen) - 12
		payload := make([]byte, payloadLen)
		if _, err := io.ReadFull(r, payload); err != nil {
			break
		}

		// Skip CRC32
		if _, err := r.Seek(4, io.SeekCurrent); err != nil {
			break
		}

		// Handle binary framing and clean up payload
		payloadStr := string(payload)
		if idx := strings.Index(payloadStr, "{"); idx != -1 {
			payloadStr = payloadStr[idx:]
		}

		// First try parsing as assistantResponseEvent
		var assistantEvt assistantResponseEvent
		if err := json.Unmarshal([]byte(payloadStr), &assistantEvt); err == nil && (assistantEvt.Content != "" || assistantEvt.ToolUseId != "" || assistantEvt.Stop) {
			appendAssistantEvent(&events, &state, assistantEvt)
			continue
		}

		// Handling 2026+ metadata events (metering, context usage)
		var metaData map[string]any
		if err := json.Unmarshal([]byte(payloadStr), &metaData); err == nil {
			if _, exists := metaData["contextUsagePercentage"]; exists {
				// Translate to a ping/metadata event for Claude Code
				events = append(events, SSEEvent{
					Event: "ping",
					Data:  map[string]any{"type": "ping", "metadata": metaData},
				})
			} else if _, exists := metaData["unit"]; exists {
				// Metering event - silently consume or log lightly
				log.Printf("Usage: %v %v", metaData["usage"], metaData["unit"])
			}
		}

		// Metrics event: Kiro reports token usage / cache accounting here. This
		// is what makes upstream cache hits visible (cache_read_input_tokens).
		if evt, ok := parseMetricsEvent(payloadStr); ok {
			events = append(events, SSEEvent{Event: "metrics", Data: evt})
		}
	}

	return events
}

func appendAssistantEvent(events *[]SSEEvent, state *parserState, evt assistantResponseEvent) {
	switch {
	case evt.Content != "":
		if state.currentBlockType != "text" {
			closeCurrentBlock(events, state)
			startTextBlock(events, state)
		}
		*events = append(*events, SSEEvent{
			Event: "content_block_delta",
			Data: map[string]interface{}{
				"type":  "content_block_delta",
				"index": state.currentBlockIndex,
				"delta": map[string]interface{}{
					"type": "text_delta",
					"text": evt.Content,
				},
			},
		})
	case evt.ToolUseId != "":
		toolName := evt.Name
		if toolName == "" {
			toolName = state.currentToolName
		}

		if state.currentBlockType != "tool_use" || state.currentToolUseID != evt.ToolUseId {
			closeCurrentBlock(events, state)
			startToolBlock(events, state, evt.ToolUseId, toolName)
		}

		if evt.Input != nil {
			*events = append(*events, SSEEvent{
				Event: "content_block_delta",
				Data: map[string]interface{}{
					"type":  "content_block_delta",
					"index": state.currentBlockIndex,
					"delta": map[string]interface{}{
						"type":         "input_json_delta",
						"id":           evt.ToolUseId,
						"name":         toolName,
						"partial_json": *evt.Input,
					},
				},
			})
		}

		if evt.Stop {
			closeCurrentBlock(events, state)
			appendMessageDelta(events, "tool_use")
		}
	case evt.Stop:
		closeCurrentBlock(events, state)
		appendMessageDelta(events, "end_turn")
	}
}

func startTextBlock(events *[]SSEEvent, state *parserState) {
	index := state.nextBlockIndex
	state.nextBlockIndex++
	state.currentBlockType = "text"
	state.currentBlockIndex = index
	state.currentToolUseID = ""
	state.currentToolName = ""

	*events = append(*events, SSEEvent{
		Event: "content_block_start",
		Data: map[string]interface{}{
			"type":  "content_block_start",
			"index": index,
			"content_block": map[string]interface{}{
				"type": "text",
				"text": "",
			},
		},
	})
}

func startToolBlock(events *[]SSEEvent, state *parserState, toolUseID, name string) {
	index := state.nextBlockIndex
	state.nextBlockIndex++
	state.currentBlockType = "tool_use"
	state.currentBlockIndex = index
	state.currentToolUseID = toolUseID
	state.currentToolName = name

	*events = append(*events, SSEEvent{
		Event: "content_block_start",
		Data: map[string]interface{}{
			"type":  "content_block_start",
			"index": index,
			"content_block": map[string]interface{}{
				"type":  "tool_use",
				"id":    toolUseID,
				"name":  name,
				"input": map[string]interface{}{},
			},
		},
	})
}

func closeCurrentBlock(events *[]SSEEvent, state *parserState) {
	if state.currentBlockType == "" || state.currentBlockIndex < 0 {
		return
	}

	*events = append(*events, SSEEvent{
		Event: "content_block_stop",
		Data: map[string]interface{}{
			"type":  "content_block_stop",
			"index": state.currentBlockIndex,
		},
	})

	state.currentBlockType = ""
	state.currentBlockIndex = -1
	state.currentToolUseID = ""
	state.currentToolName = ""
}

func appendMessageDelta(events *[]SSEEvent, stopReason string) {
	*events = append(*events, SSEEvent{
		Event: "message_delta",
		Data: map[string]interface{}{
			"type": "message_delta",
			"delta": map[string]interface{}{
				"stop_reason":   stopReason,
				"stop_sequence": nil,
			},
			"usage": map[string]interface{}{"output_tokens": 0},
		},
	})
}

// parseMetricsEvent extracts a Kiro metricsEvent frame into a normalized
// Anthropic-style usage object. The payload may be flat or nested under a
// "metricsEvent" key, and token fields come in either camelCase or snake_case
// spelling (mirrors 9router executors/kiro.js). Returns ok=false when the
// payload is not a metrics event.
func parseMetricsEvent(payload string) (map[string]any, bool) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return nil, false
	}
	obj, ok := raw["metricsEvent"]
	if !ok {
		obj = raw
	}
	m, ok := obj.(map[string]any)
	if !ok {
		return nil, false
	}
	// Must look like a metrics event: at least one recognized token field.
	if _, h := m["inputTokens"]; !h {
		if _, h := m["input_tokens"]; !h {
			if _, h := m["outputTokens"]; !h {
				if _, h := m["output_tokens"]; !h {
					return nil, false
				}
			}
		}
	}
	return map[string]any{
		"type":                        "metrics",
		"input_tokens":                numMetric(m, "inputTokens", "input_tokens"),
		"output_tokens":               numMetric(m, "outputTokens", "output_tokens"),
		"cache_read_input_tokens":     numMetric(m, "cacheReadInputTokens", "cache_read_input_tokens"),
		"cache_creation_input_tokens": numMetric(m, "cacheCreationInputTokens", "cache_creation_input_tokens"),
	}, true
}

func numMetric(m map[string]any, camel, snake string) int {
	if v, ok := m[camel]; ok {
		return intFromAny(v)
	}
	if v, ok := m[snake]; ok {
		return intFromAny(v)
	}
	return 0
}

func intFromAny(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0
		}
		return int(i)
	}
	return 0
}
