package session

import (
	"fmt"
	"net"
	"time"
)

// ConnectToSocket connects to a daemon at a specific socket path.
// Like Connect but allows specifying the socket path instead of using GetSocketPath().
func (c *TUIClient) ConnectToSocket(socketPath, version string, width, height int) error {
	return c.ConnectToSocketWithCapabilities(socketPath, version, width, height, nil)
}

// ConnectToSocketWithCapabilities connects to a daemon at a specific socket path with capabilities.
func (c *TUIClient) ConnectToSocketWithCapabilities(socketPath, version string, width, height int, caps *ClientCapabilities) error {
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon at %s: %w", socketPath, err)
	}
	c.conn = conn
	c.connected = true

	hello := &HelloPayload{
		Version:        version,
		Width:          width,
		Height:         height,
		PreferredCodec: "gob",
	}
	if caps != nil {
		hello.PixelWidth = caps.PixelWidth
		hello.PixelHeight = caps.PixelHeight
		hello.CellWidth = caps.CellWidth
		hello.CellHeight = caps.CellHeight
		hello.KittyGraphics = caps.KittyGraphics
		hello.SixelGraphics = caps.SixelGraphics
		hello.TerminalName = caps.TerminalName
	}

	msg, err := NewMessageWithCodec(MsgHello, hello, c.codec)
	if err != nil {
		_ = conn.Close()
		return err
	}
	if err := c.send(msg); err != nil {
		_ = conn.Close()
		return err
	}

	resp, err := c.recv()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if resp.Type != MsgWelcome {
		_ = conn.Close()
		return fmt.Errorf("expected welcome, got %d", resp.Type)
	}

	var welcome WelcomePayload
	if err := resp.ParsePayloadWithCodec(&welcome, c.codec); err != nil {
		_ = conn.Close()
		return fmt.Errorf("failed to parse welcome: %w", err)
	}

	c.codec = NegotiateCodec(welcome.Codec)
	c.availableSessionNames = welcome.SessionNames
	return nil
}
