package handler

import (
	"wallet-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		services: service,
	}
}

func (h *Handler) GetRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		v1 := api.Group("v1")
		{
			wallet := v1.Group("/wallet")
			{
				wallet.POST("", h.UpdateWallet)
			}

			wallets := v1.Group("/wallets")
			{
				wallets.GET("/:id", h.GetWallet)
			}
		}
	}

	return r
}
