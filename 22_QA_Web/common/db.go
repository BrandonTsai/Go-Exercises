package common

//Table: labels
// id, name

//Table: qa
// id, question, answer

//Table: labelandqa
// lid, qid

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const dbfile string = "qa.db"

func InitDB() {
	db, err := sql.Open("sqlite3", dbfile)
	checkErr(err)
	//Create Catalogue Table if not exist
	const createLabels string = `
	CREATE TABLE IF NOT EXISTS labels (
	id INTEGER NOT NULL PRIMARY KEY,
	name STRING UNIQUE NOT NULL
	);`
	_, err = db.Exec(createLabels)
	checkErr(err)
	log.Info("Table labels created")

	const createQuestions string = `
	CREATE TABLE IF NOT EXISTS questions (
	id INTEGER NOT NULL PRIMARY KEY,
	question STRING UNIQUE NOT NULL,
	answer STRING NOT NULL
	);`
	_, err = db.Exec(createLabels)
	checkErr(err)
	log.Info("Table questions created")

}
