// Package protocol re-exports the TUIOS daemon wire protocol types for external use.
package protocol

import (
	is "github.com/Gaurav-Gosain/tuios/internal/session"
)

// Core protocol types.
type (
	MessageType = is.MessageType
	Message     = is.Message
	CodecType   = is.CodecType
	Codec       = is.Codec
)

// Message type constants -- Client -> Server.
const (
	MsgHello            = is.MsgHello
	MsgAttach           = is.MsgAttach
	MsgDetach           = is.MsgDetach
	MsgNew              = is.MsgNew
	MsgList             = is.MsgList
	MsgKill             = is.MsgKill
	MsgInput            = is.MsgInput
	MsgResize           = is.MsgResize
	MsgPing             = is.MsgPing
	MsgCreatePTY        = is.MsgCreatePTY
	MsgClosePTY         = is.MsgClosePTY
	MsgListPTYs         = is.MsgListPTYs
	MsgFocusPTY         = is.MsgFocusPTY
	MsgGetState         = is.MsgGetState
	MsgUpdateState      = is.MsgUpdateState
	MsgSubscribePTY     = is.MsgSubscribePTY
	MsgUnsubscribePTY   = is.MsgUnsubscribePTY
	MsgGetTerminalState = is.MsgGetTerminalState
	MsgExecuteCommand   = is.MsgExecuteCommand
	MsgSendKeys         = is.MsgSendKeys
	MsgSetConfig        = is.MsgSetConfig
)

// Message type constants -- Server -> Client.
const (
	MsgWelcome       = is.MsgWelcome
	MsgAttached      = is.MsgAttached
	MsgDetached      = is.MsgDetached
	MsgSessionList   = is.MsgSessionList
	MsgOutput        = is.MsgOutput
	MsgError         = is.MsgError
	MsgPong          = is.MsgPong
	MsgSessionEnded  = is.MsgSessionEnded
	MsgWindowChanged = is.MsgWindowChanged
	MsgPTYList       = is.MsgPTYList
	MsgPTYCreated    = is.MsgPTYCreated
	MsgPTYClosed     = is.MsgPTYClosed
	MsgPTYOutput     = is.MsgPTYOutput
	MsgStateData     = is.MsgStateData
	MsgTerminalState = is.MsgTerminalState
	MsgCommandResult = is.MsgCommandResult
	MsgRemoteCommand = is.MsgRemoteCommand
	MsgGetLogs       = is.MsgGetLogs
	MsgLogsData      = is.MsgLogsData
	MsgQueryWindows  = is.MsgQueryWindows
	MsgWindowList    = is.MsgWindowList
	MsgQuerySession  = is.MsgQuerySession
	MsgSessionInfo   = is.MsgSessionInfo
)

// Multi-client message constants.
const (
	MsgStateSync       = is.MsgStateSync
	MsgClientJoined    = is.MsgClientJoined
	MsgClientLeft      = is.MsgClientLeft
	MsgSessionResize   = is.MsgSessionResize
	MsgForceRefresh    = is.MsgForceRefresh
	MsgRequestFullSync = is.MsgRequestFullSync
)

// Codec constants.
const (
	CodecGob  = is.CodecGob
	CodecJSON = is.CodecJSON
)

// Error code constants.
const (
	ErrCodeUnknown         = is.ErrCodeUnknown
	ErrCodeSessionNotFound = is.ErrCodeSessionNotFound
	ErrCodeSessionExists   = is.ErrCodeSessionExists
	ErrCodeInvalidMessage  = is.ErrCodeInvalidMessage
	ErrCodeInternal        = is.ErrCodeInternal
	ErrCodeNotAttached     = is.ErrCodeNotAttached
	ErrCodePTYNotFound     = is.ErrCodePTYNotFound
	ErrCodeNoTUIAttached   = is.ErrCodeNoTUIAttached
	ErrCodeCommandFailed   = is.ErrCodeCommandFailed
)

// ProtocolVersion is the current wire protocol version.
const ProtocolVersion = is.ProtocolVersion

// Payload types -- Handshake.
type (
	HelloPayload   = is.HelloPayload
	WelcomePayload = is.WelcomePayload
)

// Payload types -- Session lifecycle.
type (
	AttachPayload      = is.AttachPayload
	AttachedPayload    = is.AttachedPayload
	NewPayload         = is.NewPayload
	KillPayload        = is.KillPayload
	SessionListPayload = is.SessionListPayload
	SessionInfoType    = is.SessionInfo
	ErrorPayload       = is.ErrorPayload
	ResizePayload      = is.ResizePayload
)

// Payload types -- PTY operations.
type (
	CreatePTYPayload      = is.CreatePTYPayload
	PTYCreatedPayload     = is.PTYCreatedPayload
	ClosePTYPayload       = is.ClosePTYPayload
	FocusPTYPayload       = is.FocusPTYPayload
	InputPayload          = is.InputPayload
	PTYOutputPayload      = is.PTYOutputPayload
	ResizePTYPayload      = is.ResizePTYPayload
	SubscribePTYPayload   = is.SubscribePTYPayload
	UnsubscribePTYPayload = is.UnsubscribePTYPayload
	PTYInfo               = is.PTYInfo
	PTYListPayload        = is.PTYListPayload
)

// Payload types -- Terminal state.
type (
	GetTerminalStatePayload = is.GetTerminalStatePayload
	TerminalStatePayload    = is.TerminalStatePayload
	TerminalState           = is.TerminalState
	CellState               = is.CellState
)

// Payload types -- Session state.
type (
	SessionState     = is.SessionState
	WindowState      = is.WindowState
	SerializedBSPNode = is.SerializedBSPNode
	SerializedBSPTree = is.SerializedBSPTree
)

// Payload types -- Commands.
type (
	ExecuteCommandPayload = is.ExecuteCommandPayload
	SendKeysPayload       = is.SendKeysPayload
	SetConfigPayload      = is.SetConfigPayload
	CommandResultPayload  = is.CommandResultPayload
	RemoteCommandPayload  = is.RemoteCommandPayload
)

// Payload types -- Logs and queries.
type (
	GetLogsPayload      = is.GetLogsPayload
	LogsDataPayload     = is.LogsDataPayload
	LogEntry            = is.LogEntry
	QueryWindowsPayload = is.QueryWindowsPayload
	WindowInfo          = is.WindowInfo
	WindowListPayload   = is.WindowListPayload
	QuerySessionPayload = is.QuerySessionPayload
	SessionInfoPayload  = is.SessionInfoPayload
)

// Payload types -- Multi-client.
type (
	StateSyncPayload     = is.StateSyncPayload
	ClientJoinedPayload  = is.ClientJoinedPayload
	ClientLeftPayload    = is.ClientLeftPayload
	SessionResizePayload = is.SessionResizePayload
	ForceRefreshPayload  = is.ForceRefreshPayload
)

// Wire format helpers.
var (
	WriteMessageWithCodec = is.WriteMessageWithCodec
	ReadMessageWithCodec  = is.ReadMessageWithCodec
	WriteMessage          = is.WriteMessage
	ReadMessage           = is.ReadMessage
	NewMessageWithCodec   = is.NewMessageWithCodec
	NewMessage            = is.NewMessage
	NewRawMessage         = is.NewRawMessage
	WritePTYOutput        = is.WritePTYOutput
	WritePTYInput         = is.WritePTYInput
	ParseBinaryPTYMessage = is.ParseBinaryPTYMessage
	NegotiateCodec        = is.NegotiateCodec
	DefaultCodec          = is.DefaultCodec
	GetCodec              = is.GetCodec
	ParseCodecType        = is.ParseCodecType
)
