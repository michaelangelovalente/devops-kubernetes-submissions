package api

import (
	"log_output/internal/logger"
	"log_output/internal/utils"
	"net/http"
)

type LoggerEntryHandler struct {
	loggerStore logger.LogStorage // Use interface, not concrete type
}

func NewLoggerEntryHandler(loggerStore logger.LogStorage) *LoggerEntryHandler {
	return &LoggerEntryHandler{
		loggerStore: loggerStore,
	}
}

func (leh *LoggerEntryHandler) GetAllLogs(w http.ResponseWriter, r *http.Request) {
	logs := leh.loggerStore.GetAll()
	
	utils.WriteJSON(w, http.StatusOK,
		utils.Envelope{
			"logs": logs,
		},
	)
}
