package redisutils

import (
	"testing"
	"time"

	"github.com/sandrolain/go-utilities/pkg/testredisutils"
	"github.com/stretchr/testify/assert"
)

const (
	TestPassword = "development.password"
)

type TestStruct struct {
	Foo string
	Bar int
}

func TestSetGetDel(t *testing.T) {
	redisMock := testredisutils.NewMockServer(t, TestPassword)

	red, err := NewClient(redisMock.Addr(), TestPassword, nil, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	val := TestStruct{"hello", 123}

	{
		red.Set(Key{"testing", "setGet"}, &val, 0)
		redisMock.FastForward(time.Second)
		var res TestStruct
		ok, err := red.Get(Key{"testing", "setGet"}, &res)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("value should exist")
		}
		assert.Equal(t, res, val)
	}

	{
		red.Set(Key{"testing", "setGet"}, &val, time.Millisecond*1500)

		redisMock.FastForward(time.Second)
		var res TestStruct
		ok, err := red.Get(Key{"testing", "setGet"}, &res)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("value should exist")
		}
		assert.Equal(t, res, val)

		redisMock.FastForward(time.Second)
		var res2 TestStruct
		ok, err = red.Get(Key{"testing", "setGet"}, &res2)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal("value should not exist")
		}
	}
}
