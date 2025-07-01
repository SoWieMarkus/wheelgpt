# Get Leaderboard Position by Display Name

Returns the leaderboard position and details for a player with the specified display name in a specific Twitch channel.

```http
GET /api/leaderboard/:channelId/user/:displayName
```

## Parameters

- `channelId` (string, required) - The Twitch channel ID
- `displayName` (string, required) - The exact display name to search for

## Authentication

No authentication is required.

## Response

- Status `200 OK`

When player is found

```js
{
  channelId: "123456789",
  userId: "987654321",
  displayName: "PlayerName",
  points: 1500,
  position: 15
}
```

When player is not found

```js
null;
```
