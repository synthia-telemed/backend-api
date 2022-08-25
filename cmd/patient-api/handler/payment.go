package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrFailedToAddCreditCard          = server.NewErrorResponse("Failed to add credit card")
	ErrLimitNumberOfCreditCardReached = server.NewErrorResponse("Limited number of credit cards reached")
)

type PaymentHandler struct {
	paymentClient       payment.Client
	patientDataStore    datastore.PatientDataStore
	creditCardDataStore datastore.CreditCardDataStore
	logger              *zap.SugaredLogger
	server.GinHandler
}

func NewPaymentHandler(paymentClient payment.Client, pds datastore.PatientDataStore, cds datastore.CreditCardDataStore, logger *zap.SugaredLogger) *PaymentHandler {
	return &PaymentHandler{
		paymentClient:       paymentClient,
		patientDataStore:    pds,
		logger:              logger,
		creditCardDataStore: cds,
		GinHandler:          server.GinHandler{Logger: logger},
	}
}

func (h PaymentHandler) Register(r *gin.RouterGroup) {
	paymentGroup := r.Group("/payment", middleware.ParseUserID)
	paymentGroup.POST("/credit-card", h.AddCreditCard)
	paymentGroup.GET("/credit-card", h.GetCreditCards)

}

type AddCreditCardRequest struct {
	Name      string `json:"name" binding:"required"`
	CardToken string `json:"card_token" binding:"required"`
}

// AddCreditCard godoc
// @Summary      Add new credit card
// @Tags         Payment
// @Param 	  	 AddCreditCardRequest body AddCreditCardRequest true "Token and fingerprint from Omise"
// @Success      201
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      400  {object}  server.ErrorResponse "Failed to add credit card"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card [post]
func (h PaymentHandler) AddCreditCard(c *gin.Context) {
	patientID := h.GetUserID(c)

	var req AddCreditCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	cards, err := h.creditCardDataStore.FindByPatientID(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.FindByPatientID error")
		return
	}
	if len(cards) >= 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrLimitNumberOfCreditCardReached)
		return
	}

	patient, err := h.patientDataStore.FindByID(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.patientDataStore.FindByID error")
		return
	}

	if patient.PaymentCustomerID == nil {
		cusID, err := h.paymentClient.CreateCustomer(patientID)
		if err != nil {
			h.InternalServerError(c, err, "h.paymentClient.CreateCustomer error")
			return
		}
		patient.PaymentCustomerID = &cusID
		if err := h.patientDataStore.Save(patient); err != nil {
			h.InternalServerError(c, err, "h.patientDataStore.Save error")
			return
		}
	}

	card, err := h.paymentClient.AddCreditCard(*patient.PaymentCustomerID, req.CardToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrFailedToAddCreditCard)
		return
	}

	newCard := &datastore.CreditCard{
		IsDefault:   len(cards) == 0,
		Last4Digits: card.Last4Digits,
		Brand:       card.Brand,
		PatientID:   patientID,
		CardID:      card.ID,
		Name:        req.Name,
	}
	if err := h.creditCardDataStore.Create(newCard); err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.Create error")
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

// GetCreditCards godoc
// @Summary      Get lists of saved credit cards
// @Tags         Payment
// @Success      200  {array}   payment.Card  "List of saved cards"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card [get]
func (h PaymentHandler) GetCreditCards(c *gin.Context) {
	patientID := h.GetUserID(c)
	cards, err := h.creditCardDataStore.FindByPatientID(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.paymentClient.ListCards error")
		return
	}
	c.JSON(http.StatusOK, cards)
}
