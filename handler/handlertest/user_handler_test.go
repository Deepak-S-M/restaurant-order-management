package handlertest

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"restaurant-order-management/config"
	"restaurant-order-management/controllers"
	"restaurant-order-management/handler"
	"restaurant-order-management/middlewares"
	"restaurant-order-management/pb"
	"restaurant-order-management/testutil"
	"restaurant-order-management/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetUser(t *testing.T) {
	adminRole := testutil.RandomRole()
	adminRole.Name = "test_admin"
	adminUser := testutil.RandomUser()
	adminUser.Role = adminRole

	err := config.DB.Create(&adminUser).Error
	require.Nil(t, err)

	userRole := testutil.RandomRole()
	userRole.Name = "user"
	normalUser := testutil.RandomUser()
	normalUser.Role = userRole

	err = config.DB.Create(&normalUser).Error
	require.Nil(t, err)

	t.Cleanup(func() {
		config.DB.Unscoped().Delete(&adminUser)
		config.DB.Unscoped().Delete(&adminRole)
		config.DB.Unscoped().Delete(&normalUser)
		config.DB.Unscoped().Delete(&userRole)
	})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &controllers.GrpcUserServer{})

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	t.Cleanup(func() {
		grpcServer.Stop()
		if err := <-serveErrCh; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Errorf("unexpected grpc serve error: %v", err)
		}
	})

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	t.Cleanup(func() {
		conn.Close()
	})

	userHandler := handler.NewUserServiceHandler(conn)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users/:id", middlewares.RoleMiddleware("test_admin"), userHandler.GetUser)
	}

	t.Run("Unauthorized - No Token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		url := "/api/users/" + adminUser.Id
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.Nil(t, err)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("Forbidden - Not Admin", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		url := "/api/users/" + adminUser.Id
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.Nil(t, err)

		token, err := utils.GenerateToken(normalUser.Id, normalUser.Email, normalUser.Role.Name)
		require.Nil(t, err)
		request.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("Success - Admin User", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		url := "/api/users/" + adminUser.Id
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.Nil(t, err)

		token, err := utils.GenerateToken(adminUser.Id, adminUser.Email, adminUser.Role.Name)
		require.Nil(t, err)
		request.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)

		var responseBody map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
		require.Nil(t, err)

		require.Equal(t, "Get User API", responseBody["message"])

		data, ok := responseBody["data"].(map[string]interface{})
		require.True(t, ok)
		require.Equal(t, adminUser.Id, data["id"])
		require.Equal(t, adminUser.Email, data["email"])
		require.Equal(t, adminUser.Name, data["name"])
		require.Equal(t, adminUser.Role.Name, data["role"])
	})
}
