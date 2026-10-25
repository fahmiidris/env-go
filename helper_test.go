package env

import "testing"

func TestGetEnv(t *testing.T) {
	t.Setenv("GETENV_KEY", "GETENV_VALUE")
	t.Setenv("GETENV_EMPTY", "")

	if got := GetEnv("GETENV_KEY"); got != "GETENV_VALUE" {
		t.Errorf("GetEnv() = %q, want %q.", got, "GETENV_VALUE")
	}

	if got := GetEnv("GETENV_EMPTY"); got != "" {
		t.Errorf("GetEnv() = %q, want empty string.", got)
	}

	if got := GetEnv("GETENV_UNSET"); got != "" {
		t.Errorf("GetEnv() = %q, want empty string.", got)
	}
}

func TestEnvValueString(t *testing.T) {
	t.Setenv("ENVVALUE_SET", "ENVVALUE_SET_VALUE")

	if got := GetEnv("ENVVALUE_SET").String(); got != "ENVVALUE_SET_VALUE" {
		t.Errorf("GetEnv().String() = %q, want %q.", got, "ENVVALUE_SET_VALUE")
	}

	if got := GetEnv("ENVVALUE_UNSET").String(); got != "" {
		t.Errorf("GetEnv().String() = %q, want empty string.", got)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Setenv("GETENVDEF_KEY", "GETENVDEF_VALUE")
	t.Setenv("GETENVDEF_EMPTY", "")

	if got := GetEnv("GETENVDEF_KEY").OrDefault("DEFAULT"); got != "GETENVDEF_VALUE" {
		t.Errorf("GetEnv().OrDefault() = %q, want %q.", got, "GETENVDEF_VALUE")
	}

	if got := GetEnv("GETENVDEF_EMPTY").OrDefault("DEFAULT"); got != "DEFAULT" {
		t.Errorf("GetEnv().OrDefault() = %q, want %q.", got, "DEFAULT")
	}

	if got := GetEnv("GETENVDEF_UNSET").OrDefault("DEFAULT"); got != "DEFAULT" {
		t.Errorf("GetEnv().OrDefault() = %q, want %q.", got, "DEFAULT")
	}

	if got := GetEnvOrDefault("GETENVDEF_UNSET", "DEFAULT"); got != "DEFAULT" {
		t.Errorf("GetEnvOrDefault() = %q, want %q.", got, "DEFAULT")
	}
}
