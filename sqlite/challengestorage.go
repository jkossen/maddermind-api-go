package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type ChallengeStorage struct {
	db  *sql.DB
	dsn string
}

func (cs *ChallengeStorage) DSN(dsn string) {
	cs.dsn = dsn
}

func (cs *ChallengeStorage) Open() error {
	var err error

	/*
		if _, err = os.Stat(cs.dsn); errors.Is(err, os.ErrNotExist) {
			err = cs.createDb(cs.dbName)
			if err != nil {
				return err
			}
		}
	*/

	cs.db, err = sql.Open("sqlite3", cs.dsn)
	if err != nil {
		return err
	}

	stmt, err := cs.db.Prepare(sqlCreateTable)
	if err != nil {
		return err
	}
	stmt.Exec()
	defer stmt.Close()

	return nil
}

func (cs *ChallengeStorage) Close() error {
	return cs.db.Close()
}

func (cs *ChallengeStorage) Challenge(time int64, len int) (string, error) {
	var code string

	err := cs.Open()
	defer cs.Close()
	if err != nil {
		fmt.Println(err)
	}

	row := cs.db.QueryRow(sqlSelectTodaysChallenge, time, len)
	err = row.Scan(&code)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return "", err
	case nil:
		return code, nil
	default:
		return "", err
	}
}

func (cs *ChallengeStorage) CreateChallenge(time int64, len int, code string) error {
	stmt, err := cs.db.Prepare(sqlInsertTodaysChallenge)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(len, time, code)
	defer stmt.Close()
	return err
}

func (cs *ChallengeStorage) createDb(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	return err
}
