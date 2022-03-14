package http

import (
	"encoding/json"
	"go-midtrans/domain"
	"log"

	"github.com/gin-gonic/gin"
)

type MidtransTransactionHandler struct {
	midtransTransactionUsecase domain.MidtransTransactionUsecase
}

func NewMidtransTransactionHandler(r *gin.RouterGroup, mtu domain.MidtransTransactionUsecase) {
	handler := MidtransTransactionHandler{midtransTransactionUsecase: mtu}
	r.POST("/handler/notification", handler.Handler)
}

func (h *MidtransTransactionHandler) Handler(c *gin.Context) {
	var notifPayload map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&notifPayload)
	if err != nil {
		log.Println("error 1")
		return
	}

	orderId, exist := notifPayload["order_id"].(string)
	if !exist {
		log.Println("error 2")
		return
	}

	h.midtransTransactionUsecase.Handler(orderId)
}
