package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"strings"
	"time"
)

type Dburls struct {
	Shorturl    string `json:"short_url"`
	Originalurl string `json:"original_url"`
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

	_, err = db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS urls(user_id VARCHAR(36),short_url VARCHAR(255) PRIMARY KEY,original_url VARCHAR(255),deleted BOOLEAN);")
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
	var tmp Dburls
	for rows.Next() {
		err = rows.Scan(&tmp.Shorturl, &tmp.Originalurl)
		if err != nil {
			fmt.Println("Что-то упало на сканировании файла:", err)
		}
		urls[tmp.Shorturl] = tmp.Originalurl
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

func (d *Database) Save(shorturl, originalurl, uuid string) {
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.ExecContext(context.Background(),
		"INSERT INTO URLS (user_id,short_url,original_url)VALUES($1,$2,$3) ON CONFLICT (short_url) DO NOTHING",
		uuid, shorturl, originalurl)
	if err != nil {
		fmt.Println("Что-то упало при сохранении в базу:", err)
	}
}

func (d *Database) GetUrlsByUser(uuid string) []Dburls {
	fmt.Println("сработал метод запроса урлов пользователя в базе")
	urls := []Dburls{}
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(context.Background(),
		"SELECT SHORT_URL,ORIGINAL_URL FROM URLS WHERE user_id = $1", uuid)
	if err != nil {
		fmt.Println("Это ошибка запроса урлов пользователя", err)
		return urls
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Ошибка в чтении строк в таблице:", err)
	}
	var tmp Dburls
	for rows.Next() {
		err = rows.Scan(&tmp.Shorturl, &tmp.Originalurl)
		if err != nil {
			fmt.Println("Что-то упало на сканировании файла:", err)
		}
		urls = append(urls, Dburls{tmp.Shorturl, tmp.Originalurl})
	}
	fmt.Println("Тут должна быть структурура урлов", urls)
	return urls
}

func (d *Database) DeleteUserLinks(uid string, hashes []string) {
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dq := "UPDATE shortener SET deleted=true WHERE user_id=$1 AND hash=ANY($2::text[])"
	params := "{" + strings.Join(hashes, ",") + "}"
	_, err = db.ExecContext(context.Background(), dq,
		uid, params)
	if err != nil {
		log.Printf("DeleteUserLinks error: %#v \n", err)
	}
}
