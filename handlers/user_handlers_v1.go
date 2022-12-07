package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aZ4ziL/drug-api/auth"
	"github.com/aZ4ziL/drug-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type userHandlerV1 struct{}

func NewUserHandlerV1() userHandlerV1 {
	return userHandlerV1{}
}

func (u userHandlerV1) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "POST" {
			userRequest := &UserRequest{}

			err := ctx.ShouldBindJSON(userRequest)
			if err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
				return
			}

			validate = validator.New()
			err = validate.Struct(userRequest)
			if err != nil {
				if _, ok := err.(*validator.InvalidValidationError); ok {
					log.Println(err.Error())
					return
				}
				errorMessages := []string{}
				for _, err := range err.(validator.ValidationErrors) {
					errorMessages = append(errorMessages, fmt.Sprintf("error on field: %s, with error type: %s", err.Field(), err.ActualTag()))
				}
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": errorMessages,
				})
				return
			}

			// save
			user := models.User{
				FirstName: userRequest.FirstName,
				LastName:  userRequest.LastName,
				Username:  userRequest.Username,
				Password:  userRequest.Password,
			}
			err = models.NewUserModel().CreateNewUser(&user)
			if err != nil {
				http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
				return
			}

			ctx.JSON(http.StatusOK, user)
		}
	}
}

func (u userHandlerV1) GetToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userLogin := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := ctx.ShouldBindJSON(&userLogin)
		if err != nil {
			log.Println("error 85")
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := models.NewUserModel().GetUserByUsername(userLogin.Username)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "login_failed",
				"message": "Username or password is incorrect",
			})
			return
		}

		if !auth.DecryptionPassword(user.Password, userLogin.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "login_failed",
				"message": "Username or password is incorrect",
			})
			return
		}

		creds := auth.Credential{
			ID:         user.ID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Username:   user.Username,
			IsAdmin:    user.IsAdmin,
			DateJoined: user.DateJoined,
			LastLogin:  user.LastLogin.Time,
		}

		tokenString, err := auth.GetNewToken(creds)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Please don't share your token to anyone.",
			"token":   tokenString,
		})
	}
}
