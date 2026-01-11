package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/joho/godotenv"
)

var (
	cachedCal *ical.Calendar
	cacheTime time.Time
	cacheTTL  = 10 * time.Minute

	sourceURL   string
	accessToken string
)

func init() {
	_ = godotenv.Load()
	sourceURL = os.Getenv("ICAL_SOURCE_URL")
	accessToken = os.Getenv("ACCESS_TOKEN")
	if sourceURL == "" || accessToken == "" {
		log.Fatalf("Umgebungsvariablen ICAL_SOURCE_URL und ACCESS_TOKEN müssen gesetzt sein")
	}
}

func fetchAndFilterCalendar() (*ical.Calendar, error) {
	if cachedCal != nil && time.Since(cacheTime) < cacheTTL {
		log.Println("Cache HIT")
		return cachedCal, nil
	}

	log.Println("Cache MISS – lade Kalender neu")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Go iCal Client)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Println("Fehler beim Schließen des Response-Body:", cerr)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	srcCal, err := ical.ParseCalendar(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	newCal := ical.NewCalendar()
	newCal.SetProductId("BBZW Prüfungen + Ferien iCal Feed")
	newCal.SetMethod(ical.MethodPublish)

	for _, e := range srcCal.Events() {
		p := e.GetProperty(ical.ComponentPropertySummary)
		if p == nil {
			continue
		}

		title := p.Value
		if !strings.Contains(title, " ") {
			continue
		}

		ev := newCal.AddEvent(e.Id())
		ev.SetSummary(title)

		if startAt, err := e.GetStartAt(); err == nil {
			ev.SetStartAt(startAt)
		}
		if endAt, err := e.GetEndAt(); err == nil {
			ev.SetEndAt(endAt)
		}

		if loc := e.GetProperty(ical.ComponentPropertyLocation); loc != nil {
			ev.SetLocation(loc.Value)
		}
	}

	cachedCal = newCal
	cacheTime = time.Now()

	return newCal, nil
}
