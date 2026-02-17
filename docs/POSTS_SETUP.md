# Posts Management Setup Guide

This guide will help you set up and use notion-cli for managing content posts through a publishing pipeline in your Notion database.

## What is Posts Management?

Posts management tracks content through a structured publishing pipeline — from the initial idea through writing, review, publishing, and distribution across platforms.

**The pipeline:**
```
Idea → Outline → Draft → Review → Published → Distributed
```

## Database Setup

### 1. Create Posts Database in Notion

1. Open Notion
2. Click "+ New page" or navigate to an existing page
3. Type `/database` and select "Database - inline" or "Database - full page"
4. Name it "Content" or "Posts"

### 2. Configure Database Properties

Add the following properties:

| Property Name | Property Type | Notes |
|---------------|---------------|-------|
| **Title** | Title | Default name property (required) |
| **Status** | Status | See status options below |
| **Pillar** | Select | Your content categories |
| **Week** | Number | Week number in your content calendar |
| **Due Date** | Date | Target publish date |
| **Published Date** | Date | Actual publish date |
| **Blog URL** | URL | URL of the published post |
| **Distributed To** | Multi-select | Platforms: LinkedIn, Twitter, Dev.to, Hacker News, Reddit |
| **Distributed Date** | Date | Date distributed to platforms |
| **LinkedIn Draft** | Text | LinkedIn post copy |
| **Twitter Thread** | Text | Twitter thread copy |
| **HN Title** | Text | Hacker News submission title |
| **Reddit Title** | Text | Reddit submission title |
| **Hashtags** | Multi-select | Post hashtags |

**Status options** (add these to your Status property):
- `Idea` — Initial concept, not yet outlined
- `Outline` — Structure planned, not yet written
- `Draft` — Being written
- `Review` — Written, under review
- `Published` — Live on your blog
- `Distributed` — Shared across platforms

**Pillar options** (customize to your content focus):
- `SLURM & HPC`
- `Go Tools`
- `Infrastructure`
- `Career & AI`

**Property names are case-sensitive** and must match exactly.

### 3. Get the Database ID

1. Open your Posts database in Notion
2. Look at the URL: `https://www.notion.so/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX?v=...`
3. Copy the 32-character ID

Or use notion-cli to list all databases:
```bash
notion-cli databases list
```

### 4. Configure notion-cli

Edit `~/.notion-cli.yaml`:
```yaml
api_token: "ntn_YOUR_TOKEN_HERE"
database_id: "YOUR_POSTS_DATABASE_ID"
default_status: "Idea"
```

Or use environment variables:
```bash
export NOTION_API_TOKEN="ntn_..."
export NOTION_DATABASE_ID="YOUR_POSTS_DATABASE_ID"
```

### 5. Share Database with Integration

1. Open your Posts database in Notion
2. Click the `•••` menu in the top-right
3. Click "Connections" → "Add connections"
4. Select your notion-cli integration

### 6. Verify Setup

```bash
notion-cli databases schema
```

You should see your properties listed.

## Using Posts Commands

### Create a Post

**Quick capture of an idea:**
```bash
notion-cli posts create --title "Why SLURM scheduling matters"
```

**With full details:**
```bash
notion-cli posts create \
  --title "Introducing s9s: a terminal UI for SLURM" \
  --status "Outline" \
  --pillar "SLURM & HPC" \
  --week 1 \
  --due-date "2026-02-24"
```

**From stdin (AI workflow):**
```bash
echo '{"title":"AI Post","status":"Draft","pillar":"Go Tools","week":2}' | \
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
notion-cli posts query --status "Review"
notion-cli posts query --status "Published"
```

**Filter by content pillar:**
```bash
notion-cli posts query --pillar "Go Tools"
notion-cli posts query --pillar "SLURM & HPC"
```

**Filter by distribution platform:**
```bash
notion-cli posts query --distributed-to "LinkedIn"
```

**Combine filters:**
```bash
notion-cli posts query --status "Published" --pillar "Infrastructure"
```

**Sort results:**
```bash
notion-cli posts query --sort "last_edited_time" --order "descending"
notion-cli posts query --sort "created_time" --order "ascending"
```

### Get a Specific Post

```bash
notion-cli posts get --id "PAGE_ID"
```

### Update a Post

**Advance status:**
```bash
notion-cli posts update --id "PAGE_ID" --status "Draft"
```

**Mark as published:**
```bash
notion-cli posts update --id "PAGE_ID" \
  --status "Published" \
  --published-date "2026-02-24" \
  --blog-url "https://jontk.com/blog/introducing-s9s"
```

**Mark as distributed:**
```bash
notion-cli posts update --id "PAGE_ID" \
  --status "Distributed" \
  --distributed-to "LinkedIn,Twitter,Dev.to" \
  --distributed-date "2026-02-25"
```

**Add platform-specific copy:**
```bash
notion-cli posts update --id "PAGE_ID" \
  --linkedin-draft "Excited to share my new post about..." \
  --twitter-thread "Thread: Why every HPC admin needs s9s 🧵" \
  --hn-title "Show HN: s9s – a terminal UI for SLURM job management"
```

**Update from stdin:**
```bash
echo '{"status":"Review"}' | notion-cli posts update --id "PAGE_ID" --stdin
```

### Archive a Post

```bash
notion-cli posts archive --id "PAGE_ID"
```

## The Publishing Pipeline

### Stage 1: Idea

Capture ideas quickly before they're lost:
```bash
notion-cli posts create --title "Why Go is perfect for CLI tools"
# Status defaults to "Idea"
```

### Stage 2: Outline

When you're ready to plan the post:
```bash
notion-cli posts update --id "POST_ID" \
  --status "Outline" \
  --pillar "Go Tools" \
  --week 3 \
  --due-date "2026-03-10"
```

### Stage 3: Draft

Once writing begins:
```bash
notion-cli posts update --id "POST_ID" --status "Draft"
```

### Stage 4: Review

When writing is complete, ready for editing:
```bash
notion-cli posts update --id "POST_ID" --status "Review"
```

### Stage 5: Published

After it goes live:
```bash
notion-cli posts update --id "POST_ID" \
  --status "Published" \
  --published-date "2026-03-10" \
  --blog-url "https://jontk.com/blog/go-cli-tools"
```

### Stage 6: Distributed

After sharing across platforms:
```bash
notion-cli posts update --id "POST_ID" \
  --status "Distributed" \
  --distributed-to "LinkedIn,Twitter,Dev.to,Hacker News" \
  --distributed-date "2026-03-11"
```

## Common Workflows

### Workflow 1: Weekly Content Planning

**Check what's in flight:**
```bash
notion-cli posts query --status "Draft"
notion-cli posts query --status "Review"
```

**See this week's content:**
```bash
# Posts assigned to week 5
notion-cli posts query | jq '.[] | select(.week == 5)'
```

**Review backlog of ideas:**
```bash
notion-cli posts query --status "Idea"
```

### Workflow 2: Pre-Publish Checklist

Before publishing, prepare distribution copy:
```bash
notion-cli posts update --id "POST_ID" \
  --linkedin-draft "3 things I learned building a SLURM UI..." \
  --twitter-thread "1/ Building a TUI in Go taught me a lot about..." \
  --hn-title "Show HN: s9s – Terminal UI for SLURM clusters" \
  --reddit-title "I built a terminal UI for SLURM job management" \
  --hashtags "golang,hpc,slurm,opensource"
```

### Workflow 3: Post-Distribution Tracking

Track where your content has been shared:
```bash
# Who's seen your Go content on LinkedIn?
notion-cli posts query --pillar "Go Tools" --distributed-to "LinkedIn"

# What's been on Hacker News?
notion-cli posts query --distributed-to "Hacker News"
```

### Workflow 4: Claude AI Integration

**Capture ideas from conversation:**
```
User: "Claude, save a post idea about using Go generics for CLI tools"
Claude: notion-cli posts create \
          --title "Using Go generics to build type-safe CLI tools" \
          --status "Idea" \
          --pillar "Go Tools"
```

**Review pipeline status:**
```
User: "Claude, what posts are ready for review?"
Claude: notion-cli posts query --status "Review"
```

**Prepare distribution copy:**
```
User: "Claude, write LinkedIn and Twitter copy for my post about s9s"
Claude: [Generates copy, then:]
        notion-cli posts update --id "POST_ID" \
          --linkedin-draft "..." \
          --twitter-thread "..."
```

## Tips

### 1. Content Pillars

Use pillars to maintain focus and balance across your content:
- Plan a mix of pillars each week
- Query by pillar to see if one area is getting neglected
- Ensures your audience gets consistent value in each topic area

### 2. Week Numbers

Using week numbers lets you plan a content calendar:
```bash
# What's scheduled for week 8?
notion-cli posts query | jq '.[] | select(.week == 8)'

# What weeks have content planned?
notion-cli posts query | jq '[.[] | .week] | sort | unique'
```

### 3. Platform-Specific Copy

Write platform copy before publishing — it forces you to distill your message and you can share immediately after going live without scrambling.

### 4. Batch Processing

Move multiple posts through a stage:
```bash
# Get all posts in review
REVIEW=$(notion-cli posts query --status "Review" | jq -r '.[].id')

# Mark them published (after actually publishing each one)
for id in $REVIEW; do
  notion-cli posts update --id "$id" --status "Published"
done
```

## Troubleshooting

### "database ID is required"

Check `~/.notion-cli.yaml` has `database_id` set.

### "Could not find database"

Share the database with your integration (see step 5 in setup).

### "Property not found"

Verify property names match exactly — they're case-sensitive:
```bash
notion-cli databases schema
```

## Related Documentation

- [Main README](../README.md) - Overview of notion-cli
- [Tasks Setup](./TASKS_SETUP.md) - Task management
- [Events Setup](./EVENTS_SETUP.md) - Calendar management

---

**Quick Start:**
```bash
notion-cli posts create \
  --title "My First Post" \
  --status "Idea" \
  --pillar "Go Tools"
```
