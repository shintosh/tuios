// Package session re-exports the TUIOS daemon session management types for external use.
package session

import (
	is "github.com/Gaurav-Gosain/tuios/internal/session"
)

// Daemon types.
type (
	Daemon       = is.Daemon
	DaemonConfig = is.DaemonConfig
)

// Session types.
type (
	Session       = is.Session
	SessionConfig = is.SessionConfig
	SessionState  = is.SessionState
	SessionInfo   = is.SessionInfo
	WindowState   = is.WindowState
)

// BSP tiling types.
type (
	SerializedBSPNode = is.SerializedBSPNode
	SerializedBSPTree = is.SerializedBSPTree
)

// PTY types.
type (
	PTY           = is.PTY
	PTYInfo       = is.PTYInfo
	TerminalState = is.TerminalState
	CellState     = is.CellState
)

// Manager types.
type (
	Manager = is.Manager
)

// Constructors.
var (
	NewDaemon  = is.NewDaemon
	NewManager = is.NewManager
	NewSession = is.NewSession
)

// Utility functions.
var (
	GetSocketPath   = is.GetSocketPath
	GetPidFilePath  = is.GetPidFilePath
	GetDaemonPID    = is.GetDaemonPID
	IsDaemonRunning = is.IsDaemonRunning
	StateToCell     = is.StateToCell
)
