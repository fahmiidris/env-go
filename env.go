package env

import (
	"maps"
	"os"
)

// Load reads the specified files and sets environment variables if they
// are not already set. When no filenames are given, ".env" is used; a
// missing default ".env" is silently ignored, but a missing
// explicitly-named file is an error.
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		if err := loadFile(filename, false); err != nil {
			return err
		}
	}

	return nil
}

// Overload reads the specified files and sets environment variables,
// overwriting any existing ones. When no filenames are given, ".env" is
// used.
func Overload(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		if err := loadFile(filename, true); err != nil {
			return err
		}
	}

	return nil
}

// Read reads the specified files and returns a map of key-value pairs,
// with later files taking precedence over earlier ones. When no filenames
// are given, ".env" is used.
func Read(filenames ...string) (map[string]string, error) {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	combined := make(map[string]string)

	for _, filename := range filenames {
		envMap, err := readFile(filename)

		if err != nil {
			return nil, err
		}

		maps.Copy(combined, envMap)
	}

	return combined, nil
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)

	if err != nil {
		if os.IsNotExist(err) && filename == ".env" {
			return nil // Ignore a missing default ".env" file.
		}

		return err
	}

	for key, value := range envMap {
		if !overload {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}

func readFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return Parse(file)
}
