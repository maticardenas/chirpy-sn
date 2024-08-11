package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type errorResponseBody struct {
		Error string `json:"error"`
	}
	type successResponseBody struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		respBody := errorResponseBody{
			Error: "Something went wrong",
		}
		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	if len(reqBody.Body) > 140 {
		fmt.Println("Chirp is too long")
		respBody := errorResponseBody{
			Error: "Chirp is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	fmt.Println("Chirp is valid")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	respBody := successResponseBody{
		Valid: true,
	}
	dat, _ := json.Marshal(respBody)
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) hitsCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	file, err := os.Open("metrics.html")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	metricsHTML := fmt.Sprintf(string(fileContent), cfg.fileServerHits)
	w.Write([]byte(metricsHTML))
}

func (cfg *apiConfig) resetCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits = 0
}

func main() {
	serveMux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	cfg := &apiConfig{}
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	// server file in ./assets/logo.png
	serveMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.hitsCountHandler)
	serveMux.HandleFunc("/api/reset", cfg.resetCountHandler)

	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}
