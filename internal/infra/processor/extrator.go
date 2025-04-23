package processor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type FFmpegProcessor struct{}

var ExecCommand = exec.Command

func (f FFmpegProcessor) ExtractFrames(videoPath string) ([]string, error) {
	outputDir := filepath.Dir(videoPath)
	outputPattern := filepath.Join(outputDir, "frame-%03d.jpg")

	cmd := ExecCommand("./ffmpeg", "-i", videoPath, "-vf", "fps=1", outputPattern)
	fmt.Println("Executando comando:", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erro ao executar ffmpeg: %v\nSaída: %s\n", err, string(output))
		return nil, fmt.Errorf("erro ao extrair frames com ffmpeg: %w", err)
	}
	fmt.Printf("Saída do ffmpeg: %s\n", string(output))

	var frames []string
	for i := 1; i <= 3; i++ {
		framePath := fmt.Sprintf(filepath.Join(outputDir, "frame-%03d.jpg"), i)

		// Verifica se o arquivo realmente foi criado
		if _, err := os.Stat(framePath); os.IsNotExist(err) {
			fmt.Printf("Frame não encontrado: %s\n", framePath)
			continue
		}

		frames = append(frames, framePath)
		fmt.Println("Frame adicionado:", framePath)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("nenhum frame foi extraído de %s", videoPath)
	}

	return frames, nil
}
