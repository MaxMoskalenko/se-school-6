package scanner

import (
	"context"
	"log"
	"time"
)

func (a *App) Run(ctx context.Context) error {
	log.Printf("scanner started with interval %s", a.cfg.Interval)

	ticker := time.NewTicker(a.cfg.Interval)
	defer ticker.Stop()

	// run immediately on start
	a.runScan(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("scanner stopped")
			return ctx.Err()
		case <-ticker.C:
			a.runScan(ctx)
		}
	}
}

func (a *App) runScan(ctx context.Context) {
	log.Printf("scanner: starting scan")
	if err := a.Scan(ctx); err != nil {
		log.Printf("scanner: scan failed: %v", err)
	}
	log.Printf("scanner: scan completed")
}
