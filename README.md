# personal-microservices iCal Service

Dieses Repository stellt einen kleinen Go-Service bereit, der einen iCal-Feed filtert und über einen HTTP-Endpunkt ausliefert.

Wichtig: Die Anwendung liest folgende Umgebungsvariablen beim Start ein und beendet sich, falls sie nicht gesetzt sind:

- ICAL_SOURCE_URL – die URL zum Original iCal
- ACCESS_TOKEN – das Access-Token, das als query parameter `token` übergeben werden muss

Lokale Nutzung (mit docker-compose):

1. Erstelle eine `.env` Datei neben `docker-compose.yml` mit folgendem Inhalt:

ICAL_SOURCE_URL=https://example.com/my.ics
ACCESS_TOKEN=mein-geheimes-token

2. Starte mit:

```powershell
docker compose up --build -d
```

Der Service ist dann unter http://localhost:3411/calendar/bbzw?token=mein-geheimes-token erreichbar.

GitHub Actions:

Beim Push auf `main` oder beim Erstellen eines Tags `v*` baut die Action ein Docker-Image und pusht es in die GitHub Container Registry (ghcr.io) unter dem Namen `personal-microservices`.

Beispiel-Image-URL nach erfolgreichem Build:

- ghcr.io/<dein-github-user-or-org>/personal-microservices:latest
