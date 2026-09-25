package security_test

import (
	"testing"

	"amidesk-api-server/config"
	"amidesk-api-server/helper/security"
)

func TestLoginFailuresLockIdentityAcrossIPs(t *testing.T) {
	cfg := &config.SecurityConfig{MaxLoginFailures: 5, LoginLockoutMinute: 15}
	identity := "login-attempts-test@example.com"
	security.Clear("api", identity)
	defer security.Clear("api", identity)

	for attempt := 1; attempt <= 5; attempt++ {
		if !security.Allow(cfg, "api", identity) {
			t.Fatalf("expected login to be allowed before failure %d", attempt)
		}
		security.RecordFailure(cfg, "api", identity)
	}

	if security.Allow(cfg, "api", identity) {
		t.Fatal("expected identity to be locked after five failures")
	}
}