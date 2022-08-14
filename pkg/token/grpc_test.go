package token_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/token"
	pb "github.com/synthia-telemed/backend-api/pkg/token/proto"
	"github.com/synthia-telemed/backend-api/test/mock_token_service"
)

var _ = Describe("Token gRPC Service", func() {

	var (
		mockCtrl        *gomock.Controller
		mockTokenClient *mock_token_service.MockTokenClient
		tokenService    *token.GRPCTokenService
		req             *pb.GenerateTokenRequest
		res             *pb.TokenResponse
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTokenClient = mock_token_service.NewMockTokenClient(mockCtrl)
		tokenService = token.NewGRPCTokenServiceWithClient(mockTokenClient)

		req = &pb.GenerateTokenRequest{
			Role:   "doctor",
			UserID: 99,
		}
		res = &pb.TokenResponse{Token: "signed_token"}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	When("Token generation success", func() {
		BeforeEach(func() {
			mockTokenClient.EXPECT().GenerateToken(gomock.Any(), req).Return(res, nil).Times(1)
		})

		It("should return token", func() {
			token, err := tokenService.GenerateToken(req.UserID, req.Role)
			Expect(err).To(BeNil())
			Expect(token).To(Equal(res.Token))
		})
	})

	When("Token generation fails", func() {
		BeforeEach(func() {
			mockTokenClient.EXPECT().GenerateToken(gomock.Any(), req).Return(nil, errors.New("generation error")).Times(1)
		})

		It("should return error", func() {
			token, err := tokenService.GenerateToken(req.UserID, req.Role)
			Expect(err).ToNot(BeNil())
			Expect(token).To(BeEmpty())
		})
	})
})
