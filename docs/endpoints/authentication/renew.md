# Renew Plugin Token

Whenever you call the endpoint `/api/authentication/token` endpoint you will receive a different token. However the old tokens still remain valid.
If you need to rotate your plugin token, this endpoint.

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
