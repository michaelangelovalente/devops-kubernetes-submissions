package api

import (
	"common/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	client "log_output/internal/client/pingpong"
	"log_output/internal/store"
)

type LoggerEntryHandler struct {
	loggerStore    store.LogStorage // Use interface, not concrete type
	logger         *log.Logger
	pingpongClient client.Client
	fileInfoPath   string
}

func NewLoggerEntryHandler(loggerMemoryStore store.LogStorage, logger *log.Logger, pingpongClient client.Client, fileInfoPath string) *LoggerEntryHandler {
	return &LoggerEntryHandler{
		loggerStore:    loggerMemoryStore,
		logger:         logger,
		pingpongClient: pingpongClient,
		fileInfoPath:   fileInfoPath,
	}
}

func (leh *LoggerEntryHandler) GetLatestData(w http.ResponseWriter, r *http.Request) {

	msg := os.Getenv("MESSAGE")
	if msg == "" {
		msg = "no message found for env variable MESSAGE"
	}
	envVarMsg := fmt.Sprintf("env variable: %s=%s", "MESSAGE", msg)

	fileContent, err := os.ReadFile(leh.fileInfoPath)
	var fileContentTxt string
	if err != nil {
		fileContentTxt = fmt.Sprintf("file content: ERROR - Could not read file from path %s: %v", leh.fileInfoPath, err)
	} else {
		fileContentTxt = fmt.Sprintf("file content: %s", strings.TrimSpace(string(fileContent)))
	}

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
	logLine := fmt.Sprintf("%s: %s", logs[0].Timestamp.Format(time.RFC3339), logs[0].Value)
	ppsLine := fmt.Sprintf("Ping / Pongs: %d\n", pingpongCount)

	fullResponse := fmt.Sprintf("%s\n%s\n%s\n%s\n", fileContentTxt, envVarMsg, logLine, ppsLine)
	utils.Write(w, http.StatusOK, fullResponse)
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
