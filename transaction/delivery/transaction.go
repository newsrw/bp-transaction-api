package handler

import (
	"bp-transaction-api/configs"
	"bp-transaction-api/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// TransactionHandler http handler interface for transaction
type TransactionHandler struct {
	config   *configs.Config
	logger   *zap.Logger
	TUsecase domain.TransactionUsecase
}

// NewTransactionHandler init transactions endpoints
func NewTransactionHandler(r *gin.Engine, config *configs.Config, uc domain.TransactionUsecase) {
	handler := &TransactionHandler{
		config:   config,
		TUsecase: uc,
	}

	transactionGroup := r.Group("/transactions")
	{
		transactionGroup.POST("/broadcast", handler.BroadcastTransaction)
	}
}

func isRequestValid(m *domain.BroadcastTransactionRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// BroadcastTransaction is an exported handler function for broadcasting transaction
func (h *TransactionHandler) BroadcastTransaction(c *gin.Context) {
	var req domain.BroadcastTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ok, err := isRequestValid(&req); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	broadCastRsp, err := h.TUsecase.BroadcastTransaction(ctx, &req)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	checkStatusReq := domain.MonitoringTransactionStatusRequest{
		BroadcastTransactionResponse: broadCastRsp,
	}
	checkStatusRsp, err := h.TUsecase.MonitoringStatus(ctx, &checkStatusReq)
	if err != nil {
		c.JSON(h.getStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, checkStatusRsp)
}

func (h *TransactionHandler) getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
