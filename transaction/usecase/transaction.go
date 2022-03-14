package usecase

import (
	"encoding/json"
	"go-midtrans/domain"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type transactionUsecase struct {
	ca                            coreapi.Client
	transactionRepository         domain.TransactionRepository
	itemRepository                domain.ItemRepository
	midtransTransactionRepository domain.MidtransTrasactionRepository
}

func NewTransactionUsecase(tr domain.TransactionRepository, ir domain.ItemRepository, mtr domain.MidtransTrasactionRepository) *transactionUsecase {
	return &transactionUsecase{ca: coreapi.Client{}, transactionRepository: tr, itemRepository: ir, midtransTransactionRepository: mtr}
}

func (u *transactionUsecase) Order(itemId uint, paymentType string, user domain.User) (domain.Transaction, error) {
	item, err := u.itemRepository.GetById(itemId)
	if err != nil {
		return domain.Transaction{}, domain.ErrInternalServer
	}

	if item.ID == 0 {
		return domain.Transaction{}, domain.ErrNotFound
	}

	u.ca.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	orderId := uuid.NewString()
	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: item.Price,
		},
		Items: &[]midtrans.ItemDetails{
			midtrans.ItemDetails{
				ID:    strconv.Itoa(int(item.ID)),
				Price: item.Price,
				Qty:   1,
				Name:  item.Name,
			},
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	if paymentType == "gopay" {
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
	} else if strings.HasPrefix(paymentType, "va-") {
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		payment := strings.Split(paymentType, "-")
		if payment[1] == "bni" {
			chargeReq.BankTransfer = &coreapi.BankTransferDetails{
				Bank: midtrans.BankBni,
			}
		} else {
			return domain.Transaction{}, domain.ErrNotFound
		}
	} else {
		return domain.Transaction{}, domain.ErrNotFound
	}

	coreApiRes, _ := u.ca.ChargeTransaction(chargeReq)

	type paymentDataStruct struct {
		Key string `json:"key"`
		Qr  string `json:"qr"`
	}

	paymentData := &paymentDataStruct{}
	if paymentType == "gopay" {
		paymentData.Key = coreApiRes.Actions[1].URL
		paymentData.Qr = coreApiRes.Actions[0].URL
	} else if paymentType == "va-bni" {
		paymentData.Key = coreApiRes.VaNumbers[0].VANumber
	}

	paymentJson, _ := json.Marshal(paymentData)

	mtTrx := domain.MidtransTrasaction{
		OrderCode:   orderId,
		PaymentType: paymentType,
		Amount:      item.Price,
		Status:      "pending",
		PaymentData: string(paymentJson),
		MidtransId:  coreApiRes.TransactionID,
	}

	newMtTrx, err := u.midtransTransactionRepository.Create(mtTrx)
	if err != nil {
		return domain.Transaction{}, domain.ErrInternalServer
	}

	trx := domain.Transaction{
		ItemID:               item.ID,
		MidtransTrasactionID: newMtTrx.OrderCode,
		UserID:               user.ID,
		Price:                item.Price,
	}

	newTrx, err := u.transactionRepository.Create(trx)
	if err != nil {
		return newTrx, domain.ErrInternalServer
	}

	newTrx.MidtransTrasaction = newMtTrx
	newTrx.Item = item

	return newTrx, nil
}
