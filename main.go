package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"	
	"strings"

	"github.com/robfig/cron"
	irc "github.com/thoj/go-ircevent"
)

const server = "bnc.snoonet.org:5457"
const nick = "Moter8"
const password = ""

func main() {
	updateTopic()
	c := cron.New()
	c.AddFunc("@hourly", updateTopic)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}

func updateTopic() {
	ircnick := nick
	irccon := irc.IRC(ircnick, ircnick)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.Password = password
	irccon.UseTLS = true


	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join("#dach")
	})

	irccon.AddCallback("332", func(e *irc.Event) {
		fmt.Println(e.Raw)
		channel := strings.Split(e.Raw, " ")[3]

		if (channel != "#dach") {
			fmt.Println("wrong channel", channel)
			return
		}
		fmt.Println("MESSAGE" + e.Message())

		topic := e.Message()
		index := strings.LastIndex(topic, "Gewicht:")
		if(index < 0) {
			return
		}

		topic = topic[:index+1] + " "
		topic += getFormattedWeight(getCurrentWeight())
		irccon.SendRaw("TOPIC " + channel + " " + topic)

		irccon.Quit()
	})

	err := irccon.Connect(server)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	irccon.Loop()
}

func getFormattedWeight(weight string) string {
	formattedNick := weight + " kg"
	fmt.Println(formattedNick)
	return formattedNick
}

func getCurrentWeight() string {
	path, err := exec.LookPath("python3")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(path, "main.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(string(slurp))

}
