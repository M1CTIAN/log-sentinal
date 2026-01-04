package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Alert struct matches our DB table
type Alert struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

var db *sql.DB

func main() {
	// 1. Connect to DB (Same as Consumer)
	connStr := "postgres://admin:secretpassword@localhost:5432/logsentinel?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Setup Routes
	http.HandleFunc("/", serveHome)           // Serve the HTML Dashboard
	http.HandleFunc("/api/alerts", getAlerts) // The JSON API

	// 3. Start Server
	log.Println("🌐 Dashboard running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// Serve the HTML file
func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("cmd/dashboard/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.Execute(w, nil)
}

// Fetch alerts from DB and return as JSON
func getAlerts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, message, created_at FROM alerts ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		// Scan matches the columns in the SQL query
		if err := rows.Scan(&a.ID, &a.Message, &a.CreatedAt); err != nil {
			continue
		}
		alerts = append(alerts, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}
