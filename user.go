package main

import (
	"database/sql"
	"fmt"
)

type user struct {
	tgId      int64  `db:"tg_id"`
	tgNick    string `db:"tg_nick"`
	jokeName  string `db:"joke_name, omitempty"`
	firstName string `db:"first_name, omitempty"`
	lastName  string `db:"last_name, omitempty"`
}

// a.trofimenko #2
// на примере этой же функции. Вопрос по работе с БД в проекте
// допустима ли такая работа с базой данных. Без универсальной функции добавления в БД.
// пример: сейчас у меня есть 2 функции addUser() и addOneBeer(). Обе добавляют запись в базу
// а могла быть функция в db.go addRow(), которой передавались бы нужные параметры. Улучшит ли это код?

//P.S. Любые замечания не по темам вопросов приветствуются )

func addUser(u user) {
	db, err := connectDB()
	if err != nil { //проверяем, что коннект успешен
		panic(err)
	}

	_, err = db.Exec("insert into users (tg_id, tg_nick, first_name, last_name) values ($1, $2 , $3, $4)",
		u.tgId, u.tgNick, u.firstName, u.lastName)
	if err != nil { //проверяем, что не было ошибок при добавлении юзера
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
}

func getJokeName(tgId int64) string {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	var jokename string

	row, err := db.Query("select joke_name from users where tg_id = $1", tgId)
	if err != nil {
		panic(err)
	}

	row.Next()
	err = row.Scan(&jokename)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	if jokename == "" {
		return "Ноу Нейм"
	}

	return jokename
}

func checkUser(tgId int64) int64 {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	var ch int64
	row, err := db.Query("select tg_id from users where tg_id = $1", tgId)
	if err != nil {
		panic(err)
	}

	row.Next()
	err = row.Scan(&ch)

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	return ch
}
