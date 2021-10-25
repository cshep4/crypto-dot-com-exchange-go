package id

import "math/rand"

type (
	IDGenerator interface {
		Generate() int64
	}
	Generator struct{}
)

func (Generator) Generate() int64 {
	return rand.Int63()
}
