package worker

import (
	"context"
	"log"
	"time"

	"final-by-me/internal/models"
	"final-by-me/internal/repository"
)

type EventWorker struct {
	ch   <-chan models.EventLog
	repo *repository.EventRepo
	stop context.CancelFunc
}

func StartEventWorker(repo *repository.EventRepo, buffer int) (chan<- models.EventLog, func()) {
	ch := make(chan models.EventLog, buffer)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println("[WORKER] event worker started")
		for {
			select {
			case e, ok := <-ch:
				if !ok {
					log.Println("[WORKER] channel closed, worker stopped")
					return
				}

				// simulate async work (optional)
				time.Sleep(50 * time.Millisecond)

				// write to DB in background
				dbCtx, c := context.WithTimeout(context.Background(), 3*time.Second)
				_ = repo.Insert(dbCtx, e)
				c()

				log.Println("[WORKER] saved event:", e.Type, e.MatchKey)

			case <-ctx.Done():
				log.Println("[WORKER] cancel signal, worker stopped")
				return
			}
		}
	}()

	// return sender channel + stop function
	stopFn := func() {
		cancel()
		close(ch)
	}

	return ch, stopFn
}
