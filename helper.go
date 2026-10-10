package env

import "os"

// EnvValue is a string-like environment variable value.
type EnvValue string

// OrDefault returns the value, or defaultValue if it is empty.
func (value EnvValue) OrDefault(defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return string(value)
}

// GetEnv returns the value of the environment variable named by key as an
// EnvValue, or an empty value if it is not set.
func GetEnv(key string) EnvValue {
	return EnvValue(os.Getenv(key))
}

// GetEnvOrDefault returns the value of the environment variable named by
// key, or defaultValue if the variable is not set or is empty.
//
// It is equivalent to GetEnv(key).OrDefault(defaultValue).
func GetEnvOrDefault(key string, defaultValue string) string {
	return GetEnv(key).OrDefault(defaultValue)
}
