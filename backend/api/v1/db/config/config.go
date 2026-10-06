package config

import (
"fmt"
"log"
"os"
"github.com/Marcel-dev2009/cadence/db/connections/models"
"github.com/gin-contrib/cors"
"github.com/joho/godotenv"
"gorm.io/driver/postgres"
"gorm.io/gorm"
)
func Load(){
 _ = godotenv.Load()
 connectDB();	
}
var DB *gorm.DB
func connectDB(){
 dsn := os.Getenv("DATABASE_URL");
 db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
 if err != nil{
  log.Fatal("Database Connection Failed", err)	
 }
 DB = db
 err = db.AutoMigrate(&models.User{}, &models.ScheduledEvents{});
 if err != nil{
  log.Fatal("Migration Failed", err)	
 }
 fmt.Println("✅ Models migrated")
 fmt.Println("✅ Database connected")
 sqlDb, err := db.DB();
 sqlDb.SetMaxIdleConns(20);
 sqlDb.SetMaxOpenConns(100);
}
func CORSConfig() cors.Config {
 config := cors.DefaultConfig()
 config.AllowOrigins = []string{"http://localhost:8080"}  
 config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
 config.AllowHeaders = []string {"Origin", "Content-Type", "Authorization"}
 config.AllowCredentials = true
 return  config
}