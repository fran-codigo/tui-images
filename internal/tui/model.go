package tui

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fran-codigo/tui-images/internal/compressor"
)

type State int

const (
	StateSelectMode State = iota
	StateFilePick
	StateDirPick
	StateQualityInput
	StateProcessing
	StateDone
	StateError
)

type Mode int

const (
	ModeFile Mode = iota
	ModeDir
)

type Model struct {
	state        State
	mode         Mode
	quality      int
	qualityInput string

	filePicker  filepicker.Model
	fpConfirmed bool

	inputPath  string
	outputPath string

	progress     progress.Model
	progCurrent  int
	progTotal    int
	currentFile  string
	totalSrcSize int64
	totalDstSize int64
	errors       []string

	done bool
	err  string
}

func InitialModel(quality int) Model {
	fp := filepicker.New()
	fp.ShowPermissions = true
	fp.ShowSize = true
	fp.AutoHeight = true

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		state:        StateSelectMode,
		quality:      quality,
		qualityInput: "",
		filePicker:   fp,
		progress:     p,
	}
}

func (m Model) Init() tea.Cmd {
	return m.filePicker.Init()
}

// ProgressMsg carries a single directory-compression progress update.
type ProgressMsg struct {
	Progress compressor.Progress
	ch       <-chan compressor.Progress
}

// CompressDirMsg signals directory compression is done.
type CompressDirMsg struct{}

// CompressSingleMsg signals single file compression is done
type CompressSingleMsg struct {
	SrcSize    int64
	DstSize    int64
	OutputPath string
	Err        error
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func waitForProgress(ch <-chan compressor.Progress) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return CompressDirMsg{}
		}
		return ProgressMsg{Progress: p, ch: ch}
	}
}

// runSingleCompression compresses one file and reports the final sizes.
func runSingleCompression(inputPath, outputPath string, quality int) tea.Cmd {
	return func() tea.Msg {
		srcSize := getFileSize(inputPath)
		err := compressor.CompressImage(inputPath, outputPath, quality)
		if err != nil {
			return CompressSingleMsg{Err: err}
		}
		outputDir := filepath.Dir(outputPath)
		relPath := filepath.Base(outputPath)
		actualPath, dstSize, _ := compressor.FindOutputPath(outputDir, relPath)
		return CompressSingleMsg{
			SrcSize:    srcSize,
			DstSize:    dstSize,
			OutputPath: actualPath,
		}
	}
}
