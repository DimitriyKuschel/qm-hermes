package providers

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"queue-manager/internal/structures"
	"strings"
)

func NewConfigProvider(flags *structures.CliFlags) (*structures.Config, error) {
	var conf structures.Config

	filename := filepath.Base(flags.ConfigPath)
	viper.AddConfigPath(strings.Replace(flags.ConfigPath, filename, "", 1))
	viper.SetConfigName(strings.Replace(filename, ".yml", "", 1))
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	cnfValidator := NewCnfValidator(&conf)
	err = cnfValidator.Validate()
	if err != nil {
		return nil, err
	}

	conf.AppName = "QM_Hermes"
	conf.Path = flags.ConfigPath
	conf.Debug = flags.DebugMode

	return &conf, nil
}
