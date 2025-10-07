package reader

type PingPongReader struct {
	path string
}

func NewPingPongReader(path string) *PingPongReader {
	return &PingPongReader{
		path: path,
	}
}

func (pr *PingPongReader) GetPingPong(n int) (string, error) {
	pingpongLogs, err := getLogsFromFile(pr.path)
	if err != nil {
		return "", err
	}

	return pingpongLogs, nil
}
