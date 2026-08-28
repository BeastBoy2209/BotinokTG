# BotinokTG

Telegram bot for group chats built with Go. Provides utilities for video downloading, expense splitting, and random member selection.

## Features

- **Video Downloader**: Automatically detects and downloads videos from YouTube, TikTok, Instagram, Twitter, and Pinterest when a link is sent to the chat.
- **Expense Manager**: Allows groups to track shared expenses and calculates the simplest way to settle debts.
- **Randomizer**: Selects a random member from the chat or from a specified list of mentions.
- **Mass Mention**: Mentions all active chat members using `@all` or `@все` triggers.
- **Auto-Registration**: Automatically registers users and chats in the database upon their first message.

## Requirements

- Go 1.21+
- PostgreSQL
- `yt-dlp` binary (placed in `./bin/` or specified in `.env`)
- `ffmpeg` installed on the system (required by yt-dlp)

## Setup

1. Clone the repository
2. Create a `.env` file based on `.env.example`
3. Run database migrations:
   ```bash
   make migrate-up
   ```
4. Build and run:
   ```bash
   make run
   ```

## Commands

### Event / Expense Management
- `/event create <title> [@user1 @user2 ...]` - Create a new shared expense event.
- `/event list` - List active events.
- `/event history` - List closed events.
- `/event info <ID>` - Show details and expenses of a specific event.
- `/event add <ID> <amount> [description]` - Add an expense to an event.
- `/event debts <ID>` - Calculate who owes whom to settle the event.
- `/event close <ID>` - Mark an event as closed.

### Utilities
- `/who [@user1 @user2 ...]` - Pick a random user from the chat or from the provided mentions.
- `@all` or `@все` - Mention all active participants in the group.
