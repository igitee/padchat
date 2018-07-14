package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdp/qrterminal"
	"github.com/tuotoo/padchat"
)

func main() {
	bot, err := padchat.NewBot("ws://52.80.34.207:7777")
	if err != nil {
		panic(err)
	}
	if !bot.Init().Success {
		panic("init failed")
	}
	bot.OnQRURL(func(s string) {
		qrterminal.Generate(s, qrterminal.H, os.Stdout)
	})
	data := bot.QRLogin()
	fmt.Printf("login resp data: %+v\n", data)
	bot.OnLogin(func() {
		fmt.Println(string(bot.SyncContact().Data))
	})
	bot.OnMsg(func(msg padchat.Msg) {
		fmt.Printf("%+v\n", msg)
		fmt.Println(strings.Repeat("=", 100))
	})
	select {}
}