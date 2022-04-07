package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strconv"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
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

			if checkUser(u.tgId) != u.tgId { //Проверит есть ли такой пользователь в базе

				//a.trofimenko #1
				//след. функция добавляет юзера в таблицу. Не возвращает ошибку.
				//вопрос: по-хорошему должна ли она возвращать ее, если все возможные ошибки уже учтены в самой функции и во вложенной?
				//или я не все учел?
				//(исключая момент, когда мне нужно в связи с ошибкой что-то сообщить пользователю. Здесь нужно было бы вернуть err)
				addUser(*u)

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
					addPriceLastBeer(u.tgId, p)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
						"цена последнего пива изменена на "+strconv.Itoa(p)))
					continue
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
				addOneBeer(u.tgId, 0.5)

				f := tgbotapi.FileURL("https://i.imgur.com/unQLJIb.jpg") //отсылем какую-то картинку
				bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, f))

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Нольпяшка добавлена. Всосано литров: "+fmt.Sprintf("%.1f", getBeerCount(u.tgId))))

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Чтобы указать цену этого пива напиши **/цена ххх**"))

			case "Добавь литрушку":
				addOneBeer(u.tgId, 1)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Литрушка добавлена. Всосано литров: "+fmt.Sprintf("%.1f", getBeerCount(u.tgId))))

			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			case `^/цена...`:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Укажи цену педыдущей бутылки в рублях."))

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Слышь, "+getJokeName(update.Message.From.ID)+", нет такой команды в боте "))
			}
		}
	}
}

func addBeerButtons(msg *tgbotapi.MessageConfig) { //добавляет кнопки к следующему сообщению

	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			//tgbotapi.NewKeyboardButton("Добавь 0.3"), //пушто это объем детской соски
			tgbotapi.NewKeyboardButton("Добавь 0.5"),
			tgbotapi.NewKeyboardButton("Добавь литрушку"),
			//tgbotapi.NewKeyboardButton("Статистика"),
		),
	)

	msg.ReplyMarkup = keyboard
}
