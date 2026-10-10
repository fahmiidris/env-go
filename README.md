# env

Load `.env` files into environment variables.

`env` is a small, dependency-free Go package that parses dotenv-style files
and applies them to the process environment.

## Features

- Parse `KEY=VALUE` lines, skipping blank lines and `#` comment lines
- Strip a leading `export ` prefix
- Strip inline comments from unquoted values (`KEY=VALUE # comment`)
- Single-quoted values (`'...'`) are taken literally
- Double-quoted values (`"..."`) support `\n`, `\r`, `\t`, `\\`, and `\"` escapes
- `$VAR` and `${VAR}` expansion against previously parsed keys, then the
  process environment
- Three file-level modes: `Load` (never overwrite), `Overload` (always
  overwrite), and `Read` (merge into a map without touching the environment)
- A missing default `.env` is silently ignored by `Load` and `Overload`
- `GetEnv` returns an `EnvValue` with an `OrDefault` fallback method;
  `GetEnvOrDefault` remains available for compatibility
- Standard library only

## Installation

Install the package with `go get`:

```sh
go get github.com/fahmiidris/env-go
```

Then import it in your code (the package name is `env`):

```go
import env "github.com/fahmiidris/env-go"
```

## Usage

Load `.env` into the process environment at startup:

```go
package main

import (
	"fmt"
	"log"
	"os"

	env "github.com/fahmiidris/env-go"
)

func main() {
	// Reads ".env" by default; pass specific files to override.
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	fmt.Println(databaseURL)

	// Read a value with.
	host := env.GetEnv("HOST")
	fmt.Println(host)

	// Read a value with a fallback when it is missing or empty.
	port := env.GetEnv("PORT").OrDefault("3000")
	fmt.Println(port)

	// Merge files into a map without touching the environment.
	vars, err := env.Read(".env", ".env.local")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(vars["DATABASE_URL"])
}
```

## API

| Function                                           | Description                                                                |
| -------------------------------------------------- | -------------------------------------------------------------------------- |
| `Load(files ...string) error`                      | Sets environment variables from files, skipping keys that are already set. |
| `Overload(files ...string) error`                  | Sets environment variables from files, overwriting existing values.        |
| `Read(files ...string) (map[string]string, error)` | Merges files into a map without touching the environment.                  |
| `Parse(r io.Reader) (map[string]string, error)`    | Parses dotenv content from any reader into a map.                          |
| `GetEnv(key string) EnvValue`                      | Returns an environment value, or an empty value if unset.                  |
| `EnvValue.OrDefault(defaultValue string) string`   | Returns the value, or `defaultValue` if empty.                             |
| `GetEnvOrDefault(key, defaultValue string) string` | Compatibility helper equivalent to `GetEnv(key).OrDefault(defaultValue)`.  |

An empty file list means `.env` for every file-level function. `Load` and
`Overload` ignore a missing default `.env`, while an explicitly named
missing file is an error.

## Parsing rules

- Lines are trimmed of surrounding whitespace before parsing.
- Lines starting with `#` are ignored; a `#` in an unquoted value starts an
  inline comment.
- Values are trimmed, then interpreted by their first character: single
  quotes are taken verbatim, double quotes are unescaped, and `$VAR` /
  `${VAR}` is expanded in double-quoted and unquoted values alike.
- Expansion resolves against keys parsed earlier in the same run, falling
  back to the process environment.
- Later definitions win: a key redefined later in a file, or in a later
  file, overrides the earlier value.
