package handlers

import (
	"net/http"
)

func JumpHandler(w http.ResponseWriter, r *http.Request) {
	gs.Jump()
	w.WriteHeader(http.StatusNoContent)
}
