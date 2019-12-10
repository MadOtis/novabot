package main

import (
	"github.com/madotis/novabot/bot"
	"github.com/madotis/novabot/config"
)

func main() {
    err := config.ReadConfig()
    if err != nil {
    	panic(err.Error())
	}

	bot.Start()

    <-make(chan struct{})
    return
}

