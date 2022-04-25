package main

type user struct {
	tgId      int64  `db:"tg_id"`
	tgNick    string `db:"tg_nick"`
	jokeName  string `db:"joke_name, omitempty"`
	firstName string `db:"first_name, omitempty"`
	lastName  string `db:"last_name, omitempty"`
}

func addUser(u user) error {
	_, err := db.Exec("insert into users (tg_id, tg_nick, first_name, last_name) values ($1, $2 , $3, $4)",
		u.tgId, u.tgNick, u.firstName, u.lastName)
	if err != nil {
		return err
	}
	return nil
}

func getJokeName(tgId int64) (string, error) {
	var jokename string

	row, err := db.Query("select joke_name from users where tg_id = $1", tgId)
	if err == nil {
		row.Next()
		err = row.Scan(&jokename)
		if err == nil {
			if jokename == "" {
				return "Ноу Нейм", nil
			}

			return jokename, nil
		}
	}
	return "", err
}

func checkUser(tgId int64) (int64, error) {
	var ch int64

	row, err := db.Query("select tg_id from users where tg_id = $1", tgId)
	if err != nil {
		return 0, err
	}

	row.Next()
	err = row.Scan(&ch)

	if err != nil {
		return 0, err
	}

	return ch, nil
}
