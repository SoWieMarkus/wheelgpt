# Get Public Channels

Returns a list of `30` Twitch Channels that are using WheelGPT. The Channels need to enable the `Public Channel` setting to be visibile.

```http
GET /api/channel
```

## Authentication

No authentication is required.

## Response

- Status `200 OK`

```js
[
  {
    id: "channel_id",
    login: "channel_login",
    displayName: "Channel Display Name",
    profileImage: "https://example.com/image.png",
    isLive: true,
  },
  // ...
];
```
