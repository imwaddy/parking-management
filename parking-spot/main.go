package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var db *sql.DB

type ParkingSpot struct {
	ID          int    `json:"id"`
	SpotNumber  string `json:"spot_number"`
	Type        string `json:"type"`
	IsAvailable string `json:"is_available"`
}

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
}

// func formatDateTime(t time.Time) string {
// 	return t.Format("02-01-2006 15:04:05")
// }

func validateSpotType(spotType string) bool {
	return len(spotType) != 0
}

// func validateIsAvailable(isAvailable string) bool {
// 	if isAvailable == "yes" || isAvailable == "no" {
// 		return true
// 	}
// 	return false
// }

func AddParkingSpot(w http.ResponseWriter, r *http.Request) {
	var spot ParkingSpot
	if err := json.NewDecoder(r.Body).Decode(&spot); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !validateSpotType(spot.Type) {
		http.Error(w, "Invalid parking spot", http.StatusInternalServerError)
		return
	}

	err := db.QueryRow("INSERT INTO parking_spots (spot_number, type, is_available) VALUES ($1, $2, $3) returning id", spot.SpotNumber, spot.Type, spot.IsAvailable).Scan(&spot.ID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spot)
}

func GetAllParkingSpots(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, spot_number, type, is_available FROM parking_spots")
	if err != nil {
		http.Error(w, "Failed to scan parking spot", http.StatusInternalServerError)
		return
	}

	var spots []ParkingSpot
	for rows.Next() {
		var spot ParkingSpot
		if err := rows.Scan(&spot.ID, &spot.SpotNumber, &spot.Type, &spot.IsAvailable); err != nil {
			http.Error(w, "Failed to scan parking spot", http.StatusInternalServerError)
			return
		}

		spots = append(spots, spot)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spots)
}

func GetParkingSpot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var spot ParkingSpot
	err := db.QueryRow("SELECT id, spot_number, type, is_available FROM parking_spots WHERE id = $1", id).Scan(&spot.ID, &spot.SpotNumber, &spot.Type, &spot.IsAvailable)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spot)
}

func UpdateParkingSpot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var spot ParkingSpot
	if err := json.NewDecoder(r.Body).Decode(&spot); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if !validateSpotType(spot.Type) {
		http.Error(w, "Invalid parking spot", http.StatusInternalServerError)
		return
	}

	_, err := db.Exec("UPDATE parking_spots SET spot_number = $1, type = $2, is_available = $3 WHERE id = $4", spot.SpotNumber, spot.Type, spot.IsAvailable, id)
	if err != nil {
		http.Error(w, "Invalid parking spot", http.StatusInternalServerError)
		return
	}

	spot.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spot)
}

func DeleteParkingSpot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var spot ParkingSpot
	err := db.QueryRow("SELECT id, spot_number, type, is_available FROM parking_spots WHERE id = $1", id).Scan(&spot.ID, &spot.SpotNumber, &spot.Type, &spot.IsAvailable)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM parking_spots WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Parking spot not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Parking spot has been deleted successfully.",
	})
}

func RegisterRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/api/parking-spots", AddParkingSpot).Methods("POST")
	router.HandleFunc("/api/parking-spots/all", GetAllParkingSpots).Methods("GET")
	router.HandleFunc("/api/parking-spots/{id}", GetParkingSpot).Methods("GET")
	router.HandleFunc("/api/parking-spots/{id}", UpdateParkingSpot).Methods("PUT")
	router.HandleFunc("/api/parking-spots/{id}", DeleteParkingSpot).Methods("DELETE")

	log.Println("Server is starting  on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	InitDB()
	RegisterRoutes()
}
