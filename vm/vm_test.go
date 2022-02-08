package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVm(t *testing.T) {
	v := NewVm(nil)
	check := func(script string, expect string) {
		b, err := v.Exec(script)
		assert.NoError(t, err)
		assert.Equal(t, expect, string(b))
	}
	check(`x={
		a:1,
		b:2,
	}`, `{"a":1,"b":2}`)
	check(`x=1`, `1`)
	check(`result="1"`, `"1"`)
	check(`result={a:"a",b:{c:"c"}}`, `{"a":"a","b":{"c":"c"}}`)
}
