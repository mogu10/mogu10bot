package main

import (
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

func addOneBeer(userId int64, size float32) error {
	_, err := db.Exec("insert into beers (user_id, price, created_at, size) values ($1, $2 , $3, $4)",
		userId, 100, time.Now(), size)
	if err != nil {
		return err
	}
	return nil
}

func getBeersOfUser(userId int64) ([]beer, error) {
	rows, err := db.Query("select id, user_id, price, size, created_at from beers where user_id = $1", userId)
	if err != nil {
		panic(err)
	}

	var beers []beer
	for rows.Next() {
		b := beer{}

		err := rows.Scan(&b.id, &b.userId, &b.price, &b.size, &b.createdAt)
		if err != nil {
			return nil, err
		}
		beers = append(beers, b)
	}

	return beers, nil
}

func getBeerCount(userId int64) (float32, error) {
	var sum float32 = 0
	beers, err := getBeersOfUser(userId)

	if err == nil {
		for _, b := range beers {
			sum += b.size
		}
		return sum, nil
	}
	return 0, err
}

func addPriceLastBeer(userId int64, p int) error {
	beers, err := getBeersOfUser(userId)

	if err == nil {
		lastBeer := beers[len(beers)-1]

		_, err = db.Exec("update beers set price = $1 where id = $2", p, lastBeer.id)

		if err != nil {
			return err
		}
		return nil
	}

	return err
}
