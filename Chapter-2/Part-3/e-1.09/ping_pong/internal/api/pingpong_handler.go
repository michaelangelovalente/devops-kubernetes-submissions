package handler

import (
	"net/http"
	"ping_pong/internal/store"
	"ping_pong/internal/utils"
)

type PingPongHandler struct {
	store *store.Store
	//service..
	// logger       *log.Logger
}

func NewPingPongHandler(store *store.Store) *PingPongHandler {
	return &PingPongHandler{
		store: store,
	}
}

func (ph *PingPongHandler) Get(w http.ResponseWriter, r *http.Request) {
	count, _ := ph.store.PingPongStore.GetCurr()
	utils.WriteJSON(
		w, http.StatusOK,
		utils.Envelope{
			"count": count,
		},
	)
}

// func (ph *PingPongHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	count, err := ph.store.PingPongStore.Update(&ctx, ph.store.PingPongStore.(*store.PingPongStore).PingPongModel)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
// 		return
// 	}
// 	utils.WriteJSON(
// 		w, http.StatusOK,
// 		utils.Envelope{
// 			"count": count,
// 		},
// 	)
// }

// func (ph *PingPongHandler) Reset(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.Envelope{"error": "method not allowed"})
// 		return
// 	}
//
// 	ctx := r.Context()
// 	pingPongStore := ph.store.PingPongStore.(*store.PingPongStore)
//
// 	if err := ph.store.PingPongStore.Reset(&ctx, pingPongStore.PingPongModel); err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
// 		return
// 	}
//
// 	utils.WriteJSON(
// 		w, http.StatusOK,
// 		utils.Envelope{
// 			"message": "counter reset successfully",
// 			"count":   0,
// 		},
// 	)
// }
//
// func (ph *PingPongHandler) Ping(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	count, err := ph.store.PingPongStore.Update(&ctx, ph.store.PingPongStore.(*store.PingPongStore).PingPongModel)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
// 		return
// 	}
//
// 	utils.WriteJSON(
// 		w, http.StatusOK,
// 		utils.Envelope{
// 			"message": "pong",
// 			"count":   count,
// 		},
// 	)
// }
