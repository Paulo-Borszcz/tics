package ui

import (
	"fmt"
	"log"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
)

// ShowTemplateDialog opens a dialog to select and send a template response.
func ShowTemplateDialog(parent *gtk.Window, client *glpi.Client, ticketID int, onSent func()) {
	dialog := gtk.NewWindow()
	dialog.SetTitle("Responder chamado")
	dialog.SetDefaultSize(400, 300)
	dialog.SetModal(true)
	if parent != nil {
		dialog.SetTransientFor(parent)
	}

	box := gtk.NewBox(gtk.OrientationVertical, 12)
	box.SetMarginTop(16)
	box.SetMarginBottom(16)
	box.SetMarginStart(16)
	box.SetMarginEnd(16)

	titleLabel := gtk.NewLabel(fmt.Sprintf("Responder chamado #%d", ticketID))
	titleLabel.AddCSSClass("detail-title")
	titleLabel.SetHAlign(gtk.AlignStart)
	box.Append(titleLabel)

	infoLabel := gtk.NewLabel("Selecione uma resposta template:")
	infoLabel.SetHAlign(gtk.AlignStart)
	box.Append(infoLabel)

	// Template list
	listBox := gtk.NewListBox()
	listBox.SetSelectionMode(gtk.SelectionSingle)

	for _, tmpl := range config.Templates {
		row := gtk.NewListBoxRow()
		label := gtk.NewLabel(tmpl)
		label.SetHAlign(gtk.AlignStart)
		label.SetWrap(true)
		label.SetMarginTop(8)
		label.SetMarginBottom(8)
		label.SetMarginStart(8)
		label.SetMarginEnd(8)
		row.SetChild(label)
		listBox.Append(row)
	}

	scrolled := gtk.NewScrolledWindow()
	scrolled.SetChild(listBox)
	scrolled.SetVExpand(true)
	box.Append(scrolled)

	// Status label
	statusLabel := gtk.NewLabel("")
	statusLabel.SetHAlign(gtk.AlignStart)
	box.Append(statusLabel)

	// Buttons
	btnBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	btnBox.SetHAlign(gtk.AlignEnd)

	cancelBtn := gtk.NewButtonWithLabel("Cancelar")
	cancelBtn.ConnectClicked(func() {
		dialog.Close()
	})
	btnBox.Append(cancelBtn)

	sendBtn := gtk.NewButtonWithLabel("Enviar")
	sendBtn.AddCSSClass("suggested-action")
	sendBtn.ConnectClicked(func() {
		selected := listBox.SelectedRow()
		if selected == nil {
			statusLabel.SetText("Selecione uma resposta primeiro.")
			return
		}
		idx := selected.Index()
		if idx < 0 || idx >= len(config.Templates) {
			return
		}
		content := config.Templates[idx]
		sendBtn.SetSensitive(false)
		statusLabel.SetText("Enviando...")

		go func() {
			err := client.AddFollowup(ticketID, content)
			glib.IdleAdd(func() {
				if err != nil {
					statusLabel.SetText(fmt.Sprintf("Erro: %v", err))
					sendBtn.SetSensitive(true)
					log.Printf("Failed to send followup: %v", err)
					return
				}
				statusLabel.SetText("Enviado com sucesso!")
				if onSent != nil {
					onSent()
				}
				// Close after a brief moment
				glib.TimeoutAdd(1000, func() bool {
					dialog.Close()
					return false
				})
			})
		}()
	})
	btnBox.Append(sendBtn)
	box.Append(btnBox)

	dialog.SetChild(box)
	dialog.Show()
}
