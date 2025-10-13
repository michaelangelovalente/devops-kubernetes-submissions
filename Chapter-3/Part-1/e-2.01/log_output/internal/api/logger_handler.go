package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"common/utils"
	client "log_output/internal/client/pingpong"
	"log_output/internal/store"
)

type LoggerEntryHandler struct {
	loggerStore    store.LogStorage // Use interface, not concrete type
	logger         *log.Logger
	pingpongClient client.Client
}

func NewLoggerEntryHandler(loggerMemoryStore store.LogStorage, logger *log.Logger, pingpongClient client.Client) *LoggerEntryHandler {
	return &LoggerEntryHandler{
		loggerStore:    loggerMemoryStore,
		logger:         logger,
		pingpongClient: pingpongClient,
	}
}

func (leh *LoggerEntryHandler) GetLastLogAndCount(w http.ResponseWriter, r *http.Request) {
	logs := leh.loggerStore.GetLatest(1)
	pingpongCount, err := leh.pingpongClient.GetCount()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError,
			utils.Envelope{
				"error": fmt.Sprintf("internal server error: %v", err),
			},
		)
		return
	}

	//TODO!: fix log format
	utils.Write(w, http.StatusOK, fmt.Sprintf("%s: %s\nPing / Pongs: %d\n", logs[0].Timestamp.Format(time.RFC3339), logs[0].Value, pingpongCount))

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
