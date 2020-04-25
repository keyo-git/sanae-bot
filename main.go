package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/keyo-git/sanae-bot/tagset"
	_ "github.com/keyo-git/sanae-bot/viewer"

	"github.com/keyo-git/sanae-bot/bot"
)

var configPath = flag.String("config", "", "path to sanae config file")

func main() {
	flag.Parse()

	cfg := &bot.Config{}
	var err error
	if j, err := ioutil.ReadFile(*configPath); err == nil {
		err = json.Unmarshal(j, cfg)
	}
	if err != nil {
		log.Fatal(err)
	}

	sanae, err := bot.NewSanae(cfg)
	if err != nil {
		log.Fatal(err)
	}
	sanae.Open()
	defer sanae.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
