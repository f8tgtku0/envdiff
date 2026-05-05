# envdiff

> Diff `.env` files across environments and flag missing or mismatched variables.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

Compare two `.env` files and see what's different:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
MISSING in .env.production:
  - DEBUG
  - SMTP_HOST

MISMATCHED values:
  - DATABASE_URL  (values differ across environments)
  - LOG_LEVEL     dev="debug"  prod="warn"

OK: 12 variables match across both files.
```

You can also compare multiple files at once:

```bash
envdiff .env.development .env.staging .env.production
```

Use the `--keys-only` flag to suppress value output and only report missing keys:

```bash
envdiff --keys-only .env.development .env.production
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Only report missing keys, skip value comparison |
| `--quiet` | Exit with non-zero status if differences found (useful in CI) |
| `--version` | Print version information |

---

## License

MIT © [yourusername](https://github.com/yourusername)