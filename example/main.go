package main

import (
	"fmt"
	"bytes"

	"github.com/Baozisoftware/qrcode-terminal-go"
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
		obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.NormalBlack,
			qrcodeTerminal.ConsoleColors.BrightWhite,
			qrcodeTerminal.QRCodeRecoveryLevels.Low)
		obj.Get([]byte(s)).Print()
	})
	bot.QRLogin()
	bot.OnLogin(func() {
		fmt.Println("login success")
	})
	bot.OnMsg(func(msg padchat.Msg) {
		fmt.Println(msg.MType, msg.FromUser, msg.ToUser)
		if msg.MType == 49 {
			if bytes.Contains(msg.Content, []byte("<![CDATA[微信红包]]>")) {
				rec, err := bot.ReceiveRedPacket(msg)
				if err != nil {
					fmt.Println("receive red packet", err)
					return
				}
				fmt.Println(bot.OpenRedPacket(msg, rec.Key))
			} else if bytes.Contains(msg.Content, []byte("<![CDATA[微信转账]]>")) {
				rec, err := bot.QueryTransfer(msg)
				if err != nil {
					fmt.Println("query tranfer", err)
					return
				}
				fmt.Println("query success", rec)
				fmt.Println("accept transfer")
				fmt.Println(bot.AcceptTransfer(msg))
			}
		}
	})
	select {}
}
