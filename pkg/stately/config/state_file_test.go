package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveManagedSection_RemovesSection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	content := "user-before\n# BEGIN MANAGED\nmanaged-content\n# END MANAGED\nuser-after\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := removeManagedSection(path, "# BEGIN MANAGED", "# END MANAGED"); err != nil {
		t.Fatalf("removeManagedSection returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	want := "user-before\nuser-after\n"
	if string(got) != want {
		t.Errorf("file content mismatch\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}

func TestRemoveManagedSection_DeletesFileIfOnlyManagedContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	content := "# BEGIN MANAGED\nmanaged-content\n# END MANAGED\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := removeManagedSection(path, "# BEGIN MANAGED", "# END MANAGED"); err != nil {
		t.Fatalf("removeManagedSection returned error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, but it still exists")
	}
}

func TestRemoveManagedSection_NoMarkersIsNoop(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	content := "user-content-only\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := removeManagedSection(path, "# BEGIN MANAGED", "# END MANAGED"); err != nil {
		t.Fatalf("removeManagedSection returned error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(got) != content {
		t.Errorf("file content should be unchanged\ngot:\n%q\nwant:\n%q", string(got), content)
	}
}
