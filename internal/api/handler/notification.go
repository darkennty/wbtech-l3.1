package handler

import (
	"net/http"
	"strconv"
	"time"

	"WBTech_L3.1/internal/service"
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) handleCreate(ctx *ginext.Context) {
	var req createRequest
	if err := ctx.BindJSON(&req); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Channel == "" || req.Message == "" || req.SendAt == "" {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "missing fields")
		return
	}

	sendAt, err := time.Parse(time.RFC3339, req.SendAt)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "invalid send_at, use RFC3339")
		return
	}

	id, err := h.services.Create(ctx, service.CreateRequest{
		Channel:   req.Channel,
		Recipient: req.Recipient,
		Message:   req.Message,
		SendAt:    sendAt,
	})

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create notification")
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}

	ReturnResultResponse(ctx, ginext.H{"id": id})
}

func (h *Handler) handleGet(ctx *ginext.Context) {
	id := ctx.Param("id")
	n, err := h.services.GetNotificationByID(ctx, id)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusNotFound, "not found")
		return
	}

	ReturnResultResponse(ctx, ginext.H{"notification": n})
}

func (h *Handler) handleCancel(ctx *ginext.Context) {
	id := ctx.Param("id")
	if err := h.services.CancelNotification(ctx, id); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "failed to cancel notification")
		return
	}

	ReturnResultResponse(ctx, ginext.H{"status": "ok"})
}

func (h *Handler) handleDelete(ctx *ginext.Context) {
	id := ctx.Param("id")
	if err := h.services.Delete(ctx, id); err != nil {
		ReturnErrorResponse(ctx, http.StatusBadRequest, "failed to cancel notification")
		return
	}

	ReturnResultResponse(ctx, ginext.H{"status": "ok"})
}

func (h *Handler) handleList(ctx *ginext.Context) {
	limitStr := ctx.Query("limit")
	limit := 50

	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}

	items, err := h.services.GetRecentNotifications(ctx, limit)
	if err != nil {
		ReturnErrorResponse(ctx, http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}

	ReturnResultResponse(ctx, ginext.H{"notifications": items})
}
