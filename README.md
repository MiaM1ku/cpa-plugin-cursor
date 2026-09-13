# Cursor plugin for CLIProxyAPI

Native CLIProxyAPI plugin that uses an authorized Cursor account as an upstream.

It exposes **base models only**. Reasoning level is passed with OpenAI `reasoning_effort` (`none` / `low` / `medium` / `high` / `xhigh`) instead of advertising every Cursor variant as a separate model.
Thinking families are listed as `<base>-thinking`. The plugin sends Cursor sibling IDs (`<base>-thinking-<effort>` or `<base>-<effort>`) and ignores `-fast` variants.
Max mode is `max_mode` or a `-1m` model suffix.

Quota is shown on the plugin **Cursor 额度** page. CPA Manager Plus original quota cards still only know Codex / Claude / Antigravity / Kimi / xAI, so this plugin does **not** masquerade as Codex.

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
- `cursor/<base-id>` — non-thinking family
- `cursor/<base-id>-thinking` — thinking family

Send effort on the chat request:

```json
{
  "model": "cursor/claude-fable-5-thinking",
  "reasoning_effort": "high",
  "messages": [{"role": "user", "content": "hello"}]
}
```

That request is sent to Cursor as `claude-fable-5-thinking-high`. A non-thinking request uses `claude-fable-5-high`. If `reasoning_effort` is omitted, the plugin defaults to `medium`.

## Quota

Open the plugin management resource **Cursor 额度**. It reads Cursor dashboard usage with the account session cookie:

- Cursor Models: `autoPercentUsed` plus aggregated spend
- Other Models: `apiPercentUsed` plus guaranteed / used spend

- `Cookie: WorkosCursorSessionToken=<account_id>::<access_token>`
- `Origin: https://cursor.com`
- `Referer: https://cursor.com/dashboard`

`quota.provider` reports identifier `cursor` only. Chat still uses provider `cursor`. Do not send a Bearer token to dashboard URLs.

## License

MIT. Unofficial community plugin. Not affiliated with Cursor, Anysphere, CLIProxyAPI, or CPA Manager Plus.
