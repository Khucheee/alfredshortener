package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type Dburls struct {
	Shorturl    string `json:"short_url"`
	Originalurl string `json:"original_url"`
}
type Database struct {
	link string
	db   *sql.DB
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

func (d *Database) CreateTabledb() {
	db, err := sql.Open("pgx", d.link)
	if err != nil {
		fmt.Println(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	m, err := migrate.NewWithDatabaseInstance("file://../../internal/app/migrations/", "postgres", driver)
	fmt.Println(err)
	err = m.Up()
	if err != nil {
		fmt.Println("Упала миграция", err)
	}
	if err != nil {
		panic(err)
	}
	d.db = db
	//_, err = db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS urls(user_id VARCHAR(36),short_url VARCHAR(255) PRIMARY KEY,original_url VARCHAR(255),deleted BOOLEAN DEFAULT false);")

}

func (d *Database) Restore() map[string]URLData {
	fmt.Println("сработал метод рестор базы")
	urls := make(map[string]URLData)
	rows, err := d.db.QueryContext(context.Background(),
		"SELECT USER_ID,SHORT_URL,ORIGINAL_URL,DELETED FROM URLS")
	if err != nil {
		fmt.Println("Это ошибка запроса урлов,вернется пустая мапа при восстановлении:", err)
		return urls
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Ошибка в чтении строк в таблице:", err)
	}
	var tmp Dburls
	var isdeleted bool
	var uuid string
	for rows.Next() {
		err = rows.Scan(&uuid, &tmp.Shorturl, &tmp.Originalurl, &isdeleted)
		if err != nil {
			fmt.Println("Что-то упало на сканировании файла:", err)
		}
		urls[tmp.Shorturl] = URLData{
			originalurl: tmp.Originalurl,
			uuid:        uuid,
			isdeleted:   isdeleted,
		}
	}
	return urls
}

func (d *Database) GetShortUrldb(originalurl string) string {
	row := d.db.QueryRowContext(context.Background(),
		"SELECT SHORT_URL FROM URLS WHERE ORIGINAL_URL=$1", originalurl)
	var result string
	row.Scan(&result)
	return result
}

func (d *Database) Save(shorturl, originalurl, uuid string) {
	_, err := d.db.ExecContext(context.Background(),
		"INSERT INTO URLS (user_id,short_url,original_url)VALUES($1,$2,$3) ON CONFLICT (short_url) DO NOTHING",
		uuid, shorturl, originalurl)
	if err != nil {
		fmt.Println("Что-то упало при сохранении в базу:", err)
	}
}

func (d *Database) GetUrlsByUser(uuid string) []Dburls {
	fmt.Println("сработал метод запроса урлов пользователя в базе")
	urls := []Dburls{}
	rows, err := d.db.QueryContext(context.Background(),
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

func (d *Database) DeleteUserLink(uid string, hash string) {
	dq := "UPDATE urls SET deleted=true WHERE user_id=$1 AND short_url=$2"
	_, err := d.db.ExecContext(context.Background(), dq, uid, hash)
	if err != nil {
		log.Printf("DeleteUserLinks error: %#v \n", err)
	}

}

//старая реализация
/*func (d *Database) DeleteUserLinks(uid string, hashes []string) {
	dq := "UPDATE urls SET deleted=true WHERE user_id=$1 AND short_url=ANY($2::text[])"
	params := "{" + strings.Join(hashes, ",") + "}"
	_, err := d.db.ExecContext(context.Background(), dq, uid, params)
	if err != nil {
		log.Printf("DeleteUserLinks error: %#v \n", err)
	}

}
*/
