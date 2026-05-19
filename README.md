# env

A Go package to load environment variables from `.env` files.

## Installation

```bash
go get github.com/fahmiidris/env
```

## Usage

Create a `.env` file in the root of your project:

```env
SECRET_KEY=YOURSECRETKEYVALUE
```

### Loading variables into the environment

Use `env.Load()` to load variables from a `.env` file. By default, it looks for a file named `.env` in the current directory. This will **not** overwrite variables that are already set in the environment.

```go
package main

import (
	"log"
	"os"

	"github.com/fahmiidris/env"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")

	// Now do something with s3Bucket and secretKey
}
```

If you want to overwrite existing environment variables, use `env.Overload()`:

```go
err := env.Overload()
```

You can also specify multiple files:

```go
err := env.Load(".env", ".env.local")
```

### Reading variables without setting them

If you just want to read the variables into a map without modifying the environment, use `env.Read()`:

```go
myEnv, err := env.Read()
if err != nil {
	log.Fatal(err)
}

s3Bucket := myEnv["S3_BUCKET"]
```

### Features

- **Standard format**: `KEY=VALUE`
- **Comments**: Lines starting with `#` and inline comments are supported.
- **Exports**: Supports `export KEY=VALUE` syntax.
- **Quoting**: 
  - Single quotes (`'`) for literal values.
  - Double quotes (`"`) for values with escape sequences (e.g., `\n`, `\t`).
- **Multiple Files**: Load from multiple files sequentially.
