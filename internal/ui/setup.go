package ui

import (
	"log"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pauloborszcz/tics/internal/config"
	"github.com/pauloborszcz/tics/internal/glpi"
)

type SetupPage struct {
	box        *gtk.Box
	urlEntry   *gtk.Entry
	tokenEntry *gtk.Entry
	appEntry   *gtk.Entry
	statusLabel *gtk.Label
	saveBtn    *gtk.Button
	onSuccess  func(cfg *config.Config)
}

func NewSetupPage(cfg *config.Config, onSuccess func(cfg *config.Config)) *SetupPage {
	sp := &SetupPage{onSuccess: onSuccess}

	sp.box = gtk.NewBox(gtk.OrientationVertical, 16)
	sp.box.SetMarginTop(48)
	sp.box.SetMarginBottom(48)
	sp.box.SetMarginStart(48)
	sp.box.SetMarginEnd(48)
	sp.box.SetVAlign(gtk.AlignCenter)
	sp.box.SetHAlign(gtk.AlignCenter)
	sp.box.SetSizeRequest(420, -1)

	// Title
	title := gtk.NewLabel("Tics")
	title.AddCSSClass("setup-title")
	sp.box.Append(title)

	subtitle := gtk.NewLabel("Configure a conexao com o GLPI")
	subtitle.AddCSSClass("setup-subtitle")
	subtitle.SetMarginBottom(16)
	sp.box.Append(subtitle)

	// URL field
	urlLabel := gtk.NewLabel("URL da API GLPI")
	urlLabel.SetHAlign(gtk.AlignStart)
	urlLabel.AddCSSClass("setup-field-label")
	sp.box.Append(urlLabel)

	sp.urlEntry = gtk.NewEntry()
	sp.urlEntry.SetPlaceholderText("https://nexus.lojasmm.com.br/apirest.php")
	if cfg.GLPIURL != "" {
		sp.urlEntry.SetText(cfg.GLPIURL)
	}
	sp.box.Append(sp.urlEntry)

	// User Token field
	tokenLabel := gtk.NewLabel("User Token")
	tokenLabel.SetHAlign(gtk.AlignStart)
	tokenLabel.AddCSSClass("setup-field-label")
	tokenLabel.SetMarginTop(8)
	sp.box.Append(tokenLabel)

	sp.tokenEntry = gtk.NewEntry()
	sp.tokenEntry.SetPlaceholderText("Seu user_token do GLPI")
	if cfg.UserToken != "" {
		sp.tokenEntry.SetText(cfg.UserToken)
	}
	sp.box.Append(sp.tokenEntry)

	// App Token field
	appLabel := gtk.NewLabel("App Token (opcional)")
	appLabel.SetHAlign(gtk.AlignStart)
	appLabel.AddCSSClass("setup-field-label")
	appLabel.SetMarginTop(8)
	sp.box.Append(appLabel)

	sp.appEntry = gtk.NewEntry()
	sp.appEntry.SetPlaceholderText("App-Token (se necessario)")
	if cfg.AppToken != "" {
		sp.appEntry.SetText(cfg.AppToken)
	}
	sp.box.Append(sp.appEntry)

	// Status
	sp.statusLabel = gtk.NewLabel("")
	sp.statusLabel.AddCSSClass("setup-status")
	sp.statusLabel.SetMarginTop(8)
	sp.box.Append(sp.statusLabel)

	// Save button
	sp.saveBtn = gtk.NewButtonWithLabel("Validar e Salvar")
	sp.saveBtn.AddCSSClass("suggested-action")
	sp.saveBtn.AddCSSClass("setup-save-btn")
	sp.saveBtn.SetMarginTop(16)
	sp.saveBtn.ConnectClicked(sp.onValidate)
	sp.box.Append(sp.saveBtn)

	return sp
}

func (sp *SetupPage) onValidate() {
	url := sp.urlEntry.Text()
	token := sp.tokenEntry.Text()
	appToken := sp.appEntry.Text()

	if url == "" || token == "" {
		sp.statusLabel.SetText("URL e User Token sao obrigatorios.")
		sp.statusLabel.AddCSSClass("setup-error")
		sp.statusLabel.RemoveCSSClass("setup-success")
		return
	}

	sp.saveBtn.SetSensitive(false)
	sp.statusLabel.SetText("Validando conexao...")
	sp.statusLabel.RemoveCSSClass("setup-error")
	sp.statusLabel.RemoveCSSClass("setup-success")

	cfg := &config.Config{
		GLPIURL:   url,
		UserToken: token,
		AppToken:  appToken,
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
				sp.statusLabel.AddCSSClass("setup-error")
				sp.statusLabel.RemoveCSSClass("setup-success")
				sp.saveBtn.SetSensitive(true)
				return
			}

			if err := cfg.Save(); err != nil {
				log.Printf("setup: save failed: %v", err)
				sp.statusLabel.SetText("Erro ao salvar: " + err.Error())
				sp.statusLabel.AddCSSClass("setup-error")
				sp.saveBtn.SetSensitive(true)
				return
			}

			sp.statusLabel.SetText("Conexao validada!")
			sp.statusLabel.AddCSSClass("setup-success")
			sp.statusLabel.RemoveCSSClass("setup-error")

			if sp.onSuccess != nil {
				sp.onSuccess(cfg)
			}
		})
	}()
}
