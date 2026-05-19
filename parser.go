package env

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

// Parse reads from an io.Reader and returns a map of key-value pairs.
func Parse(r io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)

	result := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip "export " prefix
		line = strings.TrimPrefix(line, "export ")

		// Find the first '='
		keyPart, valuePart, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key := strings.TrimSpace(keyPart)
		value := strings.TrimSpace(valuePart)

		// Handle quoting
		if len(value) > 0 {
			switch value[0] {
			case '\'':
				// Single quoted: no escape sequences
				v, err := parseSingleQuoted(value)

				if err != nil {
					return nil, err
				}

				value = v
			case '"':
				// Double quoted: handle escapes
				v, err := parseDoubleQuoted(value)

				if err != nil {
					return nil, err
				}

				value = v
			default:
				// Unquoted: handle inline comments
				value = parseUnquoted(value)
			}
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseSingleQuoted(s string) (string, error) {
	if len(s) < 2 || s[0] != '\'' {
		return "", errors.New("invalid single quoted string")
	}

	content, _, found := strings.Cut(s[1:], "'")

	if !found {
		return "", errors.New("unclosed single quote")
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

	return "", errors.New("unclosed double quote")
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

			// If it's something like VAR=val#comment (no space),
			// some parsers treat it as part of value, others as comment.
			// Let's be consistent and treat any # as start of comment if not quoted.
			return strings.TrimSpace(s[:i])
		}
	}

	return strings.TrimSpace(s)
}
