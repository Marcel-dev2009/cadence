package middlewares

import (
	"net/http"
	"github.com/Marcel-dev2009/cadence/db/config"
	"github.com/gin-gonic/gin"
)
func ReadAuth(c *gin.Context) {
  cookie, err := c.Cookie("session_token")
  if err != nil {
  c.JSON(http.StatusUnauthorized, gin.H{"error":"session not found"})	
  c.Abort()
  return 
  }	
  userID, found := config.SessionCache.Get(cookie)
  if !found {
  c.JSON(http.StatusUnauthorized, gin.H{"error":"sesssion expired, please login again"})	
  c.Abort()
  return 
  } 
  
  c.Set("userID", userID)
  c.Next()
}