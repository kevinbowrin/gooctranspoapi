package main

import (
	"context"
	"flag"
	"fmt"
	api "github.com/transitreport/gooctranspoapi"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	id   = flag.String("id", "", "appID")
	key  = flag.String("key", "", "apiKey")
	stop = flag.String("stop", "", "stop number")
)

func main() {

	// Process the flags.
	flag.Parse()

	// If any of the required flags are not set, exit.
	if *id == "" {
		log.Fatalln("FATAL: An appID for the OC Transpo API is required.")
	} else if *key == "" {
		log.Fatalln("FATAL: An apiKey for the OC Transpo API is required.")
	} else if *stop == "" {
		log.Fatalln("FATAL: An stop number is required.")
	}

	// Create a new connection to the API, with a rate limit of 1 request per second,
	// with bursts of size 1.
	// Connections can also be created without a rate limit by using NewConnection()
	c := api.NewConnectionWithRateLimit(*id, *key, 1, 1)

	// Requests to the API have a context which can be canceled or timed out.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// trap Ctrl+C and call cancel on the context
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)
	go func() {
		select {
		case <-sigChan:
			log.Println("Canceling requests...")
			cancel()
			log.Println("Done, bye!")
		case <-ctx.Done():
		}
	}()

	nextTripsAllRoutes, err := c.GetNextTripsForStopAllRoutes(ctx, *stop)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Stop %v, \"%v\":\n", nextTripsAllRoutes.StopNo, nextTripsAllRoutes.StopDescription)
	for _, route := range nextTripsAllRoutes.Routes {
		fmt.Printf("  Route %v, \"%v\", going %v:\n", route.RouteNo, route.RouteHeading, route.Direction)
		for _, trip := range route.Trips {
			fmt.Printf("    %v (%v minutes old), %v\n", trip.AdjustedScheduleTime, trip.AdjustmentAge, trip.TripDestination)
		}
	}
}
