package main

import "net/http"

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	// server file in ./assets/logo.png
	serveMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	serveMux.HandleFunc("/healthz", readinessHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
