# ADR 0001: Anonymous hbbs Public-Key Bootstrap

- Status: Accepted
- Date: 2026-08-28

## Context

Users must be able to use a stock RustDesk client for ordinary ID-and-password connections without creating or signing in to an API account. The client needs the hbbs public key and the Amitronic server addresses. Manually distributing the public key is error-prone.

The RustDesk heartbeat endpoint is called without an account token. It can therefore serve both logged-out clients and clients which have logged into account features, but the API cannot reliably distinguish those two states on this endpoint.

## Decision

The API may return the hbbs public key and infrastructure options to anonymous heartbeats over HTTPS:

- `key`
- `custom-rendezvous-server`
- `relay-server`
- `api-server`

The key is read only from a mounted `id_ed25519.pub` file. The private `id_ed25519` key must never be mounted into, read by, logged by, or exposed through the API service.

Bootstrap is enabled for anonymous clients. It does not grant account access, authorize a peer connection, expose an API token, or replace RustDesk's ordinary peer authentication. Account login remains optional and only enables account-scoped functions such as address-book synchronization.

All bootstrap traffic uses the configured HTTPS API address. The TLS certificate must be valid and managed by the reverse proxy. The client receives a stable bootstrap revision and receives configuration again only after a server address or public-key change.

If the public-key file is missing or invalid, the API omits `key`; it must never return an empty key that could erase a known-good client setting.

## Security Rationale

The hbbs public key is intentionally public. Its confidentiality is not a security property. Authenticity is the relevant property: HTTPS prevents an on-path attacker from substituting a different key while the client first receives bootstrap configuration. The client subsequently uses the key to validate the hbbs identity and secure protocol flows supported by stock RustDesk.

The main residual risks are compromise of the HTTPS endpoint or its TLS termination, and an administrator configuring an incorrect key or server address. The API mitigates accidental exposure by reporting only a public-key fingerprint in the administrator view and retaining no private key material.

## Consequences

- Fresh clients need only enter the Amitronic API Server initially; the documented ID Server value remains a routing hint.
- Normal RustDesk connections work without API login.
- `pushToAnonymous` must remain enabled because stock heartbeat does not carry account identity.
- This decision does not make every anonymous TCP relay flow use secure TCP. That depends on stock-client KeyExchange behavior and is specified separately in `hbbs-secure-tcp-integration.md`.