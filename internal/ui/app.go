package ui

import (
	"log"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
	"github.com/pauloborszcz/tics/internal/service"
)

const AppID = "com.github.pauloborszcz.Tics"

type App struct {
	gtkApp *adw.Application
	cfg    *config.Config
	client *glpi.Client
	sync   *service.SyncService
	notify *service.Notifier
	window *Window
}

func NewApp(cfg *config.Config) *App {
	gtkApp := adw.NewApplication(AppID, gio.ApplicationFlagsNone)

	app := &App{
		gtkApp: gtkApp,
		cfg:    cfg,
	}

	gtkApp.ConnectActivate(func() {
		log.Println("GTK Activate signal received")
		app.onActivate()
	})

	gtkApp.ConnectShutdown(func() {
		app.onShutdown()
	})

	showAction := gio.NewSimpleAction("show-window", nil)
	showAction.ConnectActivate(func(_ *glib.Variant) {
		if app.window != nil {
			app.window.win.Present()
		}
	})
	gtkApp.AddAction(showAction)

	return app
}

func (a *App) onActivate() {
	log.Println("onActivate called")
	if a.window != nil {
		a.window.win.Present()
		return
	}

	if !a.cfg.IsConfigured() {
		a.showSetup()
		return
	}

	a.startMainApp()
}

func (a *App) showSetup() {
	win := adw.NewApplicationWindow(&a.gtkApp.Application)
	win.SetTitle("Tics - Configuracao")
	win.SetDefaultSize(500, 600)

	toolbar := adw.NewToolbarView()
	header := adw.NewHeaderBar()
	toolbar.AddTopBar(header)

	setupPage := NewSetupPage(a.cfg, func(cfg *config.Config) {
		a.cfg = cfg
		win.Close()
		a.startMainApp()
	})

	toolbar.SetContent(setupPage.box)
	win.SetContent(toolbar)
	win.Show()
}

func (a *App) startMainApp() {
	a.client = glpi.NewClient(a.cfg)

	log.Println("Initializing GLPI session...")
	if err := a.client.InitSession(); err != nil {
		log.Printf("Failed to init GLPI session: %v", err)
		a.showSetup()
		return
	}
	log.Println("GLPI session initialized successfully")

	a.sync = service.NewSyncService(a.client, a.cfg)
	a.notify = service.NewNotifier(&a.gtkApp.Application.Application)

	log.Println("Creating window...")
	a.window = NewWindow(&a.gtkApp.Application, a.client, a.sync, a.cfg, func() {
		a.window = nil
		a.showSetup()
	})
	log.Println("Showing window...")
	a.window.win.Show()
	log.Println("Window shown")

	go a.loadUserInfo()

	a.sync.OnUpdate(func(tickets []glpi.Ticket) {
		glib.IdleAdd(func() {
			a.window.UpdateTickets(tickets)
		})
	})
	a.sync.OnNewTickets(func(tickets []glpi.Ticket) {
		glib.IdleAdd(func() {
			a.notify.NotifyNewTickets(tickets)
		})
	})

	a.sync.Start()
}

func (a *App) loadUserInfo() {
	session, err := a.client.GetFullSession()
	if err != nil {
		log.Printf("Failed to get user session: %v", err)
		return
	}

	glib.IdleAdd(func() {
		a.window.SetUserName(session.Name)
	})

	pictureData, err := a.client.GetUserPicture(session.ID)
	if err != nil {
		log.Printf("Failed to get user picture: %v", err)
		return
	}

	glib.IdleAdd(func() {
		a.window.SetUserAvatar(pictureData)
	})
}

func (a *App) onShutdown() {
	if a.sync != nil {
		a.sync.Stop()
	}
	if a.client != nil {
		a.client.KillSession()
	}
}

func (a *App) Run() int {
	return a.gtkApp.Run(os.Args)
}
