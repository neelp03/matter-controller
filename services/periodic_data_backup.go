package services

import (
	"context"
	"log"
	"os"
	"time"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"github.com/neelp03/matter-controller/utils"
	"google.golang.org/api/option"
)

func initFirebase() *firestore.Client {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	serviceAccountKeyPath := os.Getenv("FIREBASE_SERVER_CREDENTIALS")

	// Initialize Firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func backupDataToFirestore(client *firestore.Client, data map[string]interface{}) {
	ctx := context.Background()

	//		Log the data to a history for ML training
	_, _, err := client.Collection("periodic_data").Add(ctx, data)
	if err != nil {
		log.Println("Error adding document: ", err)
	}

	//		Update the current data
	// Get subset of data (indoor temp, outdoor temp, window status)
	data_subset := map[string]interface{}{
		"indoor-temp":  data["indoor-temp"],
		"outdoor-temp": data["outdoor-temp"],
		"window-open":  data["window-open"],
	}
	_, err = client.Collection("current_data").Doc("latest").Set(ctx, data_subset)
	if err != nil {
		log.Println("Error updating document: ", err)
	}

	log.Println("Data backed up successfully to Firestore")
}

func compileData() (map[string]interface{}, []error) {
	indoor_temp, err := ReadTemperature()
	errors := make([]error, 0)
	if err != nil {
		log.Println("Error reading indoor temperature: ", err)
		errors = append(errors, err)
	}
	outdoor_temp, err := FetchOutdoorTemperature()
	if err != nil {
		log.Println("Error fetching outdoor weather: ", err)
		errors = append(errors, err)
	}
	window_open, is_window_event := GetWindowStatus()
	if err != nil {
		log.Println("Error checking window status: ", err)
		errors = append(errors, err)
	}

	time := utils.TimeToMinutesSinceYearStart(time.Now())

	data := map[string]interface{}{
		"indoor-temp":  indoor_temp,
		"outdoor-temp": outdoor_temp,
		// "humidity":    60,
		"time": time,
		// "longitude":  0,
		// "latitude":   0,
		// "weather":    "sunny",
		"window-open":  window_open,
		"window-event": is_window_event,
	}

	return data, errors
}

// Launch go routine to backup data every x seconds
func Interval_backup(seconds int) {
	client := initFirebase()
	defer client.Close()

	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data, errors := compileData()
			if len(errors) > 0 {
				log.Println("Errors occurred during data compilation: ", errors)
				log.Println("Skipping backup due to errors. The data may be incomplete.\nDATA:", data)
				continue
			}
			backupDataToFirestore(client, data)
			log.Println("Data backed up successfully to Firestore")
		}
	}
}
