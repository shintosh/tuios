package vt_test

import (
	"strings"
	"testing"

	"github.com/Gaurav-Gosain/tuios/pkg/vt"
)

func TestNewEmulator(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	if emu == nil {
		t.Fatal("NewEmulator returned nil")
	}
	if emu.Width() != 80 {
		t.Errorf("Width = %d, want 80", emu.Width())
	}
	if emu.Height() != 24 {
		t.Errorf("Height = %d, want 24", emu.Height())
	}
}

func TestEmulatorWrite(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	n, err := emu.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != 5 {
		t.Errorf("Write returned %d, want 5", n)
	}
}

func TestEmulatorResize(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	emu.Resize(120, 40)
	if emu.Width() != 120 {
		t.Errorf("Width after resize = %d, want 120", emu.Width())
	}
	if emu.Height() != 40 {
		t.Errorf("Height after resize = %d, want 40", emu.Height())
	}
}

func TestEmulatorString(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	_, _ = emu.Write([]byte("hello world"))
	s := emu.String()
	// String() returns the full screen buffer; verify it contains our text.
	if !strings.Contains(s, "hello world") {
		t.Errorf("String() = %q, does not contain %q", s, "hello world")
	}
}

func TestNewScreen(t *testing.T) {
	scr := vt.NewScreen(40, 10)
	if scr == nil {
		t.Fatal("NewScreen returned nil")
	}
	if scr.Width() != 40 {
		t.Errorf("Width = %d, want 40", scr.Width())
	}
	if scr.Height() != 10 {
		t.Errorf("Height = %d, want 10", scr.Height())
	}
}

func TestNewScrollback(t *testing.T) {
	sb := vt.NewScrollback(100)
	if sb == nil {
		t.Fatal("NewScrollback returned nil")
	}
	if sb.MaxLines() != 100 {
		t.Errorf("MaxLines = %d, want 100", sb.MaxLines())
	}
	if sb.Len() != 0 {
		t.Errorf("Len = %d, want 0", sb.Len())
	}
}

func TestEmulatorScrollback(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	sb := emu.Scrollback()
	if sb == nil {
		t.Fatal("Scrollback() returned nil")
	}
}

func TestEmulatorCallbacks(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	var bellCalled bool
	emu.SetCallbacks(vt.Callbacks{
		Bell: func() { bellCalled = true },
	})
	// BEL character triggers the bell callback
	_, _ = emu.Write([]byte{0x07})
	if !bellCalled {
		t.Error("Bell callback was not called")
	}
}
