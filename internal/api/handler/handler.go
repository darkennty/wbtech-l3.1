package handler

import (
	"WBTech_L3.1/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type createRequest struct {
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	SendAt    string `json:"send_at"`
}

type Handler struct {
	services *service.Service
	logger   zlog.Zerolog
}

func NewHandler(services *service.Service, logger zlog.Zerolog) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *ginext.Engine {
	r := ginext.New("")

	r.POST("/notify", handlerFunc(h.handleCreate))
	r.GET("/notify", handlerFunc(h.handleList))
	r.GET("/notify/:id", handlerFunc(h.handleGet))
	r.PUT("/notify/:id", handlerFunc(h.handleCancel))
	r.DELETE("/notify/:id", handlerFunc(h.handleDelete))

	r.Static("/static", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	return r
}
