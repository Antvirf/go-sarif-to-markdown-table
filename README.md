# SARIF to Markdown Table

Utility to convert SARIF files to Markdown tables for use in CI/CD pipelines.

## Usage

Assuming a built binary exists (and is made executable with `chmod +x sarif-to-markdown-table`):

```bash
cat results.sarif | sarif-to-markdown-table > results.md
```
