package plot

import (
	"math"
	"os"

	"github.com/goplotlib/goplotlib/render"
	"github.com/goplotlib/goplotlib/render/svg"
	"github.com/goplotlib/goplotlib/theme"
)

// Auto is a sentinel value for SetXLim/SetYLim meaning "compute this bound from the data".
var Auto = math.NaN()

// FigureOption configures a Figure.
type FigureOption func(*Figure)

// WithWidth sets the figure width in pixels.
func WithWidth(w float64) FigureOption { return func(f *Figure) { f.width = w } }

// WithHeight sets the figure height in pixels.
func WithHeight(h float64) FigureOption { return func(f *Figure) { f.height = h } }

// WithTheme sets the figure theme.
func WithTheme(t theme.Theme) FigureOption { return func(f *Figure) { f.theme = t } }

// WithTitle sets the figure title.
func WithTitle(title string) FigureOption { return func(f *Figure) { f.title = title } }

// WithSpacing sets the gap in pixels between subplot cells (default 10).
func WithSpacing(px float64) FigureOption { return func(f *Figure) { f.spacing = px } }

// AxesOption configures an Axes at creation time (e.g. its grid position).
type AxesOption func(*Axes)

// At positions the axes at a specific row and column in the subplot grid (0-indexed).
func At(row, col int) AxesOption {
	return func(ax *Axes) { ax.row = row; ax.col = col }
}

// Figure is the top-level container for one or more Axes.
type Figure struct {
	width, height float64
	theme         theme.Theme
	axes          []*Axes
	title         string
	rows, cols    int     // subplot grid dimensions (default 1×1)
	spacing       float64 // gap between subplot cells in pixels
}

// New creates a new Figure with the given options.
// Default: 900×550, Light theme, 1×1 layout.
func New(opts ...FigureOption) *Figure {
	f := &Figure{
		width:   900,
		height:  550,
		theme:   theme.Light,
		rows:    1,
		cols:    1,
		spacing: 10,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// SetLayout sets the subplot grid to rows×cols equal cells.
func (f *Figure) SetLayout(rows, cols int) *Figure {
	f.rows = rows
	f.cols = cols
	return f
}

// SubPlots sets the grid to rows×cols and returns a pre-allocated [rows][cols]*Axes grid.
func (f *Figure) SubPlots(rows, cols int) [][]*Axes {
	f.rows = rows
	f.cols = cols
	grid := make([][]*Axes, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]*Axes, cols)
		for c := 0; c < cols; c++ {
			grid[r][c] = f.AddAxes(At(r, c))
		}
	}
	return grid
}

// AddAxes adds a new Axes to the figure and returns it.
// Pass At(row, col) to position it in a subplot grid.
func (f *Figure) AddAxes(opts ...AxesOption) *Axes {
	ax := newAxes()
	for _, opt := range opts {
		opt(ax)
	}
	f.axes = append(f.axes, ax)
	return ax
}

// render lays out and draws all axes, returning the SVG canvas.
func (f *Figure) render() *svg.Canvas {
	canvas := svg.New(f.width, f.height)

	// Figure background
	canvas.DrawRect(0, 0, f.width, f.height, 0, 0, render.Style{
		Fill:        f.theme.Background,
		FillOpacity: 1.0,
	})

	// Assign colors to each axes' charts
	for _, ax := range f.axes {
		ax.assignColors(f.theme.Palette)
	}

	// Reserve space at the top for an optional figure-level title.
	topOffset := 0.0
	if f.title != "" {
		topOffset = 32.0
	}

	// Compute subplot cell dimensions.
	rows := f.rows
	cols := f.cols
	if rows < 1 {
		rows = 1
	}
	if cols < 1 {
		cols = 1
	}
	sp := f.spacing
	cellW := (f.width - sp*float64(cols-1)) / float64(cols)
	cellH := (f.height - topOffset - sp*float64(rows-1)) / float64(rows)

	// Draw each axes in its cell.
	for i, ax := range f.axes {
		r := ax.row
		c := ax.col
		cellX := float64(c) * (cellW + sp)
		cellY := topOffset + float64(r)*(cellH+sp)
		ax.render(canvas, cellX, cellY, cellW, cellH, f.theme, i+1)
	}

	// Figure title (drawn on top of everything).
	if f.title != "" {
		canvas.DrawText(f.width/2, 18, f.title, render.TextStyle{
			Color:      f.theme.TitleColor,
			FontSize:   f.theme.TitleFontSize + 2,
			FontFamily: f.theme.FontFamily,
			Anchor:     "middle",
			Baseline:   "middle",
			Bold:       true,
		})
	}

	return canvas
}

// SVG renders the figure and returns the SVG bytes.
func (f *Figure) SVG() []byte {
	return f.render().Bytes()
}

// SaveSVG renders the figure and saves it to the given file path.
func (f *Figure) SaveSVG(path string) error {
	data := f.SVG()
	return os.WriteFile(path, data, 0644)
}
