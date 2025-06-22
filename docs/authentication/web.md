# Web Token Authentication

The `web token` is used to authenticate users (channel owners) in the web interface and for endpoints that require user-level authentication.

## How it works

- The web token is a JWT signed with your backend's `JWT_SECRET_WEB`.
- It contains the channel's ID and an expiration timestamp.
- The backend validates the token, checks expiration, and ensures the channel exists.

## Obtaining a Web Token

1. **Login via Twitch** on the web interface.
2. After successful login, the backend returns a web token:

```json
{ "webToken": "<your-web-token>" }
```

3. Use this token in the `Authorization header` for web-authenticated endpoints.

```http
GET /api/authentication/token
Authorization: <web token>
```

## Token Expiration

- Web tokens are valid for 1 day by default.
- If your token expires, you must log in again to obtain a new one.
