package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			MarginBottom(1)

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#555555")).
			Padding(0, 2).
			MarginRight(1)

	optionSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#04B575")).
				Padding(0, 2).
				MarginRight(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			MarginTop(1)

	fileInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BBBBBB"))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(2)
)

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("  Image Compressor  "))
	s.WriteString("\n\n")

	switch m.state {
	case StateSelectMode:
		s.WriteString(subtitleStyle.Render("Select compression mode:"))
		s.WriteString("\n\n")

		fileOpt := optionStyle.Render("[F] File")
		dirOpt := optionStyle.Render("[D] Directory")
		if m.mode == ModeFile {
			fileOpt = optionSelectedStyle.Render("[F] File")
		} else {
			dirOpt = optionSelectedStyle.Render("[D] Directory")
		}
		s.WriteString(fileOpt + dirOpt)
		s.WriteString("\n\n")
		s.WriteString(promptStyle.Render("Press Enter to confirm, Tab to switch, Q to quit"))

	case StateFilePick, StateDirPick:
		label := "Select an image file"
		if m.state == StateDirPick {
			label = "Select a directory"
		}
		s.WriteString(subtitleStyle.Render(label))
		s.WriteString("\n")
		s.WriteString(m.filePicker.View())
		s.WriteString(helpStyle.Render("↑↓ navigate  Enter select  Esc back  Q quit"))

	case StateQualityInput:
		s.WriteString(subtitleStyle.Render("Set compression quality (1-100):"))
		s.WriteString("\n\n")
		s.WriteString(promptStyle.Render("Quality: "))
		s.WriteString(inputStyle.Render(m.qualityInput + "█"))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("Lower = smaller file, less quality. Default: 75"))
		s.WriteString(helpStyle.Render("Enter confirm  Esc back  Q quit"))

	case StateProcessing:
		s.WriteString(subtitleStyle.Render("Compressing images..."))
		s.WriteString("\n\n")

		if m.progTotal > 0 {
			pct := float64(m.progCurrent) / float64(m.progTotal)
			s.WriteString(m.progress.ViewAs(pct))
			s.WriteString("\n")
			s.WriteString(fileInfoStyle.Render(
				fmt.Sprintf("  %d/%d  %s", m.progCurrent, m.progTotal, m.currentFile),
			))
		} else {
			s.WriteString(fileInfoStyle.Render("  Processing..."))
		}

	case StateDone:
		s.WriteString(successStyle.Render("  Compression Complete!  "))
		s.WriteString("\n\n")

		if m.mode == ModeFile {
			saved := float64(m.totalSrcSize - m.totalDstSize)
			pct := float64(0)
			if m.totalSrcSize > 0 {
				pct = saved / float64(m.totalSrcSize) * 100
			}
			s.WriteString(fmt.Sprintf("  Original:  %s\n", formatBytes(m.totalSrcSize)))
			s.WriteString(fmt.Sprintf("  Compressed: %s\n", formatBytes(m.totalDstSize)))
			s.WriteString(fmt.Sprintf("  Saved:      %s (%.1f%%)\n", formatBytes(int64(saved)), pct))
			s.WriteString(fmt.Sprintf("\n  Output: %s\n", m.outputPath))
		} else {
			saved := m.totalSrcSize - m.totalDstSize
			pct := float64(0)
			if m.totalSrcSize > 0 {
				pct = float64(saved) / float64(m.totalSrcSize) * 100
			}
			s.WriteString(fmt.Sprintf("  Files processed: %d\n", m.progTotal))
			s.WriteString(fmt.Sprintf("  Original size:   %s\n", formatBytes(m.totalSrcSize)))
			s.WriteString(fmt.Sprintf("  Compressed size: %s\n", formatBytes(m.totalDstSize)))
			s.WriteString(fmt.Sprintf("  Total saved:     %s (%.1f%%)\n", formatBytes(saved), pct))
			s.WriteString(fmt.Sprintf("\n  Output directory: %s\n", m.outputPath))
		}

		if len(m.errors) > 0 {
			s.WriteString("\n")
			s.WriteString(errorStyle.Render(fmt.Sprintf("  %d error(s):", len(m.errors))))
			s.WriteString("\n")
			for _, e := range m.errors {
				s.WriteString(fmt.Sprintf("    - %s\n", e))
			}
		}

		s.WriteString(helpStyle.Render("Press Q or Esc to exit"))

	case StateError:
		s.WriteString(errorStyle.Render("  Error  "))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("  %s\n", m.err))
		s.WriteString(helpStyle.Render("Press Esc to go back, Q to quit"))
	}

	s.WriteString("\n")
	return s.String()
}
