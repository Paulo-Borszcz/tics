package ui

import (
	"log"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
	"github.com/pauloborszcz/tics/internal/service"
)

const AppID = "com.github.pauloborszcz.Tics"

type App struct {
	gtkApp *gtk.Application
	cfg    *config.Config
	client *glpi.Client
	sync   *service.SyncService
	notify *service.Notifier
	window *Window
}

func NewApp(cfg *config.Config) *App {
	gtkApp := gtk.NewApplication(AppID, gio.ApplicationFlagsNone)

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

	// Register show-window action for notifications
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
	win := gtk.NewApplicationWindow(a.gtkApp)
	win.SetTitle("Tics - Configuracao")
	win.SetDefaultSize(500, 500)

	header := gtk.NewHeaderBar()
	win.SetTitlebar(header)

	setupPage := NewSetupPage(a.cfg, func(cfg *config.Config) {
		a.cfg = cfg
		win.Close()
		a.startMainApp()
	})

	win.SetChild(setupPage.box)
	applySetupCSS(win)
	win.Show()
}

func (a *App) startMainApp() {
	a.client = glpi.NewClient(a.cfg)

	log.Println("Initializing GLPI session...")
	if err := a.client.InitSession(); err != nil {
		log.Printf("Failed to init GLPI session: %v", err)
		// Show setup again if session fails
		a.showSetup()
		return
	}
	log.Println("GLPI session initialized successfully")

	a.sync = service.NewSyncService(a.client, a.cfg)
	a.notify = service.NewNotifier(&a.gtkApp.Application)

	log.Println("Creating window...")
	a.window = NewWindow(a.gtkApp, a.client, a.sync, a.cfg, func() {
		a.window = nil
		a.showSetup()
	})
	log.Println("Showing window...")
	a.window.win.Show()
	log.Println("Window shown")

	// Load user info in background
	go a.loadUserInfo()

	// Wire sync callbacks
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

	// Start background sync
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

func applySetupCSS(win *gtk.ApplicationWindow) {
	css := gtk.NewCSSProvider()
	css.LoadFromData(`
		.setup-title {
			font-size: 2em;
			font-weight: bold;
		}
		.setup-subtitle {
			font-size: 1.1em;
			color: alpha(@theme_fg_color, 0.6);
		}
		.setup-field-label {
			font-weight: bold;
			font-size: 0.9em;
		}
		.setup-save-btn {
			padding: 8px 24px;
			font-size: 1.05em;
		}
		.setup-error {
			color: #cc0000;
		}
		.setup-success {
			color: #4e9a06;
		}
		.setup-status {
			font-size: 0.9em;
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		win.Window.Widget.Display(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
