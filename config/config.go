package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ParseConfig(path string) error {
	config := viper.New()
	fmt.Println(config)
	return nil
}
