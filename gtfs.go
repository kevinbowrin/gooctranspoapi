package gooctranspoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// ID will setup the request to return a specific row in a table by the id value.
func ID(id string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("id", id)
		return nil
	}
}

// ColumnAndValue will setup the request to return data from a specific column and value.
func ColumnAndValue(column, value string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("column", column)
		v.Set("value", value)
		return nil
	}
}

// OrderBy will setup the request to sort the data by a specific column.
func OrderBy(orderBy string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("orderBy", orderBy)
		return nil
	}
}

// Direction will setup the request to direction of sorted records, asc and desc.
func Direction(direction string) func(url.Values) error {
	return func(v url.Values) error {
		if direction != "asc" && direction != "desc" {
			return errors.New("direction only accepts asc or desc as parameters")
		}
		v.Set("direction", direction)
		return nil
	}
}

// Limit will setup the request to only return a maximum number of records.
func Limit(limit int) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("limit", strconv.Itoa(limit))
		return nil
	}
}

func setTable(table string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("table", table)
		return nil
	}
}

func (c Connection) setupGTFSURL(options ...func(url.Values) error) (*url.URL, error) {
	u, err := url.Parse(APIURLPrefix + "Gtfs")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.ID)
	v.Set("apiKey", c.Key)
	v.Set("format", "json")
	for _, opt := range options {
		err := opt(v)
		if err != nil {
			return nil, err
		}
	}
	u.RawQuery = v.Encode()
	return u, nil
}

func (c Connection) performGTFSRequest(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Close = true

	err = c.Limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("Non 200 HTTP response from API. %v %v", resp.Status, u.String())
	}

	return resp.Body, nil
}

// GTFSAgency is the GTFS agency table.
type GTFSAgency struct {
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

// GetGTFSAgency returns the GTFS agency table.
func (c Connection) GetGTFSAgency(ctx context.Context, options ...func(url.Values) error) (*GTFSAgency, error) {
	options = append(options, setTable("agency"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSAgency{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSCalendar is the GTFS calendar table.
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

// GetGTFSCalendar returns the GTFS calendar table.
func (c Connection) GetGTFSCalendar(ctx context.Context, options ...func(url.Values) error) (*GTFSCalendar, error) {
	options = append(options, setTable("calendar"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSCalendar{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSCalendarDates is GTFS calendar_dates table.
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

// GetGTFSCalendarDates returns the GTFS calendar_dates table
func (c Connection) GetGTFSCalendarDates(ctx context.Context, options ...func(url.Values) error) (*GTFSCalendarDates, error) {
	options = append(options, setTable("calendar_dates"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSCalendarDates{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSRoutes is the GTFS routes table.
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

// GetGTFSRoutes returns the GTFS routes table.
func (c Connection) GetGTFSRoutes(ctx context.Context, options ...func(url.Values) error) (*GTFSRoutes, error) {
	options = append(options, setTable("routes"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSRoutes{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSStops is the GTFS stops table.
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

// GetGTFSStops returns the GTFS stops table.
// It requires a stop_id, stop_code or id value specified, using ColumnAndValue() or ID() options.
func (c Connection) GetGTFSStops(ctx context.Context, options ...func(url.Values) error) (*GTFSStops, error) {
	options = append(options, setTable("stops"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	if v.Get("column") != "stop_id" && v.Get("column") != "stop_code" && v.Get("id") == "" {
		return nil, errors.New("a stop_id, stop_code or id value must be specified")
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSStops{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSStopTimes is the GTFS stop_times table.
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

// GetGTFSStopTimes returns the GTFS stop_times table.
// It requires a trip_id, stop_code or id value specified, using ColumnAndValue() or ID() options.
func (c Connection) GetGTFSStopTimes(ctx context.Context, options ...func(url.Values) error) (*GTFSStopTimes, error) {
	options = append(options, setTable("stop_times"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	if v.Get("column") != "trip_id" && v.Get("column") != "stop_id" && v.Get("id") == "" {
		return nil, errors.New("a trip_id, stop_id or id value must be specified")
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSStopTimes{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}

// GTFSTrips is the GTFS trips table.
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

// GetGTFSTrips returns the GTFS trips table.
// It requires a route_id or id value specified, using ColumnAndValue() or ID() options.
func (c Connection) GetGTFSTrips(ctx context.Context, options ...func(url.Values) error) (*GTFSTrips, error) {
	options = append(options, setTable("trips"))
	u, err := c.setupGTFSURL(options...)
	if err != nil {
		return nil, err
	}
	v, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	if v.Get("column") != "route_id" && v.Get("id") == "" {
		return nil, errors.New("a route_id or id value must be specified")
	}
	respBody, err := c.performGTFSRequest(ctx, u)
	if err != nil {
		return nil, err
	}
	data := &GTFSTrips{}
	err = json.NewDecoder(respBody).Decode(data)
	respBody.Close()
	return data, err
}
