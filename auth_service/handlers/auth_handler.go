package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krishnachowdaryvanam/authboard/auth_service/client"
	"github.com/krishnachowdaryvanam/authboard/auth_service/utils"
	"github.com/krishnachowdaryvanam/authboard/proto/eventpb"
	"github.com/krishnachowdaryvanam/authboard/rbac_service/rbacpb"
	"github.com/krishnachowdaryvanam/authboard/tenant_service/tenantpb"
	"github.com/krishnachowdaryvanam/authboard/user_service/userspb"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	TenantId string `json:"tenant_id"`
}

func SignUp(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userClient := client.GetUserClient()
	createUserResp, err := userClient.CreateUser(context.Background(), &userspb.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		TenantId: req.TenantId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	tenantClient := client.GetTenantClient()
	_, err = tenantClient.GetTenant(context.Background(), &tenantpb.GetTenantRequest{Id: req.TenantId})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
		return
	}

	rbacClient := client.GetRBACClient()
	_, err = rbacClient.AssignRole(context.Background(), &rbacpb.AssignRoleRequest{
		UserId: createUserResp.Id,
		Role:   req.Role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	// Publish a user created event to Kafka
	eventClient := client.GetEventClient() // Assume client has a method to get the Event service client
	eventResp, err := eventClient.PublishEvent(context.Background(), &eventpb.EventRequest{
		EventType: "user_created",
		Payload: fmt.Sprintf(`{"userId": "%s", "email": "%s", "role": "%s", "tenantId": "%s"}`,
			createUserResp.Id, createUserResp.Email, req.Role, req.TenantId),
	})
	if err != nil || !eventResp.Success {
		log.Printf("Error publishing user created event: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": createUserResp.Id,
		"email":  createUserResp.Email,
	})
}

func Login(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userClient := client.GetUserClient()
	userResp, err := userClient.GetUserByEmail(context.Background(), &userspb.GetUserByEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, userResp.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(userResp.Id, userResp.TenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
