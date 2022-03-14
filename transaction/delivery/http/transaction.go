package http

import (
	"go-midtrans/domain"
	"go-midtrans/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
}

func NewTransactionHandler(r *gin.RouterGroup, tu domain.TransactionUsecase, middleware gin.HandlerFunc) {
	handler := &TransactionHandler{transactionUsecase: tu}
	api := r.Group("/transaction")
	{
		api.POST("/new", middleware, handler.Order)
	}
}

type orderInput struct {
	ItemId      int    `binding:"required" json:"item_id"`
	PaymentType string `binding:"required" json:"payment_type"`
}

func (h *TransactionHandler) Order(c *gin.Context) {
	var input orderInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, &utils.Response{Message: err.Error()})
		return
	}

	user := c.MustGet("userLoggedin").(domain.User)

	trx, err := h.transactionUsecase.Order(uint(input.ItemId), input.PaymentType, user)
	if err != nil {
		c.JSON(utils.GetErrCode(err), &utils.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &utils.Response{Data: trx, Message: "successfully created new order"})
}
