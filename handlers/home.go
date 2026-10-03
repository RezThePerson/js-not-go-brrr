package handlers

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
)

var tmpl *template.Template

func InitTemplates(webFS embed.FS) (err error) {
	tmpl, err = template.ParseFS(webFS, "web/templates/*.html")
	if err != nil {
		return err
	}
	slog.Info("templates loaded")
	return nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		slog.Error("template execution error", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
