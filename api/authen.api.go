package api

import (
	"main/db"
	"main/interceptor"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthenAPI(router *gin.Engine) {
	authenAPI := router.Group("/api/v2")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
	}
}

// function - Login
func login(c *gin.Context) {
	var user model.User

	if c.ShouldBind(&user) == nil {
		var queryUser model.User
		if err := db.GetDB().First(&queryUser, "username = ?", user.Username).Error; err != nil {
			c.JSON(200, gin.H{"result": "nok", "error": err})
		} else if checkPasswordHash(user.Password, queryUser.Password) == false {
			c.JSON(200, gin.H{"result": "nok", "error": "invalid password", "result2": "LoginSuccess"})
		} else {
			token := interceptor.JwtSign(queryUser)
			c.JSON(200, gin.H{"result": "ok", "token": token, "result2": "LoginSuccess"})
		}

	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

// function - RegisterUse
func register(c *gin.Context) {
	var user model.User
	if c.ShouldBind(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		user.CreatedAt = time.Now()
		if err := db.GetDB().Create(&user).Error; err != nil {
			c.JSON(200, gin.H{"result": "nok", "error": err, "result2": "ข้อมูลซ้ำ(username/password)"})
		} else {
			c.JSON(200, gin.H{"result": "ok", "data": user, "result2": "บันทึกสำเร็จ"})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

// function - CheckPass & hashPasword
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// function - CheckPass & hashPasword
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
