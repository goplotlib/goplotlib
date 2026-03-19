package theme

import "github.com/goplotlib/goplotlib/color"

// Padding defines margins around the plot area.
type Padding struct {
	Top, Right, Bottom, Left float64
}

// Theme defines the visual style of a chart.
type Theme struct {
	Background     color.Color
	PlotBackground color.Color
	GridColor      color.Color
	GridWidth      float64
	GridDash       []float64 // nil = solid
	AxisColor      color.Color
	AxisWidth      float64
	TextColor      color.Color
	TitleColor     color.Color
	FontFamily     string
	FontSize       float64
	TitleFontSize  float64
	LabelFontSize  float64
	TickFontSize   float64
	LineWidth      float64
	SpineTop       bool
	SpineRight     bool
	SpineLeft      bool
	SpineBottom    bool
	Padding        Padding // margins around the plot area
	Palette        []color.Color
}

const defaultFontFamily = "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif"

// Classic is the default high-quality color palette used by Light, Dark, and Minimal themes.
var Classic = []color.Color{
	color.Parse("#4C72B0"), color.Parse("#DD8452"),
	color.Parse("#55A868"), color.Parse("#C44E52"),
	color.Parse("#8172B3"), color.Parse("#937860"),
	color.Parse("#DA8BC3"), color.Parse("#8C8C8C"),
	color.Parse("#CCB974"), color.Parse("#64B5CD"),
}

// FTEPalette is the color palette inspired by FiveThirtyEight.
var FTEPalette = []color.Color{
	color.Parse("#008fd5"), color.Parse("#fc4f30"),
	color.Parse("#e5ae38"), color.Parse("#6d904f"),
	color.Parse("#8b8b8b"), color.Parse("#810f7c"),
}

// Light is a clean, light theme (default).
var Light = Theme{
	Background:     color.Parse("#ffffff"),
	PlotBackground: color.Parse("#f7f7f7"),
	GridColor:      color.Parse("#e5e5e5"),
	GridWidth:      0.8,
	GridDash:       nil,
	AxisColor:      color.Parse("#888888"),
	AxisWidth:      1.2,
	TextColor:      color.Parse("#333333"),
	TitleColor:     color.Parse("#1a1a1a"),
	FontFamily:     defaultFontFamily,
	FontSize:       12,
	TitleFontSize:  16,
	LabelFontSize:  12,
	TickFontSize:   11,
	LineWidth:      2.5,
	SpineTop:       false,
	SpineRight:     false,
	SpineLeft:      true,
	SpineBottom:    true,
	Padding:        Padding{Top: 45, Right: 30, Bottom: 60, Left: 70},
	Palette:        Classic,
}

// Dark is a dark theme with deep blue tones.
var Dark = Theme{
	Background:     color.Parse("#1a1b2e"),
	PlotBackground: color.Parse("#22233a"),
	GridColor:      color.Parse("#2e3050"),
	GridWidth:      0.8,
	GridDash:       nil,
	AxisColor:      color.Parse("#5a5b7a"),
	AxisWidth:      1.0,
	TextColor:      color.Parse("#c8cae8"),
	TitleColor:     color.Parse("#e0e2ff"),
	FontFamily:     defaultFontFamily,
	FontSize:       12,
	TitleFontSize:  16,
	LabelFontSize:  12,
	TickFontSize:   11,
	LineWidth:      2.5,
	SpineTop:       false,
	SpineRight:     false,
	SpineLeft:      true,
	SpineBottom:    true,
	Padding:        Padding{Top: 45, Right: 30, Bottom: 60, Left: 70},
	Palette:        Classic,
}

// Minimal is a minimal, clean theme.
var Minimal = Theme{
	Background:     color.Parse("#ffffff"),
	PlotBackground: color.Parse("#ffffff"),
	GridColor:      color.Parse("#eeeeee"),
	GridWidth:      0.6,
	GridDash:       []float64{4, 4},
	AxisColor:      color.Parse("#aaaaaa"),
	AxisWidth:      0.8,
	TextColor:      color.Parse("#444444"),
	TitleColor:     color.Parse("#222222"),
	FontFamily:     defaultFontFamily,
	FontSize:       12,
	TitleFontSize:  16,
	LabelFontSize:  12,
	TickFontSize:   11,
	LineWidth:      2.5,
	SpineTop:       false,
	SpineRight:     false,
	SpineLeft:      false,
	SpineBottom:    true,
	Padding:        Padding{Top: 40, Right: 25, Bottom: 55, Left: 65},
	Palette:        Classic,
}

// FiveThirtyEight is inspired by the FiveThirtyEight data journalism style.
var FiveThirtyEight = Theme{
	Background:     color.Parse("#f0f0f0"),
	PlotBackground: color.Parse("#f0f0f0"),
	GridColor:      color.Parse("#ffffff"),
	GridWidth:      1.8,
	GridDash:       nil,
	AxisColor:      color.Parse("#999999"),
	AxisWidth:      0,
	TextColor:      color.Parse("#444444"),
	TitleColor:     color.Parse("#2e2e2e"),
	FontFamily:     defaultFontFamily,
	FontSize:       12,
	TitleFontSize:  18,
	LabelFontSize:  12,
	TickFontSize:   11,
	LineWidth:      3.0,
	SpineTop:       false,
	SpineRight:     false,
	SpineLeft:      false,
	SpineBottom:    false,
	Padding:        Padding{Top: 50, Right: 25, Bottom: 60, Left: 70},
	Palette:        FTEPalette,
}
