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

var (
	ErrInvalidRequestBody = server.NewErrorResponse("Invalid request body")
	ErrInvalidCredential  = server.NewErrorResponse("Invalid credential")
)

type AuthHandler struct {
	hospitalSysClient hospital.SystemClient
	tokenService      token.Service
	doctorDataStore   datastore.DoctorDataStore
	logger            *zap.SugaredLogger
	server.GinHandler
}

func NewAuthHandler(h hospital.SystemClient, t token.Service, ds datastore.DoctorDataStore, l *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		hospitalSysClient: h,
		tokenService:      t,
		doctorDataStore:   ds,
		logger:            l,
		GinHandler:        server.GinHandler{Logger: l},
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

// Signin godoc
// @Summary      Signin doctor with credential
// @Tags         Auth
// @Param 	  	 SigninRequest body SigninRequest true "Username and password of the doctor"
// @Success      201  {object}  SigninResponse 		   "Token is return when authentication is successes"
// @Failure      400  {object}  server.ErrorResponse   "Invalid request body"
// @Failure      401  {object}  server.ErrorResponse   "Provided credential is not in the hospital system"
// @Failure      500  {object}  server.ErrorResponse   "Internal server error"
// @Router       /auth/signin [post]
func (h AuthHandler) Signin(c *gin.Context) {
	var req SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	isCredValid, err := h.hospitalSysClient.AssertDoctorCredential(context.Background(), req.Username, req.Password)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalSysClient.AssertDoctorCredential error")
		return
	}
	if !isCredValid {
		c.JSON(http.StatusUnauthorized, ErrInvalidCredential)
		return
	}

	d, err := h.hospitalSysClient.FindDoctorByUsername(context.Background(), req.Username)
	if err != nil {
		h.InternalServerError(c, err, "h.hospitalSysClient.FindDoctorByUsername error")
		return
	}

	doctor := &datastore.Doctor{RefID: d.Id}
	if err := h.doctorDataStore.FindOrCreate(doctor); err != nil {
		h.InternalServerError(c, err, "h.doctorDataStore.FindOrCreate error")
		return
	}

	jws, err := h.tokenService.GenerateToken(uint64(doctor.ID), "Doctor")
	if err != nil {
		h.InternalServerError(c, err, "h.tokenService.GenerateToken error")
		return
	}
	c.JSON(http.StatusCreated, SigninResponse{Token: jws})
}
