
package db

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/packr/v2"
  "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var migrations *migrate.PackrMigrationSource

func init() {
	migrations = &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./migrations"),
	}

	ms, err := migrations.FindMigrations()
	if err != nil {
		panic(err)
	}

	if len(ms) == 0 {
		panic(fmt.Errorf("no SQL migrations found"))
	}
}

func MigrationsUp(db *sql.DB) error {
	nMigratnios, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	fmt.Println("Successfully ran", nMigratnios, "Migrations Up")
	return nil
}

func MigrationsDown(db *sql.DB, migrationsToRun int) error {
	nMigrations, err := migrate.ExecMax(db, "postgres", migrations, migrate.Down, int(migrationsToRun))
	if err != nil {
		return err
	}
	if migrationsToRun != 0 && nMigrations != int(migrationsToRun) {
		return err
	}
	fmt.Println("successfully ran ", nMigrations, "Migrations Down")
	return nil
}

func ConnectSQLDB(port int, host, user, password, name string) (*sqlx.DB, error)  {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		name,
	)
	db, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitSQLDB(port int, host, user, password , name string) (*sqlx.DB, error) {
  db, err := ConnectSQLDB(port,host, user,password, name)
  if err != nil {
    return nil, err
  }

	// Run DB migrations
	if err := MigrationsUp(db.DB); err != nil {
		return nil, err
	}

  return db, nil
}
