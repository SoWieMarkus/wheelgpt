# Update Personal Best

Update the personal best time for the current channel and notify the bot.

```http
POST /api/trackmania/update/pb
Authorization: <plugin token>
```

## Authentication

Requires [Plugin Authentication](../../authentication/plugin.md).

## Request Body

```js
{
  "time": 123456 // time in ms
}
```

## Response

- Status `200 OK`

```js
{ "message": "Personal best updated successfully." }
```

If the channel is currently offline and the bot is set to be disabled while offline, the request will not be processed.

- Status `200 OK`

```js
{ "message": "Channel is offline." }
```

- Status `400 Unauthorized`

```js
{ "error": "Validation error message." }
```

- Status `401 Unauthorized`

```js
{ "error": "Authentication required." }
```
