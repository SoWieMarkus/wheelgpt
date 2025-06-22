# Get Channel Information

Returns information about the authenticated channel.

```http
GET /api/channel/me
Authorization: <web token>
```

## Authentication

Requires [Web Authentication](../../authentication/web.md).

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

- Status `401 Unauthorized`

```js
{ "error": "Unauthorized: Authentication required." }
```
