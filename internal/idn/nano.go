package idn

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var _ Provider = (*NanoIDProvider)(nil)

type NanoIDProvider struct{}

func NewNanoIDProvider() *NanoIDProvider {
	return &NanoIDProvider{}
}

func (idp NanoIDProvider) ID() (id string) {
	id, err := gonanoid.New()
	if err != nil {
		return ""
	}

	return id
}
