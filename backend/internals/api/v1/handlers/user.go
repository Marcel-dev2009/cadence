package handlers

import (
"net/http"
"github.com/Marcel-dev2009/cadence/database/config"
"github.com/Marcel-dev2009/cadence/repository"
"github.com/gin-gonic/gin"
)

func GetUserMetaData(c *gin.Context) {
  db := config.DB
  dataCache := config.DataListCache	
  userID := c.GetString("userID")
 if userID == ""{
 c.JSON(http.StatusUnauthorized, gin.H{"error":"User session not found or is expired"}) 
  return
 }
 userRepo := repository.NewUserRepository(db, dataCache)

data, err := userRepo.GetUserData(userID)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
return
}

c.JSON(http.StatusOK, gin.H{"data": data})
}