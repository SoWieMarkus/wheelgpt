# Update Trackmania Room

Update or create the current Trackmania room for a channel.

```http
POST /api/trackmania/update/room
Authorization: <plugin token>
```

## Authentication

Requires [Plugin Authentication](../../authentication/plugin.md).

## Request Body

```js
{
  "login": "RoomHostLogin",
  "name": "Room Name",
  "numberOfPlayers": 12, //current number of players in room
  "maxPlayers": 64 // max players possible in room
}
```

## Response

- Status `200 OK`  
  If the room was successfully updated or created.

```js
{ "message": "Room updated successfully." }
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
