# Get Channel by ID

Returns details about a channel that is using WheelGPT.

```http
GET /api/channel/:channelId
```

## Authentication

No authentication is required.

## Response

- Status `200 OK`

```js
  {
    login: "channel_login",
    displayName: "Channel Display Name",
    profileImage: "https://example.com/image.png",
    isLive: true,
  },
```
