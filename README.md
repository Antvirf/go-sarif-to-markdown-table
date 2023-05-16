# SARIF to Markdown Table

Utility to convert SARIF files to Markdown tables for use in CI/CD pipelines.

## Usage

Assuming a built binary exists (see [bin](./bin/) folder) (and is made executable with `chmod +x sarif-to-markdown-table`):

```bash
cat results.sarif | sarif-to-markdown-table > results.md
```

## Building the executable

```bash
make builds        # build all
make build-mac     # mac/arm64
make build-linux   # linux/amd64
make build-windows # windows/amd64
```
