package dragonbook

import (
	"errors"
	"fmt"
)

type Env struct {
	table  map[string]*TokenInterface
	parent *Env
}

func (env *Env) Put(s string, sym *TokenInterface) {
	env.table[s] = sym
}

func (env *Env) Get(s string) (*TokenInterface, error) {
	for e := env; e != nil; e = e.parent {
		found := e.table[s]
		if found != nil {
			return found, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Symbol not found: %s\n", s))
}
