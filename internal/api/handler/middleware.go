package handler

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wb-go/wbf/ginext"
)

func handlerFunc(f ginext.HandlerFunc) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		logrus.Printf("%s %s %s\n", c.Request.Method, c.Request.RequestURI, time.Now().Format(time.RFC3339))
		f(c)
	}
}
