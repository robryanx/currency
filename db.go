package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func start_of_day(t time.Time) time.Time {
    year, month, day := t.Date()

    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Update(to_currency string, rate float64) {
    db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE")))
    defer db.Close()

    start_of_day_timestamp := start_of_day(time.Now().UTC()).Unix()
    
    stmtIns, _ := db.Prepare("INSERT INTO `currency_conversion` VALUES( ?, ?, ? )")
    defer stmtIns.Close()

    _, err := stmtIns.Exec(to_currency, start_of_day_timestamp, rate)
    if err != nil {
        log.Fatalln(err)
    }
}