package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/Marcel-dev2009/cadence/database/config"
	"github.com/Marcel-dev2009/cadence/database/models"
	"github.com/Marcel-dev2009/cadence/repository"
	"github.com/gin-gonic/gin"
)

type EventInput struct{
	EventName string `json:"event_name"`
    EventTime *time.Time `form:"event_time" time_format:"2006-01-02"`
}
    var db = config.DB
    var dataCache = config.DataListCache
func CreateEventHandler(c *gin.Context) {
   
    var input EventInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
      userRepo := repository.NewUserRepository(db)
      eventRepo := repository.NewEventsRepository(db, dataCache)
     userID := c.GetString("userID")
     if userID == ""{
       c.JSON(http.StatusUnauthorized, gin.H{"error":"User session not found or is expired"}) 
       return
     }
      isSynced, err := userRepo.GetgoogleCalenderSync(userID); 
   if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking sync status"})
      return
    } 
    if !isSynced {
      c.JSON(http.StatusForbidden, gin.H{"error": "Please sync Google Calendar first"})
      return
    }
    if strings.TrimSpace(input.EventName) == "" {
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
   newEvent := models.ScheduledEvents{
    EventName: input.EventName,
    EventTime: input.EventTime,
    UserID: userID,
   } 
    if err := eventRepo.CreateEvent(&newEvent); err != nil{
      c.JSON(http.StatusBadRequest, gin.H{"error":"failed to create event"})
      return
    }
    c.JSON(http.StatusCreated, gin.H{
      "message":"event created succcesfully",
      "data": newEvent,
    })
   }

   func GetUserEventsHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Pass both DB and your cache wrapper into the manager instance
	eventRepo := repository.NewEventsRepository(db, dataCache)

	events, err := eventRepo.GetEvents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to look up events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

