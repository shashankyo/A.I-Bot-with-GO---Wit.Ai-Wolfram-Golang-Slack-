package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go"
	witai "github.com/wit-ai/wit-go/v2"
)

var wolframClient *wolfram.Client

func printCommandeEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()

	}
}

func main() {

	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}
	go printCommandeEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send any question to welfram",
		Examples:    "who is the president of india",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})

			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:])
			value := gjson.Get(rough, "entities.wit$wolfram_search_wuery:wolfram_search_query")
			answer := value.String()
			res, err := wolframClient.GetSpokenAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("there is an error ")
			}

			fmt.Println(rough)
			fmt.Println(msg)

			response.Reply("received ")

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
