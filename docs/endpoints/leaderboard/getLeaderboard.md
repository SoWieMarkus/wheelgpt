# Get Leaderboard Page

Returns a paginated list of players from the guesser leaderboard for a specific Twitch channel.
The page size is `100`.

```http
GET /api/leaderboard/:channelId
```

## Parameters

- `channelId` (string, required) - The Twitch channel ID

## Query Parameters

- `page` (number, optional) - Page number for pagination. Must be a positive integer. Default: `1`

## Authentication

No authentication is required.

## Response

- Status `200 OK`

```js
[
  {
    channelId: "123456789",
    userId: "987654321",
    displayName: "PlayerName",
    points: 1500,
    position: 1,
    perfectGuessCount: 1,
  },
  // ... up to 100 entries per page
];
```
