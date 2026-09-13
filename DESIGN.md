# Cursor plugin design

- Advertise one ID per Cursor model family. Effort is `reasoning_effort` on the OpenAI request and is sent as `requested_model.parameters[{id:effort,value}]`.
- Quota is not a plugin page. CPA Manager Plus refreshes Cursor accounts through `/v0/management/api-call`; `executor.http_request` attaches the dashboard session cookie and Origin.
- Usage: OpenAI `usage` prefers Cursor `token_delta`. The host records that payload. `usage_plugin` stays off so analytics is not double-counted. Internal checkpoint retries stay inside one executor call.
