package gooctranspoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookRawRouteSummaryForStop(t *testing.T) {
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
		t.Error(err)
	}

	if routeSummary.StopNo != "7659" {
		t.Error("Unexpected StopNo in returned RouteSummaryForStop")
	}
	if routeSummary.StopDescription != "BANK / FIFTH" {
		t.Error("Unexpected StopDescription in returned RouteSummaryForStop")
	}
	if routeSummary.Error != "TestErrorStringHere" {
		t.Error("Unexpected Error in returned RouteSummaryForStop")
	}

	expectedFirstRoute := Route{
		RouteNo:      "6",
		DirectionID:  "1",
		Direction:    "Northbound",
		RouteHeading: "Rockcliffe",
	}

	if routeSummary.Routes[0] != expectedFirstRoute {
		t.Error("Unexpected first route in returned RouteSummaryForStop")
	}

	expectedSecondRoute := Route{
		RouteNo:      "7",
		DirectionID:  "1",
		Direction:    "Eastbound",
		RouteHeading: "St-Laurent",
	}

	if routeSummary.Routes[1] != expectedSecondRoute {
		t.Error("Unexpected second route in returned RouteSummaryForStop")
	}
}
