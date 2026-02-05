package ui

import (
	"fmt"
	"html"
	"log"
	"os/exec"
	"strings"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type TicketDetail struct {
	box           *gtk.Box
	client        *glpi.Client
	currentTicket *glpi.Ticket
	parentWindow  *gtk.Window

	titleLabel    *gtk.Label
	statusLabel   *gtk.Label
	priorityLabel *gtk.Label
	dateLabel     *gtk.Label
	contentLabel  *gtk.Label
	followupsBox  *gtk.Box
	actionsBox    *gtk.Box
}

func NewTicketDetail(client *glpi.Client, parentWindow *gtk.Window) *TicketDetail {
	td := &TicketDetail{client: client, parentWindow: parentWindow}

	td.box = gtk.NewBox(gtk.OrientationVertical, 12)
	td.box.SetMarginTop(16)
	td.box.SetMarginBottom(16)
	td.box.SetMarginStart(16)
	td.box.SetMarginEnd(16)

	// Placeholder
	placeholder := gtk.NewLabel("Selecione um chamado")
	placeholder.AddCSSClass("dim-label")
	td.box.Append(placeholder)

	return td
}

func (td *TicketDetail) ShowTicket(ticket glpi.Ticket) {
	td.currentTicket = &ticket

	// Clear and show loading
	td.clear()
	loading := gtk.NewLabel("Carregando detalhes...")
	loading.AddCSSClass("dim-label")
	td.box.Append(loading)

	// Fetch full ticket details + followups in background
	go func() {
		fullTicket, err := td.client.GetTicket(ticket.ID)
		if err != nil {
			log.Printf("Error fetching ticket details: %v", err)
			fullTicket = &ticket // fallback to search data
		}
		followups, fErr := td.client.GetFollowups(ticket.ID)
		if fErr != nil {
			log.Printf("Error loading followups: %v", fErr)
		}

		glib.IdleAdd(func() {
			td.currentTicket = fullTicket
			td.clear()
			td.renderTicket(fullTicket, followups)
		})
	}()
}

func (td *TicketDetail) clear() {
	for {
		child := td.box.FirstChild()
		if child == nil {
			break
		}
		td.box.Remove(child)
	}
}

func (td *TicketDetail) renderTicket(ticket *glpi.Ticket, followups []glpi.Followup) {
	// Title
	td.titleLabel = gtk.NewLabel(fmt.Sprintf("#%d - %s", ticket.ID, ticket.Name))
	td.titleLabel.AddCSSClass("detail-title")
	td.titleLabel.SetHAlign(gtk.AlignStart)
	td.titleLabel.SetWrap(true)
	td.box.Append(td.titleLabel)

	// Info card with status, date, entity
	infoCard := gtk.NewBox(gtk.OrientationVertical, 8)
	infoCard.AddCSSClass("info-card")

	// Status row
	statusRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
	statusLbl := gtk.NewLabel("Status")
	statusLbl.AddCSSClass("info-card-label")
	statusLbl.SetHAlign(gtk.AlignStart)
	statusRow.Append(statusLbl)
	td.statusLabel = gtk.NewLabel(ticket.Status)
	td.statusLabel.AddCSSClass("status-badge")
	statusRow.Append(td.statusLabel)
	infoCard.Append(statusRow)

	// Date row
	if ticket.DateCreation != "" {
		dateRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
		dateLbl := gtk.NewLabel("Criado em")
		dateLbl.AddCSSClass("info-card-label")
		dateLbl.SetHAlign(gtk.AlignStart)
		dateRow.Append(dateLbl)
		td.dateLabel = gtk.NewLabel(ticket.DateCreation)
		td.dateLabel.AddCSSClass("info-card-value")
		dateRow.Append(td.dateLabel)
		infoCard.Append(dateRow)
	}

	// Entity row
	if ticket.Entity != "" {
		entityRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
		entityLbl := gtk.NewLabel("Entidade")
		entityLbl.AddCSSClass("info-card-label")
		entityLbl.SetHAlign(gtk.AlignStart)
		entityRow.Append(entityLbl)
		entityVal := gtk.NewLabel(ticket.Entity)
		entityVal.AddCSSClass("info-card-value")
		entityVal.SetWrap(true)
		entityRow.Append(entityVal)
		infoCard.Append(entityRow)
	}

	td.box.Append(infoCard)

	// Content (description) in a card
	if ticket.Content != "" {
		contentTitle := gtk.NewLabel("Descricao")
		contentTitle.AddCSSClass("detail-section-title")
		contentTitle.SetHAlign(gtk.AlignStart)
		td.box.Append(contentTitle)

		contentCard := gtk.NewBox(gtk.OrientationVertical, 0)
		contentCard.AddCSSClass("content-card")

		td.contentLabel = gtk.NewLabel(stripHTML(ticket.Content))
		td.contentLabel.SetHAlign(gtk.AlignStart)
		td.contentLabel.SetWrap(true)
		td.contentLabel.SetSelectable(true)
		contentCard.Append(td.contentLabel)

		td.box.Append(contentCard)
	}

	// Followups section
	followupTitle := gtk.NewLabel("Acompanhamentos")
	followupTitle.AddCSSClass("detail-section-title")
	followupTitle.SetHAlign(gtk.AlignStart)
	td.box.Append(followupTitle)

	td.followupsBox = gtk.NewBox(gtk.OrientationVertical, 8)
	if len(followups) == 0 {
		emptyLabel := gtk.NewLabel("Nenhum acompanhamento")
		emptyLabel.AddCSSClass("dim-label")
		td.followupsBox.Append(emptyLabel)
	} else {
		for i, f := range followups {
			card := td.createFollowupCard(f, i+1, len(followups))
			td.followupsBox.Append(card)
		}
	}
	td.box.Append(td.followupsBox)

	// Action bar
	td.actionsBox = gtk.NewBox(gtk.OrientationHorizontal, 8)
	td.actionsBox.AddCSSClass("action-bar")
	td.actionsBox.SetHAlign(gtk.AlignFill)

	openBtn := gtk.NewButtonWithLabel("Abrir no Firefox")
	openBtn.AddCSSClass("suggested-action")
	openBtn.ConnectClicked(func() {
		td.openInBrowser()
	})
	td.actionsBox.Append(openBtn)

	replyBtn := gtk.NewButtonWithLabel("Responder")
	replyBtn.ConnectClicked(func() {
		td.showTemplateDialog()
	})
	td.actionsBox.Append(replyBtn)

	td.box.Append(td.actionsBox)
}

func (td *TicketDetail) createFollowupCard(f glpi.Followup, idx, total int) *gtk.Box {
	card := gtk.NewBox(gtk.OrientationVertical, 6)
	card.AddCSSClass("followup-card")

	// Header with number and date
	headerBox := gtk.NewBox(gtk.OrientationHorizontal, 8)

	numLabel := gtk.NewLabel(fmt.Sprintf("#%d/%d", idx, total))
	numLabel.AddCSSClass("followup-header")
	numLabel.SetHAlign(gtk.AlignStart)
	headerBox.Append(numLabel)

	dateLabel := gtk.NewLabel(f.DateCreation)
	dateLabel.AddCSSClass("followup-date")
	dateLabel.SetHAlign(gtk.AlignStart)
	dateLabel.SetHExpand(true)
	headerBox.Append(dateLabel)

	card.Append(headerBox)

	content := gtk.NewLabel(stripHTML(f.Content))
	content.SetHAlign(gtk.AlignStart)
	content.SetWrap(true)
	content.SetSelectable(true)
	card.Append(content)

	return card
}

func (td *TicketDetail) openInBrowser() {
	if td.currentTicket == nil {
		return
	}
	url := fmt.Sprintf("%s/front/ticket.form.php?id=%d", td.client.BaseURL(), td.currentTicket.ID)
	if err := exec.Command("firefox", url).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func (td *TicketDetail) showTemplateDialog() {
	if td.currentTicket == nil {
		return
	}
	ShowTemplateDialog(td.parentWindow, td.client, td.currentTicket.ID, func() {
		td.ShowTicket(*td.currentTicket)
	})
}

// stripHTML decodes HTML entities and removes HTML tags, preserving line breaks.
func stripHTML(s string) string {
	// Decode HTML entities repeatedly (GLPI double-encodes: &#38;#62; -> &#62; -> >)
	for {
		decoded := html.UnescapeString(s)
		if decoded == s {
			break
		}
		s = decoded
	}

	// Insert newlines before block-level tags
	for _, tag := range []string{"<br>", "<br/>", "<br />", "</p>", "</div>", "</h1>", "</h2>", "</h3>", "</li>", "</tr>"} {
		s = strings.ReplaceAll(s, tag, "\n")
	}

	// Remove remaining HTML tags
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}

	// Clean up excessive blank lines
	text := result.String()
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(text)
}
