package routes

import (
	"net/http"

	"github.com/Marcel-dev2009/cadence/api/v1/handlers"
	"github.com/Marcel-dev2009/cadence/database/me/middlewares"
	"github.com/gin-gonic/gin"
)
func Setup(r *gin.Engine){
 api := r.Group("/api/v1")
 auth := api.Group("/auth")
 {
  auth.POST("/register", handlers.SignUp)
  auth.POST("/login", handlers.SignIn)	
  auth.POST("/sign-out", handlers.SignOut)	
 }
 pr_route := api.Group("/mw")
 pr_route.Use(middlewares.ReadAuth);
 {
  pr_route.GET("/me", handlers.Me)	
  pr_route.GET("/test-auth", func(c *gin.Context){
   userID := c.GetString("userID")
   if userID == ""{
    c.JSON(http.StatusForbidden, gin.H{"error":"Access denied"})
    c.Abort()
    return
   }
   c.JSON(http.StatusOK, gin.H{
    "message": "middleware ran well",
    "id": userID,
   })
  })
 }
 create := api.Group("/event")
 create.Use(middlewares.ReadAuth)
 {
  create.POST("/create", handlers.CreateEventHandler)
  create.PATCH("/update-status", handlers.UpdateSyncStatus)
 }
}

