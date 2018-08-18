// Package gooctranspoapi provides a Go wrapper around the OC Transpo Live Next Bus Arrival Data Feed API.
package gooctranspoapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// APIURLPrefix is the address at which the API is available.
const APIURLPrefix = "https://api.octranspo1.com/v1.2/"

// Connection holds the Application ID and API key needed to make requests.
// It also has a rate limiter, used by the Connection's methods to limit calls on the API.
type Connection struct {
	ID      string
	Key     string
	Limiter *rate.Limiter
}

// NewConnection returns a new connection without a rate limit.
func NewConnection(id, key string) Connection {
	return Connection{
		ID:      id,
		Key:     key,
		Limiter: rate.NewLimiter(rate.Inf, 0),
	}
}

// NewConnectionWithRateLimit returns a new connection with a rate limit set.
// This is helpful for ensuring you don't go over the daily call limit, which is usually 10,000 requests per day.
// It you use the connection over 24 hours, a connection with a perSecond rate of 0.11572 would make around 9998 requests.
func NewConnectionWithRateLimit(id, key string, perSecond float64, burst int) Connection {
	return Connection{
		ID:      id,
		Key:     key,
		Limiter: rate.NewLimiter(rate.Limit(perSecond), burst),
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

// RouteSummaryForStop is a rough wrapper around the XML data returned by
// a request to GetRouteSummaryForStop. #TODO: Create a SimpleRouteSummaryForStop
type RouteSummaryForStop struct {
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

// GetRouteSummaryForStop returns the routes for a given stop number.
func (c Connection) GetRouteSummaryForStop(ctx context.Context, stopNo string) (*RouteSummaryForStop, error) {
	u, err := url.Parse(APIURLPrefix + "GetRouteSummaryForStop")
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
	data := &RouteSummaryForStop{}
	err = dec.Decode(data)
	respBody.Close()
	return data, err
}

// NextTripsForStop is a rough wrapper around the XML data returned by
// a request to GetNextTripsForStop. #TODO: Create a SimpleGetNextTripsForStop
type NextTripsForStop struct {
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
					RouteDirection struct {
						Text                  string `xml:",chardata"`
						RouteNo               string `xml:"RouteNo"`
						RouteLabel            string `xml:"RouteLabel"`
						Direction             string `xml:"Direction"`
						Error                 string `xml:"Error"`
						RequestProcessingTime string `xml:"RequestProcessingTime"`
						Trips                 struct {
							Text string `xml:",chardata"`
							Trip []struct {
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
							} `xml:"Trip"`
						} `xml:"Trips"`
					} `xml:"RouteDirection"`
				} `xml:"Route"`
			} `xml:"GetNextTripsForStopResult"`
		} `xml:"GetNextTripsForStopResponse"`
	} `xml:"Body"`
}

// GetNextTripsForStop returns the next three trips on the route for a given stop number.
func (c Connection) GetNextTripsForStop(ctx context.Context, routeNo, stopNo string) (*NextTripsForStop, error) {
	u, err := url.Parse(APIURLPrefix + "GetNextTripsForStop")
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
	data := &NextTripsForStop{}
	err = dec.Decode(data)
	respBody.Close()
	return data, err
}

// NextTripsForStopAllRoutes is a rough wrapper around the XML data returned by
// a request to GetNextTripsForStopAllRoutes. #TODO: Create a SimpleGetNextTripsForStopAllRoutes
type NextTripsForStopAllRoutes struct {
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
							Text string `xml:",chardata"`
							Trip []struct {
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
							} `xml:"Trip"`
						} `xml:"Trips"`
					} `xml:"Route"`
				} `xml:"Routes"`
			} `xml:"GetRouteSummaryForStopResult"`
		} `xml:"GetRouteSummaryForStopResponse"`
	} `xml:"Body"`
}

// GetNextTripsForStopAllRoutes returns the next three trips for all routes for a given stop number.
func (c Connection) GetNextTripsForStopAllRoutes(ctx context.Context, stopNo string) (*NextTripsForStopAllRoutes, error) {
	u, err := url.Parse(APIURLPrefix + "GetNextTripsForStopAllRoutes")
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
	data := &NextTripsForStopAllRoutes{}
	err = dec.Decode(data)
	respBody.Close()
	return data, err
}
