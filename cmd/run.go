/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"

	bb "github.com/I0HuKc/baitbotnotbytebot/internal/bot"
	"github.com/I0HuKc/baitbotnotbytebot/internal/bot/joker"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/rdstore"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/sqlstore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "r",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		botApi, err := tgbotapi.NewBotAPI(os.Getenv("APP_BOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		pg, err := db.SetPgConn(ctx, os.Getenv("APP_DB_URL"))
		if err != nil {
			log.Panic(err)
		}
		defer pg.Close()

		rclient, err := rdstore.NewClient(ctx)
		if err != nil {
			log.Panic(err)
		}
		defer rclient.Close()

		rstore := rdstore.CreateRedisStore(rclient)

		var js joker.JSchema
		if err := js.ParseSchema(os.Getenv("JOKER_SCHEMA_PATH")); err != nil {
			log.Panic(err)
		}

		baitbot := bb.CreateBaitbot(botApi, sqlstore.CreateSqlStore(pg),
			rstore, joker.CallJoker(js, rstore))
		if err := baitbot.Serve(ctx); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
