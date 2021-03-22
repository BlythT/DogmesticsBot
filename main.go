package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

// Variables used for command line parameters
var (
	Token                string
	lastDogmestic        time.Time
	lastDogmesticMessage string
	dogmesticCount       int
)

func init() {
	// parse in the discord token for authentication
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	dogmesticQuips := []string{
		"No dogmestics in chat!!\n",
		"Omg no dogmestics pls\n",
		"🚨 DOGMESTICS ALERT 🚨\n",
		"Smh dogmestics in the chat\n",
		"Can it be? A dogmestic while I'm on shift!? Not on my watch...\n",
	}

	greetings := []string{
		"Hello!\n",
		"Greetings!\n",
		"Yes?\n",
		"Howdy\n",
		"G'day\n",
		"Kia ora\n",
		"sup\n",
		"it me\n",
		"beep boop\n",
	}

	timeSinceQuips := []string{
		"The last dogmestic was %s!!\n",
		"I haven't seen a dogmestic in these parts since %s.\n",
		"If I remember correctly, the last dogmestic was... %s.\n",
	}

	// Regex is wild (puns intended)
	var dogmesticsRegex = regexp.MustCompile(`(?i).*((d+[\s.,\-_~'\|*+":;]*[o0\(aw\)]+[\s.,\-_~'\|*+":;]*[g9]+.*)|(🐕)|(🐶)).*m+[\s.,-_~'\|*+":;]*[e3]+[\s.,-_~'\|*+":;]*[s5$]+[\s.,-_~'\|*+":;]*[t+]+[\s.,-_~'\|*+":;]*[|1i!]+[\s.,-_~'\|*+":;]*[c(]+`)
	var uGoodRegex = regexp.MustCompile(`(?i)((y+o*u*)|(u+))\s*g[ou]+d+`)
	var nameRegex = regexp.MustCompile(`(?i)🤖|(.*((d+\s*[o0\(aw\)]+[\s.,\-_~'\|*+":;]*[g9]+.*)|(🐕)|(🐶)).*m+[\s.,-_~'\|*+":;]*[e3]+[\s.,-_~'\|*+":;]*[s5$]+[\s.,-_~'\|*+":;]*[t+]+[\s.,-_~'\|*+":;]*[|1i!]+[\s.,-_~'\|*+":;]*[c(]+.*)?bot|robo`)
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if uGoodRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No bro")
		time.Sleep(2 * time.Second)
		s.ChannelMessageSend(m.ChannelID, "I need nintendo switch")
		return
	}

	// Keep track of reported dogmestics (currently not restart resistant: sqllite opportunity?)
	if dogmesticsRegex.MatchString(m.Content) {
		var messageBuilder strings.Builder
		dogmesticCount++
		messageBuilder.WriteString(dogmesticQuips[rand.Intn(len(dogmesticQuips))])

		if !lastDogmestic.IsZero() {
			output := fmt.Sprintf(timeSinceQuips[rand.Intn(len(timeSinceQuips))], humanize.Time(lastDogmestic))
			messageBuilder.WriteString(output)
		}

		if lastDogmesticMessage != "" {
			output := fmt.Sprintf(`The last dogmestic I remember was "%s"`, lastDogmesticMessage)
			messageBuilder.WriteString(output)
		}

		if dogmesticCount != 0 {
			if dogmesticCount == 1 {
				output := fmt.Sprintln(`The first dogmestic today!! Congratz 🥳🎉`)
				messageBuilder.WriteString(output)
			} else {
				output := fmt.Sprintf(`That's %d dogmestics while I've been on duty!!`, dogmesticCount)
				messageBuilder.WriteString(output)
			}
		}

		s.ChannelMessageSend(m.ChannelID, messageBuilder.String())
		lastDogmesticMessage = m.Content
		lastDogmestic = time.Now()
		return
	}

	if nameRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, greetings[rand.Intn(len(greetings))])
		return
	}
}
