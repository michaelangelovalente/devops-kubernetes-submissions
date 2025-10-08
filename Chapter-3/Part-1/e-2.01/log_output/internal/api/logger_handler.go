package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"common/utils"
	"log_output/internal/store"
)

type LoggerEntryHandler struct {
	loggerStore store.LogStorage // Use interface, not concrete type
	logger      *log.Logger
}

func NewLoggerEntryHandler(loggerMemoryStore store.LogStorage, logger *log.Logger) *LoggerEntryHandler {
	return &LoggerEntryHandler{
		loggerStore: loggerMemoryStore,
		logger:      logger,
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

func (leh *LoggerEntryHandler) GetLastLogsAndStatus(w http.ResponseWriter, r *http.Request) {
	nParam := r.URL.Query().Get("n")

	var lastNLogs int64 = 10 // Default value

	if nParam != "" {
		parsedN, err := strconv.ParseInt(nParam, 10, 64)
		if err != nil {
			leh.logger.Printf("ERROR: parsing query parameter 'n': %v\n", err)
			utils.WriteJSON(w,
				http.StatusBadRequest,
				utils.Envelope{
					"error": fmt.Sprintf("invalid query parameter 'n': %v", err),
				})
			return
		}
		lastNLogs = parsedN
	}

	logs := leh.loggerStore.GetLatest(int(lastNLogs))

	response := utils.Envelope{
		"status": "ready",
	}
	response["logs"] = logs

	utils.WriteJSON(w, http.StatusOK, response)
}
