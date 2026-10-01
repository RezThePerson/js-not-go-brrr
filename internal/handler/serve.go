package routes

import "net/http"

func Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/jump", _)
	http.HandleFunc("/stream", _)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
