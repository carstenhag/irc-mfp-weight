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
	irccon.Debug = false
	irccon.Password = password
	irccon.UseTLS = true

	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join("#dach")
	})

	irccon.AddCallback("332", func(e *irc.Event) {
		fmt.Println(e.Raw)
		channel := strings.Split(e.Raw, " ")[3]
		oldTopic := strings.TrimSpace(e.Message())

		if channel != "#dach" {
			fmt.Println("wrong channel", channel)
			return
		}

		topic := oldTopic
		index := strings.LastIndex(topic, "Gewicht:")
		if index < 0 {
			return
		}

		topic = topic[:index+8] + " "
		weight := strings.TrimSpace(getCurrentWeight())
		if weight == "" {
			return
		}
		topic += getFormattedWeight(weight)

		topic = strings.TrimSpace(topic)
		if oldTopic != topic {
			fmt.Println("\n\n", oldTopic, "\n\n", topic, "\n\n")
			irccon.SendRaw("TOPIC " + channel + " " + topic)
		}

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
	formattedNick := weight + "kg"
	fmt.Println(formattedNick)
	return formattedNick
}

func getCurrentWeight() string {
	path, err := exec.LookPath("python3")
	if err != nil {
		fmt.Println("Error finding python3 executable in PATH")
		log.Fatal(err)
	}

	cmd := exec.Command(path, "main.py")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error opening Stdout")
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting the command/opening the script")
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error with the command, see https://golang.org/pkg/os/exec/#Cmd.Wait")
		fmt.Println(err) // Prints error, but doesn't exit. Just tries again in the next hour.
	}

	return strings.TrimSpace(string(slurp))

}
