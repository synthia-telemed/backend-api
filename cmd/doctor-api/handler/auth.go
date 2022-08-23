package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"github.com/synthia-telemed/backend-api/pkg/token"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	hospitalSysClient hospital.SystemClient
	logger            *zap.SugaredLogger
	tokenService      token.Service
	doctorDataStore   datastore.DoctorDataStore
}

func NewAuthHandler(h hospital.SystemClient, t token.Service, ds datastore.DoctorDataStore, l *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		hospitalSysClient: h,
		logger:            l,
		tokenService:      t,
		doctorDataStore:   ds,
	}
}

func (h AuthHandler) Register(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	authGroup.POST("/signin", h.Signin)
}

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninResponse struct {
	Token string `json:"token"`
}

func (h AuthHandler) Signin(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, server.ErrInvalidRequestBody)
		return
	}

	isCredValid, err := h.hospitalSysClient.AssertDoctorCredential(context.Background(), req.Username, req.Password)
	if err != nil {
		server.InternalServerError(c, h.logger, err, "h.hospitalSysClient.AssertDoctorCredential error")
		return
	}
	if !isCredValid {
		c.JSON(http.StatusUnauthorized, server.ErrInvalidCredential)
		return
	}

	d, err := h.hospitalSysClient.FindDoctorByUsername(context.Background(), req.Username)
	if err != nil {
		server.InternalServerError(c, h.logger, err, "h.hospitalSysClient.FindDoctorByUsername error")
		return
	}

	doctor := &datastore.Doctor{RefID: d.Id}
	if err := h.doctorDataStore.FindOrCreate(doctor); err != nil {
		server.InternalServerError(c, h.logger, err, "h.doctorDataStore.FindOrCreate error")
		return
	}

	jws, err := h.tokenService.GenerateToken(uint64(doctor.ID), "Doctor")
	if err != nil {
		server.InternalServerError(c, h.logger, err, "h.tokenService.GenerateToken error")
		return
	}
	c.JSON(http.StatusCreated, SigninResponse{Token: jws})
}
