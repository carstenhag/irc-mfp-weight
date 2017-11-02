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

	c := cron.New()
	c.AddFunc("@hourly", updateNick)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}

func updateNick() {
	ircnick := "Moter"
	irccon := irc.IRC(ircnick, ircnick)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.Password = password
	irccon.UseTLS = true

	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join("#dach")
		go irccon.Nick(getFormattedNick(getCurrentWeight()))
		irccon.Quit()
	})

	err := irccon.Connect(server)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	irccon.Loop()
}

func getFormattedNick(weight string) string {
	formattedNick := "Moter" + strings.Replace(weight, ".", "`", 1) + "kg"
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
