package rustdesk_test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"amidesk-api-server/config"
	"amidesk-api-server/helper/rustdesk"
)

func testConfig(keyFile string) *config.RustdeskBootstrapConfig {
	return &config.RustdeskBootstrapConfig{
		Enabled: true, IDServer: "rust.amitronic.pl", RelayServer: "rust.amitronic.pl",
		APIServer: "https://rust.amitronic.pl", KeyFile: keyFile, PushToAnonymous: true,
	}
}

func TestOptionsLoadsValidPublicKey(t *testing.T) {
	keyFile := filepath.Join(t.TempDir(), "id_ed25519.pub")
	key := base64.StdEncoding.EncodeToString(make([]byte, 32))
	if err := os.WriteFile(keyFile, []byte(key+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	options, status := rustdesk.Options(testConfig(keyFile))
	if !status.KeyLoaded || status.KeyStatus != "loaded" || options[rustdesk.OptionKey] != key {
		t.Fatalf("expected loaded public key, got status=%+v options=%v", status, options)
	}
}

func TestMissingKeyIsNotPushed(t *testing.T) {
	options, status := rustdesk.Options(testConfig(filepath.Join(t.TempDir(), "missing.pub")))
	if status.KeyLoaded || status.KeyStatus != "missing" {
		t.Fatalf("expected missing key status, got %+v", status)
	}
	if _, ok := options[rustdesk.OptionKey]; ok {
		t.Fatal("missing public key must not be sent as an empty key option")
	}
}

func TestDisabledBootstrapDoesNotProduceOptions(t *testing.T) {
	cfg := testConfig("")
	cfg.Enabled = false
	options, status := rustdesk.Options(cfg)
	if options != nil || status.Revision != 0 {
		t.Fatalf("disabled bootstrap must not return options or a revision: %+v", status)
	}
}

func TestKeyRotationChangesRevision(t *testing.T) {
	keyFile := filepath.Join(t.TempDir(), "id_ed25519.pub")
	cfg := testConfig(keyFile)
	if err := os.WriteFile(keyFile, []byte(base64.StdEncoding.EncodeToString(make([]byte, 32))), 0600); err != nil {
		t.Fatal(err)
	}
	_, first := rustdesk.Options(cfg)
	if err := os.WriteFile(keyFile, []byte(base64.StdEncoding.EncodeToString(append(make([]byte, 31), 1))), 0600); err != nil {
		t.Fatal(err)
	}
	_, second := rustdesk.Options(cfg)
	if first.Revision == second.Revision {
		t.Fatal("public key rotation must change revision")
	}
}

func TestMergeOptionsBootstrapOverridesInfrastructure(t *testing.T) {
	merged := rustdesk.MergeOptions(map[string]string{
		rustdesk.OptionKey: "old", "security-option": "kept",
	}, map[string]string{rustdesk.OptionKey: "new", rustdesk.OptionAPIServer: "https://rust.amitronic.pl"})
	if merged[rustdesk.OptionKey] != "new" || merged["security-option"] != "kept" {
		t.Fatalf("unexpected merged strategy: %v", merged)
	}
}