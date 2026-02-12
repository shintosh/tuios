package vt_test

import (
	"strings"
	"testing"

	"github.com/Gaurav-Gosain/tuios/internal/testutil"
	"github.com/Gaurav-Gosain/tuios/internal/vt"
)

// =============================================================================
// Basic Text Output Tests
// =============================================================================

func TestEmulator_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "multiple lines",
			input:    "Line 1\r\nLine 2\r\nLine 3",
			expected: "Line 1",
		},
		{
			name:     "text with tabs",
			input:    "Col1\tCol2\tCol3",
			expected: "Col1    Col2    Col3",
		},
		{
			name:     "unicode text",
			input:    "Hello \xe4\xb8\x96\xe7\x95\x8c", // "Hello 世界" in UTF-8
			expected: "Hello 世界",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			emu := vt.NewEmulator(80, 24)
			_, err := emu.Write([]byte(tc.input))
			if err != nil {
				t.Fatalf("Write failed: %v", err)
			}

			got := emu.String()
			if !strings.Contains(got, tc.expected) {
				t.Errorf("Expected output to contain %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// Cursor Movement Tests
// =============================================================================

func TestEmulator_CursorMovement(t *testing.T) {
	tests := []struct {
		name     string
		build    func(*testutil.ANSIBuilder) string
		expected string
	}{
		{
			name: "cursor home",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("XXXXX").CursorHome().Text("Hello").String()
			},
			expected: "Hello",
		},
		{
			name: "cursor to position",
			build: func(b *testutil.ANSIBuilder) string {
				return b.CursorTo(1, 5).Text("Test").String()
			},
			expected: "    Test",
		},
		{
			name: "cursor back and overwrite",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("Hello").CursorBackward(5).Text("World").String()
			},
			expected: "World",
		},
		{
			name: "cursor forward",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("A").CursorForward(3).Text("B").String()
			},
			expected: "A   B",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			emu := vt.NewEmulator(80, 24)
			input := tc.build(testutil.NewANSIBuilder())
			_, err := emu.Write([]byte(input))
			if err != nil {
				t.Fatalf("Write failed: %v", err)
			}

			got := emu.String()
			if !strings.Contains(got, tc.expected) {
				t.Errorf("Expected output to contain %q, got %q", tc.expected, got)
			}
		})
	}
}

// =============================================================================
// Screen Clear Tests
// =============================================================================

func TestEmulator_ScreenClear(t *testing.T) {
	tests := []struct {
		name      string
		build     func(*testutil.ANSIBuilder) string
		checkFunc func(*vt.Emulator) bool
	}{
		{
			name: "clear entire screen",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("Hello\r\nWorld").ClearScreen().CursorHome().Text("New").String()
			},
			checkFunc: func(e *vt.Emulator) bool {
				s := e.String()
				return strings.Contains(s, "New") && !strings.Contains(s, "Hello")
			},
		},
		{
			name: "clear line",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("Hello World").CR().ClearLine().Text("X").String()
			},
			checkFunc: func(e *vt.Emulator) bool {
				return strings.Contains(e.String(), "X")
			},
		},
		{
			name: "clear to end of line",
			build: func(b *testutil.ANSIBuilder) string {
				return b.Text("Hello World").CursorTo(1, 6).ClearToEndOfLine().String()
			},
			checkFunc: func(e *vt.Emulator) bool {
				s := e.String()
				return strings.Contains(s, "Hello") && !strings.Contains(s, "World")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			emu := vt.NewEmulator(80, 24)
			input := tc.build(testutil.NewANSIBuilder())
			_, err := emu.Write([]byte(input))
			if err != nil {
				t.Fatalf("Write failed: %v", err)
			}

			if !tc.checkFunc(emu) {
				t.Errorf("Check failed for %s, screen content: %q", tc.name, emu.String())
			}
		})
	}
}

// =============================================================================
// SGR (Color/Style) Tests
// =============================================================================

func TestEmulator_SGRStyles(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	// Write text with various SGR codes
	b := testutil.NewANSIBuilder()
	input := b.
		Bold().Text("Bold").Reset().Text(" ").
		Italic().Text("Italic").Reset().Text(" ").
		Underline().Text("Underline").Reset().Text(" ").
		FgColor(31).Text("Red").Reset().Text(" ").
		BgColor(44).Text("BlueBg").Reset().
		String()

	_, err := emu.Write([]byte(input))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check that text content is preserved
	got := emu.String()
	for _, word := range []string{"Bold", "Italic", "Underline", "Red", "BlueBg"} {
		if !strings.Contains(got, word) {
			t.Errorf("Expected output to contain %q, got %q", word, got)
		}
	}

	// Check rendered output contains styles
	rendered := emu.Render()
	if !strings.Contains(rendered, "\x1b[") {
		t.Error("Rendered output should contain ANSI escape sequences")
	}
}

func TestEmulator_256Colors(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	b := testutil.NewANSIBuilder()
	input := b.Fg256(196).Text("Bright Red").Reset().String()

	_, err := emu.Write([]byte(input))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := emu.String()
	if !strings.Contains(got, "Bright Red") {
		t.Errorf("Expected output to contain 'Bright Red', got %q", got)
	}
}

func TestEmulator_TrueColor(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	b := testutil.NewANSIBuilder()
	input := b.FgRGB(255, 128, 0).Text("Orange").Reset().String()

	_, err := emu.Write([]byte(input))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := emu.String()
	if !strings.Contains(got, "Orange") {
		t.Errorf("Expected output to contain 'Orange', got %q", got)
	}
}

// =============================================================================
// Alternate Screen Buffer Tests
// =============================================================================

func TestEmulator_AltScreen(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	// Track alt screen state
	var altScreenEnabled bool
	emu.SetCallbacks(vt.Callbacks{
		AltScreen: func(enabled bool) {
			altScreenEnabled = enabled
		},
	})

	b := testutil.NewANSIBuilder()

	// Write to main screen
	_, _ = emu.Write([]byte(b.Text("Main screen content").String()))
	b.Clear()

	// Switch to alt screen
	_, _ = emu.Write([]byte(b.AltScreen().String()))
	b.Clear()

	if !altScreenEnabled {
		t.Error("Expected alt screen to be enabled")
	}

	// Write to alt screen
	_, _ = emu.Write([]byte(b.ClearScreen().CursorHome().Text("Alt screen content").String()))
	b.Clear()

	// Check alt screen content
	got := emu.String()
	if !strings.Contains(got, "Alt screen content") {
		t.Errorf("Expected alt screen content, got %q", got)
	}

	// Switch back to main screen
	_, _ = emu.Write([]byte(b.MainScreen().String()))

	if altScreenEnabled {
		t.Error("Expected alt screen to be disabled")
	}

	// Main screen content should be preserved
	got = emu.String()
	if !strings.Contains(got, "Main screen content") {
		t.Errorf("Expected main screen content to be preserved, got %q", got)
	}
}

// =============================================================================
// Scrolling Tests
// =============================================================================

func TestEmulator_Scrolling(t *testing.T) {
	// Small terminal to test scrolling
	emu := vt.NewEmulator(40, 5)

	// Write more lines than the terminal height
	for i := 1; i <= 10; i++ {
		_, _ = emu.Write([]byte(testutil.NewANSIBuilder().
			Text(strings.Repeat("Line", 1)).
			Text(string(rune('0' + i%10))).
			Newline().
			String()))
	}

	// Check scrollback has content
	scrollbackLen := emu.ScrollbackLen()
	if scrollbackLen == 0 {
		t.Error("Expected scrollback to have content")
	}
}

func TestEmulator_ScrollRegion(t *testing.T) {
	emu := vt.NewEmulator(40, 10)

	b := testutil.NewANSIBuilder()

	// Set scroll region to lines 2-5
	_, _ = emu.Write([]byte(b.ScrollRegion(2, 5).String()))
	b.Clear()

	// Fill screen with content
	for i := 1; i <= 10; i++ {
		_, _ = emu.Write([]byte(b.CursorTo(i, 1).Text(strings.Repeat("=", 10)).String()))
		b.Clear()
	}

	// The scroll region should be set
	// Full testing would require checking line positions after scroll operations
}

// =============================================================================
// OSC Tests (Title, etc)
// =============================================================================

func TestEmulator_OSCTitle(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	var receivedTitle string
	emu.SetCallbacks(vt.Callbacks{
		Title: func(title string) {
			receivedTitle = title
		},
	})

	b := testutil.NewANSIBuilder()
	_, _ = emu.Write([]byte(b.OSCTitle("Test Window Title").String()))

	if receivedTitle != "Test Window Title" {
		t.Errorf("Expected title 'Test Window Title', got %q", receivedTitle)
	}
}

// =============================================================================
// Insert/Delete Character Tests
// =============================================================================

func TestEmulator_InsertDeleteChars(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	b := testutil.NewANSIBuilder()

	// Write text, move cursor, insert characters
	_, _ = emu.Write([]byte(b.Text("ABCDE").CursorTo(1, 3).InsertChars(2).String()))
	b.Clear()

	got := emu.String()
	// After inserting 2 chars at position 3, "AB  CDE" or similar
	if !strings.Contains(got, "AB") {
		t.Errorf("Insert chars test failed, got %q", got)
	}
}

func TestEmulator_InsertDeleteLines(t *testing.T) {
	emu := vt.NewEmulator(40, 10)
	b := testutil.NewANSIBuilder()

	// Fill with content
	for i := 1; i <= 5; i++ {
		_, _ = emu.Write([]byte(b.CursorTo(i, 1).Text("Line" + string(rune('0'+i))).String()))
		b.Clear()
	}

	// Insert lines at position 2
	_, _ = emu.Write([]byte(b.CursorTo(2, 1).InsertLines(2).String()))

	got := emu.String()
	// Line1 should still be at position 1
	if !strings.Contains(got, "Line1") {
		t.Errorf("Insert lines test failed, got %q", got)
	}
}

// =============================================================================
// Device Status Report Tests
// =============================================================================

func TestEmulator_DeviceStatusReport(t *testing.T) {
	emu := vt.NewEmulator(80, 24)
	b := testutil.NewANSIBuilder()

	// Move cursor to known position
	_, _ = emu.Write([]byte(b.CursorTo(5, 10).String()))

	// Note: DSR command would write a response to the emulator's pipe reader
	// which blocks if nothing reads from it. In real use, the PTY goroutine
	// reads from this pipe. For unit testing, we verify cursor positioning works.

	// Verify the cursor position was set by writing at cursor position
	b.Clear()
	_, _ = emu.Write([]byte(b.Text("X").String()))

	got := emu.String()
	// The 'X' should appear after spaces (position 10 means 9 spaces)
	if !strings.Contains(got, "X") {
		t.Errorf("Expected 'X' in output, got %q", got)
	}
}

// =============================================================================
// Resize Tests
// =============================================================================

func TestEmulator_Resize(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	// Write some content
	_, _ = emu.Write([]byte("Hello World"))

	// Resize
	emu.Resize(40, 10)

	// Check new bounds
	bounds := emu.Bounds()
	if bounds.Dx() != 40 || bounds.Dy() != 10 {
		t.Errorf("Expected bounds 40x10, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// =============================================================================
// Scrollback Tests
// =============================================================================

func TestEmulator_Scrollback(t *testing.T) {
	emu := vt.NewEmulator(40, 5)
	emu.SetScrollbackMaxLines(100)

	// Write many lines to fill scrollback
	for i := 1; i <= 20; i++ {
		_, _ = emu.Write([]byte(testutil.NewANSIBuilder().
			Text("Line " + string(rune('A'-1+i))).
			Newline().
			String()))
	}

	scrollbackLen := emu.ScrollbackLen()
	if scrollbackLen == 0 {
		t.Error("Expected scrollback to have content")
	}

	// Get a line from scrollback
	if scrollbackLen > 0 {
		line := emu.ScrollbackLine(0)
		if line == nil {
			t.Error("Expected to get scrollback line 0")
		}
	}

	// Clear scrollback
	emu.ClearScrollback()
	if emu.ScrollbackLen() != 0 {
		t.Error("Expected scrollback to be cleared")
	}
}

// =============================================================================
// Integration Test: Shell-like Output
// =============================================================================

func TestEmulator_ShellPromptSimulation(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	// Simulate typical shell output
	prompt := testutil.ShellPrompt("user", "host", "~")
	_, _ = emu.Write([]byte(prompt))
	_, _ = emu.Write([]byte("ls -la\r\n"))

	// Simulate ls output
	lsOutput := testutil.LSOutput(
		[]string{"Documents", "Downloads", "file.txt"},
		[]bool{true, true, false},
	)
	_, _ = emu.Write([]byte(lsOutput))

	// Write another prompt
	_, _ = emu.Write([]byte(prompt))

	got := emu.String()
	if !strings.Contains(got, "user@host") {
		t.Errorf("Expected shell prompt in output, got %q", got)
	}
	if !strings.Contains(got, "ls -la") {
		t.Errorf("Expected command in output, got %q", got)
	}
}

func TestEmulator_ProgressBarSimulation(t *testing.T) {
	emu := vt.NewEmulator(80, 24)

	// Simulate progress bar updates
	for i := 0; i <= 100; i += 20 {
		progress := testutil.ProgressBar(i, 20)
		_, _ = emu.Write([]byte(progress))
	}

	got := emu.String()
	if !strings.Contains(got, "100%") {
		t.Errorf("Expected 100%% in output, got %q", got)
	}
}

// =============================================================================
// Wide Character / Phantom State Tests
// =============================================================================

func TestEmulator_WideCharPhantomState(t *testing.T) {
	// In an 80-column terminal, a width-2 character placed at column 78
	// (0-indexed) fills columns 78-79. The emulator should enter phantom
	// (pending wrap) state so the next character wraps to the next line.
	emu := vt.NewEmulator(80, 24)

	// Fill columns 0-77 with 'A' (78 characters), then write a wide char.
	line := strings.Repeat("A", 78) + "世" // 世 is width-2
	_, err := emu.Write([]byte(line))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// The wide char should be at position (78, 0) with width 2.
	cell := emu.CellAt(78, 0)
	if cell == nil || cell.Content != "世" {
		t.Fatalf("Expected wide char at (78,0), got %+v", cell)
	}

	// Now write a normal ASCII char - it should wrap to (0, 1).
	_, err = emu.Write([]byte("B"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	cell = emu.CellAt(0, 1)
	if cell == nil || cell.Content != "B" {
		t.Errorf("Expected 'B' at (0,1) after wrap, got %+v", cell)
	}

	// The wide char placeholder at (79, 0) should NOT have been overwritten.
	cell = emu.CellAt(79, 0)
	if cell != nil && cell.Width != 0 {
		t.Errorf("Expected placeholder (width 0) at (79,0), got width=%d content=%q", cell.Width, cell.Content)
	}
}

func TestEmulator_WideCharOverflowWraps(t *testing.T) {
	// When a width-2 character is placed at the last column (79 in 80-col),
	// there's only 1 column remaining - not enough for a width-2 char.
	// The emulator should wrap to the next line and place it at column 0.
	emu := vt.NewEmulator(80, 24)

	// Fill columns 0-78 with 'A' (79 characters), leaving only column 79.
	line := strings.Repeat("A", 79) + "世"
	_, err := emu.Write([]byte(line))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// The wide char should NOT be at (79, 0) - it doesn't fit there.
	// It should have wrapped to (0, 1).
	cell := emu.CellAt(0, 1)
	if cell == nil || cell.Content != "世" {
		t.Errorf("Expected wide char at (0,1) after overflow wrap, got %+v", cell)
	}

	// Position (79, 0) should be empty (the last col of the first line).
	cell = emu.CellAt(79, 0)
	if cell != nil && cell.Content != "" && cell.Content != " " {
		t.Errorf("Expected empty cell at (79,0), got %+v", cell)
	}
}

func TestEmulator_WideCharSequence(t *testing.T) {
	// Fill a line with alternating wide chars. In 80 cols, we can fit
	// 40 width-2 characters per line.
	emu := vt.NewEmulator(80, 24)

	// Write 41 wide chars - 40 fit on line 0, 1 wraps to line 1.
	input := strings.Repeat("世", 41)
	_, err := emu.Write([]byte(input))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Check the last wide char on line 0 is at position 78.
	cell := emu.CellAt(78, 0)
	if cell == nil || cell.Content != "世" {
		t.Errorf("Expected wide char at (78,0), got %+v", cell)
	}

	// The 41st wide char should be at (0, 1).
	cell = emu.CellAt(0, 1)
	if cell == nil || cell.Content != "世" {
		t.Errorf("Expected wide char at (0,1), got %+v", cell)
	}
}

func TestEmulator_WideCharCursorPosition(t *testing.T) {
	// After writing a wide char that fills the last two columns, the cursor
	// should be in phantom state. We verify by checking that the next char
	// goes to the next line.
	emu := vt.NewEmulator(10, 5) // Small terminal for clarity

	// Fill 8 cols then write a width-2 char to fill cols 8-9.
	_, _ = emu.Write([]byte("AAAAAAAA世"))

	// Cursor should be at position 8 (phantom state, stays at write pos).
	pos := emu.CursorPosition()
	if pos.X != 8 || pos.Y != 0 {
		t.Errorf("Expected cursor at (8,0) in phantom state, got (%d,%d)", pos.X, pos.Y)
	}

	// Next char should trigger wrap.
	_, _ = emu.Write([]byte("X"))
	pos = emu.CursorPosition()
	if pos.Y != 1 {
		t.Errorf("Expected cursor on line 1 after wrap, got line %d", pos.Y)
	}
}

// =============================================================================
// Benchmarks
// =============================================================================

func BenchmarkEmulator_PlainTextWrite(b *testing.B) {
	emu := vt.NewEmulator(80, 24)
	data := []byte(strings.Repeat("Hello World ", 100) + "\r\n")

	b.ResetTimer()
	for b.Loop() {
		_, _ = emu.Write(data)
	}
}

func BenchmarkEmulator_ANSIColorWrite(b *testing.B) {
	emu := vt.NewEmulator(80, 24)
	builder := testutil.NewANSIBuilder()
	data := []byte(builder.
		FgColor(31).Text("Red").Reset().Text(" ").
		FgColor(32).Text("Green").Reset().Text(" ").
		FgColor(34).Text("Blue").Reset().
		Newline().
		String())

	b.ResetTimer()
	for b.Loop() {
		_, _ = emu.Write(data)
	}
}

func BenchmarkEmulator_CursorMovement(b *testing.B) {
	emu := vt.NewEmulator(80, 24)
	builder := testutil.NewANSIBuilder()
	data := []byte(builder.
		CursorTo(10, 10).Text("X").
		CursorTo(5, 5).Text("Y").
		CursorHome().Text("Z").
		String())

	b.ResetTimer()
	for b.Loop() {
		_, _ = emu.Write(data)
	}
}

func BenchmarkEmulator_Render(b *testing.B) {
	emu := vt.NewEmulator(80, 24)

	// Fill screen with styled content
	builder := testutil.NewANSIBuilder()
	for i := range 24 {
		_, _ = emu.Write([]byte(builder.
			FgColor(31 + i%6).
			Text(strings.Repeat("X", 80)).
			Newline().
			String()))
		builder.Clear()
	}

	b.ResetTimer()
	for b.Loop() {
		_ = emu.Render()
	}
}
