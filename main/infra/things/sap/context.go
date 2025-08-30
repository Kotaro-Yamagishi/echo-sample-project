package sap

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	HeaderKeyRecordModified = "X-Record-Modified"
)

type Context struct {
	context.Context
	recordsModified bool
}

func NewContextFromEcho(c echo.Context) *Context {
	req := c.Request()
	ctx := &Context{Context: req.Context()}
	ctx.setRecordModifiedFromHeader(req)
	return ctx
}

func FromContext(c context.Context) *Context {
	if ctx, ok := c.(*Context); ok {
		return ctx
	}
	return nil
}

func (c *Context) setRecordModifiedFromHeader(req *http.Request) {
	if req.Header.Get(HeaderKeyRecordModified) == "on" {
		c.SetRecordsModified()
	}
}

func (c *Context) RecordsModified() bool {
	return c.recordsModified
}

func (c *Context) SetRecordsModified() {
	c.recordsModified = true
}
