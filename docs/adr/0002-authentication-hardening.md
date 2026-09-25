# ADR 0002: Authentication Hardening for Admin and API Accounts

- Status: Proposed
- Date: 2026-08-28

## Context

Stock RustDesk clients support API account login and its existing TOTP challenge, but standard peer connections use a separate ID-and-password flow. The API must keep account login optional for normal RustDesk operation.

The administration panel is a browser application protected by HTTPS, a password, CAPTCHA, and a short-lived session. Before this ADR, its login did not enforce TOTP even where the administrator account had a configured TOTP secret. The user-list endpoint also returned TOTP secrets to every authenticated administrator.

## Decision

1. The public API bootstrap and heartbeat remain anonymous. No password, token, TOTP code, or private hbbs key is exposed through that path.
2. TOTP secrets are not returned in the administrator user-list response. A new secret is generated only during an explicit setup or rebind action and displayed once in the authenticated administration session.
3. Administrator TOTP enforcement should be introduced in two phases:
   - Phase 1: enable TOTP for every administrator through the existing user-management GUI and verify recovery access.
   - Phase 2: enable a server-side `requireAdminTOTP` configuration. Admin login must reject accounts without an enrolled TOTP secret; do not enable this before every active administrator has enrolled.
4. Apply rate limiting to `/api/login` and `/admin/auth/login`: allow five failed attempts per account, then deny attempts for 15 minutes. Successful login clears the counter. Do not include source IP in the counter, because rotating addresses would otherwise bypass the protection. Use generic error responses so usernames cannot be enumerated.
5. Keep CAPTCHA as a second abuse-control layer for the administrator panel. Deploy the panel only via HTTPS with HSTS, a valid TLS certificate, and no direct public exposure of backend port `8080`.

## Consequences

TOTP protects login to account and administration features, not the ordinary RustDesk peer password prompt. Administrators should use a password manager, unique high-entropy passwords, and a recovery procedure with two trusted operators. In-memory rate limits are acceptable for a single instance; multi-replica deployment requires a shared store such as Redis.

The `requireAdminTOTP` and rate-limit enforcement are intentionally not enabled by this ADR alone: enabling them without enrolment, persistence design, and an administrator recovery procedure risks locking out the deployment.