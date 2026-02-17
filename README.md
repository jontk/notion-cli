# notion-cli

A command-line tool for managing content in Notion databases.

> **Note:** This was built for my personal content publishing workflow. It's opinionated and designed around a specific database schema. However, the patterns and architecture are reusable if you're building something similar.

## What This Does

Programmatically manage posts in a Notion database:
- ✅ Create posts from CLI or stdin (for AI workflows)
- ✅ Query posts by status, platform, or date
- ✅ Update post properties
- ✅ Archive old content
- ✅ List databases and inspect schemas

## Why I Built This

I wanted to manage my content publishing pipeline programmatically:

```
Content Idea → Notion (Draft)
  ↓
Write Article
  ↓
Notion (Ready) → n8n automation
  ↓
Publish to platforms → Notion (Published)
```

Notion has a great UI, but I needed:
- **CLI access** for scripting and automation
- **AI integration** (Claude can manage the entire pipeline)
- **n8n workflows** that query and update posts
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
database_id: "your-database-id"
default_status: "Draft"
```

Or use environment variables:
```bash
export NOTION_API_TOKEN="secret_..."
export NOTION_DATABASE_ID="..."
```

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

### Claude Code

Claude can manage your content pipeline:

```bash
# User: "Create a post idea about Go CLI tools"
# Claude runs:
notion-cli posts create \
  --title "Building CLI Tools in Go" \
  --content "Brief outline..." \
  --status "Draft"
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
│   ├── databases/         # Database inspection
│   └── config/            # Configuration
├── internal/
│   ├── config/            # Config loading
│   ├── models/            # Domain models
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

See [docs/CLI_UX_ANALYSIS.md](docs/CLI_UX_ANALYSIS.md) for the full UX analysis.

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

See [docs/DOCS_DATABASE_SETUP.md](docs/DOCS_DATABASE_SETUP.md) for examples.

## Documentation

- [N8N Integration Guide](docs/N8N_INTEGRATION.md) - Automation workflows
- [Integration Setup](docs/INTEGRATION_SETUP.md) - Content pipeline setup
- [CLI UX Analysis](docs/CLI_UX_ANALYSIS.md) - Design decisions
- [Implementation Checklist](docs/IMPLEMENTATION_CHECKLIST.md) - Build notes

## Use Cases

This tool works well for:

✅ **Content creators** managing a publishing calendar
✅ **AI-assisted workflows** (Claude, GPT, etc.)
✅ **Automation** (n8n, Zapier, custom scripts)
✅ **Multi-platform publishing** (blog, social media)
✅ **Team workflows** with status tracking

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


