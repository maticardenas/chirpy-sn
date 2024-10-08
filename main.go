package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/maticardenas/chirpy-sn/internal/database"
)

var DbInstance *database.DB

type apiConfig struct {
	fileServerHits int
}
type chirpRequestBody struct {
	Body string `json:"body"`
}
type userRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type userResponseBody struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
type errorResponseBody struct {
	Error string `json:"error"`
}
type successResponseBody struct {
	CleanedBody string `json:"cleaned_body"`
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

func getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println("Error converting id to int:", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid ID"}`))
		return
	}

	chirp, err := DbInstance.GetChirps(id)
	if err != nil {
		fmt.Println("Error getting chirp:", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Chirp not found"}`))
		return
	}

	if chirp[0].Id == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Chirp not found"}`))
		return
	}

	dat, _ := json.Marshal(chirp[0])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := userRequestBody{}
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

	fmt.Println("User is valid")
	fmt.Println("Email:", reqBody.Email)
	fmt.Println("Password:", reqBody.Password)

	user, err := DbInstance.CreateUser(reqBody.Email, reqBody.Password)

	if err != nil {
		fmt.Println("Error creating user:", err)
		respBody := errorResponseBody{
			Error: err.Error(),
		}
		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	respBody := userResponseBody{
		Id:    user.Id,
		Email: user.Email,
	}

	dat, _ := json.Marshal(respBody)

	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := userRequestBody{}
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

	fmt.Println("Payload is valid")
	fmt.Println("Email:", reqBody.Email)
	fmt.Println("Password:", reqBody.Password)

	user, err := DbInstance.CheckUser(reqBody.Email, reqBody.Password)

	if err != nil {
		fmt.Println("Error checking user credentials:", err)
		respBody := errorResponseBody{
			Error: err.Error(),
		}
		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	respBody := userResponseBody{
		Id:    user.Id,
		Email: user.Email,
	}

	dat, _ := json.Marshal(respBody)

	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func createChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := chirpRequestBody{}
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

	chirp, err := DbInstance.CreateChirp(reqBody.Body)

	if err != nil {
		fmt.Println("Error creating chirp:", err)
		respBody := errorResponseBody{
			Error: "Something went wrong",
		}
		dat, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// respBody := successResponseBody{
	// 	CleanedBody: chirptext.ReplaceChirpInput(reqBody.Body),
	// }
	dat, _ := json.Marshal(chirp)
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := DbInstance.GetChirps()
	if err != nil {
		fmt.Println("Error getting chirps:", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	dat, _ := json.Marshal(chirps)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

func initializeDB() (*database.DB, error) {
	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return nil, err
	}
	return db, nil
}

func main() {
	DbInstance, _ = initializeDB()
	serveMux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	cfg := &apiConfig{}
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	// server file in ./assets/logo.png
	serveMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.hitsCountHandler)
	serveMux.HandleFunc("/api/reset", cfg.resetCountHandler)
	serveMux.HandleFunc("POST /api/chirps", createChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", getChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", getChirpHandler)
	serveMux.HandleFunc("POST /api/users", createUserHandler)
	serveMux.HandleFunc("POST /api/login", loginHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}
