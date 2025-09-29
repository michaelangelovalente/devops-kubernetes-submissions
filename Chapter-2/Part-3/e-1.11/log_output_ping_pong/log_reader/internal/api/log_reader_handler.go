package api

import (
	"fmt"
	"log"
	"log_reader/internal/reader"
	"log_reader/internal/utils"
	"net/http"
)

type LogReaderHandler struct {
	logReader  reader.Reader
	pingReader reader.PingPongReader
	logger     *log.Logger
}

func NewLogReaderHandler(logReader reader.Reader, pingReader reader.PingPongReader, logger *log.Logger) *LogReaderHandler {
	return &LogReaderHandler{
		logReader:  logReader,
		pingReader: pingReader,
		logger:     logger,
	}
}

func (lr *LogReaderHandler) GetAllLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := lr.logReader.GetAll()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError,
			utils.Envelope{
				"error": fmt.Sprintf("error: file read: %v", err),
			},
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK,
		utils.Envelope{
			"logs": logs,
		},
	)
}

func (lr *LogReaderHandler) GetLogsPingPong(w http.ResponseWriter, r *http.Request) {

	pingpongData, err := lr.pingReader.GetPingPong(1)
	if err != nil {
		lr.logger.Printf("ERROR: cannot read pingpong file: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError,
			utils.Envelope{
				"error": fmt.Sprintf("cannot read pingpong file: %v", err),
			},
		)
		return
	}

	logData, err := lr.logReader.GetLastLog()
	if err != nil {
		lr.logger.Printf("ERROR: cannot read log file: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError,
			utils.Envelope{
				"error": fmt.Sprintf("cannot read log file: %v", err),
			},
		)
		return
	}

	utils.Write(w, http.StatusOK, logData+"\n"+pingpongData)
}
