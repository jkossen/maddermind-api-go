package sqlite

const sqlCreateTable string = `
CREATE TABLE IF NOT EXISTS "mmdaily" (
	"id"		INTEGER,
	"length"	INTEGER,
	"date"		INTEGER,
	"challenge"		TEXT,
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
