# Update Trackmania Map

Update or delete the current Trackmania map for a channel.
If a map is provided, it will be saved and TMX data will be fetched asynchronously.
If no map is provided, the current map will be deleted.

```http
POST /api/trackmania/update/map
Authorization: <plugin token>
```

## Authentication

Requires [Plugin Authentication](../../authentication/plugin.md).

## Request Body

```js
{
  "name": "TMGL - Spring 2024 - 01",
  "uid": "abcdef1234567890",
  "author": "AuthorName",
  "authorTime": 123456, // time in ms
  "goldTime": 130000, // time in ms
  "silverTime": 140000, // time in ms
  "bronzeTime": 150000, // time in ms
  "championTime": 120000 // time in ms
}
```

## Response

- Status `200 OK`
  If a map was provided.

```js
{ "message": "Successfully updated the map." }
```

- Status `200 OK`
  If the channel is currently offline and the bot is set to be disabled while offline, the request will not be processed.

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
