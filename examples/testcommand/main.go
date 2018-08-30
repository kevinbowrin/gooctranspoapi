package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	api "github.com/transitreport/gooctranspoapi"
	"log"
	"os"
	"os/signal"
	"time"
)

const splitterString = "\n\n------------------------------------------------\n\n"

var (
	id  = flag.String("id", "", "appID")
	key = flag.String("key", "", "apiKey")
)

func main() {

	// Process the flags.
	flag.Parse()

	// If any of the required flags are not set, exit.
	if *id == "" {
		log.Fatalln("FATAL: An appID for the OC Transpo API is required.")
	} else if *key == "" {
		log.Fatalln("FATAL: An apiKey for the OC Transpo API is required.")
	}

	c := api.NewConnectionWithRateLimit(*id, *key, 1, 1)
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

	agency, err := c.GetGTFSAgency(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(agency)

	fmt.Print(splitterString)

	cal, err := c.GetGTFSCalendar(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(cal)

	fmt.Print(splitterString)

	caldates, err := c.GetGTFSCalendarDates(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(caldates)

	fmt.Print(splitterString)

	routes, err := c.GetGTFSRoutes(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(routes)

	fmt.Print(splitterString)

	stops, err := c.GetGTFSStops(ctx, api.ColumnAndValue("stop_id", "7659"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(stops)

	fmt.Print(splitterString)

	times, err := c.GetGTFSStopTimes(ctx, api.ColumnAndValue("stop_id", "7659"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(times)

	fmt.Print(splitterString)

	trips, err := c.GetGTFSTrips(ctx, api.ColumnAndValue("route_id", "1"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(trips)

	fmt.Print(splitterString)

	routeSummary, err := c.GetRouteSummaryForStop(ctx, "7659")
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(routeSummary)

	fmt.Print(splitterString)

	nextTrips, err := c.GetNextTripsForStop(ctx, "6", "7659")
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(nextTrips)

	fmt.Print(splitterString)

	nextTripsAllRoutes, err := c.GetNextTripsForStopAllRoutes(ctx, "7659")
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(nextTripsAllRoutes)

	fmt.Print(splitterString)
}
