package ui

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
	"github.com/pauloborszcz/tics/internal/service"
)

type Window struct {
	win          *gtk.ApplicationWindow
	ticketList   *TicketList
	ticketDetail *TicketDetail
	split        *gtk.Paned
	client       *glpi.Client
	sync         *service.SyncService
	cfg          *config.Config
	avatarImage  *gtk.Image
	nameLabel    *gtk.Label
	onSettings   func()
}

func NewWindow(app *gtk.Application, client *glpi.Client, syncSvc *service.SyncService, cfg *config.Config, onSettings func()) *Window {
	w := &Window{
		client:     client,
		sync:       syncSvc,
		cfg:        cfg,
		onSettings: onSettings,
	}

	w.win = gtk.NewApplicationWindow(app)
	w.win.SetTitle("Tics")
	w.win.SetDefaultSize(1000, 600)

	// Header bar
	header := gtk.NewHeaderBar()

	// Left side: avatar + name
	avatarBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	avatarBox.SetMarginStart(4)

	w.avatarImage = gtk.NewImageFromIconName("avatar-default-symbolic")
	w.avatarImage.SetPixelSize(28)
	w.avatarImage.AddCSSClass("avatar-image")
	avatarBox.Append(w.avatarImage)

	w.nameLabel = gtk.NewLabel("")
	w.nameLabel.AddCSSClass("header-username")
	avatarBox.Append(w.nameLabel)

	header.PackStart(avatarBox)

	// Auto-followup toggle
	autoBox := gtk.NewBox(gtk.OrientationHorizontal, 6)
	autoLabel := gtk.NewLabel("Auto")
	autoLabel.AddCSSClass("header-auto-label")
	autoBox.Append(autoLabel)

	autoSwitch := gtk.NewSwitch()
	autoSwitch.SetActive(false)
	autoSwitch.SetVAlign(gtk.AlignCenter)
	autoSwitch.ConnectStateSet(func(state bool) bool {
		syncSvc.SetAutoFollowup(state)
		return false
	})
	autoBox.Append(autoSwitch)
	header.PackStart(autoBox)

	// Right side: settings + refresh
	settingsBtn := gtk.NewButtonFromIconName("emblem-system-symbolic")
	settingsBtn.SetTooltipText("Configuracoes")
	settingsBtn.ConnectClicked(func() {
		if w.onSettings != nil {
			w.win.Close()
			w.onSettings()
		}
	})
	header.PackEnd(settingsBtn)

	refreshBtn := gtk.NewButtonFromIconName("view-refresh-symbolic")
	refreshBtn.SetTooltipText("Atualizar chamados")
	refreshBtn.ConnectClicked(func() {
		go syncSvc.Refresh()
	})
	header.PackEnd(refreshBtn)

	w.win.SetTitlebar(header)

	// Split view
	w.split = gtk.NewPaned(gtk.OrientationHorizontal)
	w.split.SetPosition(380)
	w.split.SetShrinkStartChild(false)
	w.split.SetShrinkEndChild(false)

	// Left: ticket list
	w.ticketList = NewTicketList()
	scrollLeft := gtk.NewScrolledWindow()
	scrollLeft.SetChild(w.ticketList.listBox)
	scrollLeft.SetSizeRequest(350, -1)
	scrollLeft.SetVExpand(true)
	w.split.SetStartChild(scrollLeft)

	// Right: ticket detail
	w.ticketDetail = NewTicketDetail(client, &w.win.Window)
	scrollRight := gtk.NewScrolledWindow()
	scrollRight.SetChild(w.ticketDetail.box)
	scrollRight.SetVExpand(true)
	w.split.SetEndChild(scrollRight)

	// Wire list selection to detail
	w.ticketList.OnSelect(func(ticket glpi.Ticket) {
		w.ticketDetail.ShowTicket(ticket)
	})

	w.win.SetChild(w.split)

	// Apply CSS
	applyCSS(w.win)

	return w
}

func (w *Window) UpdateTickets(tickets []glpi.Ticket) {
	w.ticketList.Update(tickets)
}

func (w *Window) SetUserName(name string) {
	w.nameLabel.SetText(name)
}

func (w *Window) SetUserAvatar(data []byte) {
	loader := gdkpixbuf.NewPixbufLoader()
	if err := loader.Write(data); err != nil {
		return
	}
	loader.Close()
	pixbuf := loader.Pixbuf()
	if pixbuf == nil {
		return
	}
	scaled := pixbuf.ScaleSimple(28, 28, gdkpixbuf.InterpBilinear)
	if scaled == nil {
		return
	}
	texture := gdk.NewTextureForPixbuf(scaled)
	w.avatarImage.SetFromPaintable(texture)
}

func applyCSS(win *gtk.ApplicationWindow) {
	css := gtk.NewCSSProvider()
	css.LoadFromData(`
		.ticket-row {
			padding: 8px 12px;
		}
		.ticket-id {
			font-weight: bold;
			font-size: 0.9em;
			color: @theme_selected_bg_color;
		}
		.ticket-title {
			font-weight: bold;
		}
		.ticket-meta {
			font-size: 0.85em;
			color: alpha(@theme_fg_color, 0.6);
		}
		.priority-very-low { color: #4e9a06; }
		.priority-low { color: #73d216; }
		.priority-medium { color: #f57900; }
		.priority-high { color: #cc0000; }
		.priority-very-high { color: #a40000; }
		.priority-critical { color: #a40000; font-weight: bold; }
		.detail-title {
			font-size: 1.3em;
			font-weight: bold;
		}
		.detail-section-title {
			font-weight: bold;
			font-size: 1.1em;
			margin-top: 12px;
		}
		.info-card {
			background: alpha(@theme_fg_color, 0.04);
			border-radius: 10px;
			padding: 14px;
		}
		.info-card-label {
			font-size: 0.8em;
			font-weight: bold;
			color: alpha(@theme_fg_color, 0.5);
		}
		.info-card-value {
			font-size: 0.95em;
		}
		.content-card {
			background: alpha(@theme_fg_color, 0.03);
			border-radius: 10px;
			padding: 14px;
			border: 1px solid alpha(@theme_fg_color, 0.08);
		}
		.followup-card {
			background: alpha(@theme_fg_color, 0.05);
			border-radius: 10px;
			padding: 12px;
			margin: 4px 0;
			border: 1px solid alpha(@theme_fg_color, 0.06);
		}
		.followup-date {
			font-size: 0.8em;
			color: alpha(@theme_fg_color, 0.5);
		}
		.followup-header {
			font-size: 0.82em;
			color: alpha(@theme_fg_color, 0.5);
		}
		.status-badge {
			font-size: 0.85em;
			padding: 2px 10px;
			border-radius: 6px;
			background: alpha(@theme_selected_bg_color, 0.15);
			color: @theme_selected_bg_color;
			font-weight: bold;
		}
		.action-bar {
			background: alpha(@theme_fg_color, 0.03);
			border-radius: 10px;
			padding: 12px;
			margin-top: 8px;
		}
		.action-bar button {
			padding: 6px 16px;
		}
		.header-username {
			font-weight: bold;
			font-size: 0.9em;
		}
		.header-auto-label {
			font-size: 0.85em;
			color: alpha(@theme_fg_color, 0.7);
		}
		.avatar-image {
			border-radius: 50%;
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		win.Window.Widget.Display(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
