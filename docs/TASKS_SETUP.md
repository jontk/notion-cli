# Tasks Management Setup Guide

This guide will help you set up and use notion-cli for managing tasks and TODOs in your Notion database.

## What is Tasks Management?

Tasks management helps you track TODOs, action items, and tasks with priorities, due dates, categories, and completion status.

**Perfect for:**
- Personal task lists
- Work project management
- GTD (Getting Things Done) workflows
- AI assistant task creation (Claude, GPT, etc.)

## Database Setup

### 1. Create Tasks Database in Notion

1. Open Notion
2. Click "+ New page" or navigate to an existing page
3. Type `/database` and select "Database - inline" or "Database - full page"
4. Name it "Tasks" or "My Tasks"

### 2. Configure Database Properties

Add/rename the following properties to match what notion-cli expects:

| Property Name | Property Type | Options/Notes |
|---------------|---------------|---------------|
| **Title** | Title | Default name property (required) |
| **Status** | Status | Options: "Todo", "In Progress", "Done", "Blocked" |
| **Priority** | Select | Options: "High", "Medium", "Low" |
| **Due Date** | Date | Date only (no time) |
| **Category** | Select | Add categories like: "Work", "Personal", "Home", "Health" |
| **Tags** | Multi-select | Add tags like: "urgent", "review", "research" |
| **Notes** | Text | For additional task details |

**Important**: Property names are case-sensitive and must match exactly.

### 3. Get the Database ID

**Method 1: From the URL**
1. Open your Tasks database in Notion
2. Look at the URL: `https://www.notion.so/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX?v=...`
3. Copy the 32-character ID (the X's)

**Method 2: Using notion-cli**
```bash
notion-cli databases list
```
Find your Tasks database and copy its ID.

### 4. Configure notion-cli

Add the Tasks database ID to your config:

**Edit `~/.notion-cli.yaml`:**
```yaml
api_token: "ntn_YOUR_TOKEN_HERE"
tasks_database_id: "YOUR_TASKS_DATABASE_ID"
default_task_status: "Todo"
default_priority: "Medium"
```

**Or use environment variable:**
```bash
export NOTION_TASKS_DATABASE_ID="YOUR_TASKS_DATABASE_ID"
```

### 5. Share Database with Integration

1. Open your Tasks database in Notion
2. Click the `•••` (three dots) menu in the top-right
3. Scroll down and click "Connections" or "Add connections"
4. Select your notion-cli integration

### 6. Test It Out

```bash
# Create a test task
notion-cli tasks create --title "Test Task" --priority "Medium"

# Query tasks
notion-cli tasks query

# Check today's tasks
notion-cli tasks today
```

If you see JSON output with your task, you're all set!

## Using Tasks Commands

### Create a Task

**Simple task:**
```bash
notion-cli tasks create --title "Buy milk"
```

**With priority and due date:**
```bash
notion-cli tasks create \
  --title "Review PR #123" \
  --priority "High" \
  --due "2024-03-20" \
  --category "Work"
```

**With all fields:**
```bash
notion-cli tasks create \
  --title "Prepare quarterly presentation" \
  --priority "High" \
  --due "2024-03-25" \
  --category "Work" \
  --tags "urgent,presentation,q1" \
  --notes "Include revenue metrics and team updates"
```

**From stdin (AI workflow):**
```bash
echo '{"title":"Call dentist","category":"Personal","priority":"High"}' | \
  notion-cli tasks create --stdin
```

### Query Tasks

**All tasks:**
```bash
notion-cli tasks query
```

**Filter by status:**
```bash
notion-cli tasks query --status "Todo"
notion-cli tasks query --status "In Progress"
notion-cli tasks query --status "Done"
```

**Filter by priority:**
```bash
notion-cli tasks query --priority "High"
notion-cli tasks query --priority "Medium"
```

**Filter by category:**
```bash
notion-cli tasks query --category "Work"
notion-cli tasks query --category "Personal"
```

**Combine filters:**
```bash
notion-cli tasks query --status "Todo" --priority "High" --category "Work"
```

**Limit results:**
```bash
notion-cli tasks query --limit 10
```

### Get a Specific Task

```bash
notion-cli tasks get --id "TASK_ID"
```

### Update a Task

**Update status:**
```bash
notion-cli tasks update --id "TASK_ID" --status "In Progress"
```

**Update priority:**
```bash
notion-cli tasks update --id "TASK_ID" --priority "High"
```

**Update multiple fields:**
```bash
notion-cli tasks update \
  --id "TASK_ID" \
  --status "In Progress" \
  --priority "High" \
  --due "2024-03-22" \
  --notes "Added more context after team meeting"
```

**Update from stdin:**
```bash
echo '{"status":"Done","notes":"Completed ahead of schedule"}' | \
  notion-cli tasks update --id "TASK_ID" --stdin
```

### Mark as Complete

```bash
notion-cli tasks complete --id "TASK_ID"
```

This is a shortcut for updating status to "Done".

### View Today's Tasks

```bash
notion-cli tasks today
```

Shows all tasks with status "Todo" that are due today or earlier.

### View Overdue Tasks

```bash
notion-cli tasks overdue
```

Shows all incomplete tasks that are past their due date.

## Common Workflows

### Workflow 1: Daily Task Management

**Morning:**
```bash
# Check what's due today
notion-cli tasks today

# Check what's overdue
notion-cli tasks overdue

# Review high priority items
notion-cli tasks query --priority "High" --status "Todo"
```

**During Day:**
```bash
# Start working on a task
notion-cli tasks update --id "TASK_ID" --status "In Progress"

# Mark task complete
notion-cli tasks complete --id "TASK_ID"

# Add new task as you think of it
notion-cli tasks create --title "Follow up with client" --priority "High"
```

**Evening:**
```bash
# Review what you completed
notion-cli tasks query --status "Done" --limit 20

# Plan for tomorrow
notion-cli tasks query --status "Todo" --priority "High"
```

### Workflow 2: GTD (Getting Things Done)

**Capture:**
```bash
# Quick capture of tasks
notion-cli tasks create --title "Research new framework"
notion-cli tasks create --title "Update documentation"
notion-cli tasks create --title "Schedule dentist appointment"
```

**Organize:**
```bash
# Add context and priorities
notion-cli tasks update --id "TASK_1" --category "Work" --priority "Medium" --tags "research"
notion-cli tasks update --id "TASK_2" --category "Work" --priority "Low" --tags "docs"
notion-cli tasks update --id "TASK_3" --category "Personal" --priority "High" --due "2024-03-20"
```

**Review:**
```bash
# Weekly review - what's pending?
notion-cli tasks query --status "Todo"

# What contexts do I have?
notion-cli tasks query --category "Work"
notion-cli tasks query --category "Personal"
```

### Workflow 3: Project-Based Task Management

**Create project tasks:**
```bash
notion-cli tasks create \
  --title "Design API endpoints" \
  --category "Work" \
  --tags "project-x,design" \
  --priority "High" \
  --due "2024-03-18"

notion-cli tasks create \
  --title "Implement authentication" \
  --category "Work" \
  --tags "project-x,dev" \
  --priority "High" \
  --due "2024-03-22"

notion-cli tasks create \
  --title "Write integration tests" \
  --category "Work" \
  --tags "project-x,testing" \
  --priority "Medium" \
  --due "2024-03-25"
```

**Track project progress:**
```bash
# See all project tasks (assuming tagged "project-x")
notion-cli tasks query --category "Work" | jq '.[] | select(.tags | contains(["project-x"]))'

# See what's done
notion-cli tasks query --status "Done" | jq '.[] | select(.tags | contains(["project-x"]))'
```

### Workflow 4: Claude AI Assistant

**Quick capture:**
```
User: "Claude, add a task to call the dentist tomorrow with high priority"
Claude: notion-cli tasks create --title "Call dentist" --due "2024-03-18" --priority "High"
```

**Smart suggestions:**
```
User: "Claude, what should I focus on today?"
Claude: [Runs notion-cli tasks today and notion-cli tasks overdue]
        "Here's what needs your attention:
        1. [High Priority] Review PR #123 (Work)
        2. [Overdue] Update documentation (Work)
        3. [Due Today] Call dentist (Personal)"
```

**Project breakdown:**
```
User: "Claude, I need to migrate the database. Break it down into tasks."
Claude: [Creates multiple tasks]
        - "Plan database migration strategy" (High, Work, +3 days)
        - "Backup current database" (High, Work, +4 days)
        - "Test migration on staging" (High, Work, +5 days)
        - "Execute migration" (High, Work, +7 days)
        - "Verify data integrity" (High, Work, +7 days)
```

### Workflow 5: Team Coordination

**Assign work items:**
```bash
# Create tasks with assignee in notes
notion-cli tasks create \
  --title "Review security audit" \
  --category "Work" \
  --priority "High" \
  --notes "Assigned to: John"
```

**Weekly standup prep:**
```bash
# What did I complete?
notion-cli tasks query --status "Done" --category "Work" --limit 20

# What am I working on?
notion-cli tasks query --status "In Progress" --category "Work"

# What's blocked?
notion-cli tasks query --status "Blocked"
```

## Tips and Best Practices

### 1. Use Categories Consistently

Organize tasks by life area:
- **Work**: Professional tasks
- **Personal**: Non-work tasks
- **Home**: House maintenance, chores
- **Health**: Exercise, medical appointments
- **Finance**: Bills, budgeting, taxes
- **Learning**: Courses, reading, skill development

### 2. Priority Guidelines

- **High**: Urgent and important. Do today.
- **Medium**: Important but not urgent. Schedule it. (Default)
- **Low**: Nice to have. Do when you have time.

Avoid making everything high priority - it defeats the purpose!

### 3. Tag Strategy

Use tags for:
- **Context**: @computer, @phone, @errands, @home, @office
- **Project**: project-x, migration, redesign
- **Status**: waiting-on, blocked-by, needs-review
- **Time**: quick-win (< 15 min), deep-work (needs focus)
- **Energy**: high-energy, low-energy

### 4. Due Dates

Set due dates for:
- Hard deadlines (meeting prep, bill payments)
- Self-imposed deadlines (personal goals)
- Recurring tasks (weekly review, monthly reports)

Don't set due dates for:
- Someday/maybe tasks
- Tasks without clear timelines

### 5. Status Management

- **Todo**: Not started yet
- **In Progress**: Actively working on it
- **Blocked**: Can't proceed (waiting on someone/something)
- **Done**: Completed

Review blocked tasks weekly to unblock them.

### 6. Notes Field

Use notes for:
- Context about why this task matters
- Links to related documents or tickets
- Acceptance criteria
- Dependencies
- Progress updates

## Advanced Usage

### Bash Scripts

**Daily summary:**
```bash
#!/bin/bash
# daily-summary.sh

echo "📋 Today's Tasks:"
notion-cli tasks today

echo ""
echo "⚠️  Overdue:"
notion-cli tasks overdue

echo ""
echo "🔥 High Priority:"
notion-cli tasks query --priority "High" --status "Todo" --limit 5
```

**Quick add wrapper:**
```bash
#!/bin/bash
# task-add.sh

notion-cli tasks create --title "$1" --priority "${2:-Medium}" --category "${3:-Personal}"
```

Usage:
```bash
./task-add.sh "Buy groceries" "Medium" "Personal"
```

### JSON Processing with jq

**Get task titles only:**
```bash
notion-cli tasks query | jq -r '.[] | .title'
```

**Filter by tag:**
```bash
notion-cli tasks query | jq '.[] | select(.tags | contains(["urgent"]))'
```

**Count by status:**
```bash
notion-cli tasks query | jq 'group_by(.status) | map({status: .[0].status, count: length})'
```

**This week's completed tasks:**
```bash
notion-cli tasks query --status "Done" | \
  jq --arg week "$(date -d 'last monday' +%Y-%m-%d)" \
  '.[] | select(.updated_at >= $week)'
```

### Automation with Cron

**Daily reminder (9 AM):**
```bash
# crontab -e
0 9 * * * /path/to/daily-summary.sh | mail -s "Daily Tasks" you@example.com
```

**Weekly review (Sunday 6 PM):**
```bash
0 18 * * 0 notion-cli tasks query --status "Todo" | mail -s "Week Ahead" you@example.com
```

## Troubleshooting

### "tasks database ID is required"

**Problem**: Config missing tasks_database_id

**Fix**:
```bash
# Check config
cat ~/.notion-cli.yaml

# Should have:
tasks_database_id: "your-tasks-database-id-here"
```

### "Could not find database"

**Problem**: Integration doesn't have access

**Fix**: Share database with your integration in Notion (see step 5 above)

### "Property not found" errors

**Problem**: Property names don't match

**Fix**: Verify property names in your database:
```bash
notion-cli databases schema --id "YOUR_TASKS_DATABASE_ID"

# Expected properties:
# - Title (title type)
# - Status (status type)
# - Priority (select type)
# - Due Date (date type)
# - Category (select type)
# - Tags (multi_select type)
# - Notes (rich_text type)
```

### Status property shows wrong options

**Problem**: Your database has different status options

**Note**: Notion Status properties support custom options. If you have "Blocked" as part of "Todo" state (as mentioned by user), that's fine! The code works with any status values you've configured.

## Next Steps

1. ✅ Create your Tasks database in Notion
2. ✅ Configure notion-cli with your database ID
3. ✅ Create your first task
4. ✅ Set up your daily workflow
5. 🤖 Integrate with Claude AI (see AI Assistant Integration section below)
6. 📅 Combine with Events for full productivity system: [EVENTS_SETUP.md](./EVENTS_SETUP.md)

## Related Documentation

- [Main README](../README.md) - Overview of notion-cli
- [Events Setup](./EVENTS_SETUP.md) - Calendar management
- [Posts Setup](./POSTS_SETUP.md) - Content management

---

**Quick Start Command:**
```bash
notion-cli tasks create \
  --title "Set up my productivity system" \
  --priority "High" \
  --category "Personal" \
  --notes "Configure notion-cli for task management"
```
