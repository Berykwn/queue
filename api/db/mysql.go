// db/mysql.go

package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/queue")
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging the database: %v", err)
		}

		fmt.Println("Connected to the database")
	})
	return db
}

func GetQueueData() (string, error) {
	var data string
	db := GetDB()
	err := db.QueryRow("SELECT data FROM queue LIMIT 1").Scan(&data)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching queue data from database: %v\n", err)
		}
		return "", err
	}
	return data, nil
}

func UpdateQueueProcessedAt(data string) error {
	db := GetDB()
	_, err := db.Exec("UPDATE queue SET processed_at = ? WHERE data = ?", time.Now(), data)
	if err != nil {
		log.Printf("Error updating queue processed at time in database: %v\n", err)
		return err
	}
	return nil
}
