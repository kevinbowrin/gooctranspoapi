package octranspoapi

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

func ID(id string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("id", id)
		return nil
	}
}

func ColumnAndValue(column, value string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("column", column)
		v.Set("value", value)
		return nil
	}
}

func OrderBy(orderBy string) func(url.Values) error {
	return func(v url.Values) error {
		if orderBy != "asc" && orderBy != "desc" {
			return errors.New("OrderBy only accepts asc or desc as parameters.")
		}
		v.Set("orderBy", orderBy)
		return nil
	}
}

func Direction(direction string) func(url.Values) error {
	return func(v url.Values) error {
		v.Set("direction", direction)
		return nil
	}
}

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

func (c *Connection) setupGTFSURL(options ...func(url.Values) error) (*url.URL, error) {
	u, err := url.Parse(ApiURLPrefix + "Gtfs")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.Id)
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

func (c *Connection) performGTFSRequest(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
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

// https://api.octranspo1.com/v1.2/Gtfs agency table
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

func (c *Connection) GetGTFSAgency(ctx context.Context, options ...func(url.Values) error) (*GTFSAgency, error) {
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

func (c *Connection) GetGTFSCalendar(ctx context.Context, options ...func(url.Values) error) (*GTFSCalendar, error) {
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

func (c *Connection) GetGTFSCalendarDates(ctx context.Context, options ...func(url.Values) error) (*GTFSCalendarDates, error) {
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

func (c *Connection) GetGTFSRoutes(ctx context.Context, options ...func(url.Values) error) (*GTFSRoutes, error) {
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

func (c *Connection) GetGTFSStops(ctx context.Context, options ...func(url.Values) error) (*GTFSStops, error) {
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
		return nil, errors.New("GetGTFSStops requires a stop_id, stop_code or id value specified.")
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

func (c *Connection) GetGTFSStopTimes(ctx context.Context, options ...func(url.Values) error) (*GTFSStopTimes, error) {
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
		return nil, errors.New("GetGTFSStopTimes requires a trip_id, stop_id or id value specified.")
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

func (c *Connection) GetGTFSTrips(ctx context.Context, options ...func(url.Values) error) (*GTFSTrips, error) {
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
		return nil, errors.New("GetGTFSTrips requires a route_id or id value specified.")
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
