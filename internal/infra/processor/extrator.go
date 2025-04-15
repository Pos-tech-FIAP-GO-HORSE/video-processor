package processor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExtractFrames usa o ffmpeg para extrair um frame a cada 20 segundos
func ExtractFrames(videoPath string) ([]string, error) {
	// Cria diretório temporário para os frames
	outputDir := filepath.Join("/tmp", "frames_"+filepath.Base(videoPath))
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("erro criando pasta de frames: %v", err)
	}

	// Define padrão dos arquivos gerados
	outputPattern := filepath.Join(outputDir, "frame_%04d.jpg")

	// Executa ffmpeg
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1/20", outputPattern)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar ffmpeg: %v", err)
	}

	// Lê os arquivos gerados
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("erro lendo pasta de frames: %v", err)
	}

	var paths []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".jpg") {
			paths = append(paths, filepath.Join(outputDir, f.Name()))
		}
	}

	return paths, nil
}
