package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
)

// Config holds the configuration values for the application.
type Config struct {
	Daemon struct {
		Enabled          bool          `env:"ENABLED" envDefault:"false"`
		Frequency        time.Duration `env:"FREQUENCY" envDefault:"30s"`
		ExecutionTimeout time.Duration `env:"EXECUTION_TIMEOUT" envDefault:"30m"`
	} `envPrefix:"DAEMON_"`

	Discord struct {
		Token        string `env:"TOKEN,required"`
		PublicKey    string `env:"PUBLIC_KEY,required"`
		GuildId      uint64 `env:"GUILD_ID,required"`
		ChannelId    uint64 `env:"CHANNEL_ID,required"`
		UpdateRoleId uint64 `env:"UPDATE_ROLE_ID,required"`
	} `envPrefix:"DISCORD_"`

	StatusPage struct {
		ApiKey string `env:"API_KEY,required"`
		PageId string `env:"PAGE_ID,required"`
		Url    string `env:"URL" envDefault:"status.ticketsbot.cloud"`
	} `envPrefix:"STATUSPAGE_"`

	ServerAddr string `env:"SERVER_ADDR" envDefault:":8080"`

	DatabaseUri string `env:"DATABASE_URI"`

	JsonLogs bool          `env:"JSON_LOGS" envDefault:"false"`
	LogLevel zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`
}

var Conf Config

func LoadConfig() (Config, error) {
	if _, err := os.Stat("config.toml"); err == nil {
		return fromToml()
	} else {
		return fromEnvvar()
	}
}

func fromToml() (Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		return Config{}, err
	}

	return config, nil
}

func fromEnvvar() (Config, error) {
	return env.ParseAs[Config]()
}
