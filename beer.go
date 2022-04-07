package main

import (
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

type beer struct {
	id        int
	userId    string `db:"user_id"`
	price     int
	size      float32
	createdAt pgtype.Date
}

func addOneBeer(userId int64, size float32) {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("insert into beers (user_id, price, created_at, size) values ($1, $2 , $3, $4)",
		userId, 100, time.Now(), size)
	if err != nil {
		panic(err)
	}

	defer db.Close()
}

func getBeersOfUser(userId int64) []beer {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select id, user_id, price, size, created_at from beers where user_id = $1", userId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var beers []beer
	for rows.Next() {
		b := beer{}

		err := rows.Scan(&b.id, &b.userId, &b.price, &b.size, &b.createdAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		beers = append(beers, b)
	}

	return beers
}

func getBeerCount(userId int64) float32 {
	var sum float32 = 0
	beers := getBeersOfUser(userId)

	for _, b := range beers {
		sum += b.size
	}

	return sum
}

func addPriceLastBeer(userId int64, p int) {
	beers := getBeersOfUser(userId)
	lastBeer := beers[len(beers)-1]

	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("update beers set price = $1 where id = $2", p, lastBeer.id)

	if err != nil {
		panic(err)
	}

	defer db.Close()
}
