package routes

import (
"github.com/Marcel-dev2009/cadence/db/connections/handlers"
"github.com/Marcel-dev2009/cadence/db/me/middlewares"
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
 pr_route := api.Group("/")
 pr_route.Use(middlewares.ReadAuth)
 {
  pr_route.GET("/me", handlers.Me)	
 }
}

