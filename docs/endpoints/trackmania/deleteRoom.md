# Delete Trackmania Room

Delete the current Trackmania room for a channel.

```http
DELETE /api/trackmania/room
Authorization: <plugin token>
```

## Authentication

Requires [Plugin Authentication](../../authentication/plugin.md).

## Response

- Status `200 OK`  
  If the room was successfully deleted.

```js
{ "message": "Room deleted successfully." }
```

- Status `200 OK`  
  If the channel is currently offline and the bot is set to be disabled while offline, the request will not be processed.

```js
{ "message": "Channel is offline." }
```

- Status `401 Unauthorized`

```js
{ "error": "Authentication required." }
```
