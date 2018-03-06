package main

import (
	"gopush/utils"
	"os"
	"fmt"
	"path/filepath"
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"time"
	"strings"
	"strconv"
	"sort"
)

type Migration struct {
	id int
	name string
}

func fileExist(a string, list []Migration) bool {
	for _, b := range list {
		if b.name == a {
			return true
		}
	}
	return false
}

func sortOf(migrationPath string, fileListUp []string, asc bool) []string {
	var tmpArr []string
	filesMap := make(map[int]string)
	filesKeys := make([]int, 0)
	for _, file := range fileListUp {
		filename := file[len(migrationPath)+1:]
		filename = filename[:len(filename)-4]
		filenameParts := strings.Split(filename, "_")
		i, _ := strconv.Atoi(filenameParts[0])
		filesMap[i] = file
		filesKeys = append(filesKeys, i)
	}

	sort.Ints(filesKeys)

	if !asc {
		sort.Sort(sort.Reverse(sort.IntSlice(filesKeys)))
	}

	for _, key := range filesKeys {
		tmpArr = append(tmpArr, filesMap[key])
	}

	return tmpArr
}

func executeFile(db *sql.DB, path string, complete func(db *sql.DB)) {
	content, err := ioutil.ReadFile(path)
	expressions := strings.Split(string(content), ";")
	utils.Check(err)

	tx, err := db.Begin()
	utils.Check(err)
	defer tx.Rollback()

	for _, expression := range expressions {
		stmt, errTransaction := tx.Prepare(expression)
		utils.Check(errTransaction)
		_, errTransaction = stmt.Exec()
		utils.Check(errTransaction)
		stmt.Close()
	}

	err = tx.Commit()
	utils.Check(err)

	if complete != nil {
		complete(db)
	}
}

func getMigrations(db *sql.DB) []Migration {
	var (
		result []Migration
	)

	rows, queryErr := db.Query("SELECT id, name FROM migrations ORDER BY id DESC")
	for rows.Next() {
		migration := Migration{}
		err := rows.Scan(&migration.id, &migration.name)
		utils.Check(err)
		result = append(result, migration)
	}

	queryErr = rows.Err()
	utils.Check(queryErr)

	return result
}

func saveMigration(db *sql.DB, name string) {
	tx, err := db.Begin()
	utils.Check(err)
	defer tx.Rollback()
	stmt, errTransaction := tx.Prepare("INSERT INTO migrations(name) VALUES ($1)")
	utils.Check(errTransaction)
	defer stmt.Close()

	_, errTransaction = stmt.Exec(name)
	utils.Check(errTransaction)

	errTransaction = tx.Commit()
	utils.Check(errTransaction)
}

func deleteMigration(db *sql.DB, name string) {
	tx, err := db.Begin()
	utils.Check(err)
	defer tx.Rollback()
	stmt, errTransaction := tx.Prepare("DELETE FROM migrations WHERE name = $1")
	utils.Check(errTransaction)
	defer stmt.Close()

	_, errTransaction = stmt.Exec(name)
	utils.Check(errTransaction)

	errTransaction = tx.Commit()
	utils.Check(errTransaction)
}

func main() {
	config := utils.NewLoadConfig(os.Args[1:])
	migrationPath := config.GetStr("migration_path")

	db, errConnection := sql.Open(config.GetStr("db_driver"), config.GetStr("db_url"))
	utils.Check(errConnection)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	defer db.Close()

	/**
	 * Create table
	 */
	migrationTable := fmt.Sprintf("%v/%v.sql", migrationPath, "migrations")
	if _, err := os.Stat(migrationTable); os.IsNotExist(err) {
		utils.CreateFile(migrationTable, "CREATE TABLE migrations (id SERIAL PRIMARY KEY, name VARCHAR(255) NULL);")
		executeFile(db, migrationTable,  nil)
	}

	/**
	 * Create migration
	 */
	name := config.GetParam("create")
	if len(name) > 0 {
		t := time.Now().Unix()
		utils.CreateFile(fmt.Sprintf("%v/%v_%v.sql", migrationPath, t, name), nil)
		utils.CreateFile(fmt.Sprintf("%v/%v_%v_down.sql", migrationPath, t, name), nil)
		return
	}

	up := config.GetParam("up")
	steps, _ := strconv.Atoi(up)

	down := config.GetParam("down")
	downSteps, _ := strconv.Atoi(down)

	existMigrations := getMigrations(db)

	if steps <= 0 && downSteps <= 0 {
		fmt.Println("Input up or down steps")
	}

	/**
	 * Up
	 */
	if steps > 0 {
		fileListUp := []string{}
		fileListDown := []string{}
		err := filepath.Walk(migrationPath, func(path string, f os.FileInfo, err error) error {
			if path != migrationPath {
				isDown := strings.HasSuffix(path, "down.sql")
				if isDown {
					fileListDown = append(fileListDown, path)
				} else {
					fileListUp = append(fileListUp, path)
				}
			}
			return nil
		})

		fileListUp = sortOf(migrationPath, fileListUp, true)

		utils.Check(err)

		totalSteps := 0
		for _, file := range fileListUp {
			if totalSteps > steps - 1 { break }
			filename := file[len(migrationPath)+1:]
			filename = filename[:len(filename)-4]
			filenameParts := strings.Split(filename, "_")

			i, _ := strconv.Atoi(filenameParts[0])
			if !fileExist(filename, existMigrations) && i > 0 {
				executeFile(db, file, func(db *sql.DB) {
					saveMigration(db, filename)
				})
				totalSteps++
			}
		}

		return
	}

	/**
	 * Down
	 */
	if downSteps > 0 {
		for step, file := range existMigrations {
			if step > downSteps - 1 { break }
			fileToDown := fmt.Sprintf("%v/%v_down.sql", migrationPath, file.name)
			if _, err := os.Stat(fileToDown); !os.IsNotExist(err) {
				executeFile(db, fileToDown, func(db *sql.DB) {
					deleteMigration(db, file.name)
				})
			}
		}
	}
}