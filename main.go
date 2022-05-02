package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

type Dice struct {
	nb, value int
}

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!dices") {
		rollDice(s, m)
	}
}

func getDice(str string) (Dice, error) {
	i := 0
	nb := 1
	var err error
	for ; i < len(str); i++ {
		if !(str[i] >= '0' && str[i] <= '9') {
			break
		}
	}
	if i+1 >= len(str) || str[i] != 'd' {
		return Dice{}, fmt.Errorf("Error parsing dices")
	}
	if str[0] != 'd' {
		nb, err = strconv.Atoi(str[:i])
	}
	if err != nil {
		return Dice{}, err
	}
	value, err := strconv.Atoi(str[i+1:])
	if err != nil {
		return Dice{}, err
	}
	return Dice{nb: nb, value: value}, nil
}

func getRandom(d int) int {
	min := 1
	return rand.Intn(d-min) + min
}

func rollDice(s *discordgo.Session, m *discordgo.MessageCreate) {
	dicesStr := strings.Fields(strings.Replace(m.Content, "!dices", "", 1))
	var dices []Dice

	for _, dicestr := range dicesStr {
		d, err := getDice(dicestr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		dices = append(dices, d)
	}
	alltext := ""
	line := ""
	sum := 0
	for _, d := range dices {
		line = ""
		sum = 0
		line += fmt.Sprintf("d%d: ", d.value)
		for i := 0; i < d.nb; i++ {
			r := getRandom(d.value)
			sum += r
			line += fmt.Sprintf("(%d) ", r)
		}
		line += fmt.Sprintf("sum: %d med: %d", sum, sum/d.nb)
		alltext = alltext + line + "\n"
	}
	_, err := s.ChannelMessageSend(m.ChannelID, alltext)
	if err != nil {
		fmt.Println(err)
	}
}
