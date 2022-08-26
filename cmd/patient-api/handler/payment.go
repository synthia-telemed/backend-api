package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var (
	ErrTokenHasBeenUsed               = server.NewErrorResponse("Token has been used")
	ErrLimitNumberOfCreditCardReached = server.NewErrorResponse("Limited number of credit cards is reached")
	ErrInvalidCreditCardID            = server.NewErrorResponse("Invalid credit card ID")
	ErrCreditCardOwnership            = server.NewErrorResponse("Patient doesn't own the specified credit card")
	ErrCreditCardNotFound             = server.NewErrorResponse("Credit card not found")
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
	paymentGroup.POST("/credit-card", h.CreateOrParseCustomer, h.AddCreditCard)
	paymentGroup.GET("/credit-card", h.GetCreditCards)
	paymentGroup.DELETE("/credit-card/:cardID", h.CreateOrParseCustomer, h.VerifyCreditCardOwnership, h.DeleteCreditCard)
}

type AddCreditCardRequest struct {
	Name      string `json:"name" binding:"required"`
	CardToken string `json:"card_token" binding:"required"`
}

// AddCreditCard godoc
// @Summary      Add new credit card
// @Tags         Payment
// @Param 	  	 AddCreditCardRequest body AddCreditCardRequest true "Token from Omise and name of credit card"
// @Success      201
// @Failure      400  {object}  server.ErrorResponse "Invalid request body"
// @Failure      400  {object}  server.ErrorResponse "Token has been used"
// @Failure      400  {object}  server.ErrorResponse "Limited number of credit cards is reached"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card [post]
func (h PaymentHandler) AddCreditCard(c *gin.Context) {
	patientID := h.GetUserID(c)
	customerID := h.GetCustomerID(c)

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

	card, err := h.paymentClient.AddCreditCard(customerID, req.CardToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrTokenHasBeenUsed)
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
// @Success      200  {array}   datastore.CreditCard "List of saved cards"
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

func (h PaymentHandler) DeleteCreditCard(c *gin.Context) {
	customerID := h.GetCustomerID(c)
	rawCard, ok := c.Get("CreditCard")
	if !ok {
		h.InternalServerError(c, errors.New("CreditCard not found"), "c.Get(\"CreditCard\") error")
		return
	}
	creditCard, ok := rawCard.(*datastore.CreditCard)
	if !ok {
		h.InternalServerError(c, errors.New("CreditCard type casting error"), "rawCard.(*datastore.CreditCard)")
		return
	}

	if err := h.creditCardDataStore.Delete(creditCard.ID); err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.Delete error")
		return
	}
	if err := h.paymentClient.RemoveCreditCard(customerID, creditCard.CardID); err != nil {
		h.InternalServerError(c, err, "h.paymentClient.RemoveCreditCard error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h PaymentHandler) CreateOrParseCustomer(c *gin.Context) {
	patientID := h.GetUserID(c)
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
	c.Set("CustomerID", *patient.PaymentCustomerID)
}

func (h PaymentHandler) GetCustomerID(c *gin.Context) string {
	cusID, _ := c.Get("CustomerID")
	return cusID.(string)
}

func (h PaymentHandler) VerifyCreditCardOwnership(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("cardID"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidCreditCardID)
		return
	}
	patientID := h.GetUserID(c)
	card, err := h.creditCardDataStore.FindByID(uint(cardID))
	if err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.FindByID error")
		return
	}
	if card == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrCreditCardNotFound)
		return
	}
	if card.PatientID != patientID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrCreditCardOwnership)
		return
	}
	c.Set("CreditCard", card)
}
