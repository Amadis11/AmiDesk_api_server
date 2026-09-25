package rustdesk

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"os"
	"strings"

	"amidesk-api-server/config"
)

const (
	OptionKey              = "key"
	OptionRendezvousServer = "custom-rendezvous-server"
	OptionRelayServer      = "relay-server"
	OptionAPIServer        = "api-server"
)

type BootstrapStatus struct {
	Enabled             bool   `json:"enabled"`
	KeyLoaded           bool   `json:"key_loaded"`
	KeyFingerprint      string `json:"key_fingerprint,omitempty"`
	KeyStatus           string `json:"key_status"`
	KeyFile             string `json:"key_file"`
	IDServer            string `json:"id_server"`
	RelayServer         string `json:"relay_server"`
	APIServer           string `json:"api_server"`
	Revision            int64  `json:"revision"`
	PushToAnonymous     bool   `json:"push_to_anonymous"`
	PushToAuthenticated bool   `json:"push_to_authenticated"`
}

func BootstrapStatusFor(cfg *config.RustdeskBootstrapConfig) BootstrapStatus {
	status, _ := statusAndKey(cfg)
	return status
}

func Options(cfg *config.RustdeskBootstrapConfig) (map[string]string, BootstrapStatus) {
	status, key := statusAndKey(cfg)
	if !status.Enabled {
		return nil, status
	}

	options := map[string]string{}
	if status.IDServer != "" {
		options[OptionRendezvousServer] = status.IDServer
	}
	if status.RelayServer != "" {
		options[OptionRelayServer] = status.RelayServer
	}
	if status.APIServer != "" {
		options[OptionAPIServer] = status.APIServer
	}
	if status.KeyLoaded {
		options[OptionKey] = key
	}
	return options, status
}

func statusAndKey(cfg *config.RustdeskBootstrapConfig) (BootstrapStatus, string) {
	if cfg == nil {
		return BootstrapStatus{KeyStatus: "disabled"}, ""
	}

	status := BootstrapStatus{
		Enabled:             cfg.Enabled,
		KeyFile:             cfg.KeyFile,
		IDServer:            cfg.IDServer,
		RelayServer:         cfg.RelayServer,
		APIServer:           cfg.APIServer,
		PushToAnonymous:     cfg.PushToAnonymous,
		PushToAuthenticated: cfg.PushToAuthenticated,
		KeyStatus:           "missing",
	}
	key, fingerprint, keyStatus := loadPublicKey(cfg.KeyFile)
	status.KeyStatus = keyStatus
	status.KeyLoaded = key != ""
	status.KeyFingerprint = fingerprint
	if status.Enabled {
		status.Revision = revision(cfg, key)
	}
	return status, key
}

func MergeOptions(strategy map[string]string, bootstrap map[string]string) map[string]string {
	merged := make(map[string]string, len(strategy)+len(bootstrap))
	for key, value := range strategy {
		merged[key] = value
	}
	for key, value := range bootstrap {
		merged[key] = value
	}
	return merged
}

func loadPublicKey(filename string) (string, string, string) {
	if strings.TrimSpace(filename) == "" {
		return "", "", "missing"
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", "", "missing"
	}
	key := strings.TrimSpace(string(bytes))
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(decoded) != 32 {
		return "", "", "invalid"
	}
	fingerprintHash := sha256.Sum256(decoded)
	return key, hex.EncodeToString(fingerprintHash[:8]), "loaded"
}

func revision(cfg *config.RustdeskBootstrapConfig, key string) int64 {
	values := []string{cfg.IDServer, cfg.RelayServer, cfg.APIServer, key}
	hash := sha256.Sum256([]byte(strings.Join(values, "\x00")))
	return int64(binary.BigEndian.Uint64(hash[:8]) & 0x7fffffffffffffff)
}