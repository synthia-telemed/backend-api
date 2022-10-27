package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/clock"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
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
	ErrInvalidInvoiceID               = server.NewErrorResponse("Invalid invoice ID")
	ErrInvoiceNotFound                = server.NewErrorResponse("Invoice not found")
	ErrInvoiceOwnership               = server.NewErrorResponse("Patient doesn't down the specified invoice")
	ErrInvoicePaid                    = server.NewErrorResponse("Invoice is already paid")
)

type PaymentHandler struct {
	paymentClient       payment.Client
	patientDataStore    datastore.PatientDataStore
	creditCardDataStore datastore.CreditCardDataStore
	hospitalSysClient   hospital.SystemClient
	paymentDataStore    datastore.PaymentDataStore
	clock               clock.Clock
	PatientGinHandler
}

func NewPaymentHandler(paymentClient payment.Client, pds datastore.PatientDataStore, cds datastore.CreditCardDataStore, hsc hospital.SystemClient, pay datastore.PaymentDataStore, clock clock.Clock, logger *zap.SugaredLogger) *PaymentHandler {
	return &PaymentHandler{
		paymentClient:       paymentClient,
		patientDataStore:    pds,
		creditCardDataStore: cds,
		hospitalSysClient:   hsc,
		paymentDataStore:    pay,
		clock:               clock,
		PatientGinHandler:   NewPatientGinHandler(pds, logger),
	}
}

func (h PaymentHandler) Register(r *gin.RouterGroup) {
	paymentGroup := r.Group("/payment", h.ParseUserID)
	paymentGroup.POST("/credit-card", h.CreateOrParseCustomer, h.AddCreditCard)
	paymentGroup.GET("/credit-card", h.GetCreditCards)
	paymentGroup.PATCH("/credit-card/:cardID", h.VerifyCreditCardOwnership, h.SetCreditCardIsDefault)
	paymentGroup.DELETE("/credit-card/:cardID", h.CreateOrParseCustomer, h.VerifyCreditCardOwnership, h.DeleteCreditCard)
	paymentGroup.POST("/pay/:invoiceID/credit-card/:cardID", h.ParseAndVerifyUnpaidInvoiceOwnership, h.CreateOrParseCustomer, h.VerifyCreditCardOwnership, h.PayInvoiceWithCreditCard)
}

type AddCreditCardRequest struct {
	Name      string `json:"name"`
	CardToken string `json:"card_token" binding:"required"`
	IsDefault bool   `json:"is_default"`
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

	cardCount, err := h.creditCardDataStore.Count(patientID)
	if err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.FindByPatientID error")
		return
	}
	if cardCount >= 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrLimitNumberOfCreditCardReached)
		return
	}

	card, err := h.paymentClient.AddCreditCard(customerID, req.CardToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrTokenHasBeenUsed)
		return
	}

	newCard := &datastore.CreditCard{
		Last4Digits: card.Last4Digits,
		Brand:       card.Brand,
		PatientID:   patientID,
		CardID:      card.ID,
		Name:        req.Name,
		IsDefault:   cardCount == 0 || req.IsDefault,
		Expiry:      card.Expiry,
	}
	if cardCount > 0 && req.IsDefault {
		if err := h.creditCardDataStore.SetAllToNonDefault(patientID); err != nil {
			h.InternalServerError(c, err, "h.creditCardDataStore.SetAllToNonDefault error")
			return
		}
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

type SetCreditCardIsDefaultRequest struct {
	IsDefault bool `json:"is_default"`
}

// SetCreditCardIsDefault godoc
// @Summary      Set isDefault status of credit card
// @Tags         Payment
// @Param  		 cardID 	path	 integer 	true "ID of the credit card to set isDefault"
// @Param 	  	 SetCreditCardIsDefaultRequest body SetCreditCardIsDefaultRequest true "IsDefault status of the credit card"
// @Success      200
// @Failure      400  {object}  server.ErrorResponse "Invalid credit card ID"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse "Patient doesn't own the specified credit card"
// @Failure      404  {object}  server.ErrorResponse "Credit card not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card/{cardID} [patch]
func (h PaymentHandler) SetCreditCardIsDefault(c *gin.Context) {
	var req SetCreditCardIsDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

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

	if req.IsDefault {
		if err := h.creditCardDataStore.SetAllToNonDefault(creditCard.PatientID); err != nil {
			h.InternalServerError(c, err, "h.creditCardDataStore.SetAllToNonDefault error")
			return
		}
	}
	if err := h.creditCardDataStore.SetIsDefault(creditCard.ID, req.IsDefault); err != nil {
		h.InternalServerError(c, err, "h.creditCardDataStore.SetIsDefault error")
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

// DeleteCreditCard godoc
// @Summary      Delete saved credit card
// @Tags         Payment
// @Param  		 cardID 	path	 integer 	true "ID of the credit card to delete"
// @Success      200
// @Failure      400  {object}  server.ErrorResponse "Invalid credit card ID"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse "Patient doesn't own the specified credit card"
// @Failure      404  {object}  server.ErrorResponse "Credit card not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/credit-card/{cardID} [delete]
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

type PayInvoiceWithCreditCardResponse struct {
	*datastore.Payment
	FailureMessage *string `json:"failure_message"`
}

// PayInvoiceWithCreditCard godoc
// @Summary      Pay invoice with credit card method
// @Tags         Payment
// @Param  		 cardID 	path	 integer 	true "ID of the credit card to be charged"
// @Param  		 invoiceID 	path	 integer 	true "ID of the invoice to pay"
// @Success      201  {object}	PayInvoiceWithCreditCardResponse "Payment information"
// @Failure      400  {object}  server.ErrorResponse "Invalid credit card ID or invoice ID"
// @Failure      401  {object}  server.ErrorResponse "Unauthorized"
// @Failure      403  {object}  server.ErrorResponse "Patient doesn't own the specified credit card or invoice"
// @Failure      404  {object}  server.ErrorResponse "Credit card or invoice not found"
// @Failure      500  {object}  server.ErrorResponse "Internal server error"
// @Security     UserID
// @Security     JWSToken
// @Router       /payment/pay/{invoiceID}/credit-card/{cardID} [post]
func (h PaymentHandler) PayInvoiceWithCreditCard(c *gin.Context) {
	customerID := h.GetCustomerID(c)
	rawCard, _ := c.Get("CreditCard")
	creditCard, _ := rawCard.(*datastore.CreditCard)
	rawInvoice, _ := c.Get("Invoice")
	invoice, _ := rawInvoice.(*hospital.InvoiceOverview)

	paymentCharge, err := h.paymentClient.PayWithCreditCard(customerID, creditCard.CardID, fmt.Sprintf("%d", invoice.Id), int(invoice.Total*100))
	if err != nil {
		h.InternalServerError(c, err, "h.paymentClient.PayWithCreditCard error")
		return
	}
	status := datastore.FailedPaymentStatus
	paidAt := h.clock.NowPointer()
	if paymentCharge.Success {
		status = datastore.SuccessPaymentStatus
		if err := h.hospitalSysClient.PaidInvoice(context.Background(), invoice.Id); err != nil {
			h.InternalServerError(c, err, "h.hospitalSysClient.PaidInvoice error")
			return
		}
	}
	p := &datastore.Payment{
		Method:       datastore.CreditCardPaymentMethod,
		Amount:       invoice.Total,
		PaidAt:       paidAt,
		ChargeID:     paymentCharge.ID,
		InvoiceID:    invoice.Id,
		Status:       status,
		CreditCard:   creditCard,
		CreditCardID: &creditCard.ID,
	}
	if err := h.paymentDataStore.Create(p); err != nil {
		h.InternalServerError(c, err, "h.paymentDataStore.Create error")
		return
	}
	res := &PayInvoiceWithCreditCardResponse{Payment: p, FailureMessage: paymentCharge.FailureMessage}
	c.JSON(http.StatusCreated, res)
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

func (h PaymentHandler) ParseAndVerifyUnpaidInvoiceOwnership(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("invoiceID"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidInvoiceID)
		return
	}
	invoice, err := h.hospitalSysClient.FindInvoiceByID(context.Background(), int(invoiceID))
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalSysClient.FindInvoiceByID error")
		return
	}
	if invoice == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrInvoiceNotFound)
		return
	}
	if invoice.Paid {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvoicePaid)
		return
	}
	patient, err := h.patientDataStore.FindByID(h.GetUserID(c))
	if err != nil {
		h.InternalServerError(c, err, " h.patientDataStore.FindByID error")
		return
	}
	if invoice.PatientID != patient.RefID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrInvoiceOwnership)
		return
	}
	c.Set("Invoice", invoice)
}
