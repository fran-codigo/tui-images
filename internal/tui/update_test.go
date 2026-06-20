package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSelectModeShortcutFStartsFilePicker(t *testing.T) {
	m := InitialModel(75)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'F'}})
	got := updated.(Model)

	if got.mode != ModeFile {
		t.Fatalf("mode = %v, want %v", got.mode, ModeFile)
	}
	if got.state != StateFilePick {
		t.Fatalf("state = %v, want %v", got.state, StateFilePick)
	}
	if !got.filePicker.FileAllowed || got.filePicker.DirAllowed {
		t.Fatalf("file picker permissions incorrect: file=%v dir=%v", got.filePicker.FileAllowed, got.filePicker.DirAllowed)
	}
}

func TestSelectModeShortcutDStartsDirPicker(t *testing.T) {
	m := InitialModel(75)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
	got := updated.(Model)

	if got.mode != ModeDir {
		t.Fatalf("mode = %v, want %v", got.mode, ModeDir)
	}
	if got.state != StateDirPick {
		t.Fatalf("state = %v, want %v", got.state, StateDirPick)
	}
	if got.filePicker.FileAllowed || !got.filePicker.DirAllowed {
		t.Fatalf("dir picker permissions incorrect: file=%v dir=%v", got.filePicker.FileAllowed, got.filePicker.DirAllowed)
	}
}

func TestDoneEnterContinuesFromStartDirectory(t *testing.T) {
	m := InitialModel(75)
	m.state = StateDone
	m.inputPath = "some-image.jpg"
	m.outputPath = "compressed/some-image.jpg"
	m.filePicker.CurrentDirectory = "another-dir"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := updated.(Model)

	if got.state != StateSelectMode {
		t.Fatalf("state = %v, want %v", got.state, StateSelectMode)
	}
	if got.inputPath != "" || got.outputPath != "" {
		t.Fatalf("paths were not reset: input=%q output=%q", got.inputPath, got.outputPath)
	}
	if got.filePicker.CurrentDirectory != got.startDir {
		t.Fatalf("current dir = %q, want start dir %q", got.filePicker.CurrentDirectory, got.startDir)
	}
}
