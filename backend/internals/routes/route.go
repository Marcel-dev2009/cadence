package routes

import (
"github.com/Marcel-dev2009/cadence/api/v1/handlers"
"github.com/gin-gonic/gin"
)
func Setup(r *gin.Engine){
 api := r.Group("/api/v1")
 create := api.Group("/event")

 {
  create.POST("/create", handlers.CreateEventHandler)
  create.PATCH("/update_status", handlers.UpdateSyncStatus) 
  create.PATCH("/:id/status", handlers.UpdateEventStatusHandler)
  create.PATCH("/sync_spotify", handlers.UpdateSpotifySyncStatus)
  create.GET("/upcoming_event", handlers.GetDashboardHighlightHandler)
  create.GET("/user_data", handlers.GetUserMetaData)
  create.GET("/all_events", handlers.GetUserEventsHandler)
 }
}

