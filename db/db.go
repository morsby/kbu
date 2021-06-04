package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/morsby/kbu"
	_ "modernc.org/sqlite"
)

func Connect() *sql.DB {
	os.Remove("db.sqlite") // I delete the file to avoid duplicated records. SQLite is a file based database.

	log.Println("Creating db.sqlite...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("db.sqlite created")

	db, err := sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	return db
}

func CreateTables(db *sql.DB) error {
	createRegionSQL := `CREATE TABLE region (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
	);`

	statement, err := db.Prepare(createRegionSQL)
	if err != nil {
		return err
	}
	statement.Exec()

	createUniversitySQL := `CREATE TABLE university (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
	);`

	statement, err = db.Prepare(createUniversitySQL)
	if err != nil {
		return err
	}
	statement.Exec()

	createRoundSQL := `CREATE TABLE round (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"year" integer
		"url" TEXT
		"season" TEXT
	);`

	statement, err = db.Prepare(createRoundSQL)
	if err != nil {
		return err
	}
	statement.Exec()

	createSelectionSQL := `CREATE TABLE selection (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"md5" TEXT,
		"round_id" integer, 
		"date" TEXT,
		"number" integer,
		"relnumber" REAL,
		"region_id" INTEGER,
		"start" TEXT,
		FOREIGN KEY(round_id) REFERENCES round(id),
		FOREIGN KEY(region_id) REFERENCES region(id)
	);
	
	CREATE UNIQUE INDEX idx_selections_md5 ON selection (md5);`

	statement, err = db.Prepare(createSelectionSQL)
	if err != nil {
		return err
	}
	statement.Exec()

	createPositionSQL := `CREATE TABLE position (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"location" TEXT,
		"department" TEXT,
		"specialty" TEXT,
		FOREIGN KEY(selection_id) REFERENCES selection(id)
	)`

	statement, err = db.Prepare(createPositionSQL)
	if err != nil {
		return err
	}
	statement.Exec()

	return nil
}

type Seeds struct {
	Regions      []kbu.Region
	Universities []kbu.University
}

func InsertRegions(db *sql.DB, regions []kbu.Region) (sql.Result, error) {
	if len(regions) < 1 {
		return nil, nil
	}
	query := "INSERT INTO region(name) VALUES "
	args := []interface{}{}
	for i, region := range regions {
		if i > 0 {
			query += ", "
		}
		query += "(?)"
		args = append(args, region)
	}

	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	return statement.Exec(args...)
}

func InsertUniversities(db *sql.DB, universities []kbu.University) (sql.Result, error) {
	if len(universities) < 1 {
		return nil, nil
	}
	query := "INSERT INTO university(name) VALUES "
	args := []interface{}{}
	for i, university := range universities {
		if i > 0 {
			query += ", "
		}
		query += "(?)"
		args = append(args, university)
	}

	statement, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	return statement.Exec(args...)
}

func Seed(db *sql.DB, seeds Seeds) error {
	// Seeds the DB with constants:
	_, err := InsertRegions(db, seeds.Regions)
	if err != nil {
		return err
	}

	_, err = InsertUniversities(db, seeds.Universities)
	if err != nil {
		return err
	}

	return nil
}

func InsertSelection(db *sql.DB, selections []kbu.Selection) (sql.Result, error) {
	return nil, nil
}
