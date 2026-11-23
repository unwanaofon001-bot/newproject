package main

import (
	"os"
	"strings"
	"testing"
)



// TestPrintASCII_ReturnsOutput checks that printAscii returns something when
// given a valid character. An empty result means the banner file wasn't read
// or the character wasn't rendered at all.
func TestPrintASCII_ReturnsOutput(t *testing.T) {
	result := printAscii("", "", "", false)
	if result == "" {
		t.Errorf("Expected non-empty output for input 'A', got empty string")
	}
}

// TestPrintASCII_OutputHasNineRows checks that the output contains exactly 9
// newlines, because every ASCII character in the banner is 9 rows tall.
func TestPrintASCII_OutputHasNineRows(t *testing.T) {
	result := printAscii("", "", "A", false)
	rows := strings.Count(result, "\n")
	if rows != 9 {
		t.Errorf("Expected 9 rows, got %d", rows)
	}
}

// TestPrintASCII_NoColorCode checks that when colorMode is false, no ANSI
// escape codes appear in the output at all.
func TestPrintASCII_NoColorCode(t *testing.T) {
	result := printAscii("red", "", "A", false)
	if strings.Contains(result, "\033[") {
		t.Errorf("Expected no color codes when colorMode is false, but found one")
	}
}

// TestPrintASCII_RedColorCode checks that when colorMode is true and the color
// is "red", the ANSI red code \033[31m appears in the output.
func TestPrintASCII_RedColorCode(t *testing.T) {
	result := printAscii("red", "", "A", true)
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Expected red ANSI code \\033[31m in output, but it was missing")
	}
}

// TestPrintASCII_GreenColorCode checks the same thing but for green (\033[32m).
func TestPrintASCII_GreenColorCode(t *testing.T) {
	result := printAscii("green", "", "A", true)
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("Expected green ANSI code \\033[32m in output, but it was missing")
	}
}

// TestPrintASCII_ResetCode checks that the reset code \033[0m is included
// whenever color mode is on, so colors don't bleed into the terminal.
func TestPrintASCII_ResetCode(t *testing.T) {
	result := printAscii("blue", "", "A", true)
	if !strings.Contains(result, "\033[0m") {
		t.Errorf("Expected reset ANSI code \\033[0m in output, but it was missing")
	}
}

// TestPrintASCII_UnknownColorNoCode checks that an unrecognised color name
// (one not in the colorCodes map) produces no escape code — the map lookup
// returns an empty string, so colorCode will be "".
func TestPrintASCII_UnknownColorNoCode(t *testing.T) {
	result := printAscii("magenta", "", "A", true)
	if strings.Contains(result, "\033[") {
		t.Errorf("Expected no ANSI code for unknown color 'magenta', but found one")
	}
}

// ── substring coloring tests ──────────────────────────────────────────────────

// TestPrintASCII_SubstringColored checks that when a substring is provided,
// the color code still appears (meaning the matched portion was colored).
func TestPrintASCII_SubstringColored(t *testing.T) {
	result := printAscii("green", "kit", "kitten", true)
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("Expected green color on substring 'kit' inside 'kitten'")
	}
}

// TestPrintASCII_SubstringNotInText checks that when the substring does not
// appear in the text, NO color code is produced (nothing was matched).
func TestPrintASCII_SubstringNotInText(t *testing.T) {
	result := printAscii("red", "xyz", "kitten", true)
	if strings.Contains(result, "\033[31m") {
		t.Errorf("Expected no color code when substring 'xyz' is not in 'kitten'")
	}
}

// TestPrintASCII_EmptySubstringColorsAll checks that when no substring is
// given, all characters are colored (the whole text gets the color code).
func TestPrintASCII_EmptySubstringColorsAll(t *testing.T) {
	result := printAscii("cyan", "", "Hi", true)
	// With two characters, the color code should appear at least twice
	count := strings.Count(result, "\033[36m")
	if count < 2 {
		t.Errorf("Expected color code to appear at least twice for 2-char text, got %d", count)
	}
}



// helper saves and restores os.Args automatically using defer.
// Every parseInput test calls this at the top so cleanup is guaranteed.
func setArgs(args []string) func() {
	old := os.Args
	os.Args = args
	return func() { os.Args = old } // returned function restores the original
}

// TestParseInput_AllArguments checks that all four values are parsed correctly
// when the full set of arguments is provided.
func TestParseInput_AllArguments(t *testing.T) {
	defer setArgs([]string{"cmd", "--color=red", "kit", "kitten"})()

	color, substring, text, colorMode := parseInput()

	if color != "red" {
		t.Errorf("Expected color 'red', got '%s'", color)
	}
	if substring != "kit" {
		t.Errorf("Expected substring 'kit', got '%s'", substring)
	}
	if text != "kitten" {
		t.Errorf("Expected text 'kitten', got '%s'", text)
	}
	if !colorMode {
		t.Errorf("Expected colorMode true, got false")
	}
}

// TestParseInput_NoSubstring checks the case where only color and text are
// given (no substring). substring should come back as an empty string.
func TestParseInput_NoSubstring(t *testing.T) {
	defer setArgs([]string{"cmd", "--color=blue", "hello"})()

	color, substring, text, colorMode := parseInput()

	if color != "blue" {
		t.Errorf("Expected color 'blue', got '%s'", color)
	}
	if substring != "" {
		t.Errorf("Expected empty substring, got '%s'", substring)
	}
	if text != "hello" {
		t.Errorf("Expected text 'hello', got '%s'", text)
	}
	if !colorMode {
		t.Errorf("Expected colorMode true, got false")
	}
}

// TestParseInput_ColorModeFlag checks that colorMode is true whenever the
// --color= flag is present, regardless of what color is specified.
func TestParseInput_ColorModeFlag(t *testing.T) {
	defer setArgs([]string{"cmd", "--color=yellow", "world"})()

	_, _, _, colorMode := parseInput()
	if !colorMode {
		t.Errorf("Expected colorMode true when --color= flag is present")
	}
}

// TestParseInput_DifferentColors checks several color names are all parsed
// correctly without any trimming or casing issues.
func TestParseInput_DifferentColors(t *testing.T) {
	colors := []string{"black", "red", "green", "yellow", "blue", "purple", "cyan", "white"}

	for _, c := range colors {
		defer setArgs([]string{"cmd", "--color=" + c, "test"})()
		color, _, _, _ := parseInput()
		if color != c {
			t.Errorf("Expected color '%s', got '%s'", c, color)
		}
	}
}