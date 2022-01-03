package main

import "github.com/google/uuid"

const DummyTestKey = "test_key"

type KeyBuilder interface {
	Get() (string, error)
}

type DummyKeyBuilder struct {
}

func (k DummyKeyBuilder) Get() (string, error) {
	return DummyTestKey, nil
}

type UUIDKeyBuilder struct{}

func (k UUIDKeyBuilder) Get() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func getKeyBuilder() KeyBuilder {
	return DummyKeyBuilder{}
}
