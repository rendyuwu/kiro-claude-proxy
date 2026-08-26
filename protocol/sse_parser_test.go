package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"testing"
)

func encodeTestFrame(t *testing.T, payload any) []byte {
	t.Helper()

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	var frame bytes.Buffer
	totalLen := uint32(len(data) + 12)
	if err := binary.Write(&frame, binary.BigEndian, totalLen); err != nil {
		t.Fatalf("write totalLen: %v", err)
	}
	if err := binary.Write(&frame, binary.BigEndian, uint32(0)); err != nil {
		t.Fatalf("write headerLen: %v", err)
	}
	frame.Write(data)
	frame.Write([]byte{0, 0, 0, 0})
	return frame.Bytes()
}

func TestParseEventsTextEndTurn(t *testing.T) {
	frames := bytes.Join([][]byte{
		encodeTestFrame(t, assistantResponseEvent{Content: "hello"}),
		encodeTestFrame(t, assistantResponseEvent{Stop: true}),
	}, nil)
	events := ParseEvents(frames)

	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}

	if events[0].Event != "content_block_start" {
		t.Fatalf("expected first event content_block_start, got %q", events[0].Event)
	}
	start := events[0].Data.(map[string]any)
	if got := start["index"]; got != 0 {
		t.Fatalf("expected text block index 0, got %#v", got)
	}
	block := start["content_block"].(map[string]any)
	if got := block["type"]; got != "text" {
		t.Fatalf("expected text block type, got %#v", got)
	}

	delta := events[1].Data.(map[string]any)["delta"].(map[string]any)
	if got := delta["type"]; got != "text_delta" {
		t.Fatalf("expected text_delta, got %#v", got)
	}
	if got := delta["text"]; got != "hello" {
		t.Fatalf("expected text %q, got %#v", "hello", got)
	}

	stop := events[2].Data.(map[string]any)
	if got := stop["index"]; got != 0 {
		t.Fatalf("expected text stop index 0, got %#v", got)
	}

	messageDelta := events[3].Data.(map[string]any)["delta"].(map[string]any)
	if got := messageDelta["stop_reason"]; got != "end_turn" {
		t.Fatalf("expected stop_reason end_turn, got %#v", got)
	}
}

func TestParseEventsToolUseSequence(t *testing.T) {
	input := `{"query":"drift"}`
	frames := bytes.Join([][]byte{
		encodeTestFrame(t, assistantResponseEvent{Content: "Need tool"}),
		encodeTestFrame(t, assistantResponseEvent{ToolUseId: "tool-1", Name: "search"}),
		encodeTestFrame(t, assistantResponseEvent{ToolUseId: "tool-1", Name: "search", Input: &input}),
		encodeTestFrame(t, assistantResponseEvent{ToolUseId: "tool-1", Stop: true}),
	}, nil)

	events := ParseEvents(frames)
	if len(events) != 7 {
		t.Fatalf("expected 7 events, got %d", len(events))
	}

	textStop := events[2].Data.(map[string]any)
	if got := textStop["index"]; got != 0 {
		t.Fatalf("expected text block stop index 0, got %#v", got)
	}

	if events[3].Event != "content_block_start" {
		t.Fatalf("expected tool start event, got %q", events[3].Event)
	}
	start := events[3].Data.(map[string]any)["content_block"].(map[string]any)
	if got := start["type"]; got != "tool_use" {
		t.Fatalf("expected tool_use block, got %#v", got)
	}
	if got := start["id"]; got != "tool-1" {
		t.Fatalf("expected tool id %q, got %#v", "tool-1", got)
	}

	delta := events[4].Data.(map[string]any)["delta"].(map[string]any)
	if got := delta["type"]; got != "input_json_delta" {
		t.Fatalf("expected input_json_delta, got %#v", got)
	}
	var partial string
	switch v := delta["partial_json"].(type) {
	case *string:
		partial = *v
	case string:
		partial = v
	default:
		t.Fatalf("expected partial_json as string or *string, got %T", delta["partial_json"])
	}
	if got := partial; got != input {
		t.Fatalf("expected partial_json %q, got %q", input, got)
	}

	stop := events[5].Data.(map[string]any)
	if got := stop["type"]; got != "content_block_stop" {
		t.Fatalf("expected content_block_stop, got %#v", got)
	}
	if got := stop["index"]; got != 1 {
		t.Fatalf("expected tool stop index 1, got %#v", got)
	}

	messageDelta := events[6].Data.(map[string]any)["delta"].(map[string]any)
	if got := messageDelta["stop_reason"]; got != "tool_use" {
		t.Fatalf("expected stop_reason tool_use, got %#v", got)
	}
}

func TestParseEventsMetadataPing(t *testing.T) {
	frames := encodeTestFrame(t, map[string]any{"contextUsagePercentage": 73})
	events := ParseEvents(frames)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Event != "ping" {
		t.Fatalf("expected ping event, got %q", events[0].Event)
	}
	data := events[0].Data.(map[string]any)
	if got := data["type"]; got != "ping" {
		t.Fatalf("expected ping data type, got %#v", got)
	}
}

func TestParseEventsMetricsCamelCase(t *testing.T) {
	frames := encodeTestFrame(t, map[string]any{
		"inputTokens":              100.0,
		"outputTokens":             50.0,
		"cacheReadInputTokens":     80.0,
		"cacheCreationInputTokens": 20.0,
	})
	events := ParseEvents(frames)

	metrics := findMetricsEvent(t, events)
	if got := metrics["input_tokens"]; got != 100 {
		t.Fatalf("expected input_tokens 100, got %#v", got)
	}
	if got := metrics["output_tokens"]; got != 50 {
		t.Fatalf("expected output_tokens 50, got %#v", got)
	}
	if got := metrics["cache_read_input_tokens"]; got != 80 {
		t.Fatalf("expected cache_read_input_tokens 80, got %#v", got)
	}
	if got := metrics["cache_creation_input_tokens"]; got != 20 {
		t.Fatalf("expected cache_creation_input_tokens 20, got %#v", got)
	}
}

func TestParseEventsMetricsSnakeCaseAndNested(t *testing.T) {
	frames := bytes.Join([][]byte{
		encodeTestFrame(t, map[string]any{
			"input_tokens":            7,
			"cache_read_input_tokens": 3,
		}),
		encodeTestFrame(t, map[string]any{
			"metricsEvent": map[string]any{"inputTokens": 9, "outputTokens": 4},
		}),
	}, nil)
	events := ParseEvents(frames)

	metrics := findMetricsEvent(t, events)
	if got := metrics["input_tokens"]; got != 7 {
		t.Fatalf("expected snake input_tokens 7, got %#v", got)
	}
	if got := metrics["cache_read_input_tokens"]; got != 3 {
		t.Fatalf("expected snake cache_read_input_tokens 3, got %#v", got)
	}
}

func TestParseEventsRejectsOversizedHeaderNoHugeAlloc(t *testing.T) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, uint32(0xFFFFFFFF)); err != nil {
		t.Fatalf("write totalLen: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(0xFFFFFFFF)); err != nil {
		t.Fatalf("write headerLen: %v", err)
	}

	events := ParseEvents(buf.Bytes())
	if len(events) != 0 {
		t.Fatalf("expected 0 events for oversized header, got %d", len(events))
	}
}

func TestParseEventsRejectsHeaderBeyondTotalLenNoPanic(t *testing.T) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, uint32(100)); err != nil {
		t.Fatalf("write totalLen: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(200)); err != nil {
		t.Fatalf("write headerLen: %v", err)
	}
	// Trailing bytes large enough that io.ReadFull would succeed if the bounds
	// checks were missing — this exercises the negative-payloadLen panic path.
	buf.Write(make([]byte, 200))

	events := ParseEvents(buf.Bytes())
	if len(events) != 0 {
		t.Fatalf("expected 0 events for header beyond totalLen, got %d", len(events))
	}
}

func TestParseEventsRejectsTinyTotalLen(t *testing.T) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, uint32(15)); err != nil {
		t.Fatalf("write totalLen: %v", err)
	}
	if err := binary.Write(&buf, binary.BigEndian, uint32(0)); err != nil {
		t.Fatalf("write headerLen: %v", err)
	}

	events := ParseEvents(buf.Bytes())
	if len(events) != 0 {
		t.Fatalf("expected 0 events for tiny totalLen, got %d", len(events))
	}
}

func findMetricsEvent(t *testing.T, events []SSEEvent) map[string]any {
	t.Helper()
	for _, e := range events {
		if e.Event == "metrics" {
			if data, ok := e.Data.(map[string]any); ok {
				return data
			}
			t.Fatalf("metrics event data not a map: %#v", e.Data)
		}
	}
	t.Fatalf("expected a metrics event, got %#v", events)
	return nil
}
