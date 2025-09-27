package handler

import (
	"fmt"
	"net/http"

	"ping_pong/internal/store"
	"ping_pong/internal/utils"
)

type PingPongHandler struct {
	store *store.Store
}

func NewPingPongHandler(store *store.Store) *PingPongHandler {
	return &PingPongHandler{
		store: store,
	}
}

func (ph *PingPongHandler) Get(w http.ResponseWriter, r *http.Request) {
	count, err := ph.store.PingPongStore.GetCurr()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	utils.Write(
		w, http.StatusOK,
		fmt.Sprintf("Result: %d", *count),
	)
}

func (ph *PingPongHandler) Update(w http.ResponseWriter, r *http.Request) {
	count, err := ph.store.PingPongStore.Update()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal Server Error"})
		return
	}
	utils.Write(
		w, http.StatusOK,
		fmt.Sprintf("Result: %d", *count),
	)
}
