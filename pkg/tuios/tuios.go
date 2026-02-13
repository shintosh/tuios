// Package tuios provides a reusable terminal window manager that can be
// embedded in other Bubble Tea applications or used as a standalone TUI.
//
// TUIOS (Terminal UI Operating System) is a terminal-based window manager
// that provides vim-like modal interface, workspace support, mouse interaction,
// and BSP tiling.
//
// # Basic Usage
//
// Create a new TUIOS instance with default options:
//
//	model := tuios.New()
//	p := tea.NewProgram(model)
//	if _, err := p.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// # Custom Configuration
//
// Use options to customize TUIOS behavior:
//
//	model := tuios.New(
//		tuios.WithTheme("dracula"),
//		tuios.WithShowKeys(true),
//		tuios.WithAnimations(false),
//		tuios.WithWorkspaces(9),
//	)
//
// # Using with sip (Web Terminal)
//
// TUIOS can be served through the browser using the sip library:
//
//	server := sip.NewServer(sip.DefaultConfig())
//	server.Serve(ctx, func(sess sip.Session) (tea.Model, []tea.ProgramOption) {
//		return tuios.NewForSession(sess.Pty()), nil
//	})
package tuios

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/Gaurav-Gosain/tuios/internal/app"
	"github.com/Gaurav-Gosain/tuios/internal/config"
	"github.com/Gaurav-Gosain/tuios/internal/input"
	"github.com/Gaurav-Gosain/tuios/internal/session"
	"github.com/Gaurav-Gosain/tuios/internal/terminal"
)

// Model is the main TUIOS model that implements tea.Model.
// It wraps the internal OS struct and provides a clean public API.
type Model = app.OS

// Mode represents the current interaction mode of TUIOS.
type Mode = app.Mode

// Mode constants
const (
	// WindowManagementMode allows window manipulation and navigation.
	WindowManagementMode = app.WindowManagementMode
	// TerminalMode passes input directly to the focused terminal.
	TerminalMode = app.TerminalMode
)

// Options configures a TUIOS instance.
type Options struct {
	// Theme is the color theme name (e.g., "dracula", "nord", "tokyonight").
	// Leave empty to use standard terminal colors.
	Theme string

	// ShowKeys enables the showkeys overlay to display pressed keys.
	ShowKeys bool

	// Animations enables/disables window animations.
	// When disabled, windows snap instantly to positions.
	Animations bool

	// ASCIIOnly uses ASCII characters instead of Nerd Font icons.
	ASCIIOnly bool

	// Workspaces is the number of workspaces (1-9). Default is 9.
	Workspaces int

	// BorderStyle sets the window border style.
	// Valid values: "rounded", "normal", "thick", "double", "hidden", "block", "ascii"
	BorderStyle string

	// DockbarPosition sets where the dockbar appears.
	// Valid values: "bottom", "top", "hidden"
	DockbarPosition string

	// HideWindowButtons hides the minimize/maximize/close buttons.
	HideWindowButtons bool

	// ScrollbackLines is the number of lines in scrollback buffer.
	// Default is 10000, min 100, max 1000000.
	ScrollbackLines int

	// Width is the initial width (set automatically if 0).
	Width int

	// Height is the initial height (set automatically if 0).
	Height int

	// SSHMode indicates if running over SSH.
	SSHMode bool

	// DaemonSocketPath is the path to the TUIOS daemon Unix socket.
	// When set, the TUIOS app connects to the daemon for PTY management.
	DaemonSocketPath string

	// SessionName is the daemon session name to create or attach to.
	SessionName string

	// EnableGraphicsPassthrough enables Kitty/Sixel graphics passthrough.
	// Should be true for terminal sessions, false for web sessions.
	EnableGraphicsPassthrough bool

	// HideClock hides the clock from the dockbar.
	HideClock bool

	// WindowTitlePosition sets the window title position.
	// Valid values: "left", "center", "right"
	WindowTitlePosition string

	// NoAnimations disables all UI animations.
	NoAnimations bool

	// Version is the version string sent during daemon handshake.
	Version string

	// UserConfig is a custom user configuration. If nil, defaults are used.
	UserConfig *config.UserConfig
}

// Option is a functional option for configuring TUIOS.
type Option func(*Options)

// WithTheme sets the color theme.
func WithTheme(name string) Option {
	return func(o *Options) {
		o.Theme = name
	}
}

// WithShowKeys enables the showkeys overlay.
func WithShowKeys(enabled bool) Option {
	return func(o *Options) {
		o.ShowKeys = enabled
	}
}

// WithAnimations enables or disables window animations.
func WithAnimations(enabled bool) Option {
	return func(o *Options) {
		o.Animations = enabled
	}
}

// WithASCIIOnly enables ASCII-only mode (no Nerd Font icons).
func WithASCIIOnly(enabled bool) Option {
	return func(o *Options) {
		o.ASCIIOnly = enabled
	}
}

// WithWorkspaces sets the number of workspaces (1-9).
func WithWorkspaces(n int) Option {
	return func(o *Options) {
		if n < 1 {
			n = 1
		} else if n > 9 {
			n = 9
		}
		o.Workspaces = n
	}
}

// WithBorderStyle sets the window border style.
func WithBorderStyle(style string) Option {
	return func(o *Options) {
		o.BorderStyle = style
	}
}

// WithDockbarPosition sets the dockbar position.
func WithDockbarPosition(position string) Option {
	return func(o *Options) {
		o.DockbarPosition = position
	}
}

// WithHideWindowButtons hides window control buttons.
func WithHideWindowButtons(hide bool) Option {
	return func(o *Options) {
		o.HideWindowButtons = hide
	}
}

// WithScrollbackLines sets the scrollback buffer size.
func WithScrollbackLines(lines int) Option {
	return func(o *Options) {
		if lines < 100 {
			lines = 100
		} else if lines > 1000000 {
			lines = 1000000
		}
		o.ScrollbackLines = lines
	}
}

// WithSize sets the initial terminal size.
func WithSize(width, height int) Option {
	return func(o *Options) {
		o.Width = width
		o.Height = height
	}
}

// WithSSHMode enables SSH mode.
func WithSSHMode(enabled bool) Option {
	return func(o *Options) {
		o.SSHMode = enabled
	}
}

// WithDaemonSocket sets the daemon socket path for PTY management.
func WithDaemonSocket(path string) Option {
	return func(o *Options) {
		o.DaemonSocketPath = path
	}
}

// WithSessionName sets the daemon session name.
func WithSessionName(name string) Option {
	return func(o *Options) {
		o.SessionName = name
	}
}

// WithGraphicsPassthrough enables Kitty/Sixel graphics passthrough.
func WithGraphicsPassthrough(enabled bool) Option {
	return func(o *Options) {
		o.EnableGraphicsPassthrough = enabled
	}
}

// WithHideClock hides the clock from the dockbar.
func WithHideClock(hide bool) Option {
	return func(o *Options) {
		o.HideClock = hide
	}
}

// WithWindowTitlePosition sets the window title position.
func WithWindowTitlePosition(position string) Option {
	return func(o *Options) {
		o.WindowTitlePosition = position
	}
}

// WithNoAnimations disables all UI animations.
func WithNoAnimations(enabled bool) Option {
	return func(o *Options) {
		o.NoAnimations = enabled
	}
}

// WithVersion sets the version string for daemon handshake.
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// WithUserConfig sets a custom user configuration.
func WithUserConfig(cfg *config.UserConfig) Option {
	return func(o *Options) {
		o.UserConfig = cfg
	}
}

// DefaultOptions returns the default options.
func DefaultOptions() Options {
	return Options{
		Animations:      true,
		Workspaces:      9,
		ScrollbackLines: 10000,
	}
}

// New creates a new TUIOS model with the given options.
// This is the main entry point for using TUIOS as a library.
func New(opts ...Option) *Model {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return newModel(options)
}

// NewForSession creates a new TUIOS model configured for a PTY session.
// This is useful when embedding TUIOS in web terminals or SSH servers.
//
// The pty parameter should have Width and Height fields.
type PTY interface {
	Width() int
	Height() int
}

// NewForPTY creates a new TUIOS model for a PTY session with the given options.
func NewForPTY(pty PTY, opts ...Option) *Model {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.Width = pty.Width()
	options.Height = pty.Height()

	return newModel(options)
}

// newModel creates the internal model with applied options.
func newModel(options Options) *Model {
	// Set up input handler
	app.SetInputHandler(input.HandleInput)

	// Load or create user config
	var userConfig *config.UserConfig
	if options.UserConfig != nil {
		userConfig = options.UserConfig
	} else {
		var err error
		userConfig, err = config.LoadUserConfig()
		if err != nil {
			userConfig = config.DefaultConfig()
		}
	}

	// Apply overrides using the central config system
	config.ApplyOverrides(config.Overrides{
		ASCIIOnly:           options.ASCIIOnly,
		BorderStyle:         options.BorderStyle,
		DockbarPosition:     options.DockbarPosition,
		HideWindowButtons:   options.HideWindowButtons,
		WindowTitlePosition: options.WindowTitlePosition,
		HideClock:           options.HideClock,
		ScrollbackLines:     options.ScrollbackLines,
		NoAnimations:        options.NoAnimations,
		ThemeName:           options.Theme,
	}, userConfig)

	// Backwards compat: handle Animations=false (old API)
	if !options.Animations {
		config.AnimationsEnabled = false
	}

	// Create keybind registry
	keybindRegistry := config.NewKeybindRegistry(userConfig)

	// Determine if this is a daemon session
	isDaemonSession := options.DaemonSocketPath != ""

	// Create the model using the factory function
	return app.NewOS(app.OSOptions{
		KeybindRegistry:          keybindRegistry,
		ShowKeys:                 options.ShowKeys,
		NumWorkspaces:            options.Workspaces,
		Width:                    options.Width,
		Height:                   options.Height,
		IsSSHMode:                options.SSHMode,
		IsDaemonSession:          isDaemonSession,
		SessionName:              options.SessionName,
		EnableGraphicsPassthrough: options.EnableGraphicsPassthrough,
	})
}

// Run creates and runs a full TUIOS TUI application, blocking until exit.
// This is the main entry point for running TUIOS as a standalone process.
// It handles daemon connection, signal handling, and terminal cleanup.
//
// Usage from shinto:
//
//	err := tuios.Run(
//		tuios.WithDaemonSocket("/tmp/tuios.sock"),
//		tuios.WithSessionName("main"),
//		tuios.WithGraphicsPassthrough(true),
//		tuios.WithTheme("dracula"),
//	)
func Run(opts ...Option) error {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// Set up input handler
	app.SetInputHandler(input.HandleInput)

	// Load or create user config
	var userConfig *config.UserConfig
	if options.UserConfig != nil {
		userConfig = options.UserConfig
	} else {
		var err error
		userConfig, err = config.LoadUserConfig()
		if err != nil {
			userConfig = config.DefaultConfig()
		}
	}

	// Apply overrides using the central config system
	config.ApplyOverrides(config.Overrides{
		ASCIIOnly:           options.ASCIIOnly,
		BorderStyle:         options.BorderStyle,
		DockbarPosition:     options.DockbarPosition,
		HideWindowButtons:   options.HideWindowButtons,
		WindowTitlePosition: options.WindowTitlePosition,
		HideClock:           options.HideClock,
		ScrollbackLines:     options.ScrollbackLines,
		NoAnimations:        options.NoAnimations,
		ThemeName:           options.Theme,
	}, userConfig)

	if !options.Animations {
		config.AnimationsEnabled = false
	}

	keybindRegistry := config.NewKeybindRegistry(userConfig)

	isDaemonSession := options.DaemonSocketPath != ""

	osOpts := app.OSOptions{
		KeybindRegistry:          keybindRegistry,
		ShowKeys:                 options.ShowKeys,
		NumWorkspaces:            options.Workspaces,
		Width:                    options.Width,
		Height:                   options.Height,
		IsSSHMode:                options.SSHMode,
		IsDaemonSession:          isDaemonSession,
		SessionName:              options.SessionName,
		EnableGraphicsPassthrough: options.EnableGraphicsPassthrough,
	}

	// Connect to daemon if socket path is provided
	if isDaemonSession {
		client := session.NewTUIClient()
		ver := options.Version
		if ver == "" {
			ver = "embedded"
		}
		if err := client.ConnectToSocket(options.DaemonSocketPath, ver, options.Width, options.Height); err != nil {
			return fmt.Errorf("failed to connect to daemon: %w", err)
		}
		osOpts.DaemonClient = client
	}

	model := app.NewOS(osOpts)

	p := tea.NewProgram(
		model,
		tea.WithFPS(config.NormalFPS),
		tea.WithoutSignalHandler(),
		tea.WithFilter(FilterMouseMotion),
	)

	// Handle signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		p.Send(tea.QuitMsg{})
	}()

	finalModel, err := p.Run()

	if finalOS, ok := finalModel.(*app.OS); ok {
		finalOS.Cleanup()
	}

	terminal.ResetTerminal()

	if err != nil {
		return fmt.Errorf("tuios error: %w", err)
	}

	return nil
}

// ProgramOptions returns recommended tea.ProgramOption values for running TUIOS.
// Use these when creating a tea.Program:
//
//	model := tuios.New()
//	p := tea.NewProgram(model, tuios.ProgramOptions()...)
func ProgramOptions() []tea.ProgramOption {
	return []tea.ProgramOption{
		tea.WithFPS(config.NormalFPS),
	}
}

// FilterMouseMotion is a tea.WithFilter function that reduces CPU usage
// by filtering out redundant mouse motion events.
// Only passes through mouse motion during drag/resize/selection operations.
//
// Usage:
//
//	p := tea.NewProgram(model, tea.WithFilter(tuios.FilterMouseMotion))
func FilterMouseMotion(model tea.Model, msg tea.Msg) tea.Msg {
	// Allow all non-motion events through
	if _, ok := msg.(tea.MouseMotionMsg); !ok {
		return msg
	}

	// Type assert to our OS model
	os, ok := model.(*Model)
	if !ok {
		return msg
	}

	// Allow motion events during active interactions
	if os.Dragging || os.Resizing {
		return msg
	}

	// Allow motion events during text selection
	if os.SelectionMode {
		focusedWindow := os.GetFocusedWindow()
		if focusedWindow != nil && focusedWindow.IsSelecting {
			return msg
		}
	}

	// Allow motion events when in terminal mode with alt screen apps
	if os.Mode == TerminalMode {
		focusedWindow := os.GetFocusedWindow()
		if focusedWindow != nil && focusedWindow.IsAltScreen {
			return msg
		}
	}

	// Filter out motion events when not interacting
	return nil
}

// Config re-exports the config package for customization.
// This allows users to access configuration types without importing internal packages.
var Config = struct {
	// LoadUserConfig loads the user's configuration file.
	LoadUserConfig func() (*config.UserConfig, error)
	// DefaultConfig returns the default configuration.
	DefaultConfig func() *config.UserConfig
	// GetConfigPath returns the path to the configuration file.
	GetConfigPath func() (string, error)
}{
	LoadUserConfig: config.LoadUserConfig,
	DefaultConfig:  config.DefaultConfig,
	GetConfigPath:  config.GetConfigPath,
}
