package session_test

import (
	"testing"

	"github.com/Gaurav-Gosain/tuios/pkg/session"
)

func TestNewManager(t *testing.T) {
	mgr := session.NewManager()
	if mgr == nil {
		t.Fatal("NewManager returned nil")
	}
	if mgr.SessionCount() != 0 {
		t.Errorf("SessionCount = %d, want 0", mgr.SessionCount())
	}
}

func TestNewDaemon(t *testing.T) {
	d := session.NewDaemon(&session.DaemonConfig{
		Version: "test",
	})
	if d == nil {
		t.Fatal("NewDaemon returned nil")
	}
}

func TestNewSession(t *testing.T) {
	s, err := session.NewSession("test-session", &session.SessionConfig{}, 80, 24)
	if err != nil {
		t.Fatalf("NewSession error: %v", err)
	}
	if s == nil {
		t.Fatal("NewSession returned nil")
	}

	info := s.Info()
	if info.Name != "test-session" {
		t.Errorf("Name = %q, want %q", info.Name, "test-session")
	}
	if info.Width != 80 {
		t.Errorf("Width = %d, want 80", info.Width)
	}
	if info.Height != 24 {
		t.Errorf("Height = %d, want 24", info.Height)
	}
}

func TestCreateSession(t *testing.T) {
	mgr := session.NewManager()
	s, err := mgr.CreateSession("test-session", &session.SessionConfig{}, 80, 24)
	if err != nil {
		t.Fatalf("CreateSession error: %v", err)
	}
	if s == nil {
		t.Fatal("CreateSession returned nil")
	}

	info := s.Info()
	if info.Name != "test-session" {
		t.Errorf("Name = %q, want %q", info.Name, "test-session")
	}

	if mgr.SessionCount() != 1 {
		t.Errorf("SessionCount = %d, want 1", mgr.SessionCount())
	}

	// Cleanup
	_ = mgr.DeleteSession("test-session")

	if mgr.SessionCount() != 0 {
		t.Errorf("SessionCount after delete = %d, want 0", mgr.SessionCount())
	}
}

func TestGetSocketPath(t *testing.T) {
	path, err := session.GetSocketPath()
	if err != nil {
		t.Fatalf("GetSocketPath error: %v", err)
	}
	if path == "" {
		t.Error("GetSocketPath returned empty path")
	}
}

func TestGetPidFilePath(t *testing.T) {
	path, err := session.GetPidFilePath()
	if err != nil {
		t.Fatalf("GetPidFilePath error: %v", err)
	}
	if path == "" {
		t.Error("GetPidFilePath returned empty path")
	}
}

func TestIsDaemonRunning(t *testing.T) {
	// Just verify the function is callable and returns a bool.
	// In a test environment the daemon is typically not running.
	_ = session.IsDaemonRunning()
}

func TestStateToCell(t *testing.T) {
	cs := session.CellState{
		Content: "A",
		Width:   1,
		FgColor: "#ff0000",
		Bold:    true,
	}

	cell := session.StateToCell(cs)
	if cell == nil {
		t.Fatal("StateToCell returned nil")
	}
	if cell.Content != "A" {
		t.Errorf("Content = %q, want %q", cell.Content, "A")
	}
	if cell.Width != 1 {
		t.Errorf("Width = %d, want 1", cell.Width)
	}
}
