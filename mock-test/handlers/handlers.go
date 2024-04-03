package handlers

import (
	"EXAM3/api-gateway/api/model"
	pb "EXAM3/api-gateway/genproto/user_service"
	"EXAM3/api-gateway/pkg/logger"
	"EXAM3/api-gateway/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

type MockServiceManager struct {
	s   services.IServiceManager
	log logger.Logger
}

const (
	ErrorCodeInvalidURL          = "INVALID_URL"
	ErrorCodeInvalidJSON         = "INVALID_JSON"
	ErrorCodeInvalidParams       = "INVALID_PARAMS"
	ErrorCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrorCodeUnauthorized        = "UNAUTHORIZED"
	ErrorCodeAlreadyExists       = "ALREADY_EXISTS"
	ErrorCodeNotFound            = "NOT_FOUND"
	ErrorCodeInvalidCode         = "INVALID_CODE"
	ErrorBadRequest              = "BAD_REQUEST"
	ErrorInvalidCredentials      = "INVALID_CREDENTIALS"
)

func NewMockServiceManager(s services.IServiceManager, log logger.Logger) *MockServiceManager {
	return &MockServiceManager{s: s, log: log}
}

func (s *MockServiceManager) CreateUser(c *gin.Context) {
	var (
		body        model.RegisterUserRequest
		jsbpMarshal protojson.MarshalOptions
	)
	jsbpMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		s.log.Error(err.Error())
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	exists, err := s.s.UserService().CheckField(context.Background(), &pb.CheckFieldRequest{
		Field: "email",
		Data:  body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		s.log.Error("failed to check uniqueness: ", logger.Error(err))
		return
	}

	if exists.Status {
		c.JSON(http.StatusConflict, model.ResponseError{
			Code:  ErrorCodeAlreadyExists,
			Error: "email is already exist",
		})
		s.log.Error("email is already exist in database")
		return
	}

	id := uuid.New().String()

	// s.jwtHandler = tokens.JWTHandler{
	// 	Sub:       id,
	// 	Role:      "user",
	// 	SignInKey: s.cfg.SigningKey,
	// 	Log:       s.log,
	// 	Timeout:   s.cfg.AccessTokenTimeout,
	// }
	// access, refresh, err := s.jwtHandler.GenerateAuthJWT()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, model.ResponseError{
	// 		Code:  ErrorBadRequest,
	// 		Error: err.Error(),
	// 	})
	// 	s.log.Error("cannot create access token", logger.Error(err))
	// 	return
	// }

	resp, err := s.s.UserService().CreateUser(context.Background(), &pb.User{
		Id:       id,
		Name:     body.Name,
		Age:      int64(body.Age),
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
		// RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, model.UserModel{
		// AccessToken: access,
		Id:       resp.Id,
		Name:     resp.Name,
		Age:      int(resp.Age),
		Username: resp.Username,
		Email:    resp.Email,
		Password: resp.Password,
	})
}

func (s *MockServiceManager) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user pb.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}
	user.Id = id

	result, err := s.s.UserService().UpdateUserById(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, result)
}

func (s *MockServiceManager) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := s.s.UserService().DeleteUser(context.Background(), &pb.UserId{
		UserId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}

func (s *MockServiceManager) GetAllUsers(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page := c.Param("page")
	fmt.Println(page)
	pageToInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		s.log.Error("cannot parse page query param", logger.Error(err))
		return
	}

	limit := c.Param("limit")
	fmt.Println(limit)
	LimitToInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		s.log.Error("cannot parse limit query param", logger.Error(err))
		return
	}

	response, err := s.s.UserService().ListUser(context.Background(), &pb.GetAllUserRequest{
		Page:  int64(pageToInt),
		Limit: int64(LimitToInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		s.log.Error("cannot get all users", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
