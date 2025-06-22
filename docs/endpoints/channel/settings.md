# Update Channel Settings

Update the settings for the authenticated channel.

```http
POST /api/channel/settings
Authorization: <web token>
```

## Authentication

Requires [Web Authentication](../../authentication/web.md).

## Request Body

```js
{
  "botActiveWhenOffline": true,
  "guessDelayTime": 10,
  "usagePublic": true
}
```

## Response

- Status `200 OK`

```js
{
  "id": "channel_id",
  "displayName": "Channel Display Name",
  "profileImage": "https://example.com/image.png",
  "botActiveWhenOffline": true,
  "guessDelayTime": 10, // in s
  "usagePublic": true
}
```

- Status `400 Bad Request`

```js
{ "error": "Invalid request body: <validation error message>" }
```

- Status `401 Unauthorized`

```js
{ "error": "Unauthorized: Authentication required." }
```
