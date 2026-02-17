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

// richText builds a Notion rich text slice from a plain string
func richText(content string) []notionapi.RichText {
	return []notionapi.RichText{
		{Text: &notionapi.Text{Content: content}},
	}
}

// dateProperty builds a Notion date property from a YYYY-MM-DD string
func dateProperty(dateStr string) notionapi.DateProperty {
	t := parseDate(dateStr)
	d := notionapi.Date(t)
	return notionapi.DateProperty{
		Date: &notionapi.DateObject{Start: &d},
	}
}

// multiSelect builds a Notion multi-select property from a string slice
func multiSelect(values []string) notionapi.MultiSelectProperty {
	opts := make([]notionapi.Option, 0, len(values))
	for _, v := range values {
		opts = append(opts, notionapi.Option{Name: v})
	}
	return notionapi.MultiSelectProperty{MultiSelect: opts}
}

// CreatePost creates a new post in the Notion database
func (c *Client) CreatePost(ctx context.Context, input models.PostInput, databaseID string) (*models.Post, error) {
	properties := notionapi.Properties{
		"Title": notionapi.TitleProperty{
			Title: richText(input.Title),
		},
	}

	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{Name: input.Status},
		}
	}
	if input.Week > 0 {
		properties["Week"] = notionapi.NumberProperty{
			Number: float64(input.Week),
		}
	}
	if input.Pillar != "" {
		properties["Pillar"] = notionapi.SelectProperty{
			Select: notionapi.Option{Name: input.Pillar},
		}
	}
	if input.PublishDate != "" {
		properties["Publish Date"] = dateProperty(input.PublishDate)
	}
	if input.PublishedDate != "" {
		properties["Published Date"] = dateProperty(input.PublishedDate)
	}
	if input.BlogURL != "" {
		properties["Blog URL"] = notionapi.URLProperty{URL: input.BlogURL}
	}
	if len(input.DistributedTo) > 0 {
		properties["Distributed To"] = multiSelect(input.DistributedTo)
	}
	if input.DistributedDate != "" {
		properties["Distributed Date"] = dateProperty(input.DistributedDate)
	}
	if input.LinkedInDraft != "" {
		properties["LinkedIn Draft"] = notionapi.RichTextProperty{RichText: richText(input.LinkedInDraft)}
	}
	if input.TwitterThread != "" {
		properties["Twitter Thread"] = notionapi.RichTextProperty{RichText: richText(input.TwitterThread)}
	}
	if input.HNTitle != "" {
		properties["HN Title"] = notionapi.RichTextProperty{RichText: richText(input.HNTitle)}
	}
	if input.RedditTitle != "" {
		properties["Reddit Title"] = notionapi.RichTextProperty{RichText: richText(input.RedditTitle)}
	}
	if len(input.Hashtags) > 0 {
		properties["Hashtags"] = multiSelect(input.Hashtags)
	}

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: properties,
	}

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

	if input.Title != "" {
		properties["Title"] = notionapi.TitleProperty{
			Title: richText(input.Title),
		}
	}
	if input.Status != "" {
		properties["Status"] = notionapi.StatusProperty{
			Status: notionapi.Status{Name: input.Status},
		}
	}
	if input.Week > 0 {
		properties["Week"] = notionapi.NumberProperty{
			Number: float64(input.Week),
		}
	}
	if input.Pillar != "" {
		properties["Pillar"] = notionapi.SelectProperty{
			Select: notionapi.Option{Name: input.Pillar},
		}
	}
	if input.PublishDate != "" {
		properties["Publish Date"] = dateProperty(input.PublishDate)
	}
	if input.PublishedDate != "" {
		properties["Published Date"] = dateProperty(input.PublishedDate)
	}
	if input.BlogURL != "" {
		properties["Blog URL"] = notionapi.URLProperty{URL: input.BlogURL}
	}
	if len(input.DistributedTo) > 0 {
		properties["Distributed To"] = multiSelect(input.DistributedTo)
	}
	if input.DistributedDate != "" {
		properties["Distributed Date"] = dateProperty(input.DistributedDate)
	}
	if input.LinkedInDraft != "" {
		properties["LinkedIn Draft"] = notionapi.RichTextProperty{RichText: richText(input.LinkedInDraft)}
	}
	if input.TwitterThread != "" {
		properties["Twitter Thread"] = notionapi.RichTextProperty{RichText: richText(input.TwitterThread)}
	}
	if input.HNTitle != "" {
		properties["HN Title"] = notionapi.RichTextProperty{RichText: richText(input.HNTitle)}
	}
	if input.RedditTitle != "" {
		properties["Reddit Title"] = notionapi.RichTextProperty{RichText: richText(input.RedditTitle)}
	}
	if len(input.Hashtags) > 0 {
		properties["Hashtags"] = multiSelect(input.Hashtags)
	}

	req := &notionapi.PageUpdateRequest{
		Properties: properties,
	}

	page, err := c.api.Page.Update(ctx, notionapi.PageID(pageID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

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
		Archived:   true,
		Properties: notionapi.Properties{},
	}

	page, err := c.api.Page.Update(ctx, notionapi.PageID(pageID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to archive page: %w", err)
	}

	return c.pageToPost(ctx, page)
}

// QueryOptions holds options for querying posts
type QueryOptions struct {
	Status        string
	Pillar        string
	DistributedTo string
	Sort          string
	Order         string
	Limit         int
}

// QueryPosts queries posts from a database with filters
func (c *Client) QueryPosts(ctx context.Context, databaseID string, opts QueryOptions) ([]models.Post, error) {
	var filters []notionapi.Filter

	if opts.Status != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Status",
			Status:   &notionapi.StatusFilterCondition{Equals: opts.Status},
		})
	}
	if opts.Pillar != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property: "Pillar",
			Select:   &notionapi.SelectFilterCondition{Equals: opts.Pillar},
		})
	}
	if opts.DistributedTo != "" {
		filters = append(filters, notionapi.PropertyFilter{
			Property:    "Distributed To",
			MultiSelect: &notionapi.MultiSelectFilterCondition{Contains: opts.DistributedTo},
		})
	}

	var filter notionapi.Filter
	switch len(filters) {
	case 0:
		// no filter
	case 1:
		filter = filters[0]
	default:
		compound := notionapi.AndCompoundFilter(filters)
		filter = &compound
	}

	sortField := opts.Sort
	if sortField == "" {
		sortField = "created_time"
	}
	sortOrder := notionapi.SortOrderDESC
	if opts.Order == "ascending" {
		sortOrder = notionapi.SortOrderASC
	}

	sorts := []notionapi.SortObject{
		{Timestamp: notionapi.TimestampType(sortField), Direction: sortOrder},
	}

	var allPosts []models.Post
	var cursor *string
	limit := opts.Limit
	if limit == 0 {
		limit = 100
	}

	for {
		pageSize := 100
		if remaining := limit - len(allPosts); remaining < pageSize {
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

	if prop, ok := page.Properties["Title"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
		post.Title = prop.Title[0].PlainText
	}
	if prop, ok := page.Properties["Status"].(*notionapi.StatusProperty); ok {
		post.Status = prop.Status.Name
	}
	if prop, ok := page.Properties["Week"].(*notionapi.NumberProperty); ok {
		post.Week = int(prop.Number)
	}
	if prop, ok := page.Properties["Pillar"].(*notionapi.SelectProperty); ok {
		post.Pillar = prop.Select.Name
	}
	if prop, ok := page.Properties["Publish Date"].(*notionapi.DateProperty); ok && prop.Date != nil && prop.Date.Start != nil {
		post.PublishDate = prop.Date.Start.String()
	}
	if prop, ok := page.Properties["Published Date"].(*notionapi.DateProperty); ok && prop.Date != nil && prop.Date.Start != nil {
		post.PublishedDate = prop.Date.Start.String()
	}
	if prop, ok := page.Properties["Blog URL"].(*notionapi.URLProperty); ok {
		post.BlogURL = prop.URL
	}
	if prop, ok := page.Properties["Distributed To"].(*notionapi.MultiSelectProperty); ok {
		vals := make([]string, 0, len(prop.MultiSelect))
		for _, opt := range prop.MultiSelect {
			vals = append(vals, opt.Name)
		}
		post.DistributedTo = vals
	}
	if prop, ok := page.Properties["Distributed Date"].(*notionapi.DateProperty); ok && prop.Date != nil && prop.Date.Start != nil {
		post.DistributedDate = prop.Date.Start.String()
	}
	if prop, ok := page.Properties["LinkedIn Draft"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
		post.LinkedInDraft = prop.RichText[0].PlainText
	}
	if prop, ok := page.Properties["Twitter Thread"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
		post.TwitterThread = prop.RichText[0].PlainText
	}
	if prop, ok := page.Properties["HN Title"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
		post.HNTitle = prop.RichText[0].PlainText
	}
	if prop, ok := page.Properties["Reddit Title"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
		post.RedditTitle = prop.RichText[0].PlainText
	}
	if prop, ok := page.Properties["Hashtags"].(*notionapi.MultiSelectProperty); ok {
		vals := make([]string, 0, len(prop.MultiSelect))
		for _, opt := range prop.MultiSelect {
			vals = append(vals, opt.Name)
		}
		post.Hashtags = vals
	}

	content, err := c.GetPageContent(ctx, string(page.ID))
	if err != nil {
		post.Content = ""
	} else {
		post.Content = content
	}

	return post, nil
}
