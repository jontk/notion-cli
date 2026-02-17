package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/jontk/notion-cli/internal/models"
)

// CreateTask creates a new task in the Notion database
func (c *Client) CreateTask(ctx context.Context, input models.TaskInput, databaseID string) (*models.Task, error) {
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

	// Add status
	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{
				Name: input.Status,
			},
		}
	}

	// Add priority
	if input.Priority != "" {
		properties["Priority"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Priority,
			},
		}
	}

	// Add category
	if input.Category != "" {
		properties["Category"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Category,
			},
		}
	}

	// Add tags
	if len(input.Tags) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Tags))
		for _, tag := range input.Tags {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: tag,
			})
		}
		properties["Tags"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	// Add due date
	if input.DueDate != "" {
		t := parseDate(input.DueDate)
		date := notionapi.Date(t)
		properties["Due Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
			},
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
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return c.pageToTask(ctx, page)
}

// GetTask retrieves a single task by ID
func (c *Client) GetTask(ctx context.Context, taskID string) (*models.Task, error) {
	page, err := c.api.Page.Get(ctx, notionapi.PageID(taskID))
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return c.pageToTask(ctx, page)
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(ctx context.Context, taskID string, input models.TaskInput) (*models.Task, error) {
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

	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{
				Name: input.Status,
			},
		}
	}

	if input.Priority != "" {
		properties["Priority"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Priority,
			},
		}
	}

	if input.Category != "" {
		properties["Category"] = notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: input.Category,
			},
		}
	}

	if len(input.Tags) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Tags))
		for _, tag := range input.Tags {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: tag,
			})
		}
		properties["Tags"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	if input.DueDate != "" {
		t := parseDate(input.DueDate)
		date := notionapi.Date(t)
		properties["Due Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
			},
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

	page, err := c.api.Page.Update(ctx, notionapi.PageID(taskID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return c.pageToTask(ctx, page)
}

// CompleteTask marks a task as complete
func (c *Client) CompleteTask(ctx context.Context, taskID string) (*models.Task, error) {
	return c.UpdateTask(ctx, taskID, models.TaskInput{Status: "Done"})
}

// TaskQueryOptions holds options for querying tasks
type TaskQueryOptions struct {
	Status    string
	Priority  string
	Category  string
	DueBefore string
	DueAfter  string
	Limit     int
}

// QueryTasks queries tasks from a database with filters
func (c *Client) QueryTasks(ctx context.Context, databaseID string, opts TaskQueryOptions) ([]models.Task, error) {
	var filters []notionapi.Filter

	if opts.Status != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Status",
			Status: &notionapi.StatusFilterCondition{
				Equals: opts.Status,
			},
		})
	}

	if opts.Priority != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Priority",
			Select: &notionapi.SelectFilterCondition{
				Equals: opts.Priority,
			},
		})
	}

	if opts.Category != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Category",
			Select: &notionapi.SelectFilterCondition{
				Equals: opts.Category,
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
			Property:  "Due Date",
			Direction: notionapi.SortOrderASC,
		},
	}

	var allTasks []models.Task
	var cursor *string
	limit := opts.Limit
	if limit == 0 {
		limit = 100
	}

	for {
		pageSize := 100
		remaining := limit - len(allTasks)
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
			return nil, fmt.Errorf("failed to query tasks: %w", err)
		}

		for _, page := range resp.Results {
			task, err := c.pageToTask(ctx, &page)
			if err != nil {
				return nil, err
			}
			allTasks = append(allTasks, *task)
		}

		if !resp.HasMore || len(allTasks) >= limit {
			break
		}
		cursorStr := string(resp.NextCursor)
		cursor = &cursorStr
	}

	return allTasks, nil
}

// GetTodaysTasks returns tasks due today
func (c *Client) GetTodaysTasks(ctx context.Context, databaseID string) ([]models.Task, error) {
	today := time.Now().Format("2006-01-02")
	return c.QueryTasks(ctx, databaseID, TaskQueryOptions{
		Status:    "Todo",
		DueBefore: today,
		Limit:     100,
	})
}

// GetOverdueTasks returns overdue tasks
func (c *Client) GetOverdueTasks(ctx context.Context, databaseID string) ([]models.Task, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return c.QueryTasks(ctx, databaseID, TaskQueryOptions{
		Status:    "Todo",
		DueBefore: yesterday,
		Limit:     100,
	})
}

// pageToTask converts a Notion page to our Task model
func (c *Client) pageToTask(ctx context.Context, page *notionapi.Page) (*models.Task, error) {
	task := &models.Task{
		ID:        string(page.ID),
		URL:       page.URL,
		CreatedAt: page.CreatedTime.String(),
		UpdatedAt: page.LastEditedTime.String(),
	}

	// Extract title
	if titleProp, ok := page.Properties["Title"].(*notionapi.TitleProperty); ok {
		if len(titleProp.Title) > 0 {
			task.Title = titleProp.Title[0].PlainText
		}
	}

	// Extract status
	if statusProp, ok := page.Properties["Status"].(*notionapi.StatusProperty); ok {
		task.Status = statusProp.Status.Name
	}

	// Extract priority
	if priorityProp, ok := page.Properties["Priority"].(*notionapi.SelectProperty); ok {
		if priorityProp.Select.Name != "" {
			task.Priority = priorityProp.Select.Name
		}
	}

	// Extract category
	if categoryProp, ok := page.Properties["Category"].(*notionapi.SelectProperty); ok {
		if categoryProp.Select.Name != "" {
			task.Category = categoryProp.Select.Name
		}
	}

	// Extract tags
	if tagsProp, ok := page.Properties["Tags"].(*notionapi.MultiSelectProperty); ok {
		tags := make([]string, 0, len(tagsProp.MultiSelect))
		for _, opt := range tagsProp.MultiSelect {
			tags = append(tags, opt.Name)
		}
		task.Tags = tags
	}

	// Extract due date
	if dateProp, ok := page.Properties["Due Date"].(*notionapi.DateProperty); ok {
		if dateProp.Date != nil && dateProp.Date.Start != nil {
			task.DueDate = dateProp.Date.Start.String()
		}
	}

	// Extract notes
	if notesProp, ok := page.Properties["Notes"].(*notionapi.RichTextProperty); ok {
		if len(notesProp.RichText) > 0 {
			task.Notes = notesProp.RichText[0].PlainText
		}
	}

	return task, nil
}
