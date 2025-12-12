package handlers

import (
	"net/http"
	"strconv"

	"vstats/internal/cloud/database"
	"vstats/internal/cloud/middleware"

	"github.com/gin-gonic/gin"
)

// Admin users for authorization check
var adminUsers = map[string]bool{
	"zsai001": true, // Add admin usernames here
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *gin.Context) bool {
	username := middleware.GetUsername(c)
	return adminUsers[username]
}

// AdminRequired middleware ensures only admins can access the route
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ============================================================================
// User Management Handlers (Admin Only)
// ============================================================================

// ListUsers returns a paginated list of all users
func ListUsers(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	users, total, err := database.ListAllUsers(c.Request.Context(), pageSize, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      users,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}

// GetUser returns a single user by ID
func GetUser(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	userID := c.Param("id")
	user, err := database.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Plan   string `json:"plan,omitempty"`
	Status string `json:"status,omitempty"`
}

// UpdateUser updates a user's plan or status
func UpdateUser(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	userID := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()

	// Update plan if provided
	if req.Plan != "" {
		validPlans := map[string]bool{"free": true, "pro": true, "enterprise": true}
		if !validPlans[req.Plan] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan. Must be 'free', 'pro', or 'enterprise'"})
			return
		}
		if err := database.UpdateUserPlan(ctx, userID, req.Plan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user plan"})
			return
		}
	}

	// Update status if provided
	if req.Status != "" {
		validStatuses := map[string]bool{"active": true, "suspended": true}
		if !validStatuses[req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'active' or 'suspended'"})
			return
		}
		if err := database.UpdateUserStatus(ctx, userID, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
			return
		}
	}

	// Fetch updated user
	user, err := database.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "message": "User updated successfully"})
}

// DeleteUser soft-deletes a user
func DeleteUser(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	userID := c.Param("id")

	// Don't allow deleting yourself
	currentUserID := middleware.GetUserID(c)
	if userID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	if err := database.DeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserStats returns overall user statistics
func GetUserStats(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	stats, err := database.GetUserStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

