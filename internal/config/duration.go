package config

import (
	"time"

	"github.com/pkg/errors"
)

type DurationConfig time.Duration

func (d *DurationConfig) Duration() time.Duration {
	return time.Duration(*d)
}

func (d *DurationConfig) Decode(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	*d = DurationConfig(duration)

	return nil
}

func (d *DurationConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v interface{}
	if err := unmarshal(&v); err != nil {
		return err
	}

	switch typedValue := v.(type) {
	case string:
		return d.Decode(typedValue)
	case int:
		*d = DurationConfig(typedValue)

		return nil
	default:
		return errors.Errorf("yaml: cannot unmarshal %v into value of type DurationConfig", v)
	}
}
