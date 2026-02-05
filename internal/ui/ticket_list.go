package ui

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type TicketList struct {
	listBox  *gtk.ListBox
	tickets  []glpi.Ticket
	onSelect func(glpi.Ticket)
}

func NewTicketList() *TicketList {
	tl := &TicketList{}

	tl.listBox = gtk.NewListBox()
	tl.listBox.SetSelectionMode(gtk.SelectionSingle)
	tl.listBox.SetActivateOnSingleClick(true)

	tl.listBox.ConnectRowActivated(func(row *gtk.ListBoxRow) {
		idx := row.Index()
		if idx >= 0 && idx < len(tl.tickets) && tl.onSelect != nil {
			tl.onSelect(tl.tickets[idx])
		}
	})

	placeholder := gtk.NewLabel("Nenhum chamado encontrado")
	placeholder.AddCSSClass("dim-label")
	tl.listBox.SetPlaceholder(placeholder)

	return tl
}

func (tl *TicketList) OnSelect(fn func(glpi.Ticket)) {
	tl.onSelect = fn
}

func (tl *TicketList) Update(tickets []glpi.Ticket) {
	tl.tickets = tickets

	for {
		child := tl.listBox.FirstChild()
		if child == nil {
			break
		}
		tl.listBox.Remove(child)
	}

	for _, t := range tickets {
		row := tl.createRow(t)
		tl.listBox.Append(row)
	}
}

func (tl *TicketList) createRow(t glpi.Ticket) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()

	outer := gtk.NewBox(gtk.OrientationHorizontal, 0)

	// Priority color stripe
	stripe := gtk.NewBox(gtk.OrientationVertical, 0)
	stripe.SetSizeRequest(4, -1)
	stripe.AddCSSClass("priority-stripe")
	stripe.AddCSSClass("priority-stripe-" + priorityLevel(t.Priority))
	outer.Append(stripe)

	// Content
	box := gtk.NewBox(gtk.OrientationVertical, 3)
	box.AddCSSClass("ticket-row")

	// Top: ID + date
	topBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	idLabel := gtk.NewLabel(fmt.Sprintf("#%d", t.ID))
	idLabel.AddCSSClass("ticket-id")
	idLabel.SetHAlign(gtk.AlignStart)
	topBox.Append(idLabel)

	if t.DateMod != "" {
		dateLabel := gtk.NewLabel(formatDate(t.DateMod))
		dateLabel.AddCSSClass("ticket-meta")
		dateLabel.SetHAlign(gtk.AlignEnd)
		dateLabel.SetHExpand(true)
		topBox.Append(dateLabel)
	}
	box.Append(topBox)

	// Title
	titleLabel := gtk.NewLabel(t.Name)
	titleLabel.AddCSSClass("ticket-title")
	titleLabel.SetHAlign(gtk.AlignStart)
	titleLabel.SetEllipsize(3) // PANGO_ELLIPSIZE_END
	titleLabel.SetMaxWidthChars(42)
	box.Append(titleLabel)

	// Bottom: status + priority badges
	bottomBox := gtk.NewBox(gtk.OrientationHorizontal, 6)

	statusLabel := gtk.NewLabel(t.Status)
	statusLabel.AddCSSClass("status-badge")
	statusLabel.AddCSSClass(statusCSSClass(t.StatusCode))
	statusLabel.SetHAlign(gtk.AlignStart)
	bottomBox.Append(statusLabel)

	prioLabel := gtk.NewLabel(glpi.PriorityName(t.Priority))
	prioLabel.AddCSSClass("priority-badge")
	prioLabel.AddCSSClass(priorityCSSClass(t.Priority))
	prioLabel.SetHAlign(gtk.AlignStart)
	bottomBox.Append(prioLabel)

	box.Append(bottomBox)

	outer.Append(box)
	row.SetChild(outer)
	return row
}

func priorityCSSClass(p int) string {
	switch p {
	case 1:
		return "priority-very-low"
	case 2:
		return "priority-low"
	case 3:
		return "priority-medium"
	case 4:
		return "priority-high"
	case 5:
		return "priority-very-high"
	case 6:
		return "priority-critical"
	default:
		return "priority-medium"
	}
}

func priorityLevel(p int) string {
	switch p {
	case 1:
		return "very-low"
	case 2:
		return "low"
	case 3:
		return "medium"
	case 4:
		return "high"
	case 5:
		return "very-high"
	case 6:
		return "critical"
	default:
		return "medium"
	}
}

func statusCSSClass(code int) string {
	switch code {
	case 1:
		return "status-new"
	case 2, 3:
		return "status-processing"
	case 4:
		return "status-pending"
	case 5:
		return "status-solved"
	case 6:
		return "status-closed"
	default:
		return ""
	}
}

func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return dateStr
	}
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	if t.Year() == now.Year() {
		return t.Format("02/01 15:04")
	}
	return t.Format("02/01/2006")
}
