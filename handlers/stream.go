package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

// StreamHandler serves a continuous multipart/x-mixed-replace stream of SVG frames.
// The browser loads this as <img src="/stream"> — no JS required.
// Each connected client gets its own channel; the hub broadcasts the same frame to all.
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable nginx buffering if proxied

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	id, ch := globalHub.Register()
	defer globalHub.Unregister(id)

	slog.Debug("stream client connected", "id", id, "players", globalHub.PlayerCount())

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			slog.Debug("stream client disconnected", "id", id)
			return
		case frame, ok := <-ch:
			if !ok {
				return
			}
			// Write one multipart segment with explicit Content-Length so the
			// browser can decode the frame immediately without buffering to the
			// next boundary (which caused the periodic ~1s blank-out).
			fmt.Fprintf(w, "--frame\r\nContent-Type: image/svg+xml\r\nContent-Length: %d\r\n\r\n", len(frame))
			w.Write(frame) //nolint:errcheck
			fmt.Fprintf(w, "\r\n")
			flusher.Flush()
		}
	}
}

// JumpHandler receives a POST from the jump form (submitted inside a hidden iframe).
// It applies a jump to the shared dino and returns 204 so the iframe stays blank.
func JumpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	Jump()
	w.WriteHeader(http.StatusNoContent)
}
