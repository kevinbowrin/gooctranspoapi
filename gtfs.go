package gooctranspoapi

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
