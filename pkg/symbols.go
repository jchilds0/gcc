package pkg

import (
	"errors"
	"fmt"
)

type Env struct {
	table  map[string]*Id
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{parent: parent, table: map[string]*Id{}}
}

func (env *Env) Put(s string, id *Id) {
	env.table[s] = id
}

func (env *Env) Get(s string) (*Id, error) {
	for e := env; e != nil; e = e.parent {
		found := e.table[s]
		if found != nil {
			return found, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Symbol not found: %s\n", s))
}
