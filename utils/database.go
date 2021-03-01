package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func GetDatabase() *sql.DB {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't retrieve current user, error was: %s", err))
	}

	affectumDir := filepath.Join(usr.HomeDir, "/affectum")
	database := filepath.Join(affectumDir, "affectum.sql")

	sqliteDatabase, _ := sql.Open("sqlite3", database)

	return sqliteDatabase
}

func SetupDatabase() bool {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't retrieve current user, error was: %s", err))
	}

	// @todo: merge with CreateDir().
	affectumDir := filepath.Join(usr.HomeDir, "/affectum")
	database := filepath.Join(affectumDir, "affectum.sql")
	if _, err := os.Stat(database); err != nil {
		if os.IsNotExist(err) {
			os.Create(database)
		}
	}

	sqliteDatabase, _ := sql.Open("sqlite3", database)
	defer sqliteDatabase.Close()

	// Check if we have the table already.
	_, tableCheck := sqliteDatabase.Query("select * from affectum")
	if tableCheck != nil {
		createTable(sqliteDatabase)
	}

	return true
}

func createTable(db *sql.DB) {
	createAffectumTableSQL := `CREATE TABLE affectum (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"folder" TEXT,
		"mid" integer
	  );`

	log.Println("Setting up table ...")
	statement, err := db.Prepare(createAffectumTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()

	log.Println("Table created successfully")
}

func InsertMail(folder string, mid uint32) {
	db := GetDatabase()
	defer db.Close()

	log.Println("Inserting record ...")
	insertStudentSQL := `INSERT INTO affectum(folder, mid) VALUES (?, ?)`
	statement, err := db.Prepare(insertStudentSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = statement.Exec(folder, mid)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func MailExists(folder string, mid uint32) bool {
	var output string
	db := GetDatabase()
	defer db.Close()

	query, err := db.Prepare(`SELECT count(mid) FROM affectum WHERE folder= ? AND mid= ? ORDER BY mid`)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	err = query.QueryRow(folder, mid).Scan(&output)

	if output == "0" {
		return false
	}

	return true
}
