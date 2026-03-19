package color

import (
	"testing"
)

// approxEq checks float equality within a small epsilon.
func approxEq(a, b, eps float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= eps
}

// TestParseHex6 tests parsing of 6-digit hex colors.
func TestParseHex6(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#ff0000", Color{255, 0, 0, 255}},
		{"#00ff00", Color{0, 255, 0, 255}},
		{"#0000ff", Color{0, 0, 255, 255}},
		{"#4c72b0", Color{76, 114, 176, 255}},
		{"#000000", Color{0, 0, 0, 255}},
		{"#ffffff", Color{255, 255, 255, 255}},
	}
	for _, tt := range tests {
		got := Parse(tt.input)
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// TestParseHex3 tests expansion of 3-digit shorthand hex colors.
func TestParseHex3(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#f00", Color{255, 0, 0, 255}},
		{"#0f0", Color{0, 255, 0, 255}},
		{"#00f", Color{0, 0, 255, 255}},
		{"#fff", Color{255, 255, 255, 255}},
		{"#000", Color{0, 0, 0, 255}},
		{"#abc", Color{170, 187, 204, 255}},
	}
	for _, tt := range tests {
		got := Parse(tt.input)
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// TestParseHex8 tests parsing of 8-digit hex colors with alpha.
func TestParseHex8(t *testing.T) {
	tests := []struct {
		input string
		want  Color
	}{
		{"#ff000080", Color{255, 0, 0, 128}},
		{"#ffffff00", Color{255, 255, 255, 0}},
		{"#000000ff", Color{0, 0, 0, 255}},
	}
	for _, tt := range tests {
		got := Parse(tt.input)
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// TestParseRoundtrip tests that Parse(Hex()) returns the original color.
func TestParseRoundtrip(t *testing.T) {
	colors := []Color{
		{255, 0, 0, 255},
		{76, 114, 176, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{123, 45, 67, 255},
	}
	for _, c := range colors {
		got := Parse(c.Hex())
		if got != c {
			t.Errorf("roundtrip failed for %v: Parse(Hex()) = %v", c, got)
		}
	}
}

// TestParseNamedColors tests all named colors.
func TestParseNamedColors(t *testing.T) {
	tests := []struct {
		name string
		want Color
	}{
		{"white", White},
		{"black", Black},
		{"red", Red},
		{"transparent", Transparent},
		{"steelblue", SteelBlue},
		{"gray", Gray},
		{"grey", Gray},        // alias
		{"lightgray", LightGray},
		{"lightgrey", LightGray}, // alias
	}
	for _, tt := range tests {
		got := Parse(tt.name)
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// TestParseUnknown tests that unknown names return opaque black.
func TestParseUnknown(t *testing.T) {
	cases := []string{"notacolor", "", "rgb(1,2,3)", "hsl(0,0,0)"}
	for _, s := range cases {
		got := Parse(s)
		if got.A != 255 {
			t.Errorf("Parse(%q) should return opaque color, got alpha=%d", s, got.A)
		}
	}
}

// TestHex tests the Hex() output format.
func TestHex(t *testing.T) {
	tests := []struct {
		c    Color
		want string
	}{
		{Color{255, 0, 0, 255}, "#ff0000"},
		{Color{0, 255, 0, 255}, "#00ff00"},
		{Color{0, 0, 255, 255}, "#0000ff"},
		{Color{0, 0, 0, 255}, "#000000"},
		{Color{255, 255, 255, 255}, "#ffffff"},
		{Color{76, 114, 176, 255}, "#4c72b0"},
	}
	for _, tt := range tests {
		got := tt.c.Hex()
		if got != tt.want {
			t.Errorf("%v.Hex() = %q, want %q", tt.c, got, tt.want)
		}
	}
}

// TestHexA tests the HexA() output includes alpha.
func TestHexA(t *testing.T) {
	c := Color{255, 0, 0, 128}
	got := c.HexA()
	want := "#ff000080"
	if got != want {
		t.Errorf("HexA() = %q, want %q", got, want)
	}

	// Full opacity
	c2 := Color{0, 0, 0, 255}
	if got2 := c2.HexA(); got2 != "#000000ff" {
		t.Errorf("HexA() = %q, want #000000ff", got2)
	}
}

// TestWithAlpha tests alpha replacement.
func TestWithAlpha(t *testing.T) {
	c := Red

	half := c.WithAlpha(0.5)
	if half.R != 255 || half.G != 0 || half.B != 0 {
		t.Errorf("WithAlpha changed RGB components: %v", half)
	}
	if !approxEq(float64(half.A), 127.5, 1) {
		t.Errorf("WithAlpha(0.5) A = %d, want ~127", half.A)
	}

	zero := c.WithAlpha(0)
	if zero.A != 0 {
		t.Errorf("WithAlpha(0) A = %d, want 0", zero.A)
	}

	full := c.WithAlpha(1)
	if full.A != 255 {
		t.Errorf("WithAlpha(1) A = %d, want 255", full.A)
	}

	// Clamp above 1
	clamped := c.WithAlpha(1.5)
	if clamped.A != 255 {
		t.Errorf("WithAlpha(1.5) should clamp to 255, got %d", clamped.A)
	}

	// Clamp below 0
	clampedZero := c.WithAlpha(-0.5)
	if clampedZero.A != 0 {
		t.Errorf("WithAlpha(-0.5) should clamp to 0, got %d", clampedZero.A)
	}
}

// TestLighten tests mixing toward white.
func TestLighten(t *testing.T) {
	black := Black

	// Lighten by 0 = no change
	if got := black.Lighten(0); got != black {
		t.Errorf("Lighten(0) should return original color")
	}

	// Lighten by 1 = white
	white := black.Lighten(1)
	if white.R != 255 || white.G != 255 || white.B != 255 {
		t.Errorf("Lighten(1) should return white, got %v", white)
	}

	// Lighten by 0.5 = midpoint
	mid := black.Lighten(0.5)
	if !approxEq(float64(mid.R), 127.5, 1) {
		t.Errorf("Lighten(0.5) R = %d, want ~127", mid.R)
	}

	// Alpha should be preserved
	c := Color{100, 100, 100, 128}
	if got := c.Lighten(0.5); got.A != 128 {
		t.Errorf("Lighten should preserve alpha, got %d", got.A)
	}
}

// TestDarken tests mixing toward black.
func TestDarken(t *testing.T) {
	white := White

	// Darken by 0 = no change
	if got := white.Darken(0); got != white {
		t.Errorf("Darken(0) should return original color")
	}

	// Darken by 1 = black (RGB only)
	dark := white.Darken(1)
	if dark.R != 0 || dark.G != 0 || dark.B != 0 {
		t.Errorf("Darken(1) should return black RGB, got %v", dark)
	}

	// Darken by 0.5 = midpoint
	mid := white.Darken(0.5)
	if !approxEq(float64(mid.R), 127.5, 1) {
		t.Errorf("Darken(0.5) R = %d, want ~127", mid.R)
	}

	// Alpha should be preserved
	c := Color{200, 200, 200, 64}
	if got := c.Darken(0.5); got.A != 64 {
		t.Errorf("Darken should preserve alpha, got %d", got.A)
	}
}

// TestIsTransparent tests the IsTransparent predicate.
func TestIsTransparent(t *testing.T) {
	if !Transparent.IsTransparent() {
		t.Error("Transparent should be transparent")
	}
	if Black.IsTransparent() {
		t.Error("Black should not be transparent")
	}
	if Red.WithAlpha(0).IsTransparent() == false {
		t.Error("WithAlpha(0) result should be transparent")
	}
}

// TestSVGColor tests the SVGColor output for each alpha case.
func TestSVGColor(t *testing.T) {
	// Fully opaque: use hex
	got := Red.SVGColor()
	if got != "#ff0000" {
		t.Errorf("SVGColor() for opaque = %q, want #ff0000", got)
	}

	// Fully transparent: use "none"
	got = Transparent.SVGColor()
	if got != "none" {
		t.Errorf("SVGColor() for transparent = %q, want none", got)
	}

	// Partial alpha: use rgba()
	c := Color{255, 0, 0, 128}
	got = c.SVGColor()
	if got == "#ff0000" || got == "none" {
		t.Errorf("SVGColor() for partial alpha should be rgba(), got %q", got)
	}
}
