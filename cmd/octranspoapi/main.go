package main

import (
	"context"
	"flag"
	"github.com/davecgh/go-spew/spew"
	api "github.com/transitreport/gooctranspoapi"
	"log"
	"os"
	"os/signal"
	"time"
)

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

	c, err := api.NewConnection(*id, *key, api.RateLimit(1))
	if err != nil {
		log.Fatalln(err)
	}
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

	cal, err := c.GetGTFSCalendar(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(cal)

	caldates, err := c.GetGTFSCalendarDates(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(caldates)

	routes, err := c.GetGTFSRoutes(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(routes)

	stops, err := c.GetGTFSStops(ctx, api.ColumnAndValue("stop_id", "7659"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(stops)

	times, err := c.GetGTFSStopTimes(ctx, api.ColumnAndValue("stop_id", "7659"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(times)

	trips, err := c.GetGTFSTrips(ctx, api.ColumnAndValue("route_id", "1"))
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(trips)

	x, err := c.GetRouteSummaryForStop(ctx, "7659")
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(x)
}
