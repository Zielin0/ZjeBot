package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zielin0/ZjeBot/core"
)

func main() {
	start := time.Now()

	log.Println("Initializing ZjeBot...")

	secretsLoader, err := core.NewSecretsLoader(core.SECRETS_PATH)
	if err != nil {
		log.Fatalf("Failed to load secrets: %s", err)
	}

	log.Println("Secrets loaded successfully")
	secrets := secretsLoader.GetSecrets()

	dataLoader, err := core.NewDataLoader(core.DATA_PATH)
	if err != nil {
		log.Fatalf("Failed to load data: %s", err)
	}

	log.Println("Data loaded successfully")
	data := dataLoader.GetData()

	bot := core.CreateBot(secrets, data, dataLoader, &start)
	log.Println("Bot created successfully")

	discordEnv, err := core.InitDiscord(bot)
	if err != nil {
		log.Fatalf("Failed to initialize discord environment: %s", err)
	}

	log.Println("Discord initialized successfully")
	defer discordEnv.Dg.Close()

	twitchEnv, err := core.InitTwitch(bot)
	if err != nil {
		log.Fatalf("Failed to initialize twitch environment: %s", err)
	}

	log.Println("Twitch initialized successfully")
	defer twitchEnv.Tc.Disconnect()

	log.Println("ZjeBot is running. Press Ctrl+C to stop.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case <-sc:
			log.Println("Termination signal received. Stopping ZjeBot.")
			return
		case err := <-ConnectTwitch(twitchEnv):
			if err != nil {
				log.Printf("Connection to Twitch failed: %s", err)
			}
		}
	}
}

func ConnectTwitch(te *core.TwitchEnvironment) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		err := te.Tc.Connect()
		errCh <- err
	}()

	return errCh
}
