package main

import (
	"fmt"
	_authHttpDelivery "go-midtrans/auth/delivery/http"
	_authRepository "go-midtrans/auth/repository"
	_authUsecase "go-midtrans/auth/usecase"
	"go-midtrans/domain"
	_itemRepository "go-midtrans/item/repository"
	_midtransTransactionHttpDelivery "go-midtrans/midtrans_transaction/delivery/http"
	_midtransTransactionRepository "go-midtrans/midtrans_transaction/repository"
	_midtransTransactionUsecase "go-midtrans/midtrans_transaction/usecase"
	_transactionHttpDelivery "go-midtrans/transaction/delivery/http"
	_transactionRepository "go-midtrans/transaction/repository"
	_transactionUsecase "go-midtrans/transaction/usecase"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	db, err := initDb()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	api := r.Group("/api")

	authRepository := _authRepository.NewUserRepository(db)
	authUsecase := _authUsecase.NewUserRepository(authRepository)
	middleware := _authHttpDelivery.NewMiddleware(authUsecase)
	_authHttpDelivery.NewAuthHandler(api, authUsecase)

	itemRepository := _itemRepository.NewItemRepository(db)

	midtransTransactionRepository := _midtransTransactionRepository.NewMidtransTransactionRepository(db)
	midtransTransactionUsecase := _midtransTransactionUsecase.NewMidtransTransactionUsecase(midtransTransactionRepository)
	_midtransTransactionHttpDelivery.NewMidtransTransactionHandler(api, midtransTransactionUsecase)

	transactionRepository := _transactionRepository.NewTransactionReposity(db)
	transactionUsecase := _transactionUsecase.NewTransactionUsecase(transactionRepository, itemRepository, midtransTransactionRepository)
	_transactionHttpDelivery.NewTransactionHandler(api, transactionUsecase, middleware)

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello from api",
		})
	})

	r.Run()
}

func initDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	if err := DB.AutoMigrate(&domain.MidtransTrasaction{}, &domain.Item{}, &domain.User{}, &domain.Transaction{}); err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return DB, nil
}
