package processor

type FrameProcessor interface {
	ExtractFrames(videoPath string) ([]string, error)
}
