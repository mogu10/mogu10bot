package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strconv"
)

func main() {
	var err error

	connStr := "user=postgres password=paketik26 dbname=mogu10botdb sslmode=disable"
	db, err = connectDB(connStr)
	if err != nil {
		log.Println(err)
	}

	bot, err := tgbotapi.NewBotAPI("5149822295:AAHV3IXmoKxEw0wraewx7tXfAEA12sPQaIk")
	if err != nil {
		log.Println(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var keyboard = tgbotapi.NewReplyKeyboard()

	c := tgbotapi.NewUpdate(0)
	c.Timeout = 60

	updates := bot.GetUpdatesChan(c)
	u := new(user)

	for update := range updates {
		if update.Message != nil {

			u.tgId = update.Message.Chat.ID
			u.tgNick = update.Message.Chat.UserName
			u.firstName = update.Message.Chat.FirstName
			u.lastName = update.Message.Chat.LastName

			uId, err := checkUser(u.tgId)
			if err != nil {
				log.Println(err)
			}

			if uId != u.tgId { //Проверит есть ли такой пользователь в базе

				err = addUser(*u)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"Произошла ошибка при проверке пользователя"))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Приветствую нового алкаша в нашей пати! Нас пока только *мало*. "+
						"Ну что, "+u.firstName+", по пиву?"))
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyMarkup = keyboard

			matched, err := regexp.MatchString(`^/цена *[0-9]`, update.Message.Text) //команда ввода цены
			if matched && err == nil {
				re := regexp.MustCompile("[0-9]+")
				str := re.FindString(update.Message.Text)
				p, err := strconv.Atoi(str)
				if err == nil {
					if addPriceLastBeer(u.tgId, p) == nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
							"цена последнего пива изменена на "+strconv.Itoa(p)))
						continue
					}
				}
			}

			switch update.Message.Text {
			case "я пью пиво", "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ну, тогда начнем считать, чтобы ты ничего не забыл, пьяница")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

				addBeerButtons(&msg)
				bot.Send(msg)

			//зачем 0.3? Они чо корону пьют?
			case "Добавь 0.3", "0.3", "добавь 0.3":

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Семья это, конечно, главное. Но твоя судя по всему от тебя отказалась.\nПей нормальные объемы. А то как девка"))

			case "Добавь 0.5":
				err = addOneBeer(u.tgId, 0.5)

				if err == nil {
					f := tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg") //отсылаем какую-то картинку
					bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, f))

					cBeers, err := getBeerCount(u.tgId)

					if err == nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
							"Нольпяшка добавлена. Всосано литров: "+fmt.Sprintf("%.1f", cBeers)))

						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
							"Чтобы указать цену этого пива напиши **/цена ххх**"))
					}

				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"Ошибка, бутылка не была добавлена."))
				}

			case "Добавь литрушку":
				err = addOneBeer(u.tgId, 1)
				if err == nil {
					cBeers, err := getBeerCount(u.tgId)

					if err == nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
							"Литрушка добавлена. Всосано литров: "+fmt.Sprintf("%.1f", cBeers)))

						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
							"Чтобы указать цену этого пива напиши **/цена ххх**"))
					}
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"Ошибка, бутылка не была добавлена."))
				}

			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			case `^/цена...`:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Укажи цену педыдущей бутылки в рублях."))

			default:
				jName, err := getJokeName(update.Message.From.ID)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"Слышь, "+jName+", нет такой команды в боте "))
				}
			}
		}
	}
}

func addBeerButtons(msg *tgbotapi.MessageConfig) { //добавляет кнопки к следующему сообщению

	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавь 0.3"), //пушто это объем детской соски
			tgbotapi.NewKeyboardButton("Добавь 0.5"),
			tgbotapi.NewKeyboardButton("Добавь литрушку"),
			//tgbotapi.NewKeyboardButton("Статистика"),
		),
	)

	msg.ReplyMarkup = keyboard
}
