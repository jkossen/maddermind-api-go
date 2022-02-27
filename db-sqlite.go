package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const sqlCreateTable string = `
CREATE TABLE IF NOT EXISTS "mmdaily" (
	"id"		INTEGER,
	"length"	INTEGER,
	"date"		INTEGER,
	"code"		TEXT,
	PRIMARY KEY("id" AUTOINCREMENT))`

const sqlInsertTodaysChallenge string = `
INSERT INTO mmdaily (length, date, code)
VALUES (?, ?, ?)
`

const sqlSelectTodaysChallenge string = `
SELECT code FROM mmdaily
WHERE date = ?
AND length = ?
LIMIT 1
`

func SelectTodaysChallenge(db *sql.DB, len int) (string, error) {
	var code string
	today := StartOfDayEpoch()
	row := db.QueryRow(sqlSelectTodaysChallenge, today, len)

	switch err := row.Scan(&code); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return "", err
	case nil:
		return code, nil
	default:
		checkErr(err)
		return "", err
	}
}

func CreateTodaysChallenge(db *sql.DB, len int, code string) {
	stmt, err := db.Prepare(sqlInsertTodaysChallenge)
	checkErr(err)
	s := StartOfDayEpoch()
	_, err = stmt.Exec(len, s, code)
	defer stmt.Close()
	checkErr(err)
}

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	checkErr(err)
	return file.Close()
}

func OpenDb() *sql.DB {
	if _, err := os.Stat("./maddermind.db"); errors.Is(err, os.ErrNotExist) {
		err := TouchFile("./maddermind.db")
		checkErr(err)
	}

	db, err := sql.Open("sqlite3", "./maddermind.db")
	checkErr(err)

	stmt, err := db.Prepare(sqlCreateTable)
	checkErr(err)
	stmt.Exec()
	defer stmt.Close()

	// defer db.Close() <-- why not this instead of CloseDb()?
	return db
}

func CloseDb(db *sql.DB) {
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
