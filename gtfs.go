package gooctranspoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const GTFSMethodPath = "Gtfs"

func ID(id string) func(url.Values) error {
	return func(query url.Values) error {
		query.Set("id", id)
		return nil
	}
}

func Column(column string) func(url.Values) error {
	return func(query url.Values) error {
		query.Set("column", column)
		return nil
	}
}

func Value(value string) func(url.Values) error {
	return func(query url.Values) error {
		query.Set("value", value)
		return nil
	}
}

func OrderBy(orderBy string) func(url.Values) error {
	return func(query url.Values) error {
		if orderBy != "asc" && orderBy != "desc" {
			return errors.New("OrderBy only accepts asc or desc as parameters.")
		}
		query.Set("orderBy", orderBy)
		return nil
	}
}

func Direction(direction string) func(url.Values) error {
	return func(query url.Values) error {
		query.Set("direction", direction)
		return nil
	}
}

func Limit(limit int) func(url.Values) error {
	return func(query url.Values) error {
		query.Set("limit", strconv.Itoa(limit))
		return nil
	}
}

func (c *Connection) GTFSAgency(options ...func(url.Values) error) (*GTFSAgencyData, error) {

	apiURL, err := url.Parse(ApiURLPrefix + GTFSMethodPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(apiURL)
	query := c.setupQuery()
	query.Set("table", "agency")

	for _, opt := range options {
		err := opt(query)
		if err != nil {
			return nil, err
		}
	}

	if query.Get("column") != "" && query.Get("value") == "" {
		return nil, errors.New("If a column is specified, a value must also be specified.")
	}

	apiURL.RawQuery = query.Encode()
	fmt.Println(apiURL)

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non 200 HTTP response from API. %v %v", resp.Status, apiURL)
	}

	data := &GTFSAgencyData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	return data, err
}

// https://api.octranspo1.com/v1.2/Gtfs agency table
type GTFSAgencyData struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID             string `json:"id"`
		AgencyName     string `json:"agency_name"`
		AgencyURL      string `json:"agency_url"`
		AgencyTimezone string `json:"agency_timezone"`
		AgencyLang     string `json:"agency_lang"`
		AgencyPhone    string `json:"agency_phone"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs calendar table
type GTFSCalendar struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID        string `json:"id"`
		ServiceID string `json:"service_id"`
		Monday    string `json:"monday"`
		Tuesday   string `json:"tuesday"`
		Wednesday string `json:"wednesday"`
		Thursday  string `json:"thursday"`
		Friday    string `json:"friday"`
		Saturday  string `json:"saturday"`
		Sunday    string `json:"sunday"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs calendar_dates table
type GTFSCalendarDates struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID            string `json:"id"`
		ServiceID     string `json:"service_id"`
		Date          string `json:"date"`
		ExceptionType string `json:"exception_type"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs routes table
type GTFSRoutes struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID             string `json:"id"`
		RouteID        string `json:"route_id"`
		RouteShortName string `json:"route_short_name"`
		RouteLongName  string `json:"route_long_name"`
		RouteDesc      string `json:"route_desc"`
		RouteType      string `json:"route_type"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs stops table
type GTFSStops struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Column    string `json:"column"`
		Value     string `json:"value"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID            string `json:"id"`
		StopID        string `json:"stop_id"`
		StopCode      string `json:"stop_code"`
		StopName      string `json:"stop_name"`
		StopDesc      string `json:"stop_desc"`
		StopLat       string `json:"stop_lat"`
		StopLon       string `json:"stop_lon"`
		ZoneID        string `json:"zone_id"`
		StopURL       string `json:"stop_url"`
		LocationType  string `json:"location_type"`
		ParentStation string `json:"parent_station"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs stop_times table
type GTFSStopTimes struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Column    string `json:"column"`
		Value     string `json:"value"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID            string `json:"id"`
		TripID        string `json:"trip_id"`
		ArrivalTime   string `json:"arrival_time"`
		DepartureTime string `json:"departure_time"`
		StopID        string `json:"stop_id"`
		StopSequence  string `json:"stop_sequence"`
		PickupType    string `json:"pickup_type"`
		DropOffType   string `json:"drop_off_type"`
	} `json:"Gtfs"`
}

// https://api.octranspo1.com/v1.2/Gtfs trips table
type GTFSTrips struct {
	Query struct {
		Table     string `json:"table"`
		Direction string `json:"direction"`
		Column    string `json:"column"`
		Value     string `json:"value"`
		Format    string `json:"format"`
	} `json:"Query"`
	Gtfs []struct {
		ID           string `json:"id"`
		RouteID      string `json:"route_id"`
		ServiceID    string `json:"service_id"`
		TripID       string `json:"trip_id"`
		TripHeadsign string `json:"trip_headsign"`
		DirectionID  string `json:"direction_id"`
		BlockID      string `json:"block_id"`
	} `json:"Gtfs"`
}
