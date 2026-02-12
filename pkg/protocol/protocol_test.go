package protocol_test

import (
	"bytes"
	"testing"

	"github.com/Gaurav-Gosain/tuios/pkg/protocol"
)

func TestNewMessage(t *testing.T) {
	msg, err := protocol.NewMessage(protocol.MsgPing, nil)
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if msg.Type != protocol.MsgPing {
		t.Errorf("Type = %d, want %d", msg.Type, protocol.MsgPing)
	}
}

func TestWriteReadMessage(t *testing.T) {
	msg, err := protocol.NewMessage(protocol.MsgHello, &protocol.HelloPayload{
		Version: "1.0",
		Width:   80,
		Height:  24,
	})
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}

	var buf bytes.Buffer
	if err := protocol.WriteMessage(&buf, msg); err != nil {
		t.Fatalf("WriteMessage error: %v", err)
	}

	got, err := protocol.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("ReadMessage error: %v", err)
	}

	if got.Type != protocol.MsgHello {
		t.Errorf("Read Type = %d, want %d", got.Type, protocol.MsgHello)
	}
}

func TestBinaryPTYMessage(t *testing.T) {
	var buf bytes.Buffer
	data := []byte("hello from PTY")
	ptyID := "12345678-1234-1234-1234-123456789012"

	if err := protocol.WritePTYOutput(&buf, ptyID, data); err != nil {
		t.Fatalf("WritePTYOutput error: %v", err)
	}

	raw := buf.Bytes()
	payload := raw[6:] // skip 4 (length) + 2 (type+codec)

	gotID, gotData, err := protocol.ParseBinaryPTYMessage(payload)
	if err != nil {
		t.Fatalf("ParseBinaryPTYMessage error: %v", err)
	}
	if gotID != ptyID {
		t.Errorf("PTY ID = %q, want %q", gotID, ptyID)
	}
	if !bytes.Equal(gotData, data) {
		t.Errorf("Data = %q, want %q", gotData, data)
	}
}
