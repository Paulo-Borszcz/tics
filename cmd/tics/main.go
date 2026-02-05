package main

import (
	"log"
	"os"

	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/ui"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Println("Starting Tics...")

	cfg := config.Load()
	app := ui.NewApp(cfg)

	log.Println("Launching GTK application...")
	os.Exit(app.Run())
}
