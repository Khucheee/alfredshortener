package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func DBconnect(c Configure) bool {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Dblink, `postgres`, `ALFREd2002`, `alfredshortener`)

	db, err := sql.Open("pgx", ps)
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
