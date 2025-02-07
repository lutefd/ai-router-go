package idgen_test

import (
	"testing"

	"github.com/lutefd/ai-router-go/pkg/idgen"
	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	_ = idgen.Init(1)
	id1 := idgen.Generate()
	id2 := idgen.Generate()

	assert.NotEqual(t, id1, id2)
	assert.Contains(t, id1, "chat")
}

func TestConcurrency(t *testing.T) {
	_ = idgen.Init(1)
	ids := make(chan string, 1000)
	for i := 0; i < 1000; i++ {
		go func() { ids <- idgen.Generate() }()
	}

	unique := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := <-ids
		assert.False(t, unique[id])
		unique[id] = true
	}
}
