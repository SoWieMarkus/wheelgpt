# Remove Channel

Removes the authenticated channel and all associated data.

```http
DELETE /api/authentication/remove
Authorization: <web token>
```

## Authentication

Requires [Web Authentication](../../authentication/web.md).

## Response

- Status `204 No Content`

- Status `400 Bad Request`

```js
{ "error": "Validation error message." }
```

- Status `401 Unauthorized`

```js
{ "error": "Unauthorized: Authentication required." }
```
