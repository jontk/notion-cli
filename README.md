# notion-cli

A command-line tool for managing content, tasks, and calendar events in Notion databases.

> **Note:** This was built for my personal content publishing workflow and productivity system. It's opinionated and designed around specific database schemas. However, the patterns and architecture are reusable if you're building something similar.

## What This Does

Programmatically manage your Notion workspace from the command line:

**Content Management:**
- ✅ Create posts from CLI or stdin (for AI workflows)
- ✅ Query posts by status, platform, or date
- ✅ Update post properties
- ✅ Archive old content

**Task Management:**
- ✅ Create and manage TODOs and tasks
- ✅ Query tasks by status, priority, or category
- ✅ View today's tasks and overdue items
- ✅ Mark tasks as complete

**Calendar/Events:**
- ✅ Create and manage calendar events
- ✅ View today's schedule and weekly events
- ✅ Update event details and status
- ✅ Cancel events

**Database Operations:**
- ✅ List databases and inspect schemas
- ✅ Support for multiple databases

## Why I Built This

I wanted to manage my content publishing pipeline and personal productivity system programmatically:

```
Content Idea → Notion (Draft)              Task Created → Notion (Todo)
  ↓                                          ↓
Write Article                               Work on Task
  ↓                                          ↓
Notion (Ready) → n8n automation            Notion (Done)
  ↓
Publish to platforms → Notion (Published)
```

Notion has a great UI, but I needed:
- **CLI access** for scripting and automation
- **AI integration** (Claude can manage tasks, calendar, and content)
- **n8n workflows** that query and update posts/tasks/events
- **Personal assistant** - "Claude, what's on my calendar today?"
- **Version control** for my content strategy

## Installation

### From Source

```bash
git clone https://github.com/jontk/notion-cli
cd notion-cli
go build -o notion-cli .
```

### Quick Start

```bash
# Configure
./notion-cli config init

# Create a post
./notion-cli posts create \
  --title "My Post" \
  --status "Draft" \
  --platform "Blog"

# Query posts
./notion-cli posts query --status "Draft"

# Update status
./notion-cli posts update --id "POST_ID" --status "Ready"
```

## My Database Schema

This CLI expects a Notion database with these properties:

| Property | Type | Description |
|----------|------|-------------|
| Name | Title | Post title |
| Status | Status | Draft, Ready, Published |
| Platforms | Multi-select | Blog, Twitter, LinkedIn, etc. |
| Publish Date | Date | Scheduled publish date |

**You can adapt this to your schema** by modifying `internal/notion/pages.go`.

## Commands

### Posts

```bash
# Create a post
notion-cli posts create --title "Title" --status "Draft"

# Create from stdin (AI workflow)
echo '{"title":"AI Post","content":"..."}' | notion-cli posts create --stdin

# Query posts
notion-cli posts query --status "Ready"
notion-cli posts query --platform "Twitter" --limit 10

# Get a specific post
notion-cli posts get --id "PAGE_ID"

# Update a post
notion-cli posts update --id "PAGE_ID" --status "Published"

# Archive a post
notion-cli posts archive --id "PAGE_ID"
```

### Databases

```bash
# List all databases
notion-cli databases list

# Get database schema
notion-cli databases schema
notion-cli databases schema --id "DATABASE_ID"
```

### Tasks

```bash
# Create a task
notion-cli tasks create --title "Buy milk" --priority "High"

# Create from stdin (AI workflow)
echo '{"title":"Call dentist","category":"Personal"}' | notion-cli tasks create --stdin

# Query tasks
notion-cli tasks query --status "Todo"
notion-cli tasks query --priority "High" --category "Work"

# Get a specific task
notion-cli tasks get --id "TASK_ID"

# Update a task
notion-cli tasks update --id "TASK_ID" --status "In Progress"

# Mark as complete
notion-cli tasks complete --id "TASK_ID"

# View today's tasks
notion-cli tasks today

# View overdue tasks
notion-cli tasks overdue
```

### Events

```bash
# Create an event
notion-cli events create --title "Team Meeting" --date "2024-03-20 14:00"

# Create from stdin
echo '{"title":"Doctor Appointment","date":"2024-03-25 09:30"}' | notion-cli events create --stdin

# Query events
notion-cli events query --type "Work"
notion-cli events query --status "Scheduled"

# Get a specific event
notion-cli events get --id "EVENT_ID"

# Update an event
notion-cli events update --id "EVENT_ID" --location "Conference Room B"

# Cancel an event
notion-cli events cancel --id "EVENT_ID"

# View today's schedule
notion-cli events today

# View this week's events
notion-cli events week
```

### Config

```bash
# Initialize config
notion-cli config init

# Show current config
notion-cli config show
```

## Configuration

Config is stored in `~/.notion-cli.yaml`:

```yaml
api_token: "secret_YOUR_TOKEN_HERE"
database_id: "your-posts-database-id"
tasks_database_id: "your-tasks-database-id"
events_database_id: "your-events-database-id"
default_status: "Draft"
default_task_status: "Todo"
default_priority: "Medium"
```

Or use environment variables:
```bash
export NOTION_API_TOKEN="secret_..."
export NOTION_DATABASE_ID="..."
export NOTION_TASKS_DATABASE_ID="..."
export NOTION_EVENTS_DATABASE_ID="..."
```

**Setup Guides:**
- [Posts Setup](docs/POSTS_SETUP.md) - Content management
- [Tasks Setup](docs/TASKS_SETUP.md) - TODO tracking
- [Events Setup](docs/EVENTS_SETUP.md) - Calendar management

## Integration Examples

### n8n Workflow

```javascript
// Execute Command node
const posts = await $exec('notion-cli posts query --status "Ready"');
const postsArray = JSON.parse(posts);

for (const post of postsArray) {
  // Publish to platform...
  // Update status
  await $exec(`notion-cli posts update --id "${post.id}" --status "Published"`);
}
```

### Claude Code / AI Assistant

Claude can manage your entire workflow:

```bash
# Content Management
# User: "Create a post idea about Go CLI tools"
# Claude runs:
notion-cli posts create \
  --title "Building CLI Tools in Go" \
  --content "Brief outline..." \
  --status "Draft"

# Task Management
# User: "Add a task to review the PR tomorrow"
# Claude runs:
notion-cli tasks create \
  --title "Review PR #123" \
  --due "2024-03-20" \
  --priority "High" \
  --category "Work"

# Calendar Management
# User: "Schedule a meeting with John on Friday at 2pm"
# Claude runs:
notion-cli events create \
  --title "Meeting with John" \
  --date "2024-03-22 14:00" \
  --type "Work"

# User: "What's on my schedule today?"
# Claude runs:
notion-cli events today
notion-cli tasks today
```

### Bash Script

```bash
#!/bin/bash
# Publish today's scheduled posts

TODAY=$(date +%Y-%m-%d)
POSTS=$(notion-cli posts query --status "Ready" | \
  jq --arg today "$TODAY" '.[] | select(.publish_date | startswith($today)) | .id')

for id in $POSTS; do
  echo "Publishing $id..."
  # Publish to platforms...
  notion-cli posts update --id "$id" --status "Published"
done
```

## Architecture

```
notion-cli/
├── cmd/                    # Cobra commands
│   ├── posts/             # Post CRUD commands
│   ├── tasks/             # Task management commands
│   ├── events/            # Calendar/event commands
│   ├── databases/         # Database inspection
│   └── config/            # Configuration
├── internal/
│   ├── config/            # Config loading
│   ├── models/            # Domain models (Post, Task, Event)
│   ├── notion/            # Notion API wrapper
│   └── output/            # JSON/table formatting
└── main.go
```

Built with:
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration
- [jomei/notionapi](https://github.com/jomei/notionapi) - Notion API client

## AI-Friendly Design

This CLI was designed to be easy for AI assistants (like Claude) to use:

✅ **Examples in help text** - Shows common patterns
✅ **Descriptive errors** - AI can self-correct
✅ **JSON output** - Easy to parse
✅ **stdin support** - Pipe AI-generated content
✅ **Consistent structure** - Predictable behavior

## Extending for Your Workflow

### Different Database Schema

1. Update `internal/models/post.go` with your fields
2. Modify `internal/notion/pages.go` to map properties
3. Update command flags in `cmd/posts/*.go`

### New Content Types

Want to manage "projects" or "notes" instead of posts?

1. Copy `cmd/posts/` to `cmd/projects/`
2. Create `internal/models/project.go`
3. Add commands to root in `main.go`

## Documentation

### Setup Guides
- [Posts Setup](docs/POSTS_SETUP.md) - Content publishing workflow
- [Tasks Setup](docs/TASKS_SETUP.md) - Task and TODO management
- [Events Setup](docs/EVENTS_SETUP.md) - Calendar and event management

**Note**: For AI assistant integration with Claude, see the "Workflow 4: Claude AI Assistant" sections in each setup guide.

## Use Cases

This tool works well for:

✅ **Content creators** managing a publishing calendar
✅ **Personal productivity** - tasks, TODOs, and calendar management
✅ **AI-assisted workflows** (Claude as your personal assistant)
✅ **Automation** (n8n, Zapier, custom scripts)
✅ **Multi-platform publishing** (blog, social media)
✅ **Team workflows** with status tracking
✅ **CLI-first workflows** - manage everything from the terminal

It's probably **not** a good fit if:
- ❌ You need a generic Notion client (use the API directly)
- ❌ You want GUI-first workflows (use Notion's UI)
- ❌ Your schema is completely different (significant customization needed)

## My Specific Workflow

I use this as part of an integrated content pipeline:

1. **Content repo** (`../content/`) - Markdown articles and strategy
2. **n8n workflows** (`../n8n/`) - Automation
3. **notion-cli** - The bridge to Notion
4. **Notion** - Content calendar and tracking

See my blog post: [Building a Content Publishing Pipeline](https://jontk.com/blog/Building-a-Content-Publishing-Pipeline-with-Notion-n8n-and-Go)

## Contributing

This is a personal tool open-sourced for others to learn from and adapt.

**If you're building something similar:**
- Fork it and customize for your schema
- The architecture is intentionally simple and hackable
- PRs welcome for bug fixes and improvements

**Not accepting:**
- Generic Notion client features (out of scope)
- Complex configuration systems (keep it simple)
- Breaking changes to the core workflow

## Why Open Source This?

1. **It might help someone** building similar workflows
2. **The patterns are reusable** (CLI + Notion API + automation)
3. **Example of AI-friendly CLI design**
4. **Educational** - Shows real-world Go CLI structure

## License

MIT

## Author

**Jon Thor Kristinsson**

- Blog: [jontk.is](https://jontk.com) (coming soon)
- GitHub: [@jontk](https://github.com/jontk)
- Projects: [s9s](https://github.com/jontk/s9s), [slurm-client](https://github.com/jontk/slurm-client)


