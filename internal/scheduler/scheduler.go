package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"letsgo/internal/app"

	"github.com/robfig/cron/v3"
)

func Start(application *app.Application) {
	// Create a new cron scheduler
	c := cron.New()

	// Schedule a task to run every minute
	c.AddFunc("@every 1m", func() {
		SampleTask(application)
	})

	// Start the cron scheduler
	c.Start()
}

func SampleTask(application *app.Application) {
	// Get users count from the database
	count, err := application.Queries.CountUsers(context.Background())
	if err != nil {
		log.Printf("Error getting users count: %v", err)
		return
	}
	// Log the users count with a timestamp
	timezone, _ := time.LoadLocation(application.Config.AppTimezone)
	timeStamp := time.Now().In(timezone).Format("2006-01-02 15:04:05")
	println(fmt.Sprintf("[%s] Scheduler running...", timeStamp))
	println(fmt.Sprintf("[%s] Users count: %d", timeStamp, count))
}
