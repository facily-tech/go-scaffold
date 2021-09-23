package env

import (
	"embed"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ViperConfig(embeds embed.FS) (*viper.Viper, error) {
	viperSetup := viper.GetViper()
	viperSetup.SetConfigType("env")

	app, err := embeds.Open("env/application.env")
	if err != nil {
		return nil, errors.Wrap(err, "cannot open configuration file")
	}
	defer app.Close()

	if err := viperSetup.ReadConfig(app); err != nil {
		return nil, errors.Wrap(err, "viper cannot read configuration")
	}

	viperSetup.AllowEmptyEnv(false)
	viperSetup.AutomaticEnv()

	return viperSetup, nil
}
