package teamup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/byuoitav/scheduler/calendars"
)

type Calendar struct {
	APIKey   string
	Password string // optional

	// CalendarID is like a "user ID" (maps to calendarID in teamup's docs)
	CalendarID string

	RoomID string
}

func (c *Calendar) GetEvents(ctx context.Context) ([]calendars.Event, error) {
	subCalID, err := c.GetSubcalendarID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get subcalendar id: %s", err)
	}

	url := "https://api.teamup.com/" + c.CalendarID + "/events"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating get request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teamup-Token", c.APIKey)
	if len(c.Password) > 0 {
		req.Header.Set("Teamup-Password", c.Password)
	}

	query := req.URL.Query()
	query.Add("subcalendarId[]", strconv.Itoa(subCalID))
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending get request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var eventResponse eventResponse
	if err := json.Unmarshal(body, &eventResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling json event data: %w", err)
	}

	var events []calendars.Event
	for _, event := range eventResponse.Events {
		events = append(events, calendars.Event{
			Title:     event.Title,
			StartTime: event.StartDate,
			EndTime:   event.EndDate,
		})
	}

	return events, nil
}

func (c *Calendar) CreateEvent(ctx context.Context, event calendars.Event) error {
	subCalID, err := c.GetSubcalendarID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get subcalendar id: %s", err)
	}

	// Translate event to team up event
	teamUpEvent := eventSend{
		SubCalendarID: subCalID,
		Title:         event.Title,
		StartDate:     event.StartTime,
		EndDate:       event.EndTime,
	}

	reqBody, err := json.Marshal(teamUpEvent)
	if err != nil {
		return fmt.Errorf("error marshalling event data into json: %w", err)
	}

	url := "https://api.teamup.com/" + c.CalendarID + "/events"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("Error creating post request | %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Teamup-Token", c.APIKey)
	if len(c.Password) > 0 {
		req.Header.Set("Teamup-Password", c.Password)
	}

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("invalid response code %v: %s", resp.StatusCode, body)
	}

	return nil
}

// GetSubcalendarID sends a request to the teamup api to get the appropriate subcalendar id
func (c *Calendar) GetSubcalendarID(ctx context.Context) (int, error) {
	url := "https://api.teamup.com/" + c.CalendarID + "/subcalendars"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Teamup-Token", c.APIKey)
	if len(c.Password) > 0 {
		req.Header.Set("Teamup-Password", c.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("invalid response, status code %v: %s", resp.StatusCode, body)
	}

	var subcalResponse subcalendarList
	if err := json.Unmarshal(body, &subcalResponse); err != nil {
		return 0, fmt.Errorf("error unmarshalling json: %w", err)
	}

	for _, subCal := range subcalResponse.Subcalendars {
		if subCal.Name == c.RoomID {
			return subCal.ID, nil
		}
	}

	return 0, fmt.Errorf("no calendar found with roomID %q", c.RoomID)
}
