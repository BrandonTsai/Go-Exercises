package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const dbfile string = "qa.db"

func main() {
	testLog()
	catalogues := []string{
		"Software Engineer",
		"English",
	}

	// connect to database
	db, err := sql.Open("sqlite3", dbfile)
	checkErr(err)
	initDB(db)
	newid := insetLabel(db, catalogues[0], "Golang")

	if newid != -1 {
		updateLabel(db, newid, catalogues[0], "Go")
	} else {
		updateLabel(db, 1, catalogues[0], "Python")
	}

	labels := listLabel(db, catalogues[0])
	for _, value := range labels {
		fmt.Println(value)
	}
	deleteLabel(db, 1)
	db.Close()
}

func testLog() {
	log.Info("A walrus appears")
}

func initDB(db *sql.DB) {

	//Create Catalogue Table if not exist
	const createLabels string = `
	CREATE TABLE IF NOT EXISTS label (
	id INTEGER NOT NULL PRIMARY KEY,
	catalogue STRING NOT NULL,
	name STRING UNIQUE NOT NULL
	);`
	_, err := db.Exec(createLabels)
	checkErr(err)

}

func insetLabel(db *sql.DB, catalogue string, label string) int {
	fmt.Println("Insert label", catalogue, label)
	res, err := db.Exec("INSERT INTO label VALUES(NULL,?, ?);", catalogue, label)
	if err != nil {
		fmt.Printf("Label '%s' has exist in DB\n", label)
		return -1
	}
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
	return int(id)
}

func updateLabel(db *sql.DB, id int, catalogue string, label string) {
	res, err := db.Exec("UPDATE label set catalogue=?,name=?  WHERE id=? ;", catalogue, label, id)
	checkErr(err)
	affects, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affects)
}

func listLabel(db *sql.DB, catalogue string) []string {
	rows, err := db.Query("SELECT id,name from label WHERE catalogue=? ;", catalogue)
	checkErr(err)

	var labels []string
	var id int
	var name string
	for rows.Next() {
		err = rows.Scan(&id, &name)
		checkErr(err)
		fmt.Println(id, name)
		labels = append(labels, name)
	}
	rows.Close() //good habit to close
	return labels
}

func deleteLabel(db *sql.DB, id int) {
	// delete
	res, err := db.Exec("delete from label where id=?", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
