package config

import "github.com/spf13/viper"

type Config struct {
	TelegramToken string `mapstructure:"tgToken"`

	QiwiToken        string `mapstructure:"qiwiToken"`
	QiwiWallet       string `mapstructure:"qiwiWallet"`
	QiwiPaymentsPath string `mapstructure:"paymentsPath"`
	QiwiCashInPath   string `mapstructure:"cashInPath"`

	CouchConnString string `mapstructure:"connString"`
	CouchUsername   string `mapstructure:"couchUsername"`
	CouchPassword   string `mapstructure:"couchPassword"`
	CouchBucketName string `mapstructure:"bucketName"`
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
