# Cursor plugin for CLIProxyAPI

Native CLIProxyAPI plugin that uses an authorized Cursor account as an upstream.

It exposes **base models only**. Reasoning level is passed with OpenAI `reasoning_effort` (`none` / `low` / `medium` / `high` / `xhigh`) instead of advertising every Cursor variant as a separate model.

Quota refresh uses the original CPA / CPA Manager Plus account quota path (`POST /v0/management/api-call`). The plugin injects the Cursor dashboard session cookie and `Origin: https://cursor.com` so `usage-summary`, `get-current-period-usage`, and `get-aggregated-usage-events` work. There is no separate plugin quota page.

Token usage in OpenAI responses prefers Cursor `token_delta` events. The host records that payload once. This plugin does **not** enable `usage_plugin`, so CPA usage analytics is not double-counted. Checkpoint retries stay inside one executor call.

## Install from the plugin store

Use CPA Manager Plus **Plugins → Store** after this plugin is listed, or install the Linux amd64 zip from GitHub Releases:

`cursor_<version>_linux_amd64.zip`

Place `cursor.so` in `plugins/linux/amd64/` and enable it:

```yaml
plugins:
  enabled: true
  dir: plugins
  configs:
    cursor:
      enabled: true
```

Restart CLIProxyAPI.

## Models

- `cursor/auto` — Cursor Auto
- `cursor/<base-id>` — one entry per Cursor model family (effort suffixes stripped)

Send effort on the chat request:

```json
{
  "model": "cursor/claude-4.5-sonnet",
  "reasoning_effort": "high",
  "messages": [{"role": "user", "content": "hello"}]
}
```

Model IDs that still include an effort suffix are accepted and converted to the base model plus `requested_model.parameters[{id:effort}]`.

## Quota

Refresh Cursor accounts from CPA Manager Plus **Accounts / Quota**. The plugin HTTP passthrough sets:

- `Cookie: WorkosCursorSessionToken=<account_id>::<access_token>`
- `Origin: https://cursor.com`
- `Referer: https://cursor.com/dashboard`

Do not send a Bearer token to those dashboard URLs.

## License

MIT. Unofficial community plugin. Not affiliated with Cursor, Anysphere, CLIProxyAPI, or CPA Manager Plus.
