package handlers

import (
	"fmt"
	"net/http"
)

// StreamHandler serves multipart/x-mixed-replace SVG frames.
// Load it with <img src="/stream"> — no JS needed.
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable nginx proxy buffering

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	id, ch := globalHub.Register()
	defer globalHub.Unregister(id)

	for {
		select {
		case <-r.Context().Done():
			return
		case frame, ok := <-ch:
			if !ok {
				return
			}
			// Content-Length lets the browser paint the frame immediately
			// rather than buffering until the next boundary (fixes ~1s blank-out).
			fmt.Fprintf(w, "--frame\r\nContent-Type: image/svg+xml\r\nContent-Length: %d\r\n\r\n", len(frame))
			w.Write(frame) //nolint:errcheck
			fmt.Fprintf(w, "\r\n")
			flusher.Flush()
		}
	}
}

// JumpHandler handles a jump POST submitted from the hidden iframe form.
// Returns 204 so the iframe navigation is silent.
func JumpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	gameMu.Lock()
	if !gameState.Dead && gameState.DinoY == 0 {
		gameState.DinoVY = jumpImpulse
	}
	gameMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
