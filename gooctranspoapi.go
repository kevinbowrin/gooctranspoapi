// Package gooctranspoapi provides a Go wrapper around the OC Transpo
// Live Next Bus Arrival Data Feed API.
package gooctranspoapi

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIURLPrefix is the address at which the API is available.
const APIURLPrefix = "https://api.octranspo1.com/v1.2/"

// Connection holds the Application ID and API key needed to make requests.
// It also has a rate limiter, used by the Connection's methods to
// limit calls on the API. The HTTP Client is a public field, so that it
// can be swapped out with a custom HTTP Client if needed.
type Connection struct {
	ID            string
	Key           string
	Limiter       *rate.Limiter
	HTTPClient    *http.Client
	cAPIURLPrefix string
}

// NewConnection returns a new connection without a rate limit.
func NewConnection(id, key string) Connection {
	return Connection{
		ID:            id,
		Key:           key,
		Limiter:       rate.NewLimiter(rate.Inf, 0),
		HTTPClient:    http.DefaultClient,
		cAPIURLPrefix: APIURLPrefix,
	}
}

// NewConnectionWithRateLimit returns a new connection with a rate limit set.
// This is helpful for ensuring you don't go over the daily call limit,
// which is usually 10,000 requests per day.
// It you use the connection over 24 hours, a connection with a perSec rate
// of 0.11572 would make around 9998 requests.
func NewConnectionWithRateLimit(id, key string, perSec float64, burst int) Connection {
	return Connection{
		ID:            id,
		Key:           key,
		Limiter:       rate.NewLimiter(rate.Limit(perSec), burst),
		HTTPClient:    http.DefaultClient,
		cAPIURLPrefix: APIURLPrefix,
	}
}

func (c Connection) performRequest(ctx context.Context, u url.URL, v url.Values) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)
	req.Close = true

	err = c.Limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
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

// RouteSummaryForStop is a simplified version of the data returned by
// a request to GetRouteSummaryForStop.
type RouteSummaryForStop struct {
	StopNo          string
	StopDescription string
	Error           string
	Routes          []Route
}

// Route is used by RouteSummaryForStop to store route data.
type Route struct {
	RouteNo      string
	DirectionID  string
	Direction    string
	RouteHeading string
}

// RawRouteSummaryForStop is a wrapper around the XML data returned by
// a request to GetRouteSummaryForStop.
type rawRouteSummaryForStop struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                           string `xml:",chardata"`
		GetRouteSummaryForStopResponse struct {
			Text                         string `xml:",chardata"`
			Xmlns                        string `xml:"xmlns,attr"`
			GetRouteSummaryForStopResult struct {
				Text   string `xml:",chardata"`
				StopNo struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopNo"`
				StopDescription struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopDescription"`
				Error struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"Error"`
				Routes struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
					Route []struct {
						Text         string `xml:",chardata"`
						RouteNo      string `xml:"RouteNo"`
						DirectionID  string `xml:"DirectionID"`
						Direction    string `xml:"Direction"`
						RouteHeading string `xml:"RouteHeading"`
					} `xml:"Route"`
				} `xml:"Routes"`
			} `xml:"GetRouteSummaryForStopResult"`
		} `xml:"GetRouteSummaryForStopResponse"`
	} `xml:"Body"`
}

// Cook takes a raw XML RouteSummaryForStop and simplifies it.
func (d *rawRouteSummaryForStop) cook() (*RouteSummaryForStop, error) {
	cooked := &RouteSummaryForStop{}
	cooked.StopNo = d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.StopNo.Text
	cooked.StopDescription = d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.StopDescription.Text

	errorText, err := checkErrorCode(d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.Error.Text)
	if err != nil {
		return nil, err
	}
	cooked.Error = errorText

	for _, r := range d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.Routes.Route {
		cr := Route{}
		cr.RouteNo = r.RouteNo
		cr.DirectionID = r.DirectionID
		cr.Direction = r.Direction
		cr.RouteHeading = r.RouteHeading
		cooked.Routes = append(cooked.Routes, cr)
	}
	return cooked, nil
}

// GetRouteSummaryForStop returns the routes for a given stop number.
func (c Connection) GetRouteSummaryForStop(ctx context.Context, stopNo string) (*RouteSummaryForStop, error) {
	u, err := url.Parse(c.cAPIURLPrefix + "GetRouteSummaryForStop")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.ID)
	v.Set("apiKey", c.Key)
	v.Set("stopNo", stopNo)

	respBody, err := c.performRequest(ctx, *u, v)
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(respBody)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false
	data := &rawRouteSummaryForStop{}
	err = dec.Decode(data)
	respBody.Close()
	if err != nil {
		return nil, err
	}

	return data.cook()
}

// NextTripsForStop is a simplified version of the data returned by
// a request to GetNextTripsForStop
type NextTripsForStop struct {
	StopNo          string
	StopLabel       string
	Error           string
	RouteDirections []RouteDirection
}

// RouteDirection is used by NextTripsForStop to store route direction data.
type RouteDirection struct {
	RouteNo               string
	RouteLabel            string
	Direction             string
	Error                 string
	RequestProcessingTime time.Time
	Trips                 []Trip
}

// Trip stores trip data, and includes the adjusted schedule time.
type Trip struct {
	TripDestination      string
	TripStartTime        string
	AdjustedScheduleTime int
	AdjustmentAge        float64
	LastTripOfSchedule
	BusType string
	Latitude
	Longitude
	GPSSpeed
}

// LastTripOfSchedule stores both the data and if the data was set by the API
type LastTripOfSchedule struct {
	Set   bool
	Value bool
}

// Latitude stores both the data and if the data was set by the API
type Latitude struct {
	Set   bool
	Value float64
}

// Longitude stores both the data and if the data was set by the API
type Longitude struct {
	Set   bool
	Value float64
}

// GPSSpeed stores both the data and if the data was set by the API
type GPSSpeed struct {
	Set   bool
	Value float64
}

// rawNextTripsForStop is a wrapper around the XML data returned by
// a request to GetNextTripsForStop.
type rawNextTripsForStop struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                        string `xml:",chardata"`
		GetNextTripsForStopResponse struct {
			Text                      string `xml:",chardata"`
			Xmlns                     string `xml:"xmlns,attr"`
			GetNextTripsForStopResult struct {
				Text   string `xml:",chardata"`
				StopNo struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopNo"`
				StopLabel struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopLabel"`
				Error struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"Error"`
				Route struct {
					Text           string `xml:",chardata"`
					Xmlns          string `xml:"xmlns,attr"`
					RouteDirection []struct {
						Text                  string `xml:",chardata"`
						RouteNo               string `xml:"RouteNo"`
						RouteLabel            string `xml:"RouteLabel"`
						Direction             string `xml:"Direction"`
						Error                 string `xml:"Error"`
						RequestProcessingTime string `xml:"RequestProcessingTime"`
						Trips                 struct {
							Text string       `xml:",chardata"`
							Trip []rawXMLTrip `xml:"Trip"`
						} `xml:"Trips"`
					} `xml:"RouteDirection"`
				} `xml:"Route"`
			} `xml:"GetNextTripsForStopResult"`
		} `xml:"GetNextTripsForStopResponse"`
	} `xml:"Body"`
}

type rawXMLTrip struct {
	Text                 string `xml:",chardata"`
	TripDestination      string `xml:"TripDestination"`
	TripStartTime        string `xml:"TripStartTime"`
	AdjustedScheduleTime string `xml:"AdjustedScheduleTime"`
	AdjustmentAge        string `xml:"AdjustmentAge"`
	LastTripOfSchedule   string `xml:"LastTripOfSchedule"`
	BusType              string `xml:"BusType"`
	Latitude             string `xml:"Latitude"`
	Longitude            string `xml:"Longitude"`
	GPSSpeed             string `xml:"GPSSpeed"`
}

// Cook takes a raw XML NextTripsForStop and simplifies it.
func (d *rawNextTripsForStop) cook() (*NextTripsForStop, error) {
	cooked := &NextTripsForStop{}

	cooked.StopNo = d.Body.GetNextTripsForStopResponse.GetNextTripsForStopResult.StopNo.Text
	cooked.StopLabel = d.Body.GetNextTripsForStopResponse.GetNextTripsForStopResult.StopLabel.Text

	errorText, err := checkErrorCode(d.Body.GetNextTripsForStopResponse.GetNextTripsForStopResult.Error.Text)
	if err != nil {
		return nil, err
	}
	cooked.Error = errorText

	for _, rd := range d.Body.GetNextTripsForStopResponse.GetNextTripsForStopResult.Route.RouteDirection {
		crd := RouteDirection{}
		crd.RouteNo = rd.RouteNo
		crd.RouteLabel = rd.RouteLabel
		crd.Direction = rd.Direction

		errorText, err := checkErrorCode(rd.Error)
		if err != nil {
			return nil, err
		}
		crd.Error = errorText

		tz, err := time.LoadLocation("America/Toronto")
		if err != nil {
			return nil, err
		}

		parsedProcessingTime, err := time.ParseInLocation("20060102150405", rd.RequestProcessingTime, tz)
		if err != nil {
			return nil, err
		}

		crd.RequestProcessingTime = parsedProcessingTime

		for _, t := range rd.Trips.Trip {
			ct, err := t.convert()
			if err != nil {
				return nil, err
			}
			crd.Trips = append(crd.Trips, ct)
		}
		cooked.RouteDirections = append(cooked.RouteDirections, crd)
	}
	return cooked, nil
}

// GetNextTripsForStop returns the next three trips on the route for a given stop number.
func (c Connection) GetNextTripsForStop(ctx context.Context, routeNo, stopNo string) (*NextTripsForStop, error) {
	u, err := url.Parse(c.cAPIURLPrefix + "GetNextTripsForStop")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.ID)
	v.Set("apiKey", c.Key)
	v.Set("routeNo", routeNo)
	v.Set("stopNo", stopNo)

	respBody, err := c.performRequest(ctx, *u, v)
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(respBody)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false
	data := &rawNextTripsForStop{}
	err = dec.Decode(data)
	respBody.Close()
	if err != nil {
		return nil, err
	}

	return data.cook()
}

// NextTripsForStopAllRoutes is a simplified version of the data returned by
// a request to GetNextTripsForStopAllRoutes
type NextTripsForStopAllRoutes struct {
	StopNo          string
	StopDescription string
	Error           string
	Routes          []RouteWithTrips
}

// RouteWithTrips is used by NextTripsForStopAllRoutes to store route data.
type RouteWithTrips struct {
	RouteNo      string
	DirectionID  string
	Direction    string
	RouteHeading string
	Trips        []Trip
}

// NextTripsForStopAllRoutes is a wrapper around the XML data returned by
// a request to GetNextTripsForStopAllRoutes.
type rawNextTripsForStopAllRoutes struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                           string `xml:",chardata"`
		GetRouteSummaryForStopResponse struct {
			Text                         string `xml:",chardata"`
			Xmlns                        string `xml:"xmlns,attr"`
			GetRouteSummaryForStopResult struct {
				Text   string `xml:",chardata"`
				StopNo struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopNo"`
				StopDescription struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"StopDescription"`
				Error struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
				} `xml:"Error"`
				Routes struct {
					Text  string `xml:",chardata"`
					Xmlns string `xml:"xmlns,attr"`
					Route []struct {
						Text         string `xml:",chardata"`
						RouteNo      string `xml:"RouteNo"`
						DirectionID  string `xml:"DirectionID"`
						Direction    string `xml:"Direction"`
						RouteHeading string `xml:"RouteHeading"`
						Trips        struct {
							Text string       `xml:",chardata"`
							Trip []rawXMLTrip `xml:"Trip"`
						} `xml:"Trips"`
					} `xml:"Route"`
				} `xml:"Routes"`
			} `xml:"GetRouteSummaryForStopResult"`
		} `xml:"GetRouteSummaryForStopResponse"`
	} `xml:"Body"`
}

// Cook takes a raw XML NextTripsForStopAllRoutes and simplifies it.
func (d *rawNextTripsForStopAllRoutes) cook() (*NextTripsForStopAllRoutes, error) {
	cooked := &NextTripsForStopAllRoutes{}

	cooked.StopNo = d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.StopNo.Text
	cooked.StopDescription = d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.StopDescription.Text

	errorText, err := checkErrorCode(d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.Error.Text)
	if err != nil {
		return nil, err
	}
	cooked.Error = errorText

	for _, rt := range d.Body.GetRouteSummaryForStopResponse.GetRouteSummaryForStopResult.Routes.Route {
		crt := RouteWithTrips{}
		crt.RouteNo = rt.RouteNo
		crt.DirectionID = rt.DirectionID
		crt.Direction = rt.Direction
		crt.RouteHeading = rt.RouteHeading

		for _, t := range rt.Trips.Trip {
			ct, err := t.convert()
			if err != nil {
				return nil, err
			}
			crt.Trips = append(crt.Trips, ct)
		}
		cooked.Routes = append(cooked.Routes, crt)
	}
	return cooked, nil
}

// GetNextTripsForStopAllRoutes returns the next three trips for all routes for a given stop number.
func (c Connection) GetNextTripsForStopAllRoutes(ctx context.Context, stopNo string) (*NextTripsForStopAllRoutes, error) {
	u, err := url.Parse(c.cAPIURLPrefix + "GetNextTripsForStopAllRoutes")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.ID)
	v.Set("apiKey", c.Key)
	v.Set("stopNo", stopNo)

	respBody, err := c.performRequest(ctx, *u, v)
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(respBody)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false
	data := &rawNextTripsForStopAllRoutes{}
	err = dec.Decode(data)
	respBody.Close()
	if err != nil {
		return nil, err
	}

	return data.cook()
}

func checkErrorCode(errorText string) (string, error) {
	switch errorText {
	case "1":
		return "", errors.New("error returned from API - Invalid API key")
	case "2":
		return "", errors.New("error returned from API - Unable to query data source")
	case "10":
		return "", errors.New("error returned from API - Invalid stop number")
	case "11":
		return "", errors.New("error returned from API - Invalid route number")
	case "12":
		return "", errors.New("error returned from API - Stop does not service route")
	default:
		return errorText, nil
	}
}

func (t rawXMLTrip) convert() (Trip, error) {
	ct := Trip{}
	ct.TripDestination = t.TripDestination
	ct.TripStartTime = t.TripStartTime

	pAdjustedScheduleTime, err := strconv.Atoi(t.AdjustedScheduleTime)
	if err != nil {
		return ct, err
	}
	ct.AdjustedScheduleTime = pAdjustedScheduleTime

	pAdjustmentAge, err := strconv.ParseFloat(t.AdjustmentAge, 64)
	if err != nil {
		return ct, err
	}
	ct.AdjustmentAge = pAdjustmentAge

	if t.LastTripOfSchedule == "" {
		ct.LastTripOfSchedule = LastTripOfSchedule{Set: false}
	} else {
		pLastTripOfSchedule, err := strconv.ParseBool(t.LastTripOfSchedule)
		if err != nil {
			return ct, err
		}
		ct.LastTripOfSchedule = LastTripOfSchedule{Set: true, Value: pLastTripOfSchedule}
	}

	ct.BusType = t.BusType

	if t.Latitude == "" {
		ct.Latitude = Latitude{Set: false}
	} else {
		pLatitude, err := strconv.ParseFloat(t.Latitude, 64)
		if err != nil {
			return ct, err
		}
		ct.Latitude = Latitude{Set: true, Value: pLatitude}
	}

	if t.Longitude == "" {
		ct.Longitude = Longitude{Set: false}
	} else {
		pLongitude, err := strconv.ParseFloat(t.Longitude, 64)
		if err != nil {
			return ct, err
		}
		ct.Longitude = Longitude{Set: true, Value: pLongitude}
	}

	if t.GPSSpeed == "" {
		ct.GPSSpeed = GPSSpeed{Set: false}
	} else {
		pGPSSpeed, err := strconv.ParseFloat(t.GPSSpeed, 64)
		if err != nil {
			return ct, err
		}
		ct.GPSSpeed = GPSSpeed{Set: true, Value: pGPSSpeed}
	}

	return ct, nil
}
