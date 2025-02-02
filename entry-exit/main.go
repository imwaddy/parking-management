package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type VehicleRecords struct {
	ID              int        `json:"id"`
	SpotNumber      string     `json:"spot_number"`
	LicensePlate    string     `json:"license_plate"`
	EntryTime       *time.Time `json:"-"`
	EntryTimeString string     `json:"entry_time"`
	ExitTime        *time.Time `json:"-"`
	ExitTimeString  string     `json:"exit_time,omitempty"`
}

var db *sql.DB

func InitDB() {
	var err error
	connStr := "password=12345 user=postgres dbname=parking-management sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	schema := `
    CREATE TABLE IF NOT EXISTS parking_spots (
        id SERIAL PRIMARY KEY,
        spot_number TEXT,
        type TEXT,
        is_available TEXT CHECK (is_available IN ('yes', 'no'))
    );`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	schema2 := `
    CREATE TABLE IF NOT EXISTS vehicle_records (
        id SERIAL PRIMARY KEY,
        spot_number TEXT,
        license_plate TEXT,
        entry_time TIMESTAMP,
        exit_time TIMESTAMP
    );`
	_, err = db.Exec(schema2)
	if err != nil {
		log.Fatal(err)
	}
}

func formatDateTime(t time.Time) string {
	return t.Format("02-01-2006 15:04:05")
}

func CreateVehicleEntry(w http.ResponseWriter, r *http.Request) {

	var entry VehicleRecords
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the parking spot is available
	var isAvailable string
	err := db.QueryRow("SELECT is_available FROM parking_spots WHERE spot_number = $1", entry.SpotNumber).Scan(&isAvailable)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	if isAvailable == "no" {
		http.Error(w, "Parking spot is not available", http.StatusConflict)
		return
	}

	// Set the current time as entry time
	v := time.Now()
	entry.EntryTime = &v

	// Insert the vehicle entry into the database
	err = db.QueryRow(
		"INSERT INTO vehicle_records (spot_number, license_plate, entry_time) VALUES ($1, $2, $3) RETURNING id",
		entry.SpotNumber, entry.LicensePlate, entry.EntryTime,
	).Scan(&entry.ID)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	// Update parking spot availability
	_, err = db.Exec("UPDATE parking_spots SET is_available = $1 WHERE spot_number = $2", "no", entry.SpotNumber)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	entry.EntryTimeString = formatDateTime(*entry.EntryTime) // Format for JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func CreateVehicleExit(w http.ResponseWriter, r *http.Request) {
	var exit VehicleRecords
	if err := json.NewDecoder(r.Body).Decode(&exit); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"SELECT id, entry_time FROM vehicle_records WHERE spot_number = $1 AND license_plate = $2 AND exit_time IS NULL",
		exit.SpotNumber, exit.LicensePlate,
	).Scan(&exit.ID, &exit.EntryTime)
	if err != nil {
		http.Error(w, "Active vehicle entry not found", http.StatusNotFound)
		return
	}

	exitTime := time.Now()
	exit.ExitTime = &exitTime

	_, err = db.Exec("UPDATE vehicle_records SET exit_time = $1 WHERE id = $2", exit.ExitTime, exit.ID)
	if err != nil {
		http.Error(w, "Failed to update exit time", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE parking_spots SET is_available = $1 WHERE spot_number = $2", "yes", exit.SpotNumber)
	if err != nil {
		http.Error(w, "Failed to update parking spot availability", http.StatusInternalServerError)
		return
	}

	exit.EntryTimeString = formatDateTime(*exit.EntryTime)
	exit.ExitTimeString = formatDateTime(*exit.ExitTime)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exit)
}

// NOTE: vary as per question
func GetVehicleEntryBySpotNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["spot_number"]

	rows, err := db.Query(
		"SELECT id, spot_number, license_plate, entry_time, exit_time FROM vehicle_records WHERE spot_number = $1",
		id)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	var exits []VehicleRecords
	if rows.Next() {
		var exit VehicleRecords
		var time1, time2 *time.Time
		err := rows.Scan(&exit.ID, &exit.SpotNumber, &exit.LicensePlate, &time1, &time2)
		if err != nil {
			http.Error(w, "Error scanning vehicle entry", http.StatusInternalServerError)
			return
		}

		if !time1.IsZero() {
			exit.EntryTimeString = formatDateTime(*time1)
		}
		if time2 != nil && !time2.IsZero() {
			exit.ExitTimeString = formatDateTime(*time2)
		}

		exits = append(exits, exit)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exits)
}

func RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/api/vehicle-entries", CreateVehicleEntry).Methods("POST")
	router.HandleFunc("/api/vehicle-exits", CreateVehicleExit).Methods("POST")
	router.HandleFunc("/api/vehicle-exits/{spot_number}", GetVehicleEntryBySpotNumber).Methods("GET")

	log.Println("Server is starting  on port :8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	InitDB()
	RegisterRoutes()
}
