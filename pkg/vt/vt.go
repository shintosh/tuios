// Package vt re-exports the TUIOS VT100 terminal emulator for external use.
package vt

import (
	ivt "github.com/Gaurav-Gosain/tuios/internal/vt"
)

// Core types -- type aliases are transparent (same type, no conversion needed).
type (
	Emulator    = ivt.Emulator
	Screen      = ivt.Screen
	Scrollback  = ivt.Scrollback
	Callbacks   = ivt.Callbacks
	Cursor      = ivt.Cursor
	CursorStyle = ivt.CursorStyle
	Logger      = ivt.Logger
	CharSet     = ivt.CharSet

	// Input types
	Mouse         = ivt.Mouse
	MouseClick    = ivt.MouseClick
	MouseRelease  = ivt.MouseRelease
	MouseWheel    = ivt.MouseWheel
	MouseMotion   = ivt.MouseMotion
	MouseButton   = ivt.MouseButton
	KeyMod        = ivt.KeyMod
	KeyPressEvent = ivt.KeyPressEvent

	// Graphics types
	KittyState   = ivt.KittyState
	KittyCommand = ivt.KittyCommand
	KittyImage   = ivt.KittyImage
	SixelState   = ivt.SixelState
	SixelCommand = ivt.SixelCommand

	// Handler types for extending the emulator
	DcsHandler = ivt.DcsHandler
	CsiHandler = ivt.CsiHandler
	OscHandler = ivt.OscHandler
	ApcHandler = ivt.ApcHandler
	EscHandler = ivt.EscHandler
)

// Modifier key constants.
const (
	ModShift = ivt.ModShift
	ModAlt   = ivt.ModAlt
	ModCtrl  = ivt.ModCtrl
	ModMeta  = ivt.ModMeta
)

// Cursor style constants.
const (
	CursorBlock     = ivt.CursorBlock
	CursorUnderline = ivt.CursorUnderline
	CursorBar       = ivt.CursorBar
)

// Mouse button constants.
const (
	MouseNone       = ivt.MouseNone
	MouseLeft       = ivt.MouseLeft
	MouseMiddle     = ivt.MouseMiddle
	MouseRight      = ivt.MouseRight
	MouseWheelUp    = ivt.MouseWheelUp
	MouseWheelDown  = ivt.MouseWheelDown
	MouseWheelLeft  = ivt.MouseWheelLeft
	MouseWheelRight = ivt.MouseWheelRight
	MouseBackward   = ivt.MouseBackward
	MouseForward    = ivt.MouseForward
)

// Scrollback constants.
const DefaultScrollbackSize = ivt.DefaultScrollbackSize

// NewEmulator creates a new VT100 terminal emulator with the given dimensions.
var NewEmulator = ivt.NewEmulator

// NewScreen creates a new terminal screen buffer.
var NewScreen = ivt.NewScreen

// NewScrollback creates a new scrollback buffer with the given max lines.
var NewScrollback = ivt.NewScrollback

// NewKittyState creates a new Kitty graphics state.
var NewKittyState = ivt.NewKittyState

// NewSixelState creates a new Sixel graphics state.
var NewSixelState = ivt.NewSixelState
