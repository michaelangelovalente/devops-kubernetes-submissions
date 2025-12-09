package handler

import (
	"common/utils"
	"net/http"

	"ping_pong/internal/store"
)

type PingPongHandler struct {
	pingpongRepo store.PingPongRepo
	// service..
	// logger       *log.Logger
}

func NewPingPongHandler(pingpongRepo store.PingPongRepo) *PingPongHandler {
	return &PingPongHandler{
		pingpongRepo: pingpongRepo,
	}
}

func (ph *PingPongHandler) Get(w http.ResponseWriter, r *http.Request) {
	count, _ := ph.pingpongRepo.GetCurr()
	utils.WriteJSON(
		w, http.StatusOK,
		utils.Envelope{
			"count": count,
		},
	)
}

func (ph *PingPongHandler) Update(w http.ResponseWriter, r *http.Request) {
	count, err := ph.pingpongRepo.Update()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}
	utils.WriteJSON(
		w, http.StatusOK,
		utils.Envelope{
			"count": count,
		},
	)
}
