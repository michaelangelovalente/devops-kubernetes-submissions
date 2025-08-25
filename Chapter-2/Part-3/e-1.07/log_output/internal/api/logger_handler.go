package api

import (
	"log_output/internal/store"
	"log_output/internal/utils"
	"net/http"
)

type LoggerEntryHandler struct {
	loggerStore store.LogStorage // Use interface, not concrete type
}

func NewLoggerEntryHandler(loggerStore store.LogStorage) *LoggerEntryHandler {
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

func (leh *LoggerEntryHandler) GetLastLogAndStatus(w http.ResponseWriter, r *http.Request) {
	logs := leh.loggerStore.GetLatest(10)

	utils.WriteJSON(w, http.StatusOK,
		utils.Envelope{
			"logs": logs,
		},
	)
}
