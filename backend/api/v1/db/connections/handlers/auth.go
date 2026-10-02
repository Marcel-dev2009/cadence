package handlers

import (
"net/http"
"github.com/Marcel-dev2009/cadence/db/config"
"github.com/Marcel-dev2009/cadence/db/connections/models"
"github.com/gin-gonic/gin"
"golang.org/x/crypto/bcrypt"
)
type Register struct {
 Name string `json:"name" binding:"required, min=3"`
 Email string `json:"email" binding:"required, email"`	
 Password string `json:"password" bidning:"required, min=9, max=15"`
}
type Login struct {
 Email string `json:"email" binding:"required, email"`	
 Password string `json:"password" bidning:"required, min=9, max=15"`
}
func SignUp(c *gin.Context){
 var input Register
 if err := c.ShouldBindJSON(&input); err != nil{
  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  return	
 } 
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
 // creating and storing user session
}
func SingIn(c *gin.Context){
 var input Login	
 err := c.ShouldBindJSON(&user); err != nil{
  c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})	
  return
 }
var user models.User;
if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil{
 c.JSON(http.StatusUnauthorized, gin.H{"error":"no account found for this user"})	
 return
}
if err := bcrypt.CompareHashAndPassword([]byte(input.Password)); err != nil{
	c.JSON(http.StatusUnauthorized, gin.H{"error:invalid credentials"})
	return
}
// create session
}
func SignOut(c *gin.Context){
 // sigout logic	
}