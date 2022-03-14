package usecase

import (
	"go-midtrans/domain"
	"log"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type midtransTransactionUsecase struct {
	midtransTransactionRepository domain.MidtransTrasactionRepository
	ca                            coreapi.Client
}

func NewMidtransTransactionUsecase(mtr domain.MidtransTrasactionRepository) *midtransTransactionUsecase {
	return &midtransTransactionUsecase{midtransTransactionRepository: mtr, ca: coreapi.Client{}}
}

func (u *midtransTransactionUsecase) Handler(orderId string) {
	u.ca.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	transactionStatusResp, e := u.ca.CheckTransaction(orderId)
	if e != nil {
		log.Println(e)
	}

	mTrx, err := u.midtransTransactionRepository.GetByOrderId(orderId)
	if err != nil {
		log.Println(err.Error())
	}

	if transactionStatusResp != nil {
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				mTrx.Status = "challange"
			} else if transactionStatusResp.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				mTrx.Status = "success"
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			mTrx.Status = "success"
		} else if transactionStatusResp.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			mTrx.Status = "deny"
		} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			mTrx.Status = "failure"
		} else if transactionStatusResp.TransactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			mTrx.Status = "pending"
		}
	}

	_, err = u.midtransTransactionRepository.Update(mTrx)
	if err != nil {
		log.Println(err.Error())
	}
}
