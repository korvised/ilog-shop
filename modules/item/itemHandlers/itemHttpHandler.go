package itemHandlers

import (
	"context"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/item"
	"github.com/korvised/ilog-shop/modules/item/itemUsecases"
	"github.com/korvised/ilog-shop/pkg/request"
	"github.com/korvised/ilog-shop/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	ItemHttpHandlerService interface {
		CreateItem(c echo.Context) error
		GetItem(c echo.Context) error
		GetItems(c echo.Context) error
		UpdateItem(c echo.Context) error
		EnableItem(c echo.Context) error
		DisableItem(c echo.Context) error
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecases.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecases.ItemUsecaseService) ItemHttpHandlerService {
	return &itemHttpHandler{cfg, itemUsecase}
}

func (h *itemHttpHandler) CreateItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.CreateItem(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttpHandler) GetItem(c echo.Context) error {
	ctx := context.Background()

	itemID := c.Param("item_id")

	res, err := h.itemUsecase.GetItem(ctx, itemID)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) GetItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemSearchReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.GetItems(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) UpdateItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	itemID := c.Param("item_id")

	req := new(item.ItemUpdateReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.UpdateItem(ctx, itemID, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EnableItem(c echo.Context) error {
	ctx := context.Background()

	itemID := c.Param("item_id")

	res, err := h.itemUsecase.UpdateItemStatus(ctx, itemID, true)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) DisableItem(c echo.Context) error {
	ctx := context.Background()

	itemID := c.Param("item_id")

	res, err := h.itemUsecase.UpdateItemStatus(ctx, itemID, false)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
