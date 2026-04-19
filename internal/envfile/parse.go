package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse reads a .env file and returns a map of key-value pairs.
// It supports:
//   - KEY=VALUE
//   - KEY="VALUE" or KEY='VALUE'
//   - Comments starting with #
//   - Blank lines
func Parse(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("envfile: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scan %q: %w", path, err)
	}

	return env, nil
}

func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", fmt.Errorf("invalid line %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(line[:idx])
	raw := strings.TrimSpace(line[idx+1:])

	value := stripQuotes(raw)
	return key, value, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
