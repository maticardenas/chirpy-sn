package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.Mutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func readJSONFile(path string) (DBStructure, error) {
	data := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
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

	chirps.Chirps[chirp.Id] = chirp

	fmt.Println("Chirps after adding new chirp:", chirps)

	err = writeJSONFile(db.path, chirps)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetUser(email string) (User, error) {
	fmt.Println("Getting user by email")

	users, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return User{}, err
	}

	for _, user := range users.Users {
		if user.Email == email {
			fmt.Printf("USER FOUND -> User id: %v - User email: %v\n", user.Id, user.Email)
			return user, nil
		}
	}

	fmt.Println("User not found")
	return User{}, nil
}

func (db *DB) GetLastUserID() (int, error) {
	fmt.Println("Getting last user id")

	users, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return 0, err
	}

	lastUserID := users.Users[len(users.Users)].Id

	fmt.Println("Last user id:", lastUserID)

	return lastUserID, nil
}

var (
	ErrUserExists        = errors.New("User already exists")
	ErrHashingPassword   = errors.New("Error hashing password")
	ErrReadingFile       = errors.New("Error reading file")
	ErrWritingFile       = errors.New("Error writing to file")
	ErrGettingLastUserID = errors.New("Error getting last user id")
	ErrUserNotFound      = errors.New("User not found")
	ErrWrongPassword     = errors.New("Wrong password")
)

func (db *DB) CreateUser(email string, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	existingUser, err := db.GetUser(email)
	if err != nil {
		fmt.Println("Error getting user:", err)
		return User{}, ErrReadingFile
	}
	if existingUser.Email != "" {
		fmt.Println("User already exists")
		return User{}, ErrUserExists
	}

	lastUserID, err := db.GetLastUserID()
	if err != nil {
		fmt.Println("Error getting last user id:", err)
		return User{}, ErrGettingLastUserID
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return User{}, ErrHashingPassword
	}
	user := User{
		Id:       lastUserID + 1,
		Email:    email,
		Password: string(hashedPassword),
	}

	fmt.Printf("User id: %v - User email: %v\n", user.Id, user.Email)

	data, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return User{}, ErrReadingFile
	}
	data.Users[user.Id] = user

	err = writeJSONFile(db.path, data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return User{}, ErrWritingFile
	}

	return user, nil
}

func (db *DB) GetChirps(chirpId ...int) ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	chirps, err := readJSONFile(db.path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []Chirp{}, err
	}

	chirpsList := make([]Chirp, 0, len(chirps.Chirps))

	if len(chirpId) == 0 {
		for _, chirp := range chirps.Chirps {
			chirpsList = append(chirpsList, chirp)
		}
	} else {
		chirpsList = append(chirpsList, chirps.Chirps[chirpId[0]])
	}

	return chirpsList, nil
}

func (db *DB) CheckUser(email string, password string) (User, error) {
	user, err := db.GetUser(email)
	if err != nil {
		fmt.Println("Error getting user:", err)
		return User{}, err
	}

	if user.Email == "" {
		fmt.Println("User not found")
		return User{}, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Error comparing password:", err)
		return User{}, ErrWrongPassword
	}

	return user, nil
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
			Users:  make(map[int]User),
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
