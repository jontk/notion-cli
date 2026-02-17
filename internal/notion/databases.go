package notion

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/jontk/notion-cli/internal/models"
)

// ListDatabases lists all databases accessible to the integration
func (c *Client) ListDatabases(ctx context.Context) ([]models.DatabaseInfo, error) {
	filter := notionapi.SearchFilter{
		Property: "object",
		Value:    "database",
	}

	req := &notionapi.SearchRequest{
		Filter: filter,
	}

	resp, err := c.api.Search.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search databases: %w", err)
	}

	databases := make([]models.DatabaseInfo, 0, len(resp.Results))
	for _, result := range resp.Results {
		if db, ok := result.(*notionapi.Database); ok {
			title := ""
			if len(db.Title) > 0 {
				title = db.Title[0].PlainText
			}
			databases = append(databases, models.DatabaseInfo{
				ID:    string(db.ID),
				Title: title,
			})
		}
	}

	return databases, nil
}

// GetSchema retrieves the schema of a database
func (c *Client) GetSchema(ctx context.Context, databaseID string) (*models.Schema, error) {
	db, err := c.api.Database.Get(ctx, notionapi.DatabaseID(databaseID))
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	schema := &models.Schema{
		Properties: make(map[string]models.PropertyInfo),
	}

	for name, prop := range db.Properties {
		propInfo := models.PropertyInfo{
			Type: string(prop.GetType()),
		}

		// Add options for select/multi-select properties
		switch p := prop.(type) {
		case *notionapi.SelectPropertyConfig:
			options := make(map[string]any)
			optionNames := make([]string, 0, len(p.Select.Options))
			for _, opt := range p.Select.Options {
				optionNames = append(optionNames, opt.Name)
			}
			options["options"] = optionNames
			propInfo.Options = options

		case *notionapi.MultiSelectPropertyConfig:
			options := make(map[string]any)
			optionNames := make([]string, 0, len(p.MultiSelect.Options))
			for _, opt := range p.MultiSelect.Options {
				optionNames = append(optionNames, opt.Name)
			}
			options["options"] = optionNames
			propInfo.Options = options

		case *notionapi.StatusPropertyConfig:
			options := make(map[string]any)
			optionNames := make([]string, 0, len(p.Status.Options))
			for _, opt := range p.Status.Options {
				optionNames = append(optionNames, opt.Name)
			}
			options["options"] = optionNames
			propInfo.Options = options
		}

		schema.Properties[name] = propInfo
	}

	return schema, nil
}
