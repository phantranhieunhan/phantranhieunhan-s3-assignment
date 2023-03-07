package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/phantranhieunhan/s3-assignment/common/constant"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "ES microservice config path")
}

type config struct {
	Env         string
	PostgresULR string `mapstructure:"POSTGRES_URL"`
	Server      struct {
		Port string `mapstructure:"PORT"`
	}
}

var C config

func ReadConfig() error {
	Config := &C

	if configPath == "" {
		configPathFromEnv := os.Getenv(constant.CONFIG_PATH)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/pkg/config/config.yaml", getwd)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join(rootDir(), "config"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	env := os.Getenv(constant.ENV)
	if env != "" {
		C.Env = env
	}

	postgres := os.Getenv(constant.POSTGRES_URL)
	if postgres != "" {
		C.PostgresULR = postgres
	}

	port := os.Getenv(constant.PORT)
	if port != "" {
		C.Server.Port = port
	}

	spew.Dump(C)
	return nil
}

func rootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
