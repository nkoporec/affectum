package utils

import "github.com/spf13/viper"

type Config struct {
	MailHost             string `mapstructure:"MAIL_HOST"`
	MailPort             string `mapstructure:"MAIL_PORT"`
	MailUsername         string `mapstructure:"MAIL_USERNAME"`
	MailPassword         string `mapstructure:"MAIL_PASSWORD"`
	MailFolder           string `mapstructure:"MAIL_FOLDER"`
	AttachmentFolderPath string `mapstructure:"ATTACHMENT_FOLDER_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("affectum")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
