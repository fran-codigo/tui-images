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

// CompressDirMsg signals directory compression is done
type CompressDirMsg struct {
	TotalSrcSize int64
	TotalDstSize int64
	Errors       []string
}

// CompressSingleMsg signals single file compression is done
type CompressSingleMsg struct {
	SrcSize int64
	DstSize int64
	Err     error
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// runCompression runs the compression and sends progress updates
func runCompression(inputPath, outputPath string, quality int, mode Mode) tea.Cmd {
	return func() tea.Msg {
		if mode == ModeFile {
			srcSize := getFileSize(inputPath)
			err := compressor.CompressImage(inputPath, outputPath, quality)
			if err != nil {
				return CompressSingleMsg{Err: err}
			}
			outputDir := filepath.Dir(outputPath)
			relPath := filepath.Base(outputPath)
			dstSize, _ := compressor.FindOutputFile(outputDir, relPath)
			return CompressSingleMsg{
				SrcSize: srcSize,
				DstSize: dstSize,
			}
		}

		// Directory mode
		var totalSrc, totalDst int64
		var errs []string

		ch := make(chan compressor.Progress)
		go func() {
			compressor.CompressDirectory(inputPath, outputPath, quality, ch)
		}()

		for p := range ch {
			if p.Error != "" {
				errs = append(errs, p.CurrentFile+": "+p.Error)
			}
			totalSrc += p.SrcSize
			totalDst += p.DstSize
		}

		return CompressDirMsg{
			TotalSrcSize: totalSrc,
			TotalDstSize: totalDst,
			Errors:       errs,
		}
	}
}
