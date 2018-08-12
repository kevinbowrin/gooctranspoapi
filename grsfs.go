package octranspoapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"net/http"
	"net/url"
	"strings"
)

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

func (c *Connection) GetRouteSummaryForStop(ctx context.Context, stopNo string) (*RouteSummaryForStop, error) {
	u, err := url.Parse(ApiURLPrefix + "GetRouteSummaryForStop")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("appID", c.Id)
	v.Set("apiKey", c.Key)
	v.Set("stopNo", stopNo)

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

	dec := xml.NewDecoder(resp.Body)
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false

	data := &RouteSummaryForStop{}
	err = dec.Decode(data)
	resp.Body.Close()
	return data, err
}
