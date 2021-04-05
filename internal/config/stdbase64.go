package config

import (
	b64 "encoding/base64"

	"github.com/pkg/errors"
)

type StdBase64 []byte

func (b *StdBase64) String() string {
	return string(*b)
}

func (b *StdBase64) Decode(value string) error {
	sDec, err := b64.StdEncoding.DecodeString(value)
	if err != nil {
		return errors.Errorf("failed to decode %s from base64 into value of type StdBase64", value)
	}

	*b = sDec

	return nil
}

func (b *StdBase64) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var input string
	if err := unmarshal(&input); err != nil {
		return err
	}

	return b.Decode(input)
}
