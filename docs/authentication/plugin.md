# Plugin Token Authentication

The `plugin token` is used to authenticate requests from your OpenPlanet plugin to the WheelGPT backend.
This token is tied to a specific channel and is required for endpoints that update channel-specific data (such as map updates or personal bests).

## How it works

- The plugin token is a JWT signed with your backend's `JWT_SECRET_CHANNEL`.
- It contains the channel's ID (or legacy login) and a unique token.
- The backend validates the token, checks the channel's existence, and ensures the token matches the one stored in the database.

## Obtaining a Plugin Token

Requires [Web Authentication](./web.md).

1. **Login via Twitch** on the web interface.
2. Call the endpoint:

```http
GET /api/authentication/token
Authorization: <web token>
```

The response will contain your plugin token:

```js
{ "pluginToken": "<your-plugin-token>" }
```

3. Use this token in the Authorization header for plugin-authenticated endpoints.

```http
POST /api/trackmania/update/map
Authorization: <plugin token>
```

For more detailed information visit [Get Plugin Token](../endpoints/authentication/token.md).

## Token Renewal

Whenever you call the endpoint `/api/authentication/token` you will receive a different token. However the old tokens still remain valid.
If you need to rotate your plugin token, use:

```http
POST /api/authentication/renew
Authorization: <web token>
```

```js
{ "pluginToken": "<your-plugin-token>" }
```

For more detailed information visit [Renew Plugin Token](../endpoints/authentication/renew.md).
