package config

import (
	"github.com/koding/multiconfig"
)

func Load(path string, config interface{}) (err error) {

	m := multiconfig.NewWithPath(path)

	err = m.Load(config)
	if err != nil {
		config = nil
	}

	return err
}
