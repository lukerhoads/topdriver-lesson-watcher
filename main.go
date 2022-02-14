package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const retryDelayMinutes = 3

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	cookie := os.Getenv("COOKIE")
	if cookie == "" {
		log.Fatalln("Cookie is required")
	}

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	if accountSid == "" {
		log.Fatalln("Twilio account SID is required")
	}

	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if authToken == "" {
		log.Fatalln("Twilio auth token is required")
	}

	configFile := os.Getenv("CONFIG")
	if configFile == "" {
		configFile = "config.yaml"
	}

	cfg, err := GetConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Exiting...")
		os.Exit(0)
	}()

	topDriverClient := NewTopdriverClient(cookie, cfg.InstructorId, cfg.Timings, cfg.PickupLocation, cfg.DropoffLocation)
	notifier := NewNotifier(accountSid, authToken, cfg.PhoneNumber)

	for {
		log.Println("Fetching available classes...")

		rawClasses, err := topDriverClient.GetAvailableDays()
		if err != nil {
			log.Fatal(err)
		}

		// Decode rawClasses
		appointmentAvailable, err := ParseRes(rawClasses)
		if err != nil {
			log.Fatal(err)
		}

		if appointmentAvailable {
			log.Println("Detected an available appointment, sending text...")
			notifier.SendText(cfg.ReceiverPhoneNumber, "Detected an available lesson")
		}

		time.Sleep(time.Duration(retryDelayMinutes) * time.Minute)
	}
}
