package infrastructure

import "os"

type EnvInterface interface {
	Getenv(key string) string
}

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

func (e *Env) Getenv(key string) string {
	return os.Getenv(key)
}
