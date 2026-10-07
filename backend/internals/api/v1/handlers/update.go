package handlers

import (
"net/http"
"github.com/Marcel-dev2009/cadence/database/config"
"github.com/Marcel-dev2009/cadence/repository"
"github.com/gin-gonic/gin"
)
func UpdateSyncStatus(c *gin.Context){
 db := config.DB	
 userID := c.GetString("userID")
 if userID == ""{
 c.JSON(http.StatusUnauthorized, gin.H{"error":"user id not found"})	
 c.Abort()
 return
 }	
 userRepo := repository.NewUserRepository(db)
 if err := userRepo.SetgoogleSyncStatus(userID, true); err != nil{
 c.JSON(http.StatusBadRequest, gin.H{"error":"failed to update google sync status"})
 return	
 }
 c.JSON(http.StatusOK, gin.H{
"message": "Google Calendar synchronized succesfully",
})
}