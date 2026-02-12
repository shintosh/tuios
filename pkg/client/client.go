// Package client re-exports the TUIOS daemon client library for external use.
// The TUIClient connects to a running TUIOS daemon to manage PTYs and sessions.
package client

import (
	is "github.com/Gaurav-Gosain/tuios/internal/session"
)

// Client types.
type (
	TUIClient          = is.TUIClient
	ClientCapabilities = is.ClientCapabilities
)

// Handler types for multi-client support.
type (
	StateSyncHandler     = is.StateSyncHandler
	ClientJoinedHandler  = is.ClientJoinedHandler
	ClientLeftHandler    = is.ClientLeftHandler
	SessionResizeHandler = is.SessionResizeHandler
	ForceRefreshHandler  = is.ForceRefreshHandler
	RemoteCommandHandler = is.RemoteCommandHandler
	QueryWindowsHandler  = is.QueryWindowsHandler
	QuerySessionHandler  = is.QuerySessionHandler
)

// Constructors.
var NewTUIClient = is.NewTUIClient

// Also expose the simple Client (for CLI attach use).
type (
	Client       = is.Client
	ClientConfig = is.ClientConfig
)

var NewClient = is.NewClient
