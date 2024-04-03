package tests

import (
	"EXAM3/api-gateway/api_test/handlers"
	"EXAM3/api-gateway/api_test/storage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorCaseApi(t *testing.T) {
	require.NoError(t, SetupMinimumInstance(""))
	buffer, err := OpenFile("usr.json")
	require.Error(t, err)

	// User Create
	req := NewRequest(http.MethodPost, "/users/create", buffer)
	res := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/user/create", handlers.CreateUser)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)

	var user storage.User
	require.Error(t, json.Unmarshal(res.Body.Bytes(), &user))
	require.NotEqual(t, user.Email, "test@gmail.com")
	require.NotEqual(t, user.FirstName, "testname")
	require.NotEqual(t, user.Password, "testpass")

	require.Equal(t, "", user.Id)

	// User Get
	getReq := NewRequest(http.MethodGet, "/users/get", buffer)
	q := getReq.URL.Query()
	q.Add("id", "")
	getReq.URL.RawQuery = q.Encode()
	getRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/user/get", handlers.GetUser)
	r.ServeHTTP(getRes, getReq)
	require.Equal(t, http.StatusNotFound, getRes.Code)
	var getUserResp storage.User
	require.Error(t, json.Unmarshal(getRes.Body.Bytes(), &getUserResp))
	assert.Empty(t, getUserResp.Id)
	assert.Empty(t, getUserResp.FirstName)
	assert.Empty(t, getUserResp.Email)

	// User List
	listReq := NewRequest(http.MethodGet, "/users", buffer)
	listRes := httptest.NewRecorder()

	r.GET("/usrs", handlers.ListUsers)
	r.ServeHTTP(listRes, listReq)
	assert.NotEqual(t, http.StatusOK, listRes.Code)
	bodyBytes, err := io.ReadAll(listRes.Body)
	assert.NoError(t, err)
	assert.NotNil(t, bodyBytes)

	// User Delete
	delReq := NewRequest(http.MethodDelete, "/users/delete", buffer)
	q = delReq.URL.Query()
	q.Add("id", "")
	delReq.URL.RawQuery = q.Encode()
	delRes := httptest.NewRecorder()
	r.DELETE("/usrs/delete", handlers.DeleteUser)
	r.ServeHTTP(delRes, delReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, delRes.Code)
	var message storage.Message
	bodyBytes, err = io.ReadAll(delRes.Body)
	require.NoError(t, err)
	require.Error(t, json.Unmarshal(bodyBytes, &message))
	require.NotEqual(t, "user was deleted successfully", message.Message)

	// // User Register
	regReq := NewRequest(http.MethodPost, "/users/register", buffer)
	regRes := httptest.NewRecorder()
	r.POST("/usrs/register", handlers.RegisterUser)
	r.ServeHTTP(regRes, regReq)
	assert.Equal(t, http.StatusNotFound, regRes.Code)
	var resp storage.Message
	bodyBytes, err = io.ReadAll(regRes.Body)
	require.NoError(t, err)
	require.Error(t, json.Unmarshal(bodyBytes, &resp))
	require.Empty(t, resp.Message)
	require.NotEqual(t, "a verification code was sent to your email, please check it.", resp.Message)

	// User Verify
	// uri := fmt.Sprintf("/users/verify/%s", "12345")
	// verReq := NewRequest(http.MethodGet, uri, buffer)
	// verRes := httptest.NewRecorder()
	// r = gin.Default()
	// r.GET("/users/verify/:code", handlers.Verify)
	// r.ServeHTTP(verRes, verReq)
	// assert.Equal(t, http.StatusOK, verRes.Code)
	// var response *storage.Message
	// bodyBytes, err = io.ReadAll(verRes.Body)
	// require.NoError(t, err)
	// require.NoError(t, json.Unmarshal(bodyBytes, &response))
	// require.Equal(t, "Correct code", response.Message)

	//User Verify with incorrect code
	incorrectURI := fmt.Sprintf("/users/verify/%s", "11111")
	incorrectVerReq := NewRequest(http.MethodGet, incorrectURI, buffer)
	incorrectVerRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users/verify/:code", handlers.Verify)
	r.ServeHTTP(incorrectVerRes, incorrectVerReq)
	assert.Equal(t, http.StatusBadRequest, incorrectVerRes.Code)
	var incorrectResponse storage.Message
	bodyBytes, err = io.ReadAll(incorrectVerRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &incorrectResponse))
	require.Equal(t, "Incorrect code", incorrectResponse.Message)

	gin.SetMode(gin.TestMode)
	require.NoError(t, SetupMinimumInstance(""))
	buffer, err = OpenFile("prod.json")
	require.Error(t, err)

	// product Create
	req = NewRequest(http.MethodPost, "/product/create", buffer)
	// fmt.Println(string(buffer))
	res = httptest.NewRecorder()
	r.POST("/product/create", handlers.CreateProduct)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	var product storage.Product
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &product))
	require.NotEqual(t, "Test Name", product.Name)
	require.NotEqual(t, "Test Description", product.Description)

	// Product Get
	getReq = NewRequest(http.MethodGet, "/product/get", buffer)
	q = getReq.URL.Query()
	q.Add("id", "")
	getReq.URL.RawQuery = q.Encode()
	getRes = httptest.NewRecorder()
	r = gin.Default()
	r.GET("/product/get", handlers.GetProduct)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusNotFound, getRes.Code)

	var productR storage.Product
	bodyBytes, err = io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &productR))
	// pp.Println(string(bodyBytes))
	require.Empty(t, productR.Name)
	require.Empty(t, productR.Description)

	// Product List
	listReq = NewRequest(http.MethodGet, "/products", buffer)
	listRes = httptest.NewRecorder()
	r = gin.Default()
	r.GET("/prodcts", handlers.ListProducts)
	r.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusNotFound, listRes.Code)
	bodyBytes, err = io.ReadAll(listRes.Body)
	assert.NoError(t, err)
	assert.NotNil(t, bodyBytes)
	// pp.Println(string(bodyBytes))

	// Product Delete
	delReq = NewRequest(http.MethodDelete, "/products/delete", buffer)
	q = delReq.URL.Query()
	q.Add("id", "")
	delReq.URL.RawQuery = q.Encode()
	r.DELETE("/products/delete", handlers.DeleteProduct)
	r.ServeHTTP(delRes, delReq)
	assert.Equal(t, http.StatusNotFound, delRes.Code)
	var postMessage storage.Message
	bodyBytes, err = io.ReadAll(delRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &postMessage))

}
