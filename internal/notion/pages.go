package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/jontk/notion-cli/internal/models"
)

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// CreatePost creates a new post in the Notion database
func (c *Client) CreatePost(ctx context.Context, input models.PostInput, databaseID string) (*models.Post, error) {
	properties := notionapi.Properties{
		"Name": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: input.Title,
					},
				},
			},
		},
	}

	// Add status if provided
	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{
				Name: input.Status,
			},
		}
	}

	// Add platforms if provided
	if len(input.Platforms) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Platforms))
		for _, platform := range input.Platforms {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: platform,
			})
		}
		properties["Platforms"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	// Add publish date if provided
	if input.PublishDate != "" {
		t := parseDate(input.PublishDate)
		date := notionapi.Date(t)
		properties["Publish Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
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

	// Add content as children blocks if provided
	if input.Content != "" {
		req.Children = contentToBlocks(input.Content)
	}

	page, err := c.api.Page.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return c.pageToPost(ctx, page)
}

// GetPost retrieves a single post by ID
func (c *Client) GetPost(ctx context.Context, pageID string) (*models.Post, error) {
	page, err := c.api.Page.Get(ctx, notionapi.PageID(pageID))
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	return c.pageToPost(ctx, page)
}

// UpdatePost updates an existing post
func (c *Client) UpdatePost(ctx context.Context, pageID string, input models.PostInput) (*models.Post, error) {
	properties := notionapi.Properties{}

	// Only update provided fields
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

	if len(input.Platforms) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(input.Platforms))
		for _, platform := range input.Platforms {
			multiSelect = append(multiSelect, notionapi.Option{
				Name: platform,
			})
		}
		properties["Platforms"] = notionapi.MultiSelectProperty{
			MultiSelect: multiSelect,
		}
	}

	if input.PublishDate != "" {
		t := parseDate(input.PublishDate)
		date := notionapi.Date(t)
		properties["Publish Date"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: &date,
			},
		}
	}

	req := &notionapi.PageUpdateRequest{
		Properties: properties,
	}

	page, err := c.api.Page.Update(ctx, notionapi.PageID(pageID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	// Append content blocks if provided
	if input.Content != "" {
		blocks := contentToBlocks(input.Content)
		for _, block := range blocks {
			_, err := c.api.Block.AppendChildren(ctx, notionapi.BlockID(pageID), &notionapi.AppendBlockChildrenRequest{
				Children: []notionapi.Block{block},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to append content blocks: %w", err)
			}
		}
	}

	return c.pageToPost(ctx, page)
}

// ArchivePost archives a post
func (c *Client) ArchivePost(ctx context.Context, pageID string) (*models.Post, error) {
	req := &notionapi.PageUpdateRequest{
		Archived: true,
	}

	page, err := c.api.Page.Update(ctx, notionapi.PageID(pageID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to archive page: %w", err)
	}

	return c.pageToPost(ctx, page)
}

// QueryOptions holds options for querying posts
type QueryOptions struct {
	Status    string
	Platform  string
	Sort      string
	Order     string
	Limit     int
}

// QueryPosts queries posts from a database with filters
func (c *Client) QueryPosts(ctx context.Context, databaseID string, opts QueryOptions) ([]models.Post, error) {
	var filter notionapi.Filter

	// Build filters based on options
	if opts.Status != "" && opts.Platform != "" {
		// Both status and platform filters
		filter = &notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: "Status",
				Status: &notionapi.StatusFilterCondition{
					Equals: opts.Status,
				},
			},
			notionapi.PropertyFilter{
				Property: "Platforms",
				MultiSelect: &notionapi.MultiSelectFilterCondition{
					Contains: opts.Platform,
				},
			},
		}
	} else if opts.Status != "" {
		// Only status filter
		filter = notionapi.PropertyFilter{
			Property: "Status",
			Status: &notionapi.StatusFilterCondition{
				Equals: opts.Status,
			},
		}
	} else if opts.Platform != "" {
		// Only platform filter
		filter = notionapi.PropertyFilter{
			Property: "Platforms",
			MultiSelect: &notionapi.MultiSelectFilterCondition{
				Contains: opts.Platform,
			},
		}
	}

	// Set up sorting
	sortField := opts.Sort
	if sortField == "" {
		sortField = "created_time"
	}
	sortOrder := notionapi.SortOrderDESC
	if opts.Order == "ascending" {
		sortOrder = notionapi.SortOrderASC
	}

	sorts := []notionapi.SortObject{
		{
			Timestamp: notionapi.TimestampType(sortField),
			Direction: sortOrder,
		},
	}

	// Query with pagination
	var allPosts []models.Post
	var cursor *string
	limit := opts.Limit
	if limit == 0 {
		limit = 100
	}

	for {
		pageSize := 100
		remaining := limit - len(allPosts)
		if remaining < pageSize {
			pageSize = remaining
		}

		req := &notionapi.DatabaseQueryRequest{
			Filter:    filter,
			Sorts:     sorts,
			PageSize:  pageSize,
		}

		if cursor != nil {
			req.StartCursor = notionapi.Cursor(*cursor)
		}

		resp, err := c.api.Database.Query(ctx, notionapi.DatabaseID(databaseID), req)
		if err != nil {
			return nil, fmt.Errorf("failed to query database: %w", err)
		}

		for _, page := range resp.Results {
			post, err := c.pageToPost(ctx, &page)
			if err != nil {
				return nil, err
			}
			allPosts = append(allPosts, *post)
		}

		if !resp.HasMore || len(allPosts) >= limit {
			break
		}
		cursorStr := string(resp.NextCursor)
		cursor = &cursorStr
	}

	return allPosts, nil
}

// pageToPost converts a Notion page to our Post model
func (c *Client) pageToPost(ctx context.Context, page *notionapi.Page) (*models.Post, error) {
	post := &models.Post{
		ID:        string(page.ID),
		URL:       page.URL,
		CreatedAt: page.CreatedTime.String(),
		UpdatedAt: page.LastEditedTime.String(),
	}

	// Extract title
	if titleProp, ok := page.Properties["Name"].(*notionapi.TitleProperty); ok {
		if len(titleProp.Title) > 0 {
			post.Title = titleProp.Title[0].PlainText
		}
	}

	// Extract status
	if statusProp, ok := page.Properties["Status"].(*notionapi.StatusProperty); ok {
		post.Status = statusProp.Status.Name
	}

	// Extract platforms
	if platformsProp, ok := page.Properties["Platforms"].(*notionapi.MultiSelectProperty); ok {
		platforms := make([]string, 0, len(platformsProp.MultiSelect))
		for _, opt := range platformsProp.MultiSelect {
			platforms = append(platforms, opt.Name)
		}
		post.Platforms = platforms
	}

	// Extract publish date
	if dateProp, ok := page.Properties["Publish Date"].(*notionapi.DateProperty); ok {
		if dateProp.Date != nil && dateProp.Date.Start != nil {
			post.PublishDate = dateProp.Date.Start.String()
		}
	}

	// Get content from blocks
	content, err := c.GetPageContent(ctx, string(page.ID))
	if err != nil {
		// Don't fail the whole operation if content fetch fails
		post.Content = ""
	} else {
		post.Content = content
	}

	return post, nil
}
