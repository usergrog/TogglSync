package models

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func InitDb() {
	//file, err := os.Create("smart-database.db") // Create SQLite file
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//_ = file.Close()
	log.Println("Connect to DB")

	db, _ = sqlx.Open("sqlite3", "./sync-database.db") // Open the created SQLite File
	createTables()
}

func CloseDb() {
	log.Println("Close db connection")
	_ = db.Close()
}

func createTables() {
	createEntriesTableSQL := `CREATE TABLE IF NOT EXISTS entries (
		id INTEGER NOT NULL PRIMARY KEY,
		start TEXT,
		stop TEXT,
		duration INTEGER,
		description TEXT,
		synced_at TEXT
	  );` // SQL Statement for Create Table
	log.Println("Create entries table...")
	statement, err := db.Prepare(createEntriesTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, _ = statement.Exec() // Execute SQL Statements
	log.Println("Entries table created")
}

func InsertLocation(e *TogglEntry) {
	sqlStatement := `
		INSERT INTO entries
		(id, start, stop, duration, description, synced_at)
		VALUES($1, $2, $3, $4, $5, $6)`
	stmt, err := db.Prepare(sqlStatement)
	t := time.Now()
	if err == nil {
		_, err := db.Exec(sqlStatement, e.Id, e.Start, e.Stop, e.Duration, e.Description, t.Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println("Error in sql statement:", err)
		}
		_ = stmt.Close()
	} else {
		log.Println("Error in sql statement:", err)
	}
}

func GetEntry(id int) TogglEntry {
	var entry TogglEntry
	sqlStatement :=
		`select id, start, stop, duration, description
		from entries
		where id = ?`

	db.Get(&entry, sqlStatement, id)

	return entry
}
