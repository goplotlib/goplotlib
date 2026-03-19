// cmd/testview generates a self-contained HTML page for visual review of golden test outputs.
// Run with: go run ./cmd/testview > /tmp/review.html && open /tmp/review.html
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const goldenDir = "plot/testdata/golden"

func main() {
	goldens, err := filepath.Glob(filepath.Join(goldenDir, "*.svg"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "glob: %v\n", err)
		os.Exit(1)
	}
	actuals, err := filepath.Glob(filepath.Join(goldenDir, "*.svg.actual"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "glob: %v\n", err)
		os.Exit(1)
	}

	type chart struct {
		name   string
		golden string
		actual string // empty if no diff
	}
	byName := map[string]*chart{}
	for _, path := range goldens {
		name := strings.TrimSuffix(filepath.Base(path), ".svg")
		if byName[name] == nil {
			byName[name] = &chart{name: name}
		}
		byName[name].golden = path
	}
	for _, path := range actuals {
		name := strings.TrimSuffix(strings.TrimSuffix(filepath.Base(path), ".actual"), ".svg")
		if byName[name] == nil {
			byName[name] = &chart{name: name}
		}
		byName[name].actual = path
	}

	// Sort: charts with diffs first, then alphabetical within each group.
	names := make([]string, 0, len(byName))
	for n := range byName {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		ai := byName[names[i]].actual != ""
		aj := byName[names[j]].actual != ""
		if ai != aj {
			return ai // diffs first
		}
		return names[i] < names[j]
	})

	hasDiff := false
	for _, c := range byName {
		if c.actual != "" {
			hasDiff = true
			break
		}
	}

	fmt.Print(htmlHeader(hasDiff))

	for _, name := range names {
		c := byName[name]
		goldenSVG := readFile(c.golden)
		actualSVG := ""
		if c.actual != "" {
			actualSVG = readFile(c.actual)
		}
		printChartRow(name, goldenSVG, actualSVG)
	}

	fmt.Print(htmlFooter)
}

func readFile(path string) string {
	if path == "" {
		return ""
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("<!-- error reading %s: %v -->", path, err)
	}
	return string(b)
}

// namespaceSVG rewrites all id="..." and url(#...) references in an SVG so
// they are unique within the HTML document when multiple SVGs are inlined.
func namespaceSVG(svg, prefix string) string {
	// Replace id="X" → id="prefix-X"
	svg = strings.ReplaceAll(svg, `id="`, `id="`+prefix+`-`)
	// Replace url(#X) → url(#prefix-X)
	svg = strings.ReplaceAll(svg, `url(#`, `url(#`+prefix+`-`)
	return svg
}

func printChartRow(name, goldenSVG, actualSVG string) {
	diffClass := ""
	if actualSVG != "" {
		diffClass = " has-diff"
	}
	fmt.Printf(`<div class="chart-row%s">`, diffClass)
	fmt.Printf(`<h3>%s</h3>`, name)
	fmt.Print(`<div class="pair">`)

	nsGolden := namespaceSVG(goldenSVG, name+"-g")
	fmt.Print(`<div class="golden"><p>Golden</p>`)
	fmt.Print(nsGolden)
	fmt.Print(`</div>`)

	if actualSVG != "" {
		nsActual := namespaceSVG(actualSVG, name+"-a")
		fmt.Print(`<div class="actual"><p>Actual (differs)</p>`)
		fmt.Print(nsActual)
		fmt.Print(`</div>`)

		fmt.Print(`<div class="diff"><p>Diff</p>`)
		fmt.Print(`<div class="diff-overlay">`)
		fmt.Print(`<div class="diff-base">`)
		fmt.Print(namespaceSVG(goldenSVG, name+"-dg"))
		fmt.Print(`</div>`)
		fmt.Print(`<div class="diff-top">`)
		fmt.Print(namespaceSVG(actualSVG, name+"-da"))
		fmt.Print(`</div>`)
		fmt.Print(`</div></div>`)
	}

	fmt.Print(`</div></div>`)
	fmt.Println()
}

func htmlHeader(hasDiff bool) string {
	banner := ""
	if hasDiff {
		banner = `<div class="banner">⚠ Diffs detected — charts highlighted in red have changed</div>`
	}
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>goplotlib golden review</title>
<style>
  body { font-family: system-ui, sans-serif; margin: 0; padding: 16px; background: #f5f5f5; }
  h1 { margin-bottom: 4px; }
  .banner { background: #ffeeba; border: 1px solid #ffc107; padding: 10px 16px; margin-bottom: 16px; border-radius: 4px; font-weight: bold; }
  .chart-row { background: #fff; border: 1px solid #ddd; border-radius: 6px; margin-bottom: 20px; padding: 12px 16px; }
  .chart-row.has-diff { border-color: #e74c3c; background: #fff8f8; }
  .chart-row h3 { margin: 0 0 10px; font-size: 14px; color: #333; font-family: monospace; }
  .pair { display: flex; gap: 16px; flex-wrap: wrap; }
  .golden p, .actual p { margin: 0 0 4px; font-size: 12px; color: #666; }
  .actual p { color: #e74c3c; font-weight: bold; }
  .golden svg, .actual svg { display: block; max-width: 100%; height: auto; border: 1px solid #eee; }
  .diff-overlay { position: relative; display: inline-block; border: 1px solid #eee; background: #000; }
  .diff-base { display: block; }
  .diff-base svg { display: block; }
  .diff-top { position: absolute; top: 0; left: 0; mix-blend-mode: difference; }
  .diff-top svg { display: block; }
  .diff p { margin: 0 0 4px; font-size: 12px; color: #666; }
</style>
</head>
<body>
<h1>goplotlib golden review</h1>
` + banner + `
<div class="charts">
`
}

const htmlFooter = `</div>
</body>
</html>
`
