package compressor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const maxFileSize = 100 * 1024 * 1024 // 100MB

var SupportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func IsSupportedImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return SupportedExtensions[ext]
}

func CompressImage(inputPath, outputPath string, quality int) error {
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	if info.Size() > maxFileSize {
		return fmt.Errorf("file too large: %s (max %dMB)", inputPath, maxFileSize/1024/1024)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symlinks not supported: %s", inputPath)
	}

	srcFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("opening source: %w", err)
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("decoding image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	switch ext {
	case ".jpg", ".jpeg":
		outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer outFile.Close()
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return fmt.Errorf("encoding image: %w", err)
		}

	case ".png":
		var bufJpeg, bufPng bytes.Buffer

		jpeg.Encode(&bufJpeg, img, &jpeg.Options{Quality: quality})

		pngEnc := png.Encoder{CompressionLevel: png.BestCompression}
		pngEnc.Encode(&bufPng, img)

		var bestData []byte
		var bestExt string

		if bufJpeg.Len() < bufPng.Len() {
			bestData = bufJpeg.Bytes()
			bestExt = ".jpg"
		} else {
			bestData = bufPng.Bytes()
			bestExt = ".png"
		}

		finalPath := strings.TrimSuffix(outputPath, ext) + bestExt
		if err := os.WriteFile(finalPath, bestData, 0640); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}

	default:
		return fmt.Errorf("unsupported output format: %s", ext)
	}

	return nil
}

func CompressDirectory(inputDir, outputDir string, quality int, progress chan<- Progress) error {
	defer close(progress)

	absInput, err := filepath.Abs(inputDir)
	if err != nil {
		progress <- Progress{Current: 0, Total: 0, Error: fmt.Sprintf("invalid input path: %s", err), Done: true}
		return err
	}

	var files []string
	err = filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type()&os.ModeSymlink != 0 {
			progress <- Progress{Current: 0, Total: 0, CurrentFile: path, Error: "skipping symlink", Done: false}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !IsSupportedImage(path) {
			return nil
		}

		absPath, _ := filepath.Abs(path)
		if !strings.HasPrefix(absPath, absInput) {
			progress <- Progress{Current: 0, Total: 0, CurrentFile: path, Error: "file outside directory", Done: false}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.Size() > maxFileSize {
			progress <- Progress{Current: 0, Total: 0, CurrentFile: path, Error: "file too large", Done: false}
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	total := len(files)
	if total == 0 {
		progress <- Progress{Current: 0, Total: 0, Done: true}
		return nil
	}

	progress <- Progress{Current: 0, Total: total, CurrentFile: "Starting...", Done: false}

	for i, file := range files {
		relPath, err := filepath.Rel(inputDir, file)
		if err != nil {
			progress <- Progress{Current: i, Total: total, CurrentFile: file, Error: err.Error(), Done: false}
			continue
		}

		outPath := filepath.Join(outputDir, relPath)

		err = CompressImage(file, outPath, quality)
		if err != nil {
			progress <- Progress{Current: i + 1, Total: total, CurrentFile: relPath, Error: err.Error(), Done: false}
			continue
		}

		srcInfo, _ := os.Stat(file)
		dstSize, _ := findOutputFile(outputDir, relPath)
		srcSize := srcInfo.Size()
		var saved float64
		if srcSize > 0 {
			saved = float64(srcSize-dstSize) / float64(srcSize) * 100
		}

		progress <- Progress{
			Current:     i + 1,
			Total:       total,
			CurrentFile: relPath,
			SrcSize:     srcSize,
			DstSize:     dstSize,
			Saved:       saved,
			Done:        false,
		}
	}

	progress <- Progress{Current: total, Total: total, Done: true}
	return nil
}

func findOutputFile(outputDir, relPath string) (int64, error) {
	path := filepath.Join(outputDir, relPath)
	if info, err := os.Stat(path); err == nil {
		return info.Size(), nil
	}

	jpgPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".jpg"
	if info, err := os.Stat(jpgPath); err == nil {
		return info.Size(), nil
	}

	pngPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".png"
	if info, err := os.Stat(pngPath); err == nil {
		return info.Size(), nil
	}

	return 0, fmt.Errorf("output file not found")
}

type Progress struct {
	Current     int
	Total       int
	CurrentFile string
	SrcSize     int64
	DstSize     int64
	Saved       float64
	Error       string
	Done        bool
}
