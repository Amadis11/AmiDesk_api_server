package config

import (
	"os"
	"path"
	"amidesk-api-server/util"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	DebugMode         bool                    `yaml:"debugMode"`
	Db                *DbConfig               `yaml:"db"`
	SignKey           string                  `yaml:"signKey"`
	HttpConfig        *HttpConfig             `yaml:"httpConfig"`
	SmtpConfig        *SmtpConfig             `yaml:"smtpConfig"`
	JobsConfig        *JobsConfig             `yaml:"jobsConfig"`
	RustdeskBootstrap *RustdeskBootstrapConfig `yaml:"rustdeskBootstrap"`
	Security           *SecurityConfig         `yaml:"security"`
}

type DbConfig struct {
	Driver   string `yaml:"driver"`
	Dsn      string `yaml:"dsn"`
	TimeZone string `yaml:"timeZone"`
	ShowSql  bool   `yaml:"showSql"`
}

type HttpConfig struct {
	PrintRequestLog bool   `yaml:"printRequestLog"`
	Port            string `yaml:"port"`
	StaticDir       string `yaml:"staticdir"`
}

type SmtpConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Encryption string `yaml:"encryption"` // none ssl/tls starttls
	From       string `yaml:"from"`
}

type DeviceCheckJob struct {
	Duration int `yaml:"duration"`
}

type JobsConfig struct {
	DeviceCheckJob *DeviceCheckJob `yaml:"deviceCheckJob"`
}

type RustdeskBootstrapConfig struct {
	Enabled             bool   `yaml:"enabled"`
	IDServer            string `yaml:"idServer"`
	RelayServer         string `yaml:"relayServer"`
	APIServer           string `yaml:"apiServer"`
	KeyFile             string `yaml:"keyFile"`
	PushToAnonymous     bool   `yaml:"pushToAnonymous"`
	PushToAuthenticated bool   `yaml:"pushToAuthenticated"`
}

type SecurityConfig struct {
	RequireAdminTOTP   bool `yaml:"requireAdminTOTP"`
	MaxLoginFailures   int  `yaml:"maxLoginFailures"`
	LoginLockoutMinute int  `yaml:"loginLockoutMinutes"`
}

var (
	wd       = configDirectory()
	yamlFile = path.Join(wd, "server.yaml")
)

func configDirectory() string {
	if directory := os.Getenv("RUSTDESK_API_CONFIG_DIR"); directory != "" {
		return directory
	}
	directory, _ := os.Getwd()
	return directory
}

func GetDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		DebugMode: false,
		Db: &DbConfig{
			Driver:   "sqlite",
			Dsn:      "./server.db",
			ShowSql:  true,
			TimeZone: "Europe/Warsaw",
		},
		HttpConfig: &HttpConfig{
			Port:      ":8080",
			StaticDir: "dist",
		},
		SignKey: util.RandomString(32),
		JobsConfig: &JobsConfig{
			DeviceCheckJob: &DeviceCheckJob{
				Duration: 30,
			},
		},
		RustdeskBootstrap: &RustdeskBootstrapConfig{},
		Security: &SecurityConfig{
			MaxLoginFailures:   5,
			LoginLockoutMinute: 15,
		},
	}
}

func GetServerConfig() *ServerConfig {
	cfg := GetDefaultServerConfig()
	bytes, err := os.ReadFile(yamlFile)
	if err != nil {
		WriteServerConfig(cfg)
		return cfg
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		WriteServerConfig(cfg)
		return cfg
	}
	return cfg
}

func WriteServerConfig(cfg *ServerConfig) {
	bytes, _ := yaml.Marshal(cfg)
	_ = os.WriteFile(yamlFile, bytes, 0755)
}
