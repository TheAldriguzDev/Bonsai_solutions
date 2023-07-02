package protection

import (
	"os"
	"strings"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/discord"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"

	"github.com/mitchellh/go-ps"
)

// We should also add a temp ban to the key, maybe connecting the backend!

type BadApplications []string

var BadApp = BadApplications{
	"fiddler",
	"tcpview",
	"smsniff",
	"socketsniff",
	"charles",
	"mitmweb",
	"mitmdump",
	"burpsuite",
	"burp",
	"fiddle",
	"ghidra",
	"wireshark",
	"ilspy",
	"insomnia",
}

func GetProcess(licenseKey string, userName string) {
	for {
		processes, err := ps.Processes()
		if err != nil {
			panic(err)
		}

		for _, proc := range processes {
			name := proc.Executable()

			for _, appName := range BadApp {
				if strings.Contains(strings.ToLower(name), strings.ToLower(appName)) {
					log.Error("Suspect process found!", "Protection")
					NotifyCrack(appName, licenseKey, userName)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func NotifyCrack(appName string, userLicense string, userName string) {
	whData := discord.DiscordParams{
		Embeds: []discord.Embed{
			{
				Author: discord.Author{
					Name: "Bonsaio",
					//IconURL: "https://i.imgur.com/4cJK9wD.png",
					//URL:     "https://twitter.com/lollipop_ext",
				},
				Title: "Crack attemp found!",
				Color: 9616018,
				Footer: discord.Footer{
					Text: "Bonsaio",
				},
				Fields: []discord.Field{
					{
						Name:   "Application",
						Value:  appName, // Add the discord user here
						Inline: true,
					},
					{
						Name:   "User name",
						Value:  userName, // Add the discord user here
						Inline: true,
					},
					{
						Name:   "License",
						Value:  userLicense, // Add the discord user here
						Inline: true,
					},
				},
			},
		},
	}

	err := discord.PostWebhook("https://discord.com/api/webhooks/1107341933451423754/jB0AzVPBMcg8bSQmEIhVuSXca16GMExhjnZGIKCMNJKy1W_JmcM-xjTKj0z8EUvnvFiE", whData)

	if err != nil {
		panic(err)
	} else {
		log.Success("Bonsaio team has been notified!", "Protection")
		os.Exit(0)
	}
}
