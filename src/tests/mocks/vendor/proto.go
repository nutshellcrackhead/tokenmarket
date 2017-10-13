package vendor

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type Context struct {
	mock.Mock
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	args := c.Called()
	return args.Get(0).(time.Time), args.Bool(1)
}

func (c *Context) Done() <-chan struct{} {
	args := c.Called()
	return args.Get(0).(<-chan struct{})
}

func (c *Context) Err() error {
	args := c.Called()
	return args.Error(0)
}

func (c *Context) Value(key interface{}) interface{} {
	args := c.Called(key)
	return args.Get(0)
}
