package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"go.uber.org/zap"
)

type InfoHandler struct {
	hospitalClient hospital.SystemClient
	logger         *zap.SugaredLogger
	*server.GinHandler
}

func NewInfoHandler(hospitalClient hospital.SystemClient, logger *zap.SugaredLogger) *InfoHandler {
	return &InfoHandler{
		hospitalClient: hospitalClient,
		logger:         logger,
		GinHandler:     &server.GinHandler{Logger: logger},
	}
}

type GetNameResponse struct {
	EN hospital.Name `json:"EN"`
	TH hospital.Name `json:"TH"`
}

func (h InfoHandler) GetName(c *gin.Context) {

}
