package main

import "errors"

const NotFoundError = "not_found"

type Keeper interface {
	Get(key string) (string, error)
	Set(key string, message string) error
	Clean(key string) error
}

type DummyKeeper struct {
	mem map[string]string
}

func (k DummyKeeper) Get(key string) (string, error) {
	value, ok := k.mem[key]
	if !ok {
		return "", errors.New(NotFoundError)
	}
	return value, nil
}

func (k DummyKeeper) Set(key string, message string) error {
	k.mem[key] = message
	return nil
}

func (k DummyKeeper) Clean(key string) error {
	delete(k.mem, key)
	return nil
}

func getKeeper() Keeper {
	return DummyKeeper{make(map[string]string)}
}

var keeper = getKeeper()
