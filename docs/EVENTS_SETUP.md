# Events & Calendar Management Setup Guide

This guide will help you set up and use notion-cli for managing calendar events and appointments in your Notion database.

## What is Events Management?

Events management helps you track calendar events, meetings, appointments, and scheduled activities with dates, times, locations, and attendees.

**Perfect for:**
- Personal calendar management
- Meeting scheduling
- Appointment tracking
- Event planning
- AI assistant calendar management (Claude, GPT, etc.)

## Database Setup

### 1. Create Events Database in Notion

1. Open Notion
2. Click "+ New page" or navigate to an existing page
3. Type `/database` and select "Database - inline" or "Database - full page"
4. Name it "Calendar" or "Events"

### 2. Configure Database Properties

Add/rename the following properties:

| Property Name | Property Type | Options/Notes |
|---------------|---------------|---------------|
| **Title** | Title | Default name property (required) |
| **Date** | Date | **IMPORTANT**: Enable time! Use date + time format |
| **Type** | Select | Options: "Work", "Personal", "Meeting", "Appointment" |
| **Location** | Text | Event location or meeting room |
| **Attendees** | Multi-select | Email addresses or names |
| **Status** | Multi-select | Options: "Scheduled", "Completed", "Cancelled" |
| **Notes** | Text | Additional event details |

**Critical**: The Date property MUST include time. Click on the Date property settings and enable "Include time".

### 3. Get the Database ID

**Method 1: From the URL**
1. Open your Events database in Notion
2. Look at the URL: `https://www.notion.so/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX?v=...`
3. Copy the 32-character ID (the X's)

**Method 2: Using notion-cli**
```bash
notion-cli databases list
```
Find your Events database and copy its ID.

### 4. Configure notion-cli

Add the Events database ID to your config:

**Edit `~/.notion-cli.yaml`:**
```yaml
api_token: "ntn_YOUR_TOKEN_HERE"
events_database_id: "YOUR_EVENTS_DATABASE_ID"
```

**Or use environment variable:**
```bash
export NOTION_EVENTS_DATABASE_ID="YOUR_EVENTS_DATABASE_ID"
```

### 5. Share Database with Integration

1. Open your Events database in Notion
2. Click the `•••` (three dots) menu in the top-right
3. Scroll down and click "Connections" or "Add connections"
4. Select your notion-cli integration

### 6. Test It Out

```bash
# Create a test event
notion-cli events create \
  --title "Test Meeting" \
  --date "2024-03-20 14:00" \
  --type "Meeting"

# Query events
notion-cli events query

# Check today's schedule
notion-cli events today
```

If you see JSON output with your event, you're all set!

## Using Events Commands

### Create an Event

**Simple event:**
```bash
notion-cli events create --title "Team Meeting" --date "2024-03-20 14:00"
```

**With location and type:**
```bash
notion-cli events create \
  --title "Client Presentation" \
  --date "2024-03-22 15:30" \
  --type "Meeting" \
  --location "Conference Room A"
```

**With all fields:**
```bash
notion-cli events create \
  --title "Product Review Meeting" \
  --date "2024-03-25 10:00" \
  --type "Meeting" \
  --location "Zoom: https://zoom.us/j/..." \
  --attendees "john@example.com,jane@example.com" \
  --status "Scheduled" \
  --notes "Discuss Q2 roadmap and feature priorities"
```

**From stdin (AI workflow):**
```bash
echo '{"title":"Doctor Appointment","date":"2024-03-20 09:30","type":"Personal"}' | \
  notion-cli events create --stdin
```

### Query Events

**All events:**
```bash
notion-cli events query
```

**Filter by type:**
```bash
notion-cli events query --type "Meeting"
notion-cli events query --type "Personal"
notion-cli events query --type "Work"
```

**Filter by status:**
```bash
notion-cli events query --status "Scheduled"
notion-cli events query --status "Completed"
notion-cli events query --status "Cancelled"
```

**Limit results:**
```bash
notion-cli events query --limit 20
```

### Get a Specific Event

```bash
notion-cli events get --id "EVENT_ID"
```

### Update an Event

**Update location:**
```bash
notion-cli events update --id "EVENT_ID" --location "Conference Room B"
```

**Update date/time:**
```bash
notion-cli events update --id "EVENT_ID" --date "2024-03-23 16:00"
```

**Update multiple fields:**
```bash
notion-cli events update \
  --id "EVENT_ID" \
  --date "2024-03-24 14:00" \
  --location "Building 2, Room 301" \
  --notes "Updated: Added architectural discussion to agenda"
```

**Update from stdin:**
```bash
echo '{"status":"Completed","notes":"Meeting notes: ..."}' | \
  notion-cli events update --id "EVENT_ID" --stdin
```

### Cancel an Event

```bash
notion-cli events cancel --id "EVENT_ID"
```

This updates the status to "Cancelled".

### View Today's Schedule

```bash
notion-cli events today
```

Shows all events scheduled for today.

### View This Week's Events

```bash
notion-cli events week
```

Shows all events for the current week (Monday-Sunday).

## Common Workflows

### Workflow 1: Daily Calendar Management

**Morning:**
```bash
# Check today's schedule
notion-cli events today

# Check this week
notion-cli events week

# See upcoming meetings
notion-cli events query --type "Meeting" --status "Scheduled" --limit 10
```

**During Day:**
```bash
# Quick reference for next meeting
notion-cli events today | jq -r '.[] | "\(.date): \(.title)"'

# Mark meeting as completed
notion-cli events update --id "EVENT_ID" --status "Completed"
```

**Evening:**
```bash
# Review what happened today
notion-cli events today

# Preview tomorrow
# (You'd need to filter by tomorrow's date using jq)
```

### Workflow 2: Meeting Scheduling

**Schedule a meeting:**
```bash
notion-cli events create \
  --title "Sprint Planning" \
  --date "2024-03-21 10:00" \
  --type "Meeting" \
  --location "Conference Room 1" \
  --attendees "team@company.com" \
  --notes "Plan Q2 Sprint 1"
```

**Reschedule:**
```bash
notion-cli events update \
  --id "EVENT_ID" \
  --date "2024-03-21 14:00" \
  --notes "Rescheduled due to conflict"
```

**Cancel if needed:**
```bash
notion-cli events cancel --id "EVENT_ID"
```

### Workflow 3: Appointment Tracking

**Medical appointments:**
```bash
notion-cli events create \
  --title "Dentist Appointment" \
  --date "2024-03-25 09:30" \
  --type "Personal" \
  --location "123 Main St, Suite 200" \
  --notes "Annual cleaning"
```

**Personal appointments:**
```bash
notion-cli events create \
  --title "Car Service" \
  --date "2024-03-28 14:00" \
  --type "Personal" \
  --location "Auto Shop Downtown"
```

### Workflow 4: Claude AI Assistant

**Schedule meetings:**
```
User: "Claude, schedule a team meeting tomorrow at 2pm"
Claude: notion-cli events create --title "Team Meeting" --date "2024-03-18 14:00" --type "Meeting"
```

**Check schedule:**
```
User: "Claude, what's on my calendar today?"
Claude: [Runs notion-cli events today]
        "Your schedule today:
        - 9:00 AM: Team standup (Meeting)
        - 2:00 PM: Client presentation (Meeting)
        - 4:00 PM: 1:1 with manager (Meeting)"
```

**Weekly planning:**
```
User: "Claude, what meetings do I have this week?"
Claude: [Runs notion-cli events week --type "Meeting"]
        [Presents organized view by day]
```

**Conflict resolution:**
```
User: "Claude, I need to reschedule my 3pm meeting to 4pm"
Claude: [Gets today's events, finds the 3pm meeting]
        notion-cli events update --id "EVENT_ID" --date "2024-03-17 16:00"
        "Done! Your 3pm meeting has been moved to 4pm."
```

### Workflow 5: Time Blocking

**Block focus time:**
```bash
notion-cli events create \
  --title "Deep Work: Code Review" \
  --date "2024-03-20 09:00" \
  --type "Work" \
  --notes "No interruptions, review PRs #123-#130"
```

**Block personal time:**
```bash
notion-cli events create \
  --title "Lunch Break" \
  --date "2024-03-20 12:00" \
  --type "Personal"

notion-cli events create \
  --title "Exercise" \
  --date "2024-03-20 17:30" \
  --type "Personal" \
  --location "Gym"
```

**Block preparation time:**
```bash
notion-cli events create \
  --title "Prep for Client Presentation" \
  --date "2024-03-21 13:00" \
  --type "Work" \
  --notes "Review slides, prepare demos"
```

## Tips and Best Practices

### 1. Always Include Time

Events should have specific times:
- ✅ "2024-03-20 14:00"
- ✅ "2024-03-20 09:30"
- ❌ "2024-03-20" (date only - for tasks, not events)

Use 24-hour format for clarity.

### 2. Event Types

Use consistent types:
- **Meeting**: Team meetings, client meetings, 1:1s
- **Personal**: Appointments, personal commitments
- **Work**: Focused work blocks, project time
- **Appointment**: Doctor, dentist, service appointments
- **Event**: Conferences, workshops, social events

### 3. Location Best Practices

Be specific with locations:
- Physical: "Building 2, Conference Room A"
- Virtual: "Zoom: https://zoom.us/j/123456"
- Offsite: "Coffee Shop - 123 Main St"
- Phone: "Phone: +1-555-1234"

### 4. Attendees Management

Add attendees for:
- Coordination (who's expected)
- Follow-up (who to loop in)
- Accountability (who attended)

Use consistent formats:
- Email: "john@example.com"
- Name: "John Doe"
- Team: "engineering-team"

### 5. Status Management

- **Scheduled**: Confirmed, will happen
- **Completed**: Already happened, keep for history
- **Cancelled**: Won't happen, keep for records

Review completed events monthly and archive if needed.

### 6. Notes Field

Use notes for:
- Meeting agenda
- Preparation needed
- Post-meeting summary
- Action items
- Links to related documents

## Advanced Usage

### Bash Scripts

**Daily agenda email:**
```bash
#!/bin/bash
# daily-agenda.sh

echo "📅 Today's Schedule:"
notion-cli events today | jq -r '.[] | "\(.date): \(.title) (\(.type))"'

echo ""
echo "📋 This Week:"
notion-cli events week | jq -r 'length' | xargs echo "Total events:"
```

**Check for conflicts:**
```bash
#!/bin/bash
# check-conflicts.sh

# Get today's events sorted by time
notion-cli events today | \
  jq -r 'sort_by(.date) | .[] | "\(.date): \(.title)"'
```

### JSON Processing with jq

**Events by type:**
```bash
notion-cli events week | jq 'group_by(.type) | map({type: .[0].type, count: length})'
```

**Only meeting times:**
```bash
notion-cli events today | jq -r '.[] | select(.type == "Meeting") | "\(.date): \(.title)"'
```

**Events with specific attendee:**
```bash
notion-cli events query | jq '.[] | select(.attendees | contains(["john@example.com"]))'
```

**This week's meeting count:**
```bash
notion-cli events week | jq '[.[] | select(.type == "Meeting")] | length'
```

### Integration with Other Tools

**Export to calendar format:**
```bash
#!/bin/bash
# export-calendar.sh

notion-cli events week | jq -r '.[] |
  "BEGIN:VEVENT
  SUMMARY:\(.title)
  DTSTART:\(.date | gsub("[^0-9T]"; ""))
  LOCATION:\(.location // "")
  DESCRIPTION:\(.notes // "")
  END:VEVENT"'
```

**Sync with Google Calendar:**
(Requires gcalcli or similar tool)
```bash
#!/bin/bash
notion-cli events today | jq -r '.[] |
  [.title, .date, .location] | @tsv' | \
while IFS=$'\t' read title date location; do
  gcalcli add --title "$title" --when "$date" --where "$location"
done
```

## Date Format Reference

**Correct formats:**
```bash
# Date with time (use this for events!)
--date "2024-03-20 14:00"
--date "2024-03-20 09:30"
--date "2024-12-31 23:59"

# ISO format also works
--date "2024-03-20T14:00:00Z"
```

**Incorrect formats:**
```bash
# Missing time
--date "2024-03-20"  # Don't use for events

# Wrong format
--date "03/20/2024 2pm"  # Won't parse
--date "March 20, 2pm"   # Won't parse
```

## Troubleshooting

### "events database ID is required"

**Problem**: Config missing events_database_id

**Fix**:
```bash
# Check config
cat ~/.notion-cli.yaml

# Should have:
events_database_id: "your-events-database-id-here"
```

### "Could not find database"

**Problem**: Integration doesn't have access

**Fix**: Share database with your integration in Notion (see step 5 above)

### "Property not found" errors

**Problem**: Property names don't match

**Fix**: Verify property names:
```bash
notion-cli databases schema --id "YOUR_EVENTS_DATABASE_ID"

# Expected properties:
# - Title (title type)
# - Date (date type with time enabled)
# - Type (select type)
# - Location (rich_text type)
# - Attendees (multi_select type)
# - Status (multi_select type)
# - Notes (rich_text type)
```

### Events showing wrong time

**Problem**: Time zone differences

**Solution**: Times are stored in UTC but displayed in your local timezone. Ensure your system timezone is correct.

### Date property not showing time

**Problem**: Date property doesn't have time enabled

**Fix**: In Notion, click on the Date property settings and enable "Include time".

## Next Steps

1. ✅ Create your Events database in Notion
2. ✅ Configure notion-cli with your database ID
3. ✅ Create your first event
4. ✅ Set up your daily calendar workflow
5. 🤖 Integrate with Claude AI (see AI Assistant Integration section below)
6. ✅ Combine with Tasks for full productivity system: [TASKS_SETUP.md](./TASKS_SETUP.md)

## Related Documentation

- [Main README](../README.md) - Overview of notion-cli
- [Tasks Setup](./TASKS_SETUP.md) - Task management
- [Posts Setup](./POSTS_SETUP.md) - Content management

---

**Quick Start Command:**
```bash
notion-cli events create \
  --title "Welcome Meeting" \
  --date "$(date -d '+1 day' '+%Y-%m-%d') 10:00" \
  --type "Meeting" \
  --notes "Get started with notion-cli events"
```
