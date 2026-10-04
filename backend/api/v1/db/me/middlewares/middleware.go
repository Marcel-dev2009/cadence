package middlewares

import (
"net/http"
"time"
"github.com/gin-gonic/gin"
"github.com/patrickmn/go-cache"
)
var SessionCache = cache.New(2*time.Hour, 10*time.Minute)
func ReadAuth(c *gin.Context) {
  cookie, err := c.Cookie("session_token")
  if err != nil {
  c.JSON(http.StatusUnauthorized, gin.H{"error":"session not found"})	
  c.Abort()
  return 
  }	
  userID, found := SessionCache.Get(cookie)
  if !found {
  c.JSON(http.StatusUnauthorized, gin.H{"error":"sesssion expired, please login again"})	
  c.Abort()
  return 
  } 
  
  c.Set("userID", userID)
  c.Next()
}