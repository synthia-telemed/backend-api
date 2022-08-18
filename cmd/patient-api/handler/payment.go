package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler/middleware"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"go.uber.org/zap"
	"net/http"
)

type PaymentHandler struct {
	paymentClient    payment.Client
	patientDataStore datastore.PatientDataStore
	logger           *zap.SugaredLogger
}

func NewPaymentHandler(paymentClient payment.Client, pds datastore.PatientDataStore, logger *zap.SugaredLogger) *PaymentHandler {
	return &PaymentHandler{
		paymentClient:    paymentClient,
		patientDataStore: pds,
		logger:           logger,
	}
}

func (h PaymentHandler) Register(r *gin.RouterGroup) {
	paymentGroup := r.Group("/payment", middleware.ParsePatientID)
	paymentGroup.POST("/credit-card", h.AddCreditCard)
	paymentGroup.GET("/credit-card", h.GetCreditCards)

}

type AddCreditCardRequest struct {
	CardToken string `json:"card_token" binding:"required"`
}

// AddCreditCard godoc
// @Summary      Add new credit card
// @Tags         Payment
// @Param 	  	 AddCreditCardRequest body AddCreditCardRequest true "Token from Omise"
// @Success      201
// @Failure      400  {object}  ErrorResponse "Invalid request body"
// @Failure      400  {object}  ErrorResponse "Failed to add credit card"
// @Failure      401  {object}  ErrorResponse "Unauthorized"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card [post]
func (h PaymentHandler) AddCreditCard(c *gin.Context) {
	id, _ := c.Get("patientID")
	patientID := id.(uint)

	var req AddCreditCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	patient, err := h.patientDataStore.FindByID(patientID)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.patientDataStore.FindByID error")
		return
	}

	if patient.PaymentCustomerID == nil {
		cusID, err := h.paymentClient.CreateCustomer(patientID)
		if err != nil {
			InternalServerError(c, h.logger, err, "h.paymentClient.CreateCustomer error")
			return
		}
		patient.PaymentCustomerID = &cusID
		if err := h.patientDataStore.Save(patient); err != nil {
			InternalServerError(c, h.logger, err, "h.patientDataStore.Save error")
			return
		}
	}

	if err := h.paymentClient.AddCreditCard(*patient.PaymentCustomerID, req.CardToken); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrFailedToAddCreditCard)
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

// GetCreditCards godoc
// @Summary      Get lists of saved credit cards
// @Tags         Payment
// @Success      200  {array}   payment.Card  "List of saved cards"
// @Failure      401  {object}  ErrorResponse "Unauthorized"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card [get]
func (h PaymentHandler) GetCreditCards(c *gin.Context) {
	id, _ := c.Get("patientID")
	patientID := id.(uint)

	patient, err := h.patientDataStore.FindByID(patientID)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.patientDataStore.FindByID error")
		return
	}
	if patient.PaymentCustomerID == nil {
		c.JSON(http.StatusOK, []payment.Card{})
		return
	}

	cards, err := h.paymentClient.ListCards(*patient.PaymentCustomerID)
	if err != nil {
		InternalServerError(c, h.logger, err, "h.paymentClient.ListCards error")
		return
	}
	c.JSON(http.StatusOK, cards)
}
