# Switching from the TypeScript bot to the Go bot

The Go backend reads and writes the same SQLite file as the old Prisma backend. It never changes the old tables. On its first start it only adds new ones (`ChannelState`, `PluginToken`, `Session`, `GuessRound`, `GuessResult`) and records what it did in goose's table, `goose_db_version`. Prisma's `_prisma_migrations` is left alone. The old image can therefore still open the file if you need to roll back. Every streamer, setting, current map, open guess and all-time point survives the switch without a conversion step.

## What the operator does once

1. Stop the old container and make a copy of the database file. The copy is the rollback point.

2. In the Twitch developer console, add two OAuth redirect URLs to the existing application: `https://wheelgpt.dev/api/v1/auth/twitch/callback` and `https://wheelgpt.dev/api/v1/auth/bot/callback`. The client id and secret stay the same.

3. Start the two new images, `ghcr.io/sowiemarkus/wheelgpt/backend` and `ghcr.io/sowiemarkus/wheelgpt/website` (see `dev/docker-compose.yml` for how they fit together; the website's nginx serves the site and proxies `/api` to the backend on port 8080). Point the backend at the same database file with `DATABASE_PATH`, and set the following:
   - `PUBLIC_URL=https://wheelgpt.dev`
   - `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET`
   - `TWITCH_EVENTSUB_SECRET`
   - `BOT_USERNAME`
   - `JWT_SECRET_CHANNEL`, with the old value
   - `EVENTSUB_SYNC=true`, which replaces `UPDATE_WEB_HOOKS`
   - `TRUST_PROXY_HEADERS=true`, because the website's nginx sets `X-Real-IP`

   The backend runs as an unprivileged user (uid 65532). Make the database directory writable for that user, for example with `chown -R 65532:65532` on the volume. Metrics are served on a separate port (9091), not through the website.

4. Open `https://wheelgpt.dev/api/v1/auth/bot` and log in as the bot account. This grants `user:read:chat`, `user:write:chat` and `user:bot`, so the bot can read and send chat through EventSub and the Helix API instead of IRC. The backend then subscribes every channel to `stream.online`, `stream.offline` and `channel.chat.message`. The old webhook callback `/api/twitch/webhook/stream` is still accepted, and the subscriptions are moved to `/api/v1/twitch/eventsub` automatically.

5. Watch the logs for `eventsub: synced`. Check `/readyz` on the backend and point Prometheus at port 9091 (`monitoring/` has the scrape config, alerts and a Grafana dashboard).

## What every streamer does once

Nothing breaks on switchover day. Existing plugin tokens keep working because they are still verified with the old secret. Existing settings and leaderboards stay where they are. The chat bot keeps answering in every channel where the bot account is a moderator, because Twitch lets a moderator bot read chat without further permission.

To finish the move, each streamer logs in once on wheelgpt.dev with Twitch. The only permission requested is `channel:bot`, which lets the bot read and write chat in that one channel, moderator or not. The dashboard shows whether that step is still open ("Twitch refused access to your chat" plus a button to log in again). Streamers who would rather not log in can type `/mod wheelgpt` in their chat instead.

Streamers may also create a new plugin token on the dashboard and paste it into the plugin. Doing so revokes the old token. Plugin version 1.5.0 sends the token as `Authorization: Bearer …` and talks to the new `/api/v1/plugin/...` endpoints. Older plugin versions keep using `/api/trackmania/...`, which stays available and answers with `Deprecation` and `Sunset` headers until 30 June 2027.

## What changes for viewers

The commands are the same, with one addition: `!lb` now shows this month's ranking, while `!lb all` shows the all-time board and `!lb <name>` shows both positions. Monthly leaderboards start empty at the switch, because the old backend never recorded individual rounds. The all-time board keeps every point earned before. One rule became stricter: changing a guess now restarts its "minimum guess age". Before, a viewer could guess early and then change the time at the last second.

The setting "bot active when offline" is gone. Map, room and personal bests from the plugin are only accepted while the stream is live, and the end of a stream always clears the map, the room and the open guesses. Streamers who had the setting turned on should know that PBs driven off stream no longer count. The value stays in the database, so a rollback keeps it.

## Rolling back

Stop the new container and start the old image on the same file. The old backend ignores the new tables. Two things are lost: new plugin tokens, which the old backend doesn't know, and points earned after the switch. Those points stay in the all-time table and remain visible. Monthly data is kept, unused, until you switch forward again.

## Cleaning up later

Once the dashboard shows that every active streamer has logged in again and created a new token, remove `JWT_SECRET_CHANNEL`. Legacy tokens are then rejected. After the sunset date, delete `MountLegacy` in `internal/plugin/handler.go` and the `/api/twitch/webhook/stream` route in `internal/app/router.go`.
