# Cursor plugin design

- Advertise `<base>` and `<base>-thinking`. Send Cursor sibling IDs `<base>-<effort>` or `<base>-thinking-<effort>`. Ignore `-fast` catalog rows.
- Quota is the plugin **Cursor 额度** page. `executor.http_request` attaches the dashboard session cookie and Origin for Cursor usage APIs. Do not intercept Codex WHAM.
- Usage: OpenAI `usage` prefers Cursor `token_delta`. The host records that payload. `usage_plugin` stays off so analytics is not double-counted. Internal checkpoint retries stay inside one executor call.
