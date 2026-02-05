package ui

import (
	"fmt"
	"html"
	"log"
	"os/exec"
	"strings"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type TicketDetail struct {
	box           *gtk.Box
	client        *glpi.Client
	currentTicket *glpi.Ticket
	parentWindow  *gtk.Window
}

func NewTicketDetail(client *glpi.Client, parentWindow *gtk.Window) *TicketDetail {
	td := &TicketDetail{client: client, parentWindow: parentWindow}

	td.box = gtk.NewBox(gtk.OrientationVertical, 0)

	// Empty state using adw.StatusPage
	emptyPage := adw.NewStatusPage()
	emptyPage.SetIconName("mail-unread-symbolic")
	emptyPage.SetTitle("Nenhum chamado selecionado")
	emptyPage.SetDescription("Selecione um chamado na lista para ver os detalhes")
	td.box.Append(emptyPage)

	return td
}

func (td *TicketDetail) ShowTicket(ticket glpi.Ticket) {
	td.currentTicket = &ticket

	td.clear()
	loadingPage := adw.NewStatusPage()
	loadingPage.SetTitle("Carregando...")
	loadingPage.SetIconName("content-loading-symbolic")
	td.box.Append(loadingPage)

	go func() {
		fullTicket, err := td.client.GetTicket(ticket.ID)
		if err != nil {
			log.Printf("Error fetching ticket details: %v", err)
			fullTicket = &ticket
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
	content := gtk.NewBox(gtk.OrientationVertical, 12)
	content.SetMarginTop(20)
	content.SetMarginBottom(20)
	content.SetMarginStart(20)
	content.SetMarginEnd(20)

	// Title
	titleLabel := gtk.NewLabel(fmt.Sprintf("#%d - %s", ticket.ID, ticket.Name))
	titleLabel.AddCSSClass("detail-title")
	titleLabel.SetHAlign(gtk.AlignStart)
	titleLabel.SetWrap(true)
	content.Append(titleLabel)

	// Info card (2-column layout using boxes)
	infoCard := gtk.NewBox(gtk.OrientationVertical, 12)
	infoCard.AddCSSClass("card")
	infoCard.AddCSSClass("info-card")

	// Row 1: Status + Priority
	row1 := gtk.NewBox(gtk.OrientationHorizontal, 24)
	row1.SetHomogeneous(true)
	row1.Append(infoField("STATUS", ticket.Status, "status-badge", statusCSSClass(ticket.StatusCode)))
	row1.Append(infoField("PRIORIDADE", glpi.PriorityName(ticket.Priority), "priority-badge", priorityCSSClass(ticket.Priority)))
	infoCard.Append(row1)

	// Row 2: Date + Entity
	row2 := gtk.NewBox(gtk.OrientationHorizontal, 24)
	row2.SetHomogeneous(true)
	if ticket.DateCreation != "" {
		row2.Append(infoField("CRIADO EM", formatDate(ticket.DateCreation), "info-card-value"))
	}
	if ticket.Entity != "" {
		row2.Append(infoField("ENTIDADE", ticket.Entity, "info-card-value"))
	}
	if row2.FirstChild() != nil {
		infoCard.Append(row2)
	}

	content.Append(infoCard)

	// Description
	if ticket.Content != "" {
		sectionTitle := gtk.NewLabel("Descricao")
		sectionTitle.AddCSSClass("detail-section-title")
		sectionTitle.SetHAlign(gtk.AlignStart)
		content.Append(sectionTitle)

		contentCard := gtk.NewBox(gtk.OrientationVertical, 0)
		contentCard.AddCSSClass("card")
		contentCard.AddCSSClass("info-card")

		contentLabel := gtk.NewLabel(stripHTML(ticket.Content))
		contentLabel.SetHAlign(gtk.AlignStart)
		contentLabel.SetWrap(true)
		contentLabel.SetSelectable(true)
		contentCard.Append(contentLabel)

		content.Append(contentCard)
	}

	// Followups
	countStr := ""
	if len(followups) > 0 {
		countStr = fmt.Sprintf(" (%d)", len(followups))
	}
	followupTitle := gtk.NewLabel("Acompanhamentos" + countStr)
	followupTitle.AddCSSClass("detail-section-title")
	followupTitle.SetHAlign(gtk.AlignStart)
	content.Append(followupTitle)

	if len(followups) == 0 {
		emptyLabel := gtk.NewLabel("Nenhum acompanhamento")
		emptyLabel.AddCSSClass("dim-label")
		emptyLabel.SetMarginTop(4)
		content.Append(emptyLabel)
	} else {
		for i, f := range followups {
			card := td.createFollowupCard(f, i+1, len(followups))
			content.Append(card)
		}
	}

	// Actions
	actionsBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	actionsBox.SetMarginTop(12)

	openBtn := gtk.NewButtonWithLabel("Abrir no navegador")
	openBtn.AddCSSClass("suggested-action")
	openBtn.AddCSSClass("pill")
	openBtn.ConnectClicked(func() {
		td.openInBrowser()
	})
	actionsBox.Append(openBtn)

	replyBtn := gtk.NewButtonWithLabel("Responder")
	replyBtn.AddCSSClass("pill")
	replyBtn.ConnectClicked(func() {
		td.showTemplateDialog()
	})
	actionsBox.Append(replyBtn)

	content.Append(actionsBox)

	td.box.Append(content)
}

func infoField(label, value string, classes ...string) *gtk.Box {
	field := gtk.NewBox(gtk.OrientationVertical, 4)

	lbl := gtk.NewLabel(label)
	lbl.AddCSSClass("info-card-label")
	lbl.SetHAlign(gtk.AlignStart)
	field.Append(lbl)

	val := gtk.NewLabel(value)
	val.SetHAlign(gtk.AlignStart)
	val.SetWrap(true)
	for _, cls := range classes {
		if cls != "" {
			val.AddCSSClass(cls)
		}
	}
	field.Append(val)
	return field
}

func (td *TicketDetail) createFollowupCard(f glpi.Followup, idx, total int) *gtk.Box {
	card := gtk.NewBox(gtk.OrientationVertical, 6)
	card.AddCSSClass("card")
	card.AddCSSClass("followup-card")

	headerBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	numLabel := gtk.NewLabel(fmt.Sprintf("#%d/%d", idx, total))
	numLabel.AddCSSClass("followup-header")
	numLabel.SetHAlign(gtk.AlignStart)
	headerBox.Append(numLabel)

	if f.DateCreation != "" {
		dateLabel := gtk.NewLabel(formatDate(f.DateCreation))
		dateLabel.AddCSSClass("followup-date")
		dateLabel.SetHAlign(gtk.AlignEnd)
		dateLabel.SetHExpand(true)
		headerBox.Append(dateLabel)
	}

	card.Append(headerBox)

	contentText := stripHTML(f.Content)
	if contentText != "" {
		contentLabel := gtk.NewLabel(contentText)
		contentLabel.SetHAlign(gtk.AlignStart)
		contentLabel.SetWrap(true)
		contentLabel.SetSelectable(true)
		card.Append(contentLabel)
	}

	return card
}

func (td *TicketDetail) openInBrowser() {
	if td.currentTicket == nil {
		return
	}
	url := fmt.Sprintf("%s/front/ticket.form.php?id=%d", td.client.BaseURL(), td.currentTicket.ID)
	if err := exec.Command("xdg-open", url).Start(); err != nil {
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
	for {
		decoded := html.UnescapeString(s)
		if decoded == s {
			break
		}
		s = decoded
	}

	for _, tag := range []string{"<br>", "<br/>", "<br />", "</p>", "</div>", "</h1>", "</h2>", "</h3>", "</li>", "</tr>"} {
		s = strings.ReplaceAll(s, tag, "\n")
	}

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

	text := result.String()
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(text)
}
