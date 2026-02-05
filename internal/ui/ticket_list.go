package ui

import (
	"fmt"

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

	// Placeholder when empty
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

	// Remove all existing rows
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

	box := gtk.NewBox(gtk.OrientationVertical, 4)
	box.AddCSSClass("ticket-row")

	// Top line: ID + priority
	topBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	idLabel := gtk.NewLabel(fmt.Sprintf("#%d", t.ID))
	idLabel.AddCSSClass("ticket-id")
	idLabel.SetHAlign(gtk.AlignStart)
	topBox.Append(idLabel)
	box.Append(topBox)

	// Title
	titleLabel := gtk.NewLabel(t.Name)
	titleLabel.AddCSSClass("ticket-title")
	titleLabel.SetHAlign(gtk.AlignStart)
	titleLabel.SetEllipsize(3) // PANGO_ELLIPSIZE_END
	titleLabel.SetMaxWidthChars(40)
	box.Append(titleLabel)

	// Bottom: status + date
	bottomBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	statusLabel := gtk.NewLabel(t.Status)
	statusLabel.AddCSSClass("status-badge")
	statusLabel.SetHAlign(gtk.AlignStart)
	bottomBox.Append(statusLabel)

	if t.DateMod != "" {
		dateLabel := gtk.NewLabel(t.DateMod)
		dateLabel.AddCSSClass("ticket-meta")
		dateLabel.SetHAlign(gtk.AlignEnd)
		dateLabel.SetHExpand(true)
		bottomBox.Append(dateLabel)
	}
	box.Append(bottomBox)

	row.SetChild(box)
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
