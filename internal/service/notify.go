package service

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type Notifier struct {
	app *gio.Application
}

func NewNotifier(app *gio.Application) *Notifier {
	return &Notifier{app: app}
}

// NotifyNewTickets sends a desktop notification for new tickets.
func (n *Notifier) NotifyNewTickets(tickets []glpi.Ticket) {
	if len(tickets) == 0 {
		return
	}

	var title, body string
	if len(tickets) == 1 {
		title = "Novo chamado"
		body = fmt.Sprintf("#%d - %s", tickets[0].ID, tickets[0].Name)
	} else {
		title = fmt.Sprintf("%d novos chamados", len(tickets))
		body = ""
		for i, t := range tickets {
			if i > 2 {
				body += fmt.Sprintf("... e mais %d", len(tickets)-3)
				break
			}
			if i > 0 {
				body += "\n"
			}
			body += fmt.Sprintf("#%d - %s", t.ID, t.Name)
		}
	}

	notification := gio.NewNotification(title)
	notification.SetBody(body)
	notification.SetDefaultAction("app.show-window")
	n.app.SendNotification("new-tickets", notification)
}
