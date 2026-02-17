# Posts Management Setup Guide

This guide will help you set up and use notion-cli for managing content posts in your Notion database.

## What is Posts Management?

Posts management is designed for content creators who want to track blog posts, social media content, articles, or any written content through a publishing pipeline.

**Typical workflow:**
```
Idea → Draft → Ready → Published
```

## Database Setup

### 1. Create Posts Database in Notion

1. Open Notion
2. Click "+ New page" or navigate to an existing page
3. Type `/database` and select "Database - inline" or "Database - full page"
4. Name it "Content" or "Posts"

### 2. Configure Database Properties

Add/rename the following properties:

| Property Name | Property Type | Options/Notes |
|---------------|---------------|---------------|
| **Title** | Title | Default name property (required) |
| **Status** | Status | Options: "Draft", "Ready", "Published" |
| **Platforms** | Multi-select | Add platforms like: "Blog", "Twitter", "LinkedIn" |
| **Publish Date** | Date | Scheduled publish date (date only, no time) |

**Note**: The property names must match exactly (case-sensitive).

### 3. Get the Database ID

**Method 1: From the URL**
1. Open your Posts database in Notion
2. Look at the URL: `https://www.notion.so/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX?v=...`
3. Copy the 32-character ID (the X's)

**Method 2: Using notion-cli**
```bash
notion-cli databases list
```
This will show all databases - find your Posts database and copy its ID.

### 4. Configure notion-cli

Add the Posts database ID to your config:

**Edit `~/.notion-cli.yaml`:**
```yaml
api_token: "ntn_YOUR_TOKEN_HERE"
database_id: "YOUR_POSTS_DATABASE_ID"
default_status: "Draft"
```

**Or use environment variable:**
```bash
export NOTION_API_TOKEN="ntn_..."
export NOTION_DATABASE_ID="YOUR_POSTS_DATABASE_ID"
```

### 5. Share Database with Integration

1. Open your Posts database in Notion
2. Click the `•••` (three dots) menu in the top-right
3. Scroll down and click "Connections" or "Add connections"
4. Select your notion-cli integration

### 6. Verify Setup

Check the database schema:
```bash
notion-cli databases schema
```

You should see your properties listed.

## Using Posts Commands

### Create a Post

**Simple post:**
```bash
notion-cli posts create --title "My First Post" --status "Draft"
```

**With all fields:**
```bash
notion-cli posts create \
  --title "Building CLI Tools in Go" \
  --content "Learn how to build powerful CLI tools..." \
  --status "Draft" \
  --platform "Blog,Twitter" \
  --date "2024-03-25"
```

**From stdin (AI workflow):**
```bash
echo '{"title":"AI-Generated Post","content":"...","status":"Draft"}' | \
  notion-cli posts create --stdin
```

### Query Posts

**All posts:**
```bash
notion-cli posts query
```

**Filter by status:**
```bash
notion-cli posts query --status "Draft"
notion-cli posts query --status "Ready"
notion-cli posts query --status "Published"
```

**Filter by platform:**
```bash
notion-cli posts query --platform "Blog"
notion-cli posts query --platform "Twitter"
```

**Limit results:**
```bash
notion-cli posts query --status "Draft" --limit 10
```

**Sort by date:**
```bash
notion-cli posts query --sort "created_time" --order "descending"
```

### Get a Specific Post

```bash
notion-cli posts get --id "PAGE_ID"
```

This returns the full post including content blocks.

### Update a Post

**Update status:**
```bash
notion-cli posts update --id "PAGE_ID" --status "Ready"
```

**Update multiple fields:**
```bash
notion-cli posts update \
  --id "PAGE_ID" \
  --status "Published" \
  --date "2024-03-20" \
  --platform "Blog,Twitter,LinkedIn"
```

**Update from stdin:**
```bash
echo '{"status":"Published","date":"2024-03-20"}' | \
  notion-cli posts update --id "PAGE_ID" --stdin
```

### Archive a Post

```bash
notion-cli posts archive --id "PAGE_ID"
```

This marks the post as archived in Notion (doesn't delete it).

## Common Workflows

### Workflow 1: Content Publishing Pipeline

**Step 1: Create draft**
```bash
notion-cli posts create \
  --title "10 Tips for Better Productivity" \
  --status "Draft" \
  --platform "Blog"
```

**Step 2: Write content**
(Edit in Notion UI or use your favorite editor)

**Step 3: Mark as ready**
```bash
notion-cli posts update --id "POST_ID" --status "Ready"
```

**Step 4: Schedule publish date**
```bash
notion-cli posts update \
  --id "POST_ID" \
  --status "Ready" \
  --date "2024-03-25"
```

**Step 5: After publishing**
```bash
notion-cli posts update --id "POST_ID" --status "Published"
```

### Workflow 2: Multi-Platform Content

Create a post for multiple platforms:
```bash
notion-cli posts create \
  --title "Product Launch Announcement" \
  --content "We're excited to announce..." \
  --status "Ready" \
  --platform "Blog,Twitter,LinkedIn,Newsletter" \
  --date "2024-03-20"
```

Query what's ready for each platform:
```bash
notion-cli posts query --status "Ready" --platform "Twitter"
notion-cli posts query --status "Ready" --platform "Blog"
```

### Workflow 3: Claude AI Integration

**Idea Capture:**
```
User: "Claude, I have an idea for a post about Git workflows"
Claude: notion-cli posts create --title "Git Workflows Best Practices" --status "Draft"
```

**Content Generation:**
```
User: "Claude, generate an outline for my Git workflows post"
Claude: [Generates outline]
         echo '{"content":"[outline content]"}' | \
           notion-cli posts update --id "POST_ID" --stdin
```

**Publishing Prep:**
```
User: "Claude, what posts are ready to publish this week?"
Claude: notion-cli posts query --status "Ready"
        [Shows posts and helps schedule them]
```

### Workflow 4: n8n Automation

**Example: Auto-publish scheduled posts**

n8n workflow:
1. **Schedule**: Run daily at 9 AM
2. **Execute**: `notion-cli posts query --status "Ready"`
3. **Filter**: Posts with today's publish date
4. **For Each Post**:
   - Publish to platform (Medium, Dev.to, etc.)
   - Update status: `notion-cli posts update --id "POST_ID" --status "Published"`
5. **Notify**: Send Slack message with published posts


### Workflow 5: Content Calendar Review

**Weekly review:**
```bash
# What's published this week
notion-cli posts query --status "Published" --limit 50

# What's ready to go
notion-cli posts query --status "Ready"

# What's in draft
notion-cli posts query --status "Draft"
```

**Monthly planning:**
```bash
# See all scheduled content
notion-cli posts query --status "Ready" | jq '.[] | {title, publish_date}'
```

## Tips and Best Practices

### 1. Status Workflow

Use statuses to track progress:
- **Draft**: Work in progress, not ready
- **Ready**: Edited, reviewed, ready to publish
- **Published**: Already published to platforms

### 2. Platform Tags

Use consistent platform names:
- Blog
- Twitter
- LinkedIn
- Medium
- Dev.to
- Newsletter
- YouTube

### 3. Scheduling

Set publish dates for:
- Planning your content calendar
- Scheduling posts in advance
- Tracking when content should go live

### 4. Content Organization

Add additional properties for better organization:
- **Category**: Tutorial, Opinion, News, Case Study
- **Tags**: go, productivity, tutorial, beginner
- **Target Audience**: developers, managers, beginners
- **Word Count**: Track post length
- **SEO Keywords**: For search optimization

### 5. Batch Operations

Process multiple posts efficiently:
```bash
# Get all draft posts
DRAFTS=$(notion-cli posts query --status "Draft" | jq -r '.[] | .id')

# Mark them all as ready (with careful review!)
for id in $DRAFTS; do
  notion-cli posts update --id "$id" --status "Ready"
done
```

## Customization

### Adding Custom Properties

Want to track more fields? Extend the database:

1. **In Notion**: Add new properties to your database
2. **In code**: Modify `internal/notion/pages.go` to handle new properties
3. **In commands**: Add flags to `cmd/posts/*.go`

Example: Adding "Category" field
- Add Select property "Category" in Notion
- Update `internal/models/post.go` to include `Category string`
- Update `CreatePost()` to set Category property
- Add `--category` flag to create/update commands

### Different Status Options

Your workflow might use different statuses:
- Idea → Outline → Draft → Review → Published
- Backlog → Writing → Editing → Published → Promoted

Change status options in your Notion database Status property.

## Troubleshooting

### "database ID is required"

**Problem**: Config missing database ID

**Fix**:
```bash
# Check config
cat ~/.notion-cli.yaml

# Should have:
database_id: "your-database-id-here"
```

### "Could not find database"

**Problem**: Integration doesn't have access

**Fix**: Share database with your integration in Notion (see step 5 above)

### "Property not found" errors

**Problem**: Property names don't match

**Fix**: Check your property names exactly match:
```bash
notion-cli databases schema

# Expected properties:
# - Title (title type)
# - Status (status type)
# - Platforms (multi_select type)
# - Publish Date (date type)
```

### Empty content when using get

**Problem**: Content isn't included in query

**Solution**: Content is stored as page blocks, use `posts get` to retrieve it:
```bash
notion-cli posts get --id "PAGE_ID"
```

## Integration Examples

### With Claude AI

```
"Claude, create a content calendar for next month with 8 blog post ideas about Go programming"

"Claude, what blog posts are ready to publish this week?"

"Claude, generate an outline for my post about Docker optimization"
```

### With Scripts

**Daily summary email:**
```bash
#!/bin/bash
# send-content-summary.sh

READY=$(notion-cli posts query --status "Ready" | jq -r 'length')
DRAFT=$(notion-cli posts query --status "Draft" | jq -r 'length')

echo "Content Summary:"
echo "Ready to publish: $READY"
echo "In draft: $DRAFT"
```

### With Git Hooks

**Pre-commit hook to update Notion:**
```bash
#!/bin/bash
# .git/hooks/pre-commit

# When you commit a new article, create a post in Notion
if git diff --cached --name-only | grep -q "^content/"; then
  TITLE=$(grep -m 1 "^# " path/to/article.md | sed 's/# //')
  notion-cli posts create --title "$TITLE" --status "Draft"
fi
```

## Next Steps

1. ✅ Set up your Posts database in Notion
2. ✅ Configure notion-cli with your database ID
3. ✅ Create your first post
4. ✅ Set up your content workflow (Draft → Ready → Published)
5. 📚 Explore automation opportunities with your workflow tools
6. 🤖 Set up Claude as your content assistant

## Related Documentation

- [Main README](../README.md) - Overview of notion-cli
- [Tasks Setup](./TASKS_SETUP.md) - Task management
- [Events Setup](./EVENTS_SETUP.md) - Calendar management

---

**Quick Start Command:**
```bash
notion-cli posts create \
  --title "My First Post" \
  --content "This is a test post" \
  --status "Draft" \
  --platform "Blog"
```
