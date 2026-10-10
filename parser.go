package env

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"unicode"
)

// Parse reads dotenv content from an io.Reader and returns a map of
// key-value pairs. It supports blank lines and "#" comments, an "export "
// prefix, single and double quoted values, and $VAR/${VAR} expansion
// against previously parsed keys and the process environment.
func Parse(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comment-only lines.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip a leading "export " prefix.
		line = strings.TrimPrefix(line, "export ")

		// Find the first '='; lines without one are ignored.
		keyPart, valuePart, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key := strings.TrimSpace(keyPart)
		value := strings.TrimSpace(valuePart)

		// Handle the quoting style indicated by the first character.
		if len(value) > 0 {
			switch value[0] {
			case '\'':
				// Single-quoted values support neither escaping nor expansion.
				v, err := parseSingleQuoted(value)

				if err != nil {
					return nil, err
				}

				value = v

			case '"':
				// Double-quoted values support escaping and expansion.
				v, err := parseDoubleQuoted(value)

				if err != nil {
					return nil, err
				}

				value = expand(v, result)

			default:
				// Unquoted values support inline comments and expansion.
				value = expand(parseUnquoted(value), result)
			}
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func expand(s string, current map[string]string) string {
	return os.Expand(s, func(k string) string {
		if v, ok := current[k]; ok {
			return v
		}

		return os.Getenv(k)
	})
}

func parseSingleQuoted(s string) (string, error) {
	if len(s) < 2 || s[0] != '\'' {
		return "", errors.New("invalid single quoted string.")
	}

	content, _, found := strings.Cut(s[1:], "'")

	if !found {
		return "", errors.New("unclosed single quote.")
	}

	return content, nil
}

func parseDoubleQuoted(s string) (string, error) {
	var sb strings.Builder

	escaped := false

	for i := 1; i < len(s); i++ {
		c := s[i]

		if escaped {
			switch c {
			case 'n':
				sb.WriteByte('\n')

			case 'r':
				sb.WriteByte('\r')

			case 't':
				sb.WriteByte('\t')

			case '\\':
				sb.WriteByte('\\')

			case '"':
				sb.WriteByte('"')

			default:
				sb.WriteByte('\\')
				sb.WriteByte(c)

			}

			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			return sb.String(), nil
		}

		sb.WriteByte(c)
	}

	return "", errors.New("unclosed double quote.")
}

func parseUnquoted(s string) string {
	for i, c := range s {
		if c == '#' {
			if i > 0 && unicode.IsSpace(rune(s[i-1])) {
				return strings.TrimSpace(s[:i])
			}

			if i == 0 {
				return ""
			}

			return strings.TrimSpace(s[:i])
		}
	}

	return strings.TrimSpace(s)
}
