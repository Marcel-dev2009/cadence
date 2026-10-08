package main

import (
	"log"
	"os"
	"time"

	"github.com/Marcel-dev2009/cadence/database/config"
	"github.com/Marcel-dev2009/cadence/database/me/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {	
 config.Load()
 config.SessionCache = cache.New(14*25*time.Hour, 10*time.Minute)
 config.DataListCache = cache.New(2*time.Hour, 10*time.Minute)
 r := gin.Default();
 r.Use(cors.New(config.CORSConfig()))
 r.GET("/", func(c *gin.Context){
  c.JSON(200, gin.H{"status":"backend server is running smoothly"})	
 })
 routes.Setup(r);
 port := os.Getenv("PORT")
 if port == ""{
  port = "8080"	
 }
 log.Println("server running at port", port)
 r.Run(":"+port)
}