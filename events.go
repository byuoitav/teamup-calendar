package teamup

import "time"

type event struct {
	ID                 string    `json:"id"`
	RemoteID           string    `json:"remote_id"`
	SeriesID           string    `json:"series_id"`
	SubCalendarIDArray []int     `json:"subcalendar_id"`
	SubCalendarID      int       `json:"subcalendar_id"`
	StartDate          time.Time `json:"start_dt"`
	EndDate            time.Time `json:"end_dt"`
	AllDay             bool      `json:"all_day"`
	Title              string    `json:"title"`
	Who                string    `json:"who"`
	Location           string    `json:"location"`
	Notes              string    `json:"notes"`
	RecurRule          string    `json:"rrule"`
	RecurInstStart     string    `json:"ristart_dt"`
	RecurSeriesStart   string    `json:"rsstart_dt"`
	TimeZone           string    `json:"tz"`
	Version            string    `json:"version"`
	ReadOnly           bool      `json:"readonly"`
	CreationDate       time.Time `json:"creation_dt"`
	UpdateDate         time.Time `json:"update_dt"`
}

//eventResponse represents the json object returned by the TeamUp API events get request
type eventResponse struct {
	Events    []event `json:"events"`
	TimeStamp int     `json:"timestamp"`
}

//eventSend represents the event data to be sent to the TeamUp API
type eventSend struct {
	SubCalendarID int    `json:"subcalendar_id"`
	StartDate     string `json:"start_dt"`
	EndDate       string `json:"end_dt"`
	AllDay        bool   `json:"all_day"`
	RecurRule     string `json:"rrule"`
	Title         string `json:"title"`
	Who           string `json:"who"`
	Location      string `json:"location"`
	Notes         string `json:"notes"`
}

//subcalendarList holds a list of subcalendars
type subcalendarList struct {
	Subcalendars []subcalendar `json:"subcalendars"`
}

//subcalendar represents a subcalendar in json
type subcalendar struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Color        int    `json:"color"`
	Overlap      bool   `json:"overlap"`
	ReadOnly     bool   `json:"readonly"`
	CreationDate string `json:"creation_dt"`
	UpdateDate   string `json:"update_dt"`
}
