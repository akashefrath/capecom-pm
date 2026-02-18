package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
)

func Migrate(DB *sqlx.DB) {

	// ensure migrations table exists
	DB.MustExec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(50) PRIMARY KEY,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(files)

	for _, file := range files {

		version := filepath.Base(file)

		var exists string
		err := DB.Get(&exists,
			"SELECT version FROM schema_migrations WHERE version=?",
			version,
		)

		if err == nil {
			fmt.Println("skip", version)
			continue
		}

		fmt.Println("running", version)

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		tx := DB.MustBegin()

		_, err = tx.Exec(string(sqlBytes))
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			log.Fatalf("migrations failed %s: %v", version, err)
		}

		_, err = tx.Exec("INSERT INTO schema_migrations(version) VALUES(?)", version)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			log.Fatal(err)
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	fmt.Println("migrations done")
}
