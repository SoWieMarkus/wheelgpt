# Get Plugin Token

Returns the current plugin token for the authenticated channel.

```http
GET /api/authentication/token
Authorization: <web token>
```

## Authentication

Requires [Web Authentication](../../authentication/web.md).

## Response

- Status `200 OK`

```js
{ "pluginToken": "<your-plugin-token>" }
```

- Status `401 Unauthorized`

```js
{ "error": "Unauthorized: Authentication required." }
```
