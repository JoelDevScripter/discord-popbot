# Project README

# Pop Girls Card Bot 🎶🃏

A Discord bot built in Go that generates collectible cards of pop artists, with rarity levels, eras, and marketplace features.

## Setup

1. Install Go.
2. Run `go mod tidy`.
3. Create `.env` with:
   - `DISCORD_TOKEN`
   - `SPOTIFY_CLIENT_ID`
   - `SPOTIFY_CLIENT_SECRET`
4. Run the bot:
   ```bash
   go run cmd/bot/main.go
   ```
