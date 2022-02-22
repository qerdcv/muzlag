package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

//used to print errors majorly.
const appPrefix = "muzlag"

type Config struct {
	Token  string
	Prefix string `default:"!"`
}

func New() (c Config, err error) {

	err = envconfig.Process(appPrefix, &c)

	if err != nil {
		fmt.Println(err.Error())
		return c, err
	}

	return c, nil
}
