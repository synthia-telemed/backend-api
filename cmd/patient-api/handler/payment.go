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
}

type AddCreditCardRequest struct {
	CardToken string `json:"card_token" binding:"required"`
}

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

	c.Status(http.StatusCreated)
}
