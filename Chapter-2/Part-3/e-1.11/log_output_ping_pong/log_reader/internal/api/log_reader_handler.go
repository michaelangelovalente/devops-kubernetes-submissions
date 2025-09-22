package api

import (
	"fmt"
	"log"
	"log_reader/internal/reader"
	"log_reader/internal/utils"
	"net/http"
)

type LogReaderHandler struct {
	logReader reader.LogReader
	logger    *log.Logger
}

func NewLogReaderHandler(logReader reader.LogReader, logger *log.Logger) *LogReaderHandler {
	return &LogReaderHandler{
		logReader: logReader,
		logger:    logger,
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
