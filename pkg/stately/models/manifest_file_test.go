package models

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeSectionRaw_NewFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "testfile")

	f := &ManifestFile{
		Install:      MergeSection,
		Content:      "managed-line-1\nmanaged-line-2\n",
		SectionStart: "# BEGIN MANAGED",
		SectionEnd:   "# END MANAGED",
	}

	if err := f.MergeSectionRaw(dest); err != nil {
		t.Fatalf("MergeSectionRaw returned error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	want := "# BEGIN MANAGED\nmanaged-line-1\nmanaged-line-2\n# END MANAGED\n"
	if string(got) != want {
		t.Errorf("file content mismatch\ngot:\n%s\nwant:\n%s", string(got), want)
	}
}

func TestMergeSectionRaw_ExistingFileNoMarkers(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "testfile")

	existing := "user-line-1\nuser-line-2\n"
	if err := os.WriteFile(dest, []byte(existing), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	f := &ManifestFile{
		Install:      MergeSection,
		Content:      "managed-line-1\n",
		SectionStart: "# BEGIN MANAGED",
		SectionEnd:   "# END MANAGED",
	}

	if err := f.MergeSectionRaw(dest); err != nil {
		t.Fatalf("MergeSectionRaw returned error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	want := "user-line-1\nuser-line-2\n\n# BEGIN MANAGED\nmanaged-line-1\n# END MANAGED\n"
	if string(got) != want {
		t.Errorf("file content mismatch\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}

func TestMergeSectionRaw_ExistingFileWithMarkers(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "testfile")

	existing := "user-before\n# BEGIN MANAGED\nold-managed\n# END MANAGED\nuser-after\n"
	if err := os.WriteFile(dest, []byte(existing), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	f := &ManifestFile{
		Install:      MergeSection,
		Content:      "new-managed\n",
		SectionStart: "# BEGIN MANAGED",
		SectionEnd:   "# END MANAGED",
	}

	if err := f.MergeSectionRaw(dest); err != nil {
		t.Fatalf("MergeSectionRaw returned error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	want := "user-before\n# BEGIN MANAGED\nnew-managed\n# END MANAGED\nuser-after\n"
	if string(got) != want {
		t.Errorf("file content mismatch\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}

func TestMergeSectionRaw_UpdateManagedContent(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "testfile")

	// First write with initial managed content.
	existing := "header\n# BEGIN MANAGED\ninitial-content\n# END MANAGED\nfooter\n"
	if err := os.WriteFile(dest, []byte(existing), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	f := &ManifestFile{
		Install:      MergeSection,
		Content:      "updated-line-1\nupdated-line-2\n",
		SectionStart: "# BEGIN MANAGED",
		SectionEnd:   "# END MANAGED",
	}

	if err := f.MergeSectionRaw(dest); err != nil {
		t.Fatalf("MergeSectionRaw returned error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	want := "header\n# BEGIN MANAGED\nupdated-line-1\nupdated-line-2\n# END MANAGED\nfooter\n"
	if string(got) != want {
		t.Errorf("file content mismatch\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}
