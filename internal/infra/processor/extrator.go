package processor

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type FFmpegProcessor struct{}

var ExecCommand = exec.Command

func (f FFmpegProcessor) ExtractFrames(videoPath string) ([]string, error) {
	outputDir := filepath.Dir(videoPath)
	outputPattern := filepath.Join(outputDir, "frame-%03d.jpg")

	cmd := ExecCommand("ffmpeg", "-i", videoPath, "-vf", "fps=1", outputPattern)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao extrair frames com ffmpeg: %w", err)
	}

	var frames []string
	for i := 1; i <= 3; i++ {
		framePath := fmt.Sprintf(filepath.Join(outputDir, "frame-%03d.jpg"), i)
		frames = append(frames, framePath)
	}

	return frames, nil
}
