package app

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func DBconnect(c Configure) bool {
	db, err := sql.Open("pgx", c.Dblink)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return false
	}
	return true
}
func CreateTabledb(c Configure) {
	db, err := sql.Open("pgx", c.Dblink)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS urls(short_url VARCHAR(255) PRIMARY KEY,original_url VARCHAR(255));")
	if err != nil {
		panic(err)
	}
}

func GetOriginalUrldb(shorturl string, c Configure) string {
	db, err := sql.Open("pgx", c.Dblink)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRowContext(context.Background(),
		"SELECT ORIGINAL_URL FROM URLS WHERE SHORT_URL=$1", shorturl)
	var result string
	row.Scan(&result)
	return result
}
func GetShortUrldb(originalurl string, c Configure) string {
	db, err := sql.Open("pgx", c.Dblink)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRowContext(context.Background(),
		"SELECT SHORT_URL FROM URLS WHERE ORIGINAL_URL=$1", originalurl)
	var result string
	row.Scan(&result)
	return result
}
func AddURLdb(shorturl, originalurl string, c Configure) error {
	db, err := sql.Open("pgx", c.Dblink)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.ExecContext(context.Background(), "INSERT INTO URLS (short_url,original_url) VALUES($1,$2) ON CONFLICT (short_url) DO NOTHING", shorturl, originalurl)
	if err != nil {
		return err
	}
	return nil
}
