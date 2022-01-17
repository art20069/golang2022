package api

import (
	"main/db"
	"main/interceptor"
	"main/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupTransactionAPI(router *gin.Engine) {
	TransactionAPI := router.Group("/api/v2")
	{
		TransactionAPI.GET("/Transaction", getTransaction)
		TransactionAPI.POST("/Transaction", interceptor.JwtVerify, createTransaction)
	}
}

type TransactionResult struct {
	ID            uint
	Total         float64
	Paid          float64
	Change        float64
	PaymentType   string
	PaymentDetail string
	OrderList     string
	Staff         string
	CreatedAt     time.Time
}

func getTransaction(c *gin.Context) {
	// var transactions []model.Transaction
	// db.GetDB().Find(&transactions)
	// c.JSON(200, transactions)
	var result []TransactionResult
	db.GetDB().Debug().Raw("SELECT transactions.id, total, paid, change, payment_type, payment_detail, order_list, users.username as Staff, transactions.created_at FROM transactions join users on transactions.staff_id = users.id", nil).Scan(&result)
	c.JSON(200, result)
}

func createTransaction(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{"result": "createTransaction"})
	var transaction model.Transaction
	if err := c.ShouldBind(&transaction); err == nil {
		transaction.StaffID = c.GetString("jwt_staff_id")
		transaction.CreatedAt = time.Now()
		db.GetDB().Create(&transaction)
		c.JSON(http.StatusOK, gin.H{"result": "ok", "data": transaction})
	} else {
		c.JSON(404, gin.H{"result": "nok"})
	}
}
