package handler

import (
	"context"
	"errors"
	"net/http"
	"wallet-service/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

func (h *Handler) UpdateWallet(c *gin.Context) {
	var in UpdateWalletRequest

	if err := c.BindJSON(&in); err != nil {
		log.Error(err)
		return
	}

	parseID, err := uuid.Parse(in.WalletID)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorResponse{Message: ErrInvalidFormatID.Error()})
		return
	}

	var serviceCall func(ctx context.Context, id uuid.UUID, amount int64) (*domain.Wallet, error)

	switch in.OperationType {
	case "DEPOSIT":
		serviceCall = h.services.Deposit
	case "WITHDRAW":
		serviceCall = h.services.Withdraw
	}

	wallet, err := serviceCall(c, parseID, in.Amount)
	if err != nil {
		log.Error(err)

		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrInsufficientBalance):
			c.AbortWithStatus(http.StatusConflict)
			return
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, &UpdateWalletResponse{
		WalletID:   wallet.ID().String(),
		NewBalance: wallet.Balance(),
	})
}

func (h *Handler) GetWallet(c *gin.Context) {
	walletID := c.Param("id")
	if walletID == "" {
		log.Error(ErrPathParameterID)
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorResponse{Message: ErrPathParameterID.Error()})
		return
	}

	parseID, err := uuid.Parse(walletID)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, &ErrorResponse{Message: ErrInvalidFormatID.Error()})
		return
	}

	wallet, err := h.services.Wallet.Get(c, parseID)
	if err != nil {
		log.Error(err)
		if errors.Is(err, domain.ErrWalletNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &GetWalletResponse{
		WalletID: wallet.ID().String(),
		Balance:  wallet.Balance(),
	})
}
