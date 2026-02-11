package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"final-by-me/internal/db"
	"final-by-me/internal/handlers"
	"final-by-me/internal/repository"
	"final-by-me/internal/seed"
	"final-by-me/internal/worker"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "EPL-Connect"
	}

	client, err := db.Connect(mongoURI)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database(dbName)

	teamRepo := repository.NewTeamRepo(database)
	matchRepo := repository.NewMatchRepo(database)
	eventRepo := repository.NewEventRepo(database)

	// Start background worker (goroutine + channel)
	eventCh, stopWorker := worker.StartEventWorker(eventRepo, 100)
	defer stopWorker()

	// Example background heartbeat goroutine (optional extra)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Println("[BG] server alive")
		}
	}()

	// indexes + seed
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := teamRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal("team index error:", err)
	}
	if err := matchRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal("match index error:", err)
	}
	if err := teamRepo.SeedIfEmpty(ctx, seed.EPLTeams()); err != nil {
		log.Fatal("seed teams error:", err)
	}

	mux := http.NewServeMux()

	// Frontend
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Handlers
	teamH := handlers.NewTeamHandler(teamRepo)
	matchH := handlers.NewMatchMongoHandler(matchRepo, teamRepo, eventCh)
	tableH := handlers.NewTableHandler(teamRepo, matchRepo)

	// Routes
	mux.HandleFunc("GET /teams", teamH.ListTeams)

	mux.HandleFunc("GET /matches", matchH.ListMatches)
	mux.HandleFunc("POST /matches", matchH.CreateMatch)
	mux.HandleFunc("PATCH /matches/", matchH.AddEvent) // /matches/{key}/events
	mux.HandleFunc("POST /matches/", matchH.Finalize)  // /matches/{key}/finalize

	mux.HandleFunc("GET /table", tableH.GetTable)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	addr := ":" + port
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
