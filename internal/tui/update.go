package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle quality input state
	if m.state == StateQualityInput {
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "esc":
				return m.handleEsc()
			case "enter":
				return m.handleEnter()
			case "backspace":
				if len(m.qualityInput) > 0 {
					m.qualityInput = m.qualityInput[:len(m.qualityInput)-1]
				}
				return m, nil
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if len(m.qualityInput) < 3 {
					m.qualityInput += km.String()
				}
				return m, nil
			}
			return m, nil
		}
	}

	// Always update filepicker when in picker state
	if m.state == StateFilePick || m.state == StateDirPick {
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Check if user selected a path
		if m.filePicker.Path != "" && m.inputPath == "" {
			m.inputPath = m.filePicker.Path
			m.state = StateQualityInput
			m.qualityInput = strconv.Itoa(m.quality)
			return m, tea.Batch(cmds...)
		}

		// Handle keys for picker state
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = StateSelectMode
				m.filePicker.Path = ""
				return m, nil
			}
		}

		return m, tea.Batch(cmds...)
	}

	// Handle other states
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.handleEsc()
		case "enter":
			return m.handleEnter()
		case "tab":
			if m.state == StateSelectMode {
				if m.mode == ModeFile {
					m.mode = ModeDir
				} else {
					m.mode = ModeFile
				}
			}
		}

	case tea.WindowSizeMsg:
		h := msg.Height - 10
		if h < 5 {
			h = 5
		}
		m.filePicker.Height = h

	case CompressSingleMsg:
		if msg.Err != nil {
			m.state = StateError
			m.err = msg.Err.Error()
		} else {
			m.totalSrcSize = msg.SrcSize
			m.totalDstSize = msg.DstSize
			m.state = StateDone
		}

	case CompressDirMsg:
		m.totalSrcSize = msg.TotalSrcSize
		m.totalDstSize = msg.TotalDstSize
		m.errors = msg.Errors
		m.state = StateDone
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleEsc() (Model, tea.Cmd) {
	switch m.state {
	case StateQualityInput:
		m.inputPath = ""
		m.filePicker.Path = ""
		if m.mode == ModeFile {
			m.state = StateFilePick
		} else {
			m.state = StateDirPick
		}
		return m, m.filePicker.Init()
	case StateFilePick, StateDirPick:
		m.state = StateSelectMode
	case StateSelectMode, StateDone, StateError:
		return m, tea.Quit
	default:
		m.state = StateSelectMode
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.state {
	case StateSelectMode:
		if m.mode == ModeFile {
			m.state = StateFilePick
			m.filePicker.FileAllowed = true
			m.filePicker.DirAllowed = false
		} else {
			m.state = StateDirPick
			m.filePicker.FileAllowed = false
			m.filePicker.DirAllowed = true
		}
		return m, m.filePicker.Init()

	case StateQualityInput:
		q, err := strconv.Atoi(m.qualityInput)
		if err != nil || q < 1 || q > 100 {
			return m, nil
		}
		m.quality = q

		if _, err := os.Stat(m.inputPath); err != nil {
			m.state = StateError
			m.err = fmt.Sprintf("Path not found: %s", m.inputPath)
			return m, nil
		}

		if m.mode == ModeFile {
			dir := filepath.Dir(m.inputPath)
			name := filepath.Base(m.inputPath)
			ext := filepath.Ext(name)
			base := strings.TrimSuffix(name, ext)
			m.outputPath = filepath.Join(dir, "compressed", base+ext)
		} else {
			dir := filepath.Dir(m.inputPath)
			name := filepath.Base(m.inputPath)
			m.outputPath = filepath.Join(dir, name+"_compressed")
		}

		m.state = StateProcessing
		return m, runCompression(m.inputPath, m.outputPath, m.quality, m.mode)

	default:
		return m, nil
	}
}

// Run starts the TUI
func Run(quality int) error {
	m := InitialModel(quality)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("running TUI: %w", err)
	}

	fm := finalModel.(Model)
	if fm.err != "" {
		return fmt.Errorf("compression failed: %s", fm.err)
	}
	return nil
}
