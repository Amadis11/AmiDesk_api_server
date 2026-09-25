# hbbs Secure TCP Integration

Bootstrap distributes the hbbs public key, but it does not itself make anonymous peer traffic use secure TCP.

## Stock Client Behavior

In the RustDesk target code, `src/client.rs`, `Client::request_relay`, invokes `secure_tcp` only when the configured server key is non-empty and either the account `token` or a `switch_code` is non-empty:

```rust
if !key.is_empty() && (!token.is_empty() || !switch_code.is_empty()) {
    secure_tcp(&mut socket, key).await?;
}
```

`src/rendezvous_mediator.rs`, `RendezvousMediator::start_tcp`, does secure its persistent rendezvous connection with the configured key. The relay request includes `id`, `token`, `uuid`, `relay_server`, `secure`, and `switch_code`. The punch request built in `src/client.rs` includes `id`, `token`, `nat_type`, `licence_key`, `conn_type`, `version`, `udp_port`, `force_relay`, and `switch_code`.

## Requirements for the hbbs Fork

1. Implement the existing KeyExchange framing expected by `src/common.rs`, `secure_tcp_impl`; it validates that the configured Base64 key decodes to a 32-byte Ed25519 public key.
2. Protect token-bearing TCP relay requests and switch-side flows only after the client has initiated KeyExchange. Do not assume an anonymous client starts this handshake solely because it has received a bootstrap key.
3. Preserve stock handling for tokenless `PunchHoleRequest` traffic. The API must not alter the rendezvous protocol or invent an API-to-hbbs token path.
4. Verify the patched hbbs against three cases: a key without token, a key with account token, and a tokenless remote endpoint. The first may not exercise secure relay TCP under stock-client logic.

This is a client contract limitation, not something the API bootstrap endpoint can override.