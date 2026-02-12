package client_test

import (
	"testing"

	"github.com/Gaurav-Gosain/tuios/pkg/client"
)

func TestNewTUIClient(t *testing.T) {
	c := client.NewTUIClient()
	if c == nil {
		t.Fatal("NewTUIClient returned nil")
	}
	if c.IsConnected() {
		t.Error("New client should not be connected")
	}
}

func TestNewClient(t *testing.T) {
	c := client.NewClient(&client.ClientConfig{})
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
}
