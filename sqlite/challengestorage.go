package sqlite

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var errNoChal = errors.New("challengestorage: no challenge for given datetime/len")

type ChallengeStorage struct {
	db  *sql.DB
	dsn string
}

func (cs *ChallengeStorage) ErrNoChal() error {
	return errNoChal
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

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("challengestorage: Open(): stmt.Close() failed")
		}
	}(stmt)

	_, err = stmt.Exec()

	return err
}

func (cs *ChallengeStorage) Close() error {
	return cs.db.Close()
}

func (cs *ChallengeStorage) Challenge(time int64, len int) (string, error) {
	var code string

	row := cs.db.QueryRow(sqlSelectTodaysChallenge, time, len)
	err := row.Scan(&code)

	switch err {
	case sql.ErrNoRows:
		return "", errNoChal
	case nil:
		return code, nil
	default:
		return "", err
	}
}

func (cs *ChallengeStorage) Create(time int64, len int, code string) error {
	stmt, err := cs.db.Prepare(sqlInsertTodaysChallenge)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(len, time, code)

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("challengestorage: Create(): stmt.Close() failed")
		}
	}(stmt)

	return err
}

func (cs *ChallengeStorage) createDb(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("challengestorage: createDb(): file.Close() failed")
		}
	}(file)

	return err
}
