# Renew Plugin Token

Generates and returns a new plugin token for the authenticated channel.

```http
POST /api/authentication/renew
Authorization: <web token>
```

## Authentication

Requires [Web Authentication](../../authentication/web.md).

## Request Body

No request body is needed.

## Response

- Status `200 OK`

```js
{ "pluginToken": "<your-new-plugin-token>" }
```

- Status `401 Unauthorized`

```js
{ "error": "Unauthorized: Authentication required." }
```
