package service

import (
	"net/http"
	"strconv"

	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/gin-gonic/gin"
)

type userSvc struct{}

func (s *userSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	userId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be integer"})
		return
	}

	// get user from db
	user, err := db.GetUser(int32(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *userSvc) Update(c *gin.Context) {
	stringId := c.Param("id")
	userId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be integer"})
		return
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := db.UpdateUser(int32(userId), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
