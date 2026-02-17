package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/jontk/notion-cli/internal/models"
)

// parseDateTime parses a date-time string in "2006-01-02 15:04" format
func parseDateTime(dateTimeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		// Try date-only format
		t, err = time.Parse("2006-01-02", dateTimeStr)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}

// CreateEvent creates a new event in the Notion database
func (c *Client) CreateEvent(ctx context.Context, input models.EventInput, databaseID string) (*models.Event, error) {
	properties := notionapi.Properties{
		"Title": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Title,
					},
				},
			},
		},
	}

	// Add date
	if input.Date != "" {
		t := parseDateTime(input.Date)
		date := notionapi.Date(t)
		properties["Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
			},
		}
	}

	// Add type
	if input.Type != "" {
		properties["Type"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Type,
			},
		}
	}

	// Add status
	if input.Status != "" {
		properties["Status"] = notionapi.MultiSelectProperty{
			MultiSelect: []notionapi.Option{
				{
					Name: input.Status,
				},
			},
		}
	}

	// Add location
	if input.Location != "" {
		properties["Location"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Location,
					},
				},
			},
		}
	}

	// Add attendees
	if len(input.Attendees) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Attendees))
		for _, attendee := range input.Attendees {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: attendee,
			})
		}
		properties["Attendees"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	// Add notes
	if input.Notes != "" {
		properties["Notes"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Notes,
					},
				},
			},
		}
	}

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: properties,
	}

	page, err := c.api.Page.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return c.pageToEvent(ctx, page)
}

// GetEvent retrieves a single event by ID
func (c *Client) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	page, err := c.api.Page.Get(ctx, notionapi.PageID(eventID))
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return c.pageToEvent(ctx, page)
}

// UpdateEvent updates an existing event
func (c *Client) UpdateEvent(ctx context.Context, eventID string, input models.EventInput) (*models.Event, error) {
	properties := notionapi.Properties{}

	if input.Title != "" {
		properties["Name"] = notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Title,
					},
				},
			},
		}
	}

	if input.Date != "" {
		t := parseDateTime(input.Date)
		date := notionapi.Date(t)
		properties["Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
			},
		}
	}

	if input.Type != "" {
		properties["Type"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Type,
			},
		}
	}

	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{
				Name: input.Status,
			},
		}
	}

	if input.Location != "" {
		properties["Location"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Location,
					},
				},
			},
		}
	}

	if len(input.Attendees) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Attendees))
		for _, attendee := range input.Attendees {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: attendee,
			})
		}
		properties["Attendees"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	if input.Notes != "" {
		properties["Notes"] = notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Notes,
					},
				},
			},
		}
	}

	req := &notionapi.PageUpdateRequest{
		Properties: properties,
	}

	page, err := c.api.Page.Update(ctx, notionapi.PageID(eventID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return c.pageToEvent(ctx, page)
}

// CancelEvent marks an event as cancelled
func (c *Client) CancelEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return c.UpdateEvent(ctx, eventID, models.EventInput{Status: "Cancelled"})
}

// EventQueryOptions holds options for querying events
type EventQueryOptions struct {
	Type       string
	Status     string
	DateAfter  string
	DateBefore string
	Limit      int
}

// QueryEvents queries events from a database with filters
func (c *Client) QueryEvents(ctx context.Context, databaseID string, opts EventQueryOptions) ([]models.Event, error) {
	var filters []notionapi.Filter

	if opts.Type != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Type",
			Select: &notionapi.SelectFilterCondition{
				Equals: opts.Type,
			},
		})
	}

	if opts.Status != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Status",
			MultiSelect: &notionapi.MultiSelectFilterCondition{
				Contains: opts.Status,
			},
		})
	}

	var filter notionapi.Filter
	if len(filters) > 1 {
		filter = &notionapi.AndCompoundFilter{filters[0], filters[1]}
	} else if len(filters) == 1 {
		filter = filters[0]
	}

	sorts := []notionapi.SortObject{
		{
			Property:  "Date",
			Direction: notionapi.SortOrderASC,
		},
	}

	var allEvents []models.Event
	var cursor *string
	limit := opts.Limit
	if limit == 0 {
		limit = 100
	}

	for {
		pageSize := 100
		remaining := limit - len(allEvents)
		if remaining < pageSize {
			pageSize = remaining
		}

		req := &notionapi.DatabaseQueryRequest{
			Filter:   filter,
			Sorts:    sorts,
			PageSize: pageSize,
		}

		if cursor != nil {
			req.StartCursor = notionapi.Cursor(*cursor)
		}

		resp, err := c.api.Database.Query(ctx, notionapi.DatabaseID(databaseID), req)
		if err != nil {
			return nil, fmt.Errorf("failed to query events: %w", err)
		}

		for _, page := range resp.Results {
			event, err := c.pageToEvent(ctx, &page)
			if err != nil {
				return nil, err
			}
			allEvents = append(allEvents, *event)
		}

		if !resp.HasMore || len(allEvents) >= limit {
			break
		}
		cursorStr := string(resp.NextCursor)
		cursor = &cursorStr
	}

	return allEvents, nil
}

// GetTodaysEvents returns events happening today
func (c *Client) GetTodaysEvents(ctx context.Context, databaseID string) ([]models.Event, error) {
	today := time.Now().Format("2006-01-02")
	return c.QueryEvents(ctx, databaseID, EventQueryOptions{
		DateAfter:  today,
		DateBefore: today,
		Limit:      100,
	})
}

// GetWeeksEvents returns events for this week
func (c *Client) GetWeeksEvents(ctx context.Context, databaseID string) ([]models.Event, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return c.QueryEvents(ctx, databaseID, EventQueryOptions{
		DateAfter:  startOfWeek.Format("2006-01-02"),
		DateBefore: endOfWeek.Format("2006-01-02"),
		Limit:      100,
	})
}

// pageToEvent converts a Notion page to our Event model
func (c *Client) pageToEvent(ctx context.Context, page *notionapi.Page) (*models.Event, error) {
	event := &models.Event{
		ID:        string(page.ID),
		URL:       page.URL,
		CreatedAt: page.CreatedTime.String(),
		UpdatedAt: page.LastEditedTime.String(),
	}

	// Extract title
	if titleProp, ok := page.Properties["Title"].(*notionapi.TitleProperty); ok {
		if len(titleProp.Title) > 0 {
			event.Title = titleProp.Title[0].PlainText
		}
	}

	// Extract date
	if dateProp, ok := page.Properties["Date"].(*notionapi.DateProperty); ok {
		if dateProp.Date != nil && dateProp.Date.Start != nil {
			event.Date = dateProp.Date.Start.String()
		}
	}

	// Extract type
	if typeProp, ok := page.Properties["Type"].(*notionapi.SelectProperty); ok {
		if typeProp.Select.Name != "" {
			event.Type = typeProp.Select.Name
		}
	}

	// Extract status
	if statusProp, ok := page.Properties["Status"].(*notionapi.MultiSelectProperty); ok {
		if len(statusProp.MultiSelect) > 0 {
			event.Status = statusProp.MultiSelect[0].Name
		}
	}

	// Extract location
	if locationProp, ok := page.Properties["Location"].(*notionapi.RichTextProperty); ok {
		if len(locationProp.RichText) > 0 {
			event.Location = locationProp.RichText[0].PlainText
		}
	}

	// Extract attendees
	if attendeesProp, ok := page.Properties["Attendees"].(*notionapi.MultiSelectProperty); ok {
		attendees := make([]string, 0, len(attendeesProp.MultiSelect))
		for _, opt := range attendeesProp.MultiSelect {
			attendees = append(attendees, opt.Name)
		}
		event.Attendees = attendees
	}

	// Extract notes
	if notesProp, ok := page.Properties["Notes"].(*notionapi.RichTextProperty); ok {
		if len(notesProp.RichText) > 0 {
			event.Notes = notesProp.RichText[0].PlainText
		}
	}

	return event, nil
}
