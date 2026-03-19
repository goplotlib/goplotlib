# Contributing

## Build

```bash
go build ./...
```

## Test

```bash
go test ./...
```

Golden SVG tests live in `plot/testdata/golden/`. To regenerate all golden files after a rendering change:

```bash
go test ./plot/ -update
```

To add new golden files without overwriting existing ones:

```bash
go test ./plot/ -update-new
```

## Run the example

```bash
go run ./examples/basic/
```

Output SVGs are written to `/tmp/`.

## Docs

The documentation site is a [Hugo](https://gohugo.io) site under `docs/`. Source pages live in `docs/src/` and are pre-processed by `cmd/gendocs`, which executes all ` ```go run ``` ` code blocks and embeds the resulting SVGs.

**Regenerate docs content and chart images:**

```bash
go run ./cmd/gendocs
```

This writes processed Markdown to `docs/content/` and SVGs to `docs/static/img/`.

**Serve the docs locally:**

```bash
cd docs && hugo server
```

Then open [http://localhost:1313](http://localhost:1313).

**Build the static site:**

```bash
cd docs && hugo --gc --minify
```

Output is written to `docs/public/`. The site deploys automatically to GitHub Pages on every push to `main`.
