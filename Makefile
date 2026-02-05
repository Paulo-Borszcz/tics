PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
SERVICEDIR = $(HOME)/.config/systemd/user
DESKTOPDIR = $(HOME)/.local/share/applications

.PHONY: build install uninstall enable-service disable-service clean

build:
	go build -o tics ./cmd/tics

install: build
	install -Dm755 tics $(BINDIR)/tics
	install -Dm644 deploy/tics.service $(SERVICEDIR)/tics.service
	install -Dm644 deploy/com.github.pauloborszcz.Tics.desktop $(DESKTOPDIR)/com.github.pauloborszcz.Tics.desktop

uninstall:
	rm -f $(BINDIR)/tics
	rm -f $(SERVICEDIR)/tics.service
	rm -f $(DESKTOPDIR)/com.github.pauloborszcz.Tics.desktop

enable-service:
	systemctl --user daemon-reload
	systemctl --user enable --now tics.service

disable-service:
	systemctl --user disable --now tics.service

clean:
	rm -f tics
