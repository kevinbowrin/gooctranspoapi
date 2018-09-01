package gooctranspoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGTFSAgency(t *testing.T) {
	rawJSONString := `{"Query":{"table":"agency",
                                "direction":"ASC","format":"json"},
                     "Gtfs":[{"id":"1","agency_name":"Test Agency",
                              "agency_url":"http://test.com",
                              "agency_timezone":"America/Toronto",
                              "agency_lang":"","agency_phone":""}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	agency, err := c.GetGTFSAgency(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if agency.Gtfs[0].ID != "1" {
		t.Fatal("Unexpected ID in returned GTFSAgency")
	}
	if agency.Gtfs[0].AgencyName != "Test Agency" {
		t.Fatal("Unexpected AgencyName in returned GTFSAgency")
	}
	if agency.Gtfs[0].AgencyTimezone != "America/Toronto" {
		t.Fatal("Unexpected AgencyTimezone in returned GTFSAgency")
	}
}

func TestGTFSCalendar(t *testing.T) {
	rawJSONString := `{"Query": {"table":"calendar","direction":"ASC","column": 
                                 "id","value":"1","format":"json"},
                       "Gtfs": [{"id":"1",
                                 "service_id":"JUN26-JUNDA13-Weekday-01",
                                 "monday":"1","tuesday":"1","wednesday":"1",
                                 "thursday":"1","friday":"1","saturday":"0",
                                 "sunday":"0","start_date": "20130626",
                                 "end_date": "20130627"}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	calendar, err := c.GetGTFSCalendar(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if calendar.Gtfs[0].ID != "1" {
		t.Fatal("Unexpected ID in returned GTFSCalendar")
	}
	if calendar.Gtfs[0].ServiceID != "JUN26-JUNDA13-Weekday-01" {
		t.Fatal("Unexpected ServiceID in returned GTFSCalendar")
	}
	if calendar.Gtfs[0].EndDate != "20130627" {
		t.Fatal("Unexpected EndDate in returned GTFSCalendar")
	}
}

func TestGTFSCalendarDates(t *testing.T) {
	rawJSONString := `{"Query":{"table":"calendar_dates",
                                "direction":"ASC","column":"id","value":"1",
                                "format":"json"},
                        "Gtfs":[{"id":"1",
                                 "service_id":"JUN13-JUNDA13-Weekday-99",
                                 "date":"20130701","exception_type":"2"}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	dates, err := c.GetGTFSCalendarDates(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if dates.Gtfs[0].ID != "1" {
		t.Fatal("Unexpected ID in returned GTFSCalendarDates")
	}
	if dates.Gtfs[0].ServiceID != "JUN13-JUNDA13-Weekday-99" {
		t.Fatal("Unexpected ServiceID in returned GTFSCalendarDates")
	}
	if dates.Gtfs[0].ExceptionType != "2" {
		t.Fatal("Unexpected ExceptionType in returned GTFSCalendarDates")
	}
}

func TestGTFSRoutes(t *testing.T) {
	rawJSONString := `{"Query":{"table":"routes","direction":"ASC",
	                            "column":"id","value":"1","format":"json"},
	                   "Gtfs":[{"id":"1","route_id":"1-146",
	                            "route_short_name":"1","route_long_name":"",
	                            "route_desc":"","route_type":"3"}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	routes, err := c.GetGTFSRoutes(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if routes.Gtfs[0].RouteID != "1-146" {
		t.Fatal("Unexpected RouteID in returned GTFSRoutes")
	}
	if routes.Gtfs[0].RouteType != "3" {
		t.Fatal("Unexpected RouteType in returned GTFSRoutes")
	}
}

func TestGTFSStops(t *testing.T) {
	rawJSONString := `{"Query":{"table":"stops","direction":"ASC",
	                            "column":"stop_id","value":"AA010",
	                            "format":"json"},
	                   "Gtfs":[{"id":"1","stop_id":"AA010","stop_code":"8767",
	                            "stop_name":"SUSSEX / CHUTE RIDEAU FALLS",
	                            "stop_desc":"","stop_lat":"45.4399",
	                            "stop_lon":"-75.6958","stop_street":"",
	                            "stop_city":"","stop_region":"",
	                            "stop_postcode":"","stop_country":"","zone_id":""}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	stops, err := c.GetGTFSStops(context.TODO(), ID("1"))
	if err != nil {
		t.Fatal(err)
	}

	if stops.Gtfs[0].StopID != "AA010" {
		t.Fatal("Unexpected StopID in returned GTFSStops")
	}
	if stops.Gtfs[0].StopName != "SUSSEX / CHUTE RIDEAU FALLS" {
		t.Fatal("Unexpected StopName in returned GTFSStops")
	}
}

func TestGTFSStopTimes(t *testing.T) {
	rawJSONString := `{"Query":{"table":"stop_times","direction":"ASC",
	                            "column":"stop_id","value":"AA010",
	                            "format":"json"},
	                   "Gtfs":[{"id":"133436",
	                            "trip_id":"27212870-CADA13-CADA13-Sunday-71",
	                            "arrival_time":"08:29:00",
	                            "departure_time":"08:29:00","stop_id":"AA010",
	                            "stop_sequence":"20","pickup_type":"0",
	                            "drop_off_type":"0"}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	times, err := c.GetGTFSStopTimes(context.TODO(), ID("133436"))
	if err != nil {
		t.Fatal(err)
	}

	if times.Gtfs[0].TripID != "27212870-CADA13-CADA13-Sunday-71" {
		t.Fatal("Unexpected TripID in returned GTFSStopTimes")
	}
	if times.Gtfs[0].StopSequence != "20" {
		t.Fatal("Unexpected StopSequence in returned GTFSStopTimes")
	}
}

func TestGTFSTrips(t *testing.T) {
	rawJSONString := `{"Query":{"table":"trips","direction":"ASC",
	                            "column":"route_id","value":"135-147",
	                            "format":"json"},
	                   "Gtfs":[{"id":"1","route_id":"135-147",
	                            "service_id":"CADA13-CADA13-Sunday-71",
	                            "trip_id":"27210104-CADA13-CADA13-Sunday-71",
	                            "trip_headsign":"Esprit",
	                            "block_id":"3406628"}]}`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawJSONString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	trips, err := c.GetGTFSTrips(context.TODO(), ID("1"))
	if err != nil {
		t.Fatal(err)
	}

	if trips.Gtfs[0].TripHeadsign != "Esprit" {
		t.Fatal("Unexpected TripHeadsign in returned GTFSTrips")
	}
	if trips.Gtfs[0].BlockID != "3406628" {
		t.Fatal("Unexpected BlockID in returned GTFSTrips")
	}
}
