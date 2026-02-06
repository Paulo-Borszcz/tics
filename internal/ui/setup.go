package ui

import (
	"log"
	"time"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type SetupPage struct {
	box         *gtk.Box
	urlRow      *adw.EntryRow
	tokenRow    *adw.PasswordEntryRow
	appRow      *adw.EntryRow
	statusLabel *gtk.Label
	saveBtn     *gtk.Button
	onSuccess   func(cfg *config.Config)
}

func NewSetupPage(cfg *config.Config, onSuccess func(cfg *config.Config)) *SetupPage {
	sp := &SetupPage{onSuccess: onSuccess}

	sp.box = gtk.NewBox(gtk.OrientationVertical, 0)

	page := adw.NewStatusPage()
	page.SetIconName("network-server-symbolic")
	page.SetTitle("Tics")
	page.SetDescription("Configure a conexao com o GLPI")

	formBox := gtk.NewBox(gtk.OrientationVertical, 16)

	// Preferences group with entry rows
	group := adw.NewPreferencesGroup()

	sp.urlRow = adw.NewEntryRow()
	sp.urlRow.SetTitle("URL da API GLPI")
	if cfg.GLPIURL != "" {
		sp.urlRow.SetText(cfg.GLPIURL)
	}
	group.Add(sp.urlRow)

	sp.tokenRow = adw.NewPasswordEntryRow()
	sp.tokenRow.SetTitle("User Token")
	if cfg.UserToken != "" {
		sp.tokenRow.SetText(cfg.UserToken)
	}
	group.Add(sp.tokenRow)

	sp.appRow = adw.NewEntryRow()
	sp.appRow.SetTitle("App Token (opcional)")
	if cfg.AppToken != "" {
		sp.appRow.SetText(cfg.AppToken)
	}
	group.Add(sp.appRow)

	formBox.Append(group)

	// Status
	sp.statusLabel = gtk.NewLabel("")
	sp.statusLabel.SetMarginTop(4)
	formBox.Append(sp.statusLabel)

	// Save button
	sp.saveBtn = gtk.NewButtonWithLabel("Validar e Salvar")
	sp.saveBtn.AddCSSClass("suggested-action")
	sp.saveBtn.AddCSSClass("pill")
	sp.saveBtn.SetHAlign(gtk.AlignCenter)
	sp.saveBtn.ConnectClicked(sp.onValidate)
	formBox.Append(sp.saveBtn)

	page.SetChild(formBox)
	sp.box.Append(page)

	return sp
}

func (sp *SetupPage) onValidate() {
	url := sp.urlRow.Text()
	token := sp.tokenRow.Text()
	appToken := sp.appRow.Text()

	if url == "" || token == "" {
		sp.statusLabel.SetText("URL e User Token sao obrigatorios.")
		sp.statusLabel.AddCSSClass("error")
		return
	}

	sp.saveBtn.SetSensitive(false)
	sp.statusLabel.SetText("Validando conexao...")
	sp.statusLabel.RemoveCSSClass("error")
	sp.statusLabel.RemoveCSSClass("success")

	cfg := &config.Config{
		GLPIURL:      url,
		UserToken:    token,
		AppToken:     appToken,
		SyncInterval: 30 * time.Second,
	}

	go func() {
		client := glpi.NewClient(cfg)
		err := client.InitSession()
		if err == nil {
			client.KillSession()
		}

		glib.IdleAdd(func() {
			if err != nil {
				log.Printf("setup: validation failed: %v", err)
				sp.statusLabel.SetText("Falha na conexao: " + err.Error())
				sp.statusLabel.AddCSSClass("error")
				sp.statusLabel.RemoveCSSClass("success")
				sp.saveBtn.SetSensitive(true)
				return
			}

			if err := cfg.Save(); err != nil {
				log.Printf("setup: save failed: %v", err)
				sp.statusLabel.SetText("Erro ao salvar: " + err.Error())
				sp.statusLabel.AddCSSClass("error")
				sp.saveBtn.SetSensitive(true)
				return
			}

			sp.statusLabel.SetText("Conexao validada!")
			sp.statusLabel.AddCSSClass("success")
			sp.statusLabel.RemoveCSSClass("error")

			if sp.onSuccess != nil {
				sp.onSuccess(cfg)
			}
		})
	}()
}
