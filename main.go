package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"final-by-me/internal/db"
	"final-by-me/internal/handlers"
	"final-by-me/internal/middleware"
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

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// DB connect
	client, err := db.Connect(mongoURI)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	database := client.Database(dbName)

	teamRepo := repository.NewTeamRepo(database)
	matchRepo := repository.NewMatchRepo(database)
	eventRepo := repository.NewEventRepo(database)

	// background worker
	eventCh, stopWorker := worker.StartEventWorker(eventRepo, 100)
	defer stopWorker()

	// indexes + seed
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := teamRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal("team index error:", err)
	}
	if err := matchRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal("match index error:", err)
	}

	// If you use leagues, seed with TeamsAll()
	// If you still seed EPL only, change it back.
	if err := teamRepo.SeedIfEmpty(ctx, seed.TeamsAll()); err != nil {
		log.Fatal("seed teams error:", err)
	}

	// Handlers
	authH := handlers.NewAuthHandler(database, jwtSecret, teamRepo)
	teamH := handlers.NewTeamHandler(teamRepo)
	matchH := handlers.NewMatchMongoHandler(matchRepo, teamRepo, eventCh)
	tableH := handlers.NewTableHandler(teamRepo, matchRepo)
	statsH := handlers.NewStatsHandler(matchRepo)

	// Router
	mux := http.NewServeMux()

	// Frontend
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Public API
	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login", authH.Login)

	mux.HandleFunc("GET /leagues", teamH.ListLeagues)
	mux.HandleFunc("GET /teams", teamH.ListTeams)

	mux.HandleFunc("GET /matches", matchH.ListMatches)
	mux.HandleFunc("GET /table", tableH.GetTable)
	mux.HandleFunc("GET /stats", statsH.GetStats)

	// Admin-only chain
	adminChain := func(h http.Handler) http.Handler {
		return middleware.WithJSON(
			middleware.AuthJWT(jwtSecret)(
				middleware.RequireRole("admin")(h),
			),
		)
	}

	// ✅ Admin routes MUST match UI calls (NO conflicts)
	mux.Handle("POST /matches", adminChain(http.HandlerFunc(matchH.CreateMatch)))
	mux.Handle("PATCH /matches/{key}/events", adminChain(http.HandlerFunc(matchH.AddEvent)))
	mux.Handle("PATCH /matches/{key}/status", adminChain(http.HandlerFunc(matchH.SetStatus)))
	mux.Handle("POST /matches/{key}/finalize", adminChain(http.HandlerFunc(matchH.Finalize)))

	addr := ":" + port
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
