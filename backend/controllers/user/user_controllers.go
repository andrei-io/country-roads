package user

import (
	"country/dic"
	"country/domain/entity"
	"country/domain/repo/email"
	"country/domain/repo/user"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Gin controller for GET /users/:id
func ReadOne(c *gin.Context) {
	repo := dic.Container.Get(dic.UserRepo).(user.IUserRepo)
	id := c.Param("id")

	u, err := repo.Get(id)
	u.Password = "Secret"
	u.Email = "Secret"

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, entity.NotFoundError)
		return
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, entity.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

// Gin controlle for signup
func SignUp(c *gin.Context) {
	userRepo := dic.Container.Get(dic.UserRepo).(user.IUserRepo)
	emailRepo := dic.Container.Get(dic.EmailRepo).(email.IEmailRepo)
	var json entity.User

	// Validate request form
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{Message: err.Error()})
		return
	}

	json.Role = 1
	json.Verified = false
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: err.Error()})
	}
	json.Password = string(hashedPassword)

	// Create in db. Will error out when invalid
	err = userRepo.Create(&json)
	if err != nil {
		err := err.(*pgconn.PgError)
		message := err.Message
		if strings.Contains(message, "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, entity.NewDuplicateEntityErrorResponse(err.ConstraintName))
			return
		} else {
			c.JSON(http.StatusUnprocessableEntity, entity.Response{Message: err.Error()})
			return
		}
	}

	verficationCode, err := entity.GenerateOTP(6)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Response{Message: err.Error()})
	}

	if viper.GetString("ENV") == "PROD" {
		_, err := emailRepo.Send(json.Email, "Country Roads verification email", verficationCode)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Printf("Verification code for user %s is %s\n", json.Email, verficationCode)
	}

	c.JSON(http.StatusOK, entity.SuccesResponse)

}
