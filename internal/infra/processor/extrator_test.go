package processor_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"video-processor/internal/infra/processor"

	"github.com/stretchr/testify/assert"
)

func TestExtractFrames_Success(t *testing.T) {
	processor.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo")
	}

	outputDir := "/tmp"
	for i := 1; i <= 3; i++ {
		framePath := filepath.Join(outputDir, fmt.Sprintf("frame-%03d.jpg", i))
		err := os.WriteFile(framePath, []byte("fake frame content"), 0644)
		require.NoError(t, err)
		defer os.Remove(framePath) // Limpa após o teste
	}

	ffmpeg := processor.FFmpegProcessor{}
	frames, err := ffmpeg.ExtractFrames("/tmp/video.mp4")

	assert.NoError(t, err)
	assert.Len(t, frames, 3)
	assert.Contains(t, frames[0], "frame-001.jpg")
}

func TestExtractFrames_Error(t *testing.T) {
	processor.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false") // comando que sempre retorna erro
	}

	ffmpeg := processor.FFmpegProcessor{}
	_, err := ffmpeg.ExtractFrames("/tmp/video.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao extrair frames com ffmpeg")
}
