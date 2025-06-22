# Get Public Channels

Returns a list of `30` Twitch Channels that are using WheelGPT. The Channels need to enable the `Public Channel` setting to be visibile.

```http
GET /api/landing
```

## Authentication

No authentication is required.

## Response

- Status `200 OK`

```js
[
  {
    login: "channel_login",
    displayName: "Channel Display Name",
    profileImage: "https://example.com/image.png",
    isLive: true,
  },
  // ...
];
```
