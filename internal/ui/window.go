package ui

import (
	"fmt"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
	"github.com/pauloborszcz/tics/internal/service"
)

type Window struct {
	win          *adw.ApplicationWindow
	toastOverlay *adw.ToastOverlay
	titleWidget  *adw.WindowTitle
	avatar       *adw.Avatar
	ticketList   *TicketList
	ticketDetail *TicketDetail
	client       *glpi.Client
	sync         *service.SyncService
	cfg          *config.Config
	onSettings   func()
}

func NewWindow(app *gtk.Application, client *glpi.Client, syncSvc *service.SyncService, cfg *config.Config, onSettings func()) *Window {
	w := &Window{
		client:     client,
		sync:       syncSvc,
		cfg:        cfg,
		onSettings: onSettings,
	}

	w.win = adw.NewApplicationWindow(app)
	w.win.SetTitle("Tics")
	w.win.SetDefaultSize(1050, 680)

	// ToolbarView: header + content
	toolbar := adw.NewToolbarView()
	header := adw.NewHeaderBar()

	// Avatar menu button (left)
	avatarBtn := gtk.NewMenuButton()
	avatarBtn.SetHasFrame(false)
	avatarBtn.AddCSSClass("circular")

	w.avatar = adw.NewAvatar(32, "", false)
	avatarBtn.SetChild(w.avatar)

	// Popover menu
	popBox := gtk.NewBox(gtk.OrientationVertical, 2)
	popBox.SetMarginTop(6)
	popBox.SetMarginBottom(6)
	popBox.SetMarginStart(6)
	popBox.SetMarginEnd(6)
	popBox.SetSizeRequest(220, -1)

	// Auto-followup toggle row
	autoRow := gtk.NewBox(gtk.OrientationHorizontal, 10)
	autoRow.SetMarginTop(6)
	autoRow.SetMarginBottom(4)
	autoRow.SetMarginStart(8)
	autoRow.SetMarginEnd(8)

	autoIcon := gtk.NewImageFromIconName("mail-send-symbolic")
	autoIcon.SetPixelSize(16)
	autoRow.Append(autoIcon)

	autoLabel := gtk.NewLabel("Auto-followup")
	autoLabel.SetHExpand(true)
	autoLabel.SetHAlign(gtk.AlignStart)
	autoRow.Append(autoLabel)

	autoSwitch := gtk.NewSwitch()
	autoSwitch.SetVAlign(gtk.AlignCenter)
	autoSwitch.ConnectStateSet(func(state bool) bool {
		syncSvc.SetAutoFollowup(state)
		return false
	})
	autoRow.Append(autoSwitch)
	popBox.Append(autoRow)

	popBox.Append(gtk.NewSeparator(gtk.OrientationHorizontal))

	refreshBtn := popoverButton("view-refresh-symbolic", "Atualizar chamados")
	popBox.Append(refreshBtn)

	settingsItem := popoverButton("emblem-system-symbolic", "Configuracoes")
	popBox.Append(settingsItem)

	popover := gtk.NewPopover()
	popover.SetChild(popBox)
	avatarBtn.SetPopover(popover)

	refreshBtn.ConnectClicked(func() {
		popover.Popdown()
		go syncSvc.Refresh()
	})

	settingsItem.ConnectClicked(func() {
		popover.Popdown()
		if w.onSettings != nil {
			w.win.Close()
			w.onSettings()
		}
	})

	header.PackStart(avatarBtn)

	// Title with subtitle (ticket count)
	w.titleWidget = adw.NewWindowTitle("Tics", "")
	header.SetTitleWidget(w.titleWidget)

	toolbar.AddTopBar(header)

	// Split view
	split := gtk.NewPaned(gtk.OrientationHorizontal)
	split.SetPosition(380)
	split.SetShrinkStartChild(false)
	split.SetShrinkEndChild(false)

	// Left: ticket list
	w.ticketList = NewTicketList()
	scrollLeft := gtk.NewScrolledWindow()
	scrollLeft.SetChild(w.ticketList.listBox)
	scrollLeft.SetSizeRequest(350, -1)
	scrollLeft.SetVExpand(true)
	split.SetStartChild(scrollLeft)

	// Right: ticket detail with toast overlay
	w.ticketDetail = NewTicketDetail(client, &w.win.ApplicationWindow.Window)
	w.toastOverlay = adw.NewToastOverlay()
	scrollRight := gtk.NewScrolledWindow()
	scrollRight.SetChild(w.ticketDetail.box)
	scrollRight.SetVExpand(true)
	w.toastOverlay.SetChild(scrollRight)
	split.SetEndChild(w.toastOverlay)

	w.ticketList.OnSelect(func(ticket glpi.Ticket) {
		w.ticketDetail.ShowTicket(ticket)
	})

	toolbar.SetContent(split)
	w.win.SetContent(toolbar)

	applyCSS(w.win)

	return w
}

func popoverButton(iconName, label string) *gtk.Button {
	btn := gtk.NewButton()
	btn.SetHasFrame(false)

	box := gtk.NewBox(gtk.OrientationHorizontal, 8)
	box.SetMarginStart(6)
	box.SetMarginEnd(6)
	box.SetMarginTop(4)
	box.SetMarginBottom(4)

	icon := gtk.NewImageFromIconName(iconName)
	icon.SetPixelSize(16)
	box.Append(icon)

	lbl := gtk.NewLabel(label)
	lbl.SetHAlign(gtk.AlignStart)
	box.Append(lbl)

	btn.SetChild(box)
	return btn
}

func (w *Window) UpdateTickets(tickets []glpi.Ticket) {
	w.ticketList.Update(tickets)
	if len(tickets) == 1 {
		w.titleWidget.SetSubtitle("1 chamado")
	} else {
		w.titleWidget.SetSubtitle(fmt.Sprintf("%d chamados", len(tickets)))
	}
}

func (w *Window) SetUserName(name string) {
	w.avatar.SetText(name)
	w.avatar.SetShowInitials(true)
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
	texture := gdk.NewTextureForPixbuf(pixbuf)
	w.avatar.SetCustomImage(texture)
}

func (w *Window) ShowToast(message string) {
	toast := adw.NewToast(message)
	w.toastOverlay.AddToast(toast)
}

func applyCSS(win *adw.ApplicationWindow) {
	css := gtk.NewCSSProvider()
	css.LoadFromData(`
		/* Ticket list */
		.ticket-row { padding: 10px 14px; }
		.ticket-id {
			font-weight: bold;
			font-size: 0.9em;
			color: @accent_color;
		}
		.ticket-title { font-weight: 600; }
		.ticket-meta {
			font-size: 0.8em;
			color: alpha(@window_fg_color, 0.5);
		}

		/* Priority stripe */
		.priority-stripe { border-radius: 3px 0 0 3px; }
		.priority-stripe-very-low { background: #33d17a; }
		.priority-stripe-low { background: #57e389; }
		.priority-stripe-medium { background: #e5a50a; }
		.priority-stripe-high { background: #ff7800; }
		.priority-stripe-very-high { background: #ed333b; }
		.priority-stripe-critical { background: #c01c28; }

		/* Priority text */
		.priority-very-low { color: #26a269; }
		.priority-low { color: #33d17a; }
		.priority-medium { color: #e5a50a; }
		.priority-high { color: #ff7800; }
		.priority-very-high { color: #ed333b; }
		.priority-critical { color: #c01c28; font-weight: bold; }

		/* Priority badge */
		.priority-badge {
			font-size: 0.8em;
			padding: 2px 8px;
			border-radius: 4px;
			font-weight: 600;
		}

		/* Status badges */
		.status-badge {
			font-size: 0.8em;
			padding: 2px 10px;
			border-radius: 6px;
			font-weight: bold;
		}
		.status-new {
			background: alpha(#3584e4, 0.15);
			color: #3584e4;
		}
		.status-processing {
			background: alpha(#2ec27e, 0.15);
			color: #2ec27e;
		}
		.status-pending {
			background: alpha(#ff7800, 0.15);
			color: #e66100;
		}
		.status-solved {
			background: alpha(#26a269, 0.15);
			color: #26a269;
		}
		.status-closed {
			background: alpha(@window_fg_color, 0.08);
			color: alpha(@window_fg_color, 0.5);
		}

		/* Detail */
		.detail-title { font-size: 1.3em; font-weight: bold; }
		.detail-section-title {
			font-weight: bold;
			font-size: 1em;
			margin-top: 16px;
			color: alpha(@window_fg_color, 0.7);
		}
		.info-card { padding: 16px; }
		.info-card-label {
			font-size: 0.7em;
			font-weight: bold;
			color: alpha(@window_fg_color, 0.45);
			letter-spacing: 1px;
		}
		.info-card-value { font-size: 0.95em; }
		.followup-card { padding: 12px; margin: 2px 0; }
		.followup-date {
			font-size: 0.8em;
			color: alpha(@window_fg_color, 0.45);
		}
		.followup-header {
			font-size: 0.82em;
			color: alpha(@window_fg_color, 0.45);
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		win.ApplicationWindow.Window.Widget.Display(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
