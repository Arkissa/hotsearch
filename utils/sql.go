package utils

import (
	"database/sql"
	"fmt"
	"hotsearch/log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	sqli *sql.DB
}

func NewDatabase(databaseName string) *Database {
	sqlite, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		panic(err)
	}

	db := new(Database)
	db.sqli = sqlite

	return db
}

func (d *Database) CreateTable(tableName string) bool {
	log.LogPut("[INFO] Start Create Table %s\n", tableName)
	_, err := d.sqli.Exec(
		fmt.Sprintf("CREATE TABLE %s (id INTEGER PRIMARY KEY AUTOINCREMENT, hotTitle TEXT, date TEXT NOT NULL DEFAULT (datetime('now','localtime')))",
			tableName))
	if err != nil {
		log.LogOutErr(fmt.Sprintf("create table %s err", tableName), err)
		return false
	}

	return true
}

func (d *Database) FindTable(tableNames ...string) []string {
	smts, err := d.sqli.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.LogOutErr("select tables err", err)
		return nil
	}
	defer smts.Close()

	var (
		tableMap   = make(map[string]bool)
		name       string
		noFindName []string
	)

	for i := 0; i < len(tableNames); i++ {
		tableMap[tableNames[i]] = true
	}

	for smts.Next() {
		if err = smts.Scan(&name); err != nil {
			log.LogOutErr("get table name err", err)
			return nil
		}

		tableMap[name] = false
	}


	for tableName, ok := range tableMap {
		if ok {
			log.LogPut("[WARNING] Not Found %s\n", tableName)
			noFindName = append(noFindName, tableName)
		}
	}

	return noFindName
}

func (d *Database) InsertData(tableName, data string) bool {
	log.LogPut("[INFO] Start Insert %s Data %s\n", tableName, data)
	smts, err := d.sqli.Prepare(fmt.Sprintf("INSERT INTO %s (hotTitle) VALUES (?)", tableName))
	if err != nil {
		log.LogOutErr("select tables err", err)
		return false
	}
	defer smts.Close()

	_, err = smts.Exec(data)
	if err != nil {
		log.LogOutErr("select tables err", err)
		return false
	}

	return true
}

func (d *Database) Deduplication(tableName string) {
	log.LogPut("[INFO] Start Deduplicationing for %s\n", tableName)
	_, err := d.sqli.Exec(fmt.Sprintf("DELETE FROM %s WHERE id NOT in ( SELECT MIN(id) FROM %s GROUP BY hotTitle )",
		tableName, tableName))
	if err != nil {
		log.LogOutErr("delete columns err", err)
	}
}

func (d *Database) Close() {
	if err := d.sqli.Close(); err != nil {
		log.LogOutErr("close database err", err)
		return
	}
}
