package gooctranspoapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGTFSAgency(t *testing.T) {
	rawJSONString := `{"Query":{"table":"agency","direction":"ASC","format":"json"},
                     "Gtfs":[{"id":"1",
                             "agency_name":"Test Agency",
                             "agency_url":"http://test.com",
                             "agency_timezone":"America/Toronto",
                             "agency_lang":"",
                             "agency_phone":""
                             }]}`

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
		t.Fatal("Unexpected agency_name in returned GTFSAgency")
	}
	if agency.Gtfs[0].AgencyTimezone != "America/Toronto" {
		t.Fatal("Unexpected agency_timezone in returned GTFSAgency")
	}
}
