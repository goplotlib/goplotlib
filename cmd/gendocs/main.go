// cmd/gendocs processes docs/src/ → docs/content/ and generates extra SVGs.
//
// Each .md file in docs/src/ is copied to docs/content/ with two transforms:
//   - Code fences tagged  ```go run:filename.svg  are executed via go run;
//     the generated SVG lands in docs/static/img/.
//   - Lines between // @nodoc and // @doc (inclusive) are stripped from the
//     displayed code so readers see clean examples without the SaveSVG call.
//
// Run from the repository root:
//
//	go run ./cmd/gendocs
//	go run ./cmd/gendocs --check   # CI: exit non-zero if output is stale
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	srcDir  = "docs/src"
	dstDir  = "docs/content"
	runDir  = "cmd/gendocs/_run"
	imgDir  = "docs/static/img"
)

func main() {
	check := flag.Bool("check", false, "verify generated files are up to date without writing")
	flag.Parse()

	if err := os.MkdirAll(imgDir, 0755); err != nil {
		fatalf("creating img dir: %v", err)
	}

	// Process docs/src/ → docs/content/
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		rel, _ := filepath.Rel(srcDir, path)
		outPath := filepath.Join(dstDir, rel)
		return processFile(path, outPath, *check)
	})
	if err != nil {
		fatalf("%v", err)
	}

	// Generate showcase, theme, and utility SVGs not tied to doc pages.
	generateExtras()

	fmt.Println("Done.")
}

// processFile reads a source .md, runs any run: fences, strips @nodoc regions,
// and writes the result to dst. In check mode it only compares.
func processFile(src, dst string, check bool) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	out, err := processMarkdown(in)
	if err != nil {
		return fmt.Errorf("%s: %w", src, err)
	}
	if check {
		existing, _ := os.ReadFile(dst)
		if !bytes.Equal(existing, out) {
			return fmt.Errorf("%s is out of date; run go run ./cmd/gendocs", dst)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	fmt.Println("wrote", dst)
	return os.WriteFile(dst, out, 0644)
}

// processMarkdown transforms src markdown, executing run: fences and stripping
// @nodoc regions from the displayed code.
func processMarkdown(src []byte) ([]byte, error) {
	var out bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(src))

	inFence := false
	var fenceCode strings.Builder
	runFile := ""

	for scanner.Scan() {
		line := scanner.Text()

		if !inFence {
			if after, ok := strings.CutPrefix(line, "```go run:"); ok {
				runFile = strings.TrimSpace(after)
				inFence = true
				fenceCode.Reset()
			} else {
				out.WriteString(line)
				out.WriteByte('\n')
			}
		} else {
			if line == "```" {
				// End of fence: run the snippet, emit stripped display version.
				code := fenceCode.String()
				if err := runSnippet(code, runFile); err != nil {
					return nil, fmt.Errorf("running %s: %w", runFile, err)
				}
				display := stripNodoc(code)
				if strings.TrimSpace(display) != "" {
					out.WriteString("```go\n")
					out.WriteString(display)
					out.WriteString("```\n")
				}
				inFence = false
				runFile = ""
			} else {
				fenceCode.WriteString(line)
				fenceCode.WriteByte('\n')
			}
		}
	}
	return out.Bytes(), scanner.Err()
}

// runSnippet writes code to a temp file in cmd/gendocs/_run/ and executes it
// with go run from the repository root so module imports resolve correctly.
func runSnippet(code, svgName string) error {
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return err
	}
	tmpFile := filepath.Join(runDir, "main.go")
	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		return err
	}
	cmd := exec.Command("go", "run", filepath.Join(runDir, "main.go"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go run failed for %s: %w", svgName, err)
	}
	fmt.Println("  generated", svgName)
	return nil
}

// stripNodoc removes lines between // @nodoc and // @doc (inclusive) so the
// displayed code is clean. If there is no // @doc the region extends to EOF.
func stripNodoc(code string) string {
	var out strings.Builder
	hiding := false
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "// @nodoc" {
			hiding = true
			continue
		}
		if trimmed == "// @doc" {
			hiding = false
			continue
		}
		if !hiding {
			// Avoid a trailing blank line produced by Split on a newline-terminated string.
			if i == len(lines)-1 && line == "" {
				continue
			}
			out.WriteString(line)
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "gendocs: "+format+"\n", args...)
	os.Exit(1)
}
