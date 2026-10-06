package handlers

import (
"net/http"
"strings"
"time"
"github.com/Marcel-dev2009/cadence/db/config"
"github.com/Marcel-dev2009/cadence/db/connections/models"
"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

type EventInput struct{
	EventName string `json:"event_name"`
    EventTime *time.Time `form:"event_time" time_format:"2006-01-02"`
}
 func checkGoogleCalenderSync(c *gin.Context, db *gorm.DB, userID string)  bool {
  var result struct {
    SyncGoogleCalender bool 
  } 
  err := db.Table("users").Select("sync_google_calender").Where("id = ?", userID).First(&result).Error
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error":"user record not found"})
    return false
  }
  if result.SyncGoogleCalender == false {
    c.JSON(http.StatusUnauthorized, gin.H{"error":"sync to google calender to create events, go to your settings page to do so"})
    c.Abort()
    return false
  }
  return true
 }

func CreateEventHandler(c *gin.Context) {
    db := config.DB
    var input EventInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
     userID := c.GetString("userID")
     if userID == ""{
       c.JSON(http.StatusUnauthorized, gin.H{"error":"User session not found or is expired"}) 
       return
     }
      if allowed := checkGoogleCalenderSync(c, db, userID); !allowed{
        return
      }
    if strings.TrimSpace(input.EventName) == " " {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You must give the event a name"})
        return
    }
    if err := c.ShouldBindQuery(&input); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid Time format"})
    }
   if input.EventTime == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must set a specific time for the event via ?event_time=YYYY-MM-DD"})
		return
	}
    if input.EventTime.Before(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot schedule an event in the past!"})
        return
    }
    // create the event for them
   event := models.ScheduledEvents{
    EventName: input.EventName,
    EventTime: input.EventTime,
    UserID: userID,
   } 
   if err := db.Create(&event).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
    "error":"error creating events",
   })
   c.JSON(http.StatusCreated, gin.H {
    "message":"event created successfully!",
    "data": event,
   })
    return
   }
}
