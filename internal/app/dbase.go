package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type dburls struct {
	shorturl    string
	originalurl string
}
type Database struct {
	link string
}

func DBconnect(c *Configure) bool {
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

func CreateTabledb(c *Configure) {
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

func (d *Database) Restore() map[string]string {
	fmt.Println("сработал метод рестор базы")
	urls := make(map[string]string)
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(),
		"SELECT SHORT_URL,ORIGINAL_URL FROM URLS")
	if err != nil {
		fmt.Println("Это ошибка запроса урлов,вернется пустая мапа при восстановлении:", err)
		return urls
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Ошибка в чтении строк в таблице:", err)
	}
	var tmp dburls
	for rows.Next() {
		err = rows.Scan(&tmp.shorturl, &tmp.originalurl)
		if err != nil {
			fmt.Println("Что-то упало на сканировании файла:", err)
		}
		urls[tmp.shorturl] = tmp.originalurl
	}
	return urls
}

func (d *Database) GetShortUrldb(originalurl string) string {
	db, err := sql.Open("pgx", d.link)
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

func (d *Database) Save(shorturl, originalurl string) {
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.ExecContext(context.Background(), "INSERT INTO URLS (short_url,original_url) VALUES($1,$2) ON CONFLICT (short_url) DO NOTHING", shorturl, originalurl)
	if err != nil {
		fmt.Println("Что-то упало при сохранении в базу:", err)
	}
}
