package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aZ4ziL/drug-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DrugRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Stock uint   `json:"stock" validate:"required,number"`
	Price uint   `json:"price" validate:"required,number"`
}

type drugHandlerV1 struct{}

func NewDrugHandlerV1() drugHandlerV1 {
	return drugHandlerV1{}
}

// Index
// handler for drug index
func (d drugHandlerV1) Index() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if id, ok := ctx.GetQuery("id"); ok {
			idInt, _ := strconv.Atoi(id)
			drug, err := models.NewDrugModel().GetDrugByID(uint(idInt))
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "not_found",
					"message": "the drug by id is not found.",
				})
				return
			}
			ctx.JSON(http.StatusOK, drug)
			return
		}

		drugs := models.NewDrugModel().GetAllDrug()
		ctx.JSON(http.StatusOK, drugs)
	}
}

// Add
// handler for add drug
func (d drugHandlerV1) Add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		drugRequest := &DrugRequest{}

		err := ctx.ShouldBindJSON(&drugRequest)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		// validate the post request
		validate = validator.New()
		err = validate.Struct(drugRequest)
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

		drug := models.Drug{
			ID:    drugRequest.ID,
			Name:  drugRequest.Name,
			Stock: drugRequest.Stock,
			Price: drugRequest.Price,
		}
		err = models.NewDrugModel().CreateNewDrug(&drug)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Successfully to add new drug.",
			"drug":    drug,
		})
	}
}

// Edit
// handler for edit drug
func (d drugHandlerV1) Edit() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		drugRequest := &DrugRequest{}
		err := ctx.ShouldBindJSON(drugRequest)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		drug, err := models.NewDrugModel().GetDrugByID(drugRequest.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "not_found",
				"message": fmt.Sprintf("error: drug by id: %d is not found.", drugRequest.ID),
			})
			return
		}

		// validate the post request
		validate = validator.New()
		err = validate.Struct(drugRequest)
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

		// Edit and save
		drug.ID = drugRequest.ID
		drug.Name = drugRequest.Name
		drug.Stock = drugRequest.Stock
		drug.Price = drugRequest.Price

		err = models.NewDrugModel().UpdateDrug(&drug)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Successfully to update the drug by id: %d", drugRequest.ID),
			"drug":    drug,
		})
	}
}

// Delete
// delete handler for drug
func (d drugHandlerV1) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		drugRequest := struct {
			ID uint `json:"id" validate:"required,number"`
		}{}

		err := ctx.ShouldBindJSON(&drugRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		validate = validator.New()
		err = validate.Struct(&drugRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Please set the id",
			})
			return
		}

		drug, err := models.NewDrugModel().GetDrugByID(drugRequest.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		err = models.NewDrugModel().DeleteDrug(&drug)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Successfully to delete drug by id: %d", drugRequest.ID),
		})
	}
}
