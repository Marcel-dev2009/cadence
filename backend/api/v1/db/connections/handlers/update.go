package handlers

import (
"net/http"

"github.com/Marcel-dev2009/cadence/db/config"
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
 result := db.Table("users").Where("id = ?", userID).Update("sync_google_calender", true)
 if result.Error != nil{
 c.JSON(http.StatusBadRequest, gin.H{"error":"failed to update google sync status"})
 return	
 }
 c.JSON(http.StatusOK, gin.H{
"message": "Google Calendar sync enabled successfully!",
"sync":    true,
})
}