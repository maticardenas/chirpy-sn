package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.Mutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func readJSONFile(path string) (DBStructure, error) {
	data := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return data, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding file:", err)
		return data, err
	}

	return data, nil
}

func writeJSONFile(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		err = fmt.Errorf("error creating file: %v", err)
		return err
	}

	dat, _ := json.Marshal(data)
	if _, err := file.Write(dat); err != nil {
		err = fmt.Errorf("error writing to file: %v", err)
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	chirps, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Chirp{}, err
	}

	fmt.Printf("Chirps: %v and chirps.Chirps: %v\n", chirps, chirps.Chirps)
	fmt.Printf("Chirps length: %v\n", len(chirps.Chirps))
	lastID := chirps.Chirps[len(chirps.Chirps)].Id

	fmt.Println("Last ID:", lastID)

	chirp := Chirp{
		Id:   lastID + 1,
		Body: body,
	}

	fmt.Printf("Chirp id: %v - Chirp body: %v\n", chirp.Id, chirp.Body)
	fmt.Printf("Chirps: %v and chirps.Chirps: %v\n", chirps, chirps.Chirps)

	// chirps.Chirps = make(map[int]Chirp)
	chirps.Chirps[chirp.Id] = chirp

	fmt.Println("Chirps after adding new chirp:", chirps)

	err = writeJSONFile(db.path, chirps)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	chirps, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []Chirp{}, err
	}

	chirpsList := make([]Chirp, 0, len(chirps.Chirps))

	for _, chirp := range chirps.Chirps {
		chirpsList = append(chirpsList, chirp)
	}

	return chirpsList, nil
}

func NewDB(path string) (*DB, error) {
	dbExists := false
	if _, err := os.Stat(path); err == nil {
		dbExists = true
		fmt.Println("Database file exists")
	}

	if !dbExists {
		dbStructure := DBStructure{
			Chirps: make(map[int]Chirp),
		}
		fmt.Printf("DBStrucutre: %v\n", dbStructure)

		err := writeJSONFile(path, dbStructure)
		if err != nil {
			return nil, err
		}
	}

	db := &DB{
		path: path,
		mux:  &sync.Mutex{},
	}

	return db, nil
}
