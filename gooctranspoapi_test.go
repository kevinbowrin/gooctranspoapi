package gooctranspoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRouteSummaryForStop(t *testing.T) {
	rawXMLString := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetRouteSummaryForStopResponse xmlns="http://octranspo.com">
      <GetRouteSummaryForStopResult>
        <StopNo xmlns="http://tempuri.org/">7659</StopNo>
        <StopDescription xmlns="http://tempuri.org/">BANK / FIFTH</StopDescription>
        <Error xmlns="http://tempuri.org/">TestErrorStringHere</Error>
        <Routes xmlns="http://tempuri.org/">
          <Route>
            <RouteNo>6</RouteNo>
            <DirectionID>1</DirectionID>
            <Direction>Northbound</Direction>
            <RouteHeading>Rockcliffe</RouteHeading>
          </Route>
          <Route>
            <RouteNo>7</RouteNo>
            <DirectionID>1</DirectionID>
            <Direction>Eastbound</Direction>
            <RouteHeading>St-Laurent</RouteHeading>
          </Route>
        </Routes>
      </GetRouteSummaryForStopResult>
    </GetRouteSummaryForStopResponse>
  </soap:Body>
</soap:Envelope>`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawXMLString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	routeSummary, err := c.GetRouteSummaryForStop(context.TODO(), "7659")
	if err != nil {
		t.Fatal(err)
	}

	if routeSummary.StopNo != "7659" {
		t.Fatal("Unexpected StopNo in returned RouteSummaryForStop")
	}
	if routeSummary.StopDescription != "BANK / FIFTH" {
		t.Fatal("Unexpected StopDescription in returned RouteSummaryForStop")
	}
	if routeSummary.Error != "TestErrorStringHere" {
		t.Fatal("Unexpected Error in returned RouteSummaryForStop")
	}

	expectedFirstRoute := Route{
		RouteNo:      "6",
		DirectionID:  "1",
		Direction:    "Northbound",
		RouteHeading: "Rockcliffe",
	}

	if routeSummary.Routes[0] != expectedFirstRoute {
		t.Fatal("Unexpected first route in returned RouteSummaryForStop")
	}

	expectedSecondRoute := Route{
		RouteNo:      "7",
		DirectionID:  "1",
		Direction:    "Eastbound",
		RouteHeading: "St-Laurent",
	}

	if routeSummary.Routes[1] != expectedSecondRoute {
		t.Fatal("Unexpected second route in returned RouteSummaryForStop")
	}
}

func TestRouteSummaryForStopWithError(t *testing.T) {
	rawXMLString := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetRouteSummaryForStopResponse xmlns="http://octranspo.com">
      <GetRouteSummaryForStopResult>
        <Error xmlns="http://tempuri.org/">10</Error>
      </GetRouteSummaryForStopResult>
    </GetRouteSummaryForStopResponse>
  </soap:Body>
</soap:Envelope>`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawXMLString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	_, err := c.GetRouteSummaryForStop(context.TODO(), "000000")
	if err == nil {
		t.Fatal("Expected error from parsing RouteSummaryForStop with Error")
	}

}

func TestGetNextTripsForStop(t *testing.T) {
	rawXMLString := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetNextTripsForStopResponse xmlns="http://octranspo.com">
      <GetNextTripsForStopResult>
        <StopNo xmlns="http://tempuri.org/">3020</StopNo>
        <StopLabel xmlns="http://tempuri.org/">LAURIER STATION</StopLabel>
        <Error xmlns="http://tempuri.org/">TestErrorStringHere</Error>
        <Route xmlns="http://tempuri.org/">
          <RouteDirection>
            <RouteNo>94</RouteNo>
            <RouteLabel>Riverview</RouteLabel>
            <Direction>Westbound</Direction>
            <Error>TestRouteDirectionErrorHere</Error>
            <RequestProcessingTime>20180831114042</RequestProcessingTime>
            <Trips>
              <Trip>
                <TripDestination>Riverview</TripDestination>
                <TripStartTime>11:13</TripStartTime>
                <AdjustedScheduleTime>16</AdjustedScheduleTime>
                <AdjustmentAge>0.34</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>6EB - 60</BusType>
                <Latitude>45.431521</Latitude>
                <Longitude>-75.605296</Longitude>
                <GPSSpeed>63.0</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Riverview</TripDestination>
                <TripStartTime>10:59</TripStartTime>
                <AdjustedScheduleTime>17</AdjustedScheduleTime>
                <AdjustmentAge>0.32</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>4EB - DD</BusType>
                <Latitude>45.426999</Latitude>
                <Longitude>-75.600192</Longitude>
                <GPSSpeed>11.4</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Riverview</TripDestination>
                <TripStartTime>11:28</TripStartTime>
                <AdjustedScheduleTime>35</AdjustedScheduleTime>
                <AdjustmentAge>0.32</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>4EA - DD</BusType>
                <Latitude>45.455889</Latitude>
                <Longitude>-75.504171</Longitude>
                <GPSSpeed>0.5</GPSSpeed>
              </Trip>
            </Trips>
          </RouteDirection>
          <RouteDirection>
            <RouteNo>94</RouteNo>
            <RouteLabel>Millennium</RouteLabel>
            <Direction>Eastbound</Direction>
            <Error/>
            <RequestProcessingTime>20180831114042</RequestProcessingTime>
            <Trips>
              <Trip>
                <TripDestination>Millennium</TripDestination>
                <TripStartTime>11:00</TripStartTime>
                <AdjustedScheduleTime>12</AdjustedScheduleTime>
                <AdjustmentAge>0.44</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>4EB - DD</BusType>
                <Latitude>45.404710</Latitude>
                <Longitude>-75.732058</Longitude>
                <GPSSpeed>15.9</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Millennium</TripDestination>
                <TripStartTime>11:15</TripStartTime>
                <AdjustedScheduleTime>25</AdjustedScheduleTime>
                <AdjustmentAge>0.49</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>6EB - 60</BusType>
                <Latitude>45.344501</Latitude>
                <Longitude>-75.758024</Longitude>
                <GPSSpeed>19.2</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Millennium</TripDestination>
                <TripStartTime>11:30</TripStartTime>
                <AdjustedScheduleTime>42</AdjustedScheduleTime>
                <AdjustmentAge>0.25</AdjustmentAge>
                <LastTripOfSchedule>false</LastTripOfSchedule>
                <BusType>4LB - DD</BusType>
                <Latitude>45.281069</Latitude>
                <Longitude>-75.721529</Longitude>
                <GPSSpeed>33.5</GPSSpeed>
              </Trip>
            </Trips>
          </RouteDirection>
        </Route>
      </GetNextTripsForStopResult>
    </GetNextTripsForStopResponse>
  </soap:Body>
</soap:Envelope>`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawXMLString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	nextTrips, err := c.GetNextTripsForStop(context.TODO(), "94", "3020")
	if err != nil {
		t.Fatal(err)
	}

	if nextTrips.StopNo != "3020" {
		t.Fatal("Unexpected StopNo in returned NextTripsForStop")
	}
	if nextTrips.StopLabel != "LAURIER STATION" {
		t.Fatal("Unexpected StopLabel in returned NextTripsForStop")
	}
	if nextTrips.Error != "TestErrorStringHere" {
		t.Fatal("Unexpected Error in returned NextTripsForStop")
	}

	tz, err := time.LoadLocation("America/Toronto")
	if err != nil {
		t.Fatal(err)
	}

	expectedFirstRouteDirection := RouteDirection{
		RouteNo:               "94",
		RouteLabel:            "Riverview",
		Direction:             "Westbound",
		Error:                 "TestRouteDirectionErrorHere",
		RequestProcessingTime: time.Date(2018, time.August, 31, 11, 40, 42, 0, tz),
	}

	if nextTrips.RouteDirections[0].RouteNo != expectedFirstRouteDirection.RouteNo {
		t.Fatal("Unexpected routeNo in first route direction in returned NextTripsForStop")
	}
	if nextTrips.RouteDirections[0].RouteLabel != expectedFirstRouteDirection.RouteLabel {
		t.Fatal("Unexpected RouteLabel in first route direction in returned NextTripsForStop")
	}
	if nextTrips.RouteDirections[0].Direction != expectedFirstRouteDirection.Direction {
		t.Fatal("Unexpected Direction in first route direction in returned NextTripsForStop")
	}
	if nextTrips.RouteDirections[0].Error != expectedFirstRouteDirection.Error {
		t.Fatal("Unexpected Direction in first route direction in returned NextTripsForStop")
	}
	if !nextTrips.RouteDirections[0].RequestProcessingTime.Equal(expectedFirstRouteDirection.RequestProcessingTime) {
		t.Fatal("Unexpected RequestProcessingTime in first route direction in returned NextTripsForStop")
	}

	expectedTrips := []Trip{
		{
			TripDestination:      "Millennium",
			TripStartTime:        "11:00",
			AdjustedScheduleTime: 12,
			AdjustmentAge:        0.44,
			LastTripOfSchedule:   LastTripOfSchedule{Set: true, Value: false},
			BusType:              "4EB - DD",
			Latitude:             Latitude{Set: true, Value: 45.404710},
			Longitude:            Longitude{Set: true, Value: -75.732058},
			GPSSpeed:             GPSSpeed{Set: true, Value: 15.9},
		},
		{
			TripDestination:      "Millennium",
			TripStartTime:        "11:15",
			AdjustedScheduleTime: 25,
			AdjustmentAge:        0.49,
			LastTripOfSchedule:   LastTripOfSchedule{Set: true, Value: false},
			BusType:              "6EB - 60",
			Latitude:             Latitude{Set: true, Value: 45.344501},
			Longitude:            Longitude{Set: true, Value: -75.758024},
			GPSSpeed:             GPSSpeed{Set: true, Value: 19.2},
		},
		{
			TripDestination:      "Millennium",
			TripStartTime:        "11:30",
			AdjustedScheduleTime: 42,
			AdjustmentAge:        0.25,
			LastTripOfSchedule:   LastTripOfSchedule{Set: true, Value: false},
			BusType:              "4LB - DD",
			Latitude:             Latitude{Set: true, Value: 45.281069},
			Longitude:            Longitude{Set: true, Value: -75.721529},
			GPSSpeed:             GPSSpeed{Set: true, Value: 33.5},
		},
	}

	for i, trip := range nextTrips.RouteDirections[1].Trips {
		if trip != expectedTrips[i] {
			t.Fatal("Trip doesn't match expected trip in returned NextTripsForStop")
		}
	}
}

func TestGetNextTripsForStopAllRoutes(t *testing.T) {
	rawXMLString := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetRouteSummaryForStopResponse xmlns="http://octranspo.com">
      <GetRouteSummaryForStopResult>
        <StopNo xmlns="http://tempuri.org/">3020</StopNo>
        <StopDescription xmlns="http://tempuri.org/">LAURIER STATION</StopDescription>
        <Error xmlns="http://tempuri.org/"/>
        <Routes xmlns="http://tempuri.org/">
          <Route>
            <RouteNo>97</RouteNo>
            <DirectionID>0</DirectionID>
            <Direction>Eastbound</Direction>
            <RouteHeading>Airport / Aéroport</RouteHeading>
            <Trips>
              <Trip>
                <TripDestination>Airport / Aéroport</TripDestination>
                <TripStartTime>13:14</TripStartTime>
                <AdjustedScheduleTime>8</AdjustedScheduleTime>
                <AdjustmentAge>0.42</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>6EB - 60</BusType>
                <Latitude>45.413769</Latitude>
                <Longitude>-75.710547</Longitude>
                <GPSSpeed>25.7</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Airport / Aéroport</TripDestination>
                <TripStartTime>13:29</TripStartTime>
                <AdjustedScheduleTime>22</AdjustedScheduleTime>
                <AdjustmentAge>-1</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>4LB - DD</BusType>
                <Latitude/>
                <Longitude/>
                <GPSSpeed/>
              </Trip>
              <Trip>
                <TripDestination>South Keys</TripDestination>
                <TripStartTime>12:43</TripStartTime>
                <AdjustedScheduleTime>23</AdjustedScheduleTime>
                <AdjustmentAge>0.40</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType> - DD</BusType>
                <Latitude>45.365881</Latitude>
                <Longitude>-75.783288</Longitude>
                <GPSSpeed>15.9</GPSSpeed>
              </Trip>
            </Trips>
          </Route>
          <Route>
            <RouteNo>97</RouteNo>
            <DirectionID>1</DirectionID>
            <Direction>Westbound</Direction>
            <RouteHeading>Bells Corners</RouteHeading>
            <Trips>
              <Trip>
                <TripDestination>Bayshore</TripDestination>
                <TripStartTime>12:57</TripStartTime>
                <AdjustedScheduleTime>2</AdjustedScheduleTime>
                <AdjustmentAge>0.44</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>4E - DEH</BusType>
                <Latitude>45.418886</Latitude>
                <Longitude>-75.678187</Longitude>
                <GPSSpeed>22.4</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Bells Corners</TripDestination>
                <TripStartTime>13:12</TripStartTime>
                <AdjustedScheduleTime>15</AdjustedScheduleTime>
                <AdjustmentAge>0.51</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>6EB - 60</BusType>
                <Latitude>45.387714</Latitude>
                <Longitude>-75.673109</Longitude>
                <GPSSpeed>71.9</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>Tunney's Pasture</TripDestination>
                <TripStartTime>13:02</TripStartTime>
                <AdjustedScheduleTime>16</AdjustedScheduleTime>
                <AdjustmentAge>0.51</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType> - DD</BusType>
                <Latitude>45.384286</Latitude>
                <Longitude>-75.676965</Longitude>
                <GPSSpeed>18.1</GPSSpeed>
              </Trip>
            </Trips>
          </Route>
          <Route>
            <RouteNo>98</RouteNo>
            <DirectionID>1</DirectionID>
            <Direction>Northbound</Direction>
            <RouteHeading>Tunney's Pasture</RouteHeading>
            <Trips>
              <Trip>
                <TripDestination>LeBreton</TripDestination>
                <TripStartTime>12:46</TripStartTime>
                <AdjustedScheduleTime>14</AdjustedScheduleTime>
                <AdjustmentAge>0.37</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>6EB - 60</BusType>
                <Latitude>45.410505</Latitude>
                <Longitude>-75.664115</Longitude>
                <GPSSpeed>51.1</GPSSpeed>
              </Trip>
              <Trip>
                <TripDestination>LeBreton</TripDestination>
                <TripStartTime>13:01</TripStartTime>
                <AdjustedScheduleTime>26</AdjustedScheduleTime>
                <AdjustmentAge>0.44</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>6EAB - 60</BusType>
                <Latitude/>
                <Longitude/>
                <GPSSpeed/>
              </Trip>
              <Trip>
                <TripDestination>LeBreton</TripDestination>
                <TripStartTime>13:16</TripStartTime>
                <AdjustedScheduleTime>42</AdjustedScheduleTime>
                <AdjustmentAge>0.47</AdjustmentAge>
                <LastTripOfSchedule/>
                <BusType>6EB - 60</BusType>
                <Latitude>45.369819</Latitude>
                <Longitude>-75.613623</Longitude>
                <GPSSpeed>41.1</GPSSpeed>
              </Trip>
            </Trips>
          </Route>
        </Routes>
      </GetRouteSummaryForStopResult>
    </GetRouteSummaryForStopResponse>
  </soap:Body>
</soap:Envelope>`

	rawHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, rawXMLString)
	}
	ts := httptest.NewServer(http.HandlerFunc(rawHandler))
	defer ts.Close()

	c := NewConnection("", "")
	c.cAPIURLPrefix = ts.URL + "/"

	nextTripsAllRoutes, err := c.GetNextTripsForStopAllRoutes(context.TODO(), "3020")
	if err != nil {
		t.Fatal(err)
	}

	if nextTripsAllRoutes.StopNo != "3020" {
		t.Fatal("Unexpected StopNo in returned NextTripsForStopAllRoutes")
	}
	if nextTripsAllRoutes.StopDescription != "LAURIER STATION" {
		t.Fatal("Unexpected StopDescription in returned NextTripsForStopAllRoutes")
	}

	if nextTripsAllRoutes.Routes[0].RouteNo != "97" {
		t.Fatal("Unexpected RouteNo in first route in returned NextTripsForStopAllRoutes")
	}
	if nextTripsAllRoutes.Routes[0].DirectionID != "0" {
		t.Fatal("Unexpected DirectionID in first route in returned NextTripsForStopAllRoutes")
	}
	if nextTripsAllRoutes.Routes[0].Direction != "Eastbound" {
		t.Fatal("Unexpected Direction in first route in returned NextTripsForStopAllRoutes")
	}
	if nextTripsAllRoutes.Routes[0].RouteHeading != "Airport / Aéroport" {
		t.Fatal("Unexpected RouteHeading in first route in returned NextTripsForStopAllRoutes")
	}

	expectedTrips := []Trip{
		{
			TripDestination:      "LeBreton",
			TripStartTime:        "12:46",
			AdjustedScheduleTime: 14,
			AdjustmentAge:        0.37,
			LastTripOfSchedule:   LastTripOfSchedule{Set: false},
			BusType:              "6EB - 60",
			Latitude:             Latitude{Set: true, Value: 45.410505},
			Longitude:            Longitude{Set: true, Value: -75.664115},
			GPSSpeed:             GPSSpeed{Set: true, Value: 51.1},
		},
		{
			TripDestination:      "LeBreton",
			TripStartTime:        "13:01",
			AdjustedScheduleTime: 26,
			AdjustmentAge:        0.44,
			LastTripOfSchedule:   LastTripOfSchedule{Set: false},
			BusType:              "6EAB - 60",
			Latitude:             Latitude{Set: false},
			Longitude:            Longitude{Set: false},
			GPSSpeed:             GPSSpeed{Set: false},
		},
		{
			TripDestination:      "LeBreton",
			TripStartTime:        "13:16",
			AdjustedScheduleTime: 42,
			AdjustmentAge:        0.47,
			LastTripOfSchedule:   LastTripOfSchedule{Set: false},
			BusType:              "6EB - 60",
			Latitude:             Latitude{Set: true, Value: 45.369819},
			Longitude:            Longitude{Set: true, Value: -75.613623},
			GPSSpeed:             GPSSpeed{Set: true, Value: 41.1},
		},
	}

	for i, trip := range nextTripsAllRoutes.Routes[2].Trips {
		if trip != expectedTrips[i] {
			t.Fatal("Trip doesn't match expected trip in returned NextTripsForStopAllRoutes")
		}
	}

}
