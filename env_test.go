package env

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	content := "TEST_KEY=TEST_VALUE"
	tmpFile, err := os.CreateTemp("", ".env")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	os.Unsetenv("TEST_KEY")
	defer os.Unsetenv("TEST_KEY")

	if err := Load(tmpFile.Name()); err != nil {
		t.Errorf("Load() error = %v.", err)
	}

	if val := os.Getenv("TEST_KEY"); val != "TEST_VALUE" {
		t.Errorf("got %s, want TEST_VALUE.", val)
	}
}

func TestOverload(t *testing.T) {
	os.Setenv("TEST_KEY", "OLD_VALUE")
	defer os.Unsetenv("TEST_KEY")

	content := "TEST_KEY=NEW_VALUE"
	tmpFile, err := os.CreateTemp("", ".env")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	// Load should not overwrite an already-set value.
	if err := Load(tmpFile.Name()); err != nil {
		t.Errorf("Load() error = %v.", err)
	}

	if val := os.Getenv("TEST_KEY"); val != "OLD_VALUE" {
		t.Errorf("Load() overwrote: got %s, want OLD_VALUE.", val)
	}

	// Overload should overwrite the existing value.
	if err := Overload(tmpFile.Name()); err != nil {
		t.Errorf("Overload() error = %v.", err)
	}

	if val := os.Getenv("TEST_KEY"); val != "NEW_VALUE" {
		t.Errorf("Overload() failed: got %s, want NEW_VALUE.", val)
	}
}

func TestRead(t *testing.T) {
	content := "TEST_READ=READ_VALUE"
	tmpFile, err := os.CreateTemp("", ".env")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	got, err := Read(tmpFile.Name())

	if err != nil {
		t.Errorf("Read() error = %v.", err)
	}

	if val, ok := got["TEST_READ"]; !ok || val != "READ_VALUE" {
		t.Errorf("got %v, want map with TEST_READ=READ_VALUE.", got)
	}
}
