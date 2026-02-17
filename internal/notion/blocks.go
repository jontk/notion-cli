package notion

import (
	"context"
	"strings"

	"github.com/jomei/notionapi"
)

// GetPageContent retrieves all block children recursively and extracts plain text
func (c *Client) GetPageContent(ctx context.Context, pageID string) (string, error) {
	var content strings.Builder

	blocks, err := c.getAllBlocks(ctx, notionapi.BlockID(pageID))
	if err != nil {
		return "", err
	}

	for _, block := range blocks {
		text := extractTextFromBlock(block)
		if text != "" {
			content.WriteString(text)
			content.WriteString("\n")
		}
	}

	return strings.TrimSpace(content.String()), nil
}

// getAllBlocks retrieves all blocks for a page, handling pagination
func (c *Client) getAllBlocks(ctx context.Context, blockID notionapi.BlockID) ([]notionapi.Block, error) {
	var allBlocks []notionapi.Block
	var cursor *string

	for {
		pagination := &notionapi.Pagination{
			PageSize: 100,
		}
		if cursor != nil {
			pagination.StartCursor = notionapi.Cursor(*cursor)
		}

		resp, err := c.api.Block.GetChildren(ctx, blockID, pagination)
		if err != nil {
			return nil, err
		}

		allBlocks = append(allBlocks, resp.Results...)

		if !resp.HasMore {
			break
		}
		cursorStr := string(resp.NextCursor)
		cursor = &cursorStr
	}

	return allBlocks, nil
}

// extractTextFromBlock extracts plain text from various block types
func extractTextFromBlock(block notionapi.Block) string {
	switch b := block.(type) {
	case *notionapi.ParagraphBlock:
		return extractRichText(b.Paragraph.RichText)
	case *notionapi.Heading1Block:
		return extractRichText(b.Heading1.RichText)
	case *notionapi.Heading2Block:
		return extractRichText(b.Heading2.RichText)
	case *notionapi.Heading3Block:
		return extractRichText(b.Heading3.RichText)
	case *notionapi.BulletedListItemBlock:
		return "• " + extractRichText(b.BulletedListItem.RichText)
	case *notionapi.NumberedListItemBlock:
		return extractRichText(b.NumberedListItem.RichText)
	case *notionapi.QuoteBlock:
		return "> " + extractRichText(b.Quote.RichText)
	case *notionapi.CodeBlock:
		return "```\n" + extractRichText(b.Code.RichText) + "\n```"
	default:
		return ""
	}
}

// extractRichText converts RichText array to plain text
func extractRichText(richTexts []notionapi.RichText) string {
	var result strings.Builder
	for _, rt := range richTexts {
		result.WriteString(rt.PlainText)
	}
	return result.String()
}

// contentToBlocks converts a content string to paragraph blocks
func contentToBlocks(content string) []notionapi.Block {
	if content == "" {
		return nil
	}

	paragraphs := strings.Split(content, "\n\n")
	blocks := make([]notionapi.Block, 0, len(paragraphs))

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		blocks = append(blocks, &notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeParagraph,
			},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: para,
						},
					},
				},
			},
		})
	}

	return blocks
}
