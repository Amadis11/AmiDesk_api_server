# RustDesk Bootstrap Client Contract

This API server supports stock RustDesk clients targeted by this project (1.4.6). It does not require a custom client.

## Confirmed Client Behavior

The RustDesk client sends `POST /api/heartbeat` from `src/hbbs_http/sync.rs`, function `start_hbbs_sync_async`. `heartbeat_url` uses the configured `api-server` or derives an API URL from `custom-rendezvous-server`; a client must therefore be configured with an API Server before it can receive bootstrap options.

The heartbeat request has no authentication header: the call is `post_request(url.clone(), v.to_string(), "")`. It sends the locally persisted `strategy_timestamp` as `modified_at`.

The response is deserialized in the same function. Its `strategy` field uses `StrategyOptions`, whose supported configuration map is `config_options: HashMap<String, String>`. When `modified_at` differs from the sent value, the response value is persisted as `strategy_timestamp`. When a `strategy` exists, `handle_config_options` writes every `config_options` entry through `Config::set_options`.

The exact option names are confirmed by `src/hbbs_http/sync.rs` and the custom-server UI in `src/ui/index.tis`:

- `key`
- `custom-rendezvous-server`
- `relay-server`
- `api-server`

The API returns a strategy only when the supplied `modified_at` differs from the bootstrap revision. This avoids needless writes while ensuring a change to the configured servers or mounted public key is propagated.

## Operational Consequences

1. On a fresh stock client, set API Server to `https://rust.amitronic.pl` manually. Set ID Server to `rust.amitronic.pl` as the initial routing hint. Leave Key empty.
2. The next heartbeat installs the `key`, rendezvous, relay, and API options. Restarting the client retains the saved options.
3. Account login is optional for normal peer connections, address book synchronization is separate from the bootstrap path.
4. Stock heartbeat does not transmit the account token. Consequently, it cannot distinguish an authenticated heartbeat from an anonymous one. `pushToAnonymous` must be enabled for stock-client bootstrap, including after account login.
5. A missing or invalid mounted key omits the `key` option entirely; it never sends `"key": ""` and therefore does not clear a previously valid client key.

## Test Procedure

Test a fresh stock RustDesk 1.4.6 client with ID Server `rust.amitronic.pl`, API Server `https://rust.amitronic.pl`, Relay Server empty or `rust.amitronic.pl`, and an empty Key.

1. Without account login, wait for heartbeat and confirm the custom-server settings contain the hbbs public key and configured servers. Restart and confirm they remain.
2. Connect to a remote device using ID and password. The remote device need not have an API account.
3. Log in, confirm address-book operations work, then confirm heartbeat bootstrap settings remain present.
4. Replace only the read-only mounted `id_ed25519.pub`. Wait for heartbeat, confirm the revision changes in `GET /admin/bootstrap`, then confirm the client receives the replacement key.

The admin endpoint requires an administrator session and never returns the public key contents.