package security

import (
	"sync"
	"time"

	"amidesk-api-server/config"
)

type attempt struct {
	failures  int
	lockedTo  time.Time
	updatedAt time.Time
}

var attempts = struct {
	sync.Mutex
	entries map[string]attempt
}{entries: make(map[string]attempt)}

func Allow(cfg *config.SecurityConfig, scope, identity string) bool {
	if cfg == nil || cfg.MaxLoginFailures <= 0 || cfg.LoginLockoutMinute <= 0 {
		return true
	}
	key := attemptKey(scope, identity)
	now := time.Now()
	attempts.Lock()
	defer attempts.Unlock()
	entry, found := attempts.entries[key]
	if !found || entry.lockedTo.IsZero() || now.After(entry.lockedTo) {
		if found && !entry.lockedTo.IsZero() && now.After(entry.lockedTo) {
			delete(attempts.entries, key)
		}
		return true
	}
	return false
}

func RecordFailure(cfg *config.SecurityConfig, scope, identity string) {
	if cfg == nil || cfg.MaxLoginFailures <= 0 || cfg.LoginLockoutMinute <= 0 {
		return
	}
	key := attemptKey(scope, identity)
	now := time.Now()
	attempts.Lock()
	defer attempts.Unlock()
	entry := attempts.entries[key]
	if !entry.lockedTo.IsZero() && now.After(entry.lockedTo) {
		entry = attempt{}
	}
	entry.failures++
	entry.updatedAt = now
	if entry.failures >= cfg.MaxLoginFailures {
		entry.lockedTo = now.Add(time.Duration(cfg.LoginLockoutMinute) * time.Minute)
	}
	attempts.entries[key] = entry
}

func Clear(scope, identity string) {
	attempts.Lock()
	defer attempts.Unlock()
	delete(attempts.entries, attemptKey(scope, identity))
}

func attemptKey(scope, identity string) string {
	return scope + "\x00" + identity
}