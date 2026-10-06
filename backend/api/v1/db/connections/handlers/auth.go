package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"github.com/Marcel-dev2009/cadence/db/config"
	"github.com/Marcel-dev2009/cadence/db/connections/models"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)
type Register struct {
 Name string `json:"name" binding:"required,min=3"`
 Email string `json:"email" binding:"required"`	
 Password string `json:"password" binding:"required,min=9,max=15"`
}
type Login struct {
 Email string `json:"email" binding:"required,email"`	
 Password string `json:"password" binding:"required,min=9,max=15"`
}
func generateSessionID() string {
 b := make([]byte, 32)
 rand.Read(b)
 return hex.EncodeToString(b)   
}
func SignUp(c *gin.Context){
 var input Register
 if err := c.ShouldBindJSON(&input); err != nil{
 c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})         
  return          
 }
 fmt.Println("Reached here!");
 var existing models.User
 if err := config.DB.Where("email = ? ", input.Email).First(&existing).Error; err == nil{
  c.JSON(http.StatusConflict, gin.H{"error":"email already exisis"})	
  return
 }
 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
 if err != nil {
 c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to encrypt password"})
 return	
 }
 user := models.User{
 Name: input.Name,
 Email: input.Email,
 Password: string(hashedPassword),	
 }
 if err = config.DB.Create(&user).Error; err != nil{
 c.JSON(http.StatusInternalServerError, gin.H{"error":"error creating user"})	
 return
 }
 // create session
 sessionID := generateSessionID()
 fmt.Println(sessionID) // here
 config.SessionCache.Set(sessionID, user.ID, cache.DefaultExpiration)
 c.SetCookie(
  "session_token",
  sessionID,
  7200,
  "/",
  "",   //would be changed to our deployed domain when deployed
  false,  // remember when deplyed changed to https(true)
  true,
 )
 c.JSON(http.StatusCreated , gin.H{
 "message":"account created succesfully!",
 "user":gin.H{
 "name":user.Name,
 "email":user.Email,
 "id":user.ID,
 },
})
}
func SignIn(c *gin.Context){
 var input Login	
 if err := c.ShouldBindJSON(&input); err != nil{
   c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
  return
 }

var user models.User;
if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil{
 c.JSON(http.StatusUnauthorized, gin.H{"error":"no account found for this user"})	
 return
}
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil{
	c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"})
	return
}
 sessionID := generateSessionID()
  config.SessionCache.Set(sessionID, user.ID, cache.DefaultExpiration)
 c.SetCookie(
  "session_token",
  sessionID,
  7200,
  "/",
   "",   //would be changed to our deployed domain when deployed
  false,  // remember when deplyed changed to https(true)
  true,
 )
 c.JSON(http.StatusOK, gin.H{
 "message":"signed in succesfully!",
 "user": gin.H{
  "id": user.ID,
  "email":user.Email,      
 },
})
}
func SignOut(c *gin.Context){
 cookie, err := c.Cookie("session_token")
 if err == nil{
  config.SessionCache.Delete(cookie)        
 }
 c.SetCookie("session_token", "", -1, "", "/", false, true) 
 c.JSON(http.StatusOK, gin.H{"message": "successfully signed out"})
}
func Me (c *gin.Context){
  cookie, err := c.Cookie("session_token")
  if err != nil{
    c.JSON(http.StatusUnauthorized, gin.H{"error":"user not found"})      
    c.Abort()
    return 
  }
  userID, found := config.SessionCache.Get(cookie)
  if !found{
   c.JSON(http.StatusUnauthorized, gin.H{"error":" your session has expired, please login again"})
   c.Abort()
   return        
  }
  c.Set("userID", userID)
  c.Next()
 }
