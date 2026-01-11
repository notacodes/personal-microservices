package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Filtered BBZW Calendar Endpoint
	r.GET("/calendar/bbzw", func(c *gin.Context) {
		token := c.Query("token")

		if token != accessToken {
			log.Println("Unauthorized Zugriff")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		cal, err := fetchAndFilterCalendar()
		if err != nil {
			log.Println("Kalender Fehler:", err)
			c.String(http.StatusInternalServerError, "calendar error")
			return
		}

		c.Header("Content-Type", "text/calendar; charset=utf-8")
		c.Header("Content-Disposition", "inline; filename=calendar.ics")
		c.String(http.StatusOK, cal.Serialize())
	})

	log.Println("iCal Feed verfügbar unter http://localhost:3411/calendar/bbzw")
	if err := r.Run(":3411"); err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
