# Login with Twitch Code

Exchanges a Twitch OAuth code for a WheelGPT web token.

```http
POST /api/authentication/login
```

## Request Body

```js
{
  "code": "twitch_oauth_code"
}
```

## Response

- Status `200 OK`

```js
{ "webToken": "<your-web-token>" }
```

- Status `400 Bad Request`

```js
{ "error": "Validation error message." }
```
