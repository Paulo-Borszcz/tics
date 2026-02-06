# Tics

A native GNOME desktop application for managing GLPI helpdesk tickets. Built with Go, GTK4 and libadwaita for a modern Linux desktop experience.

* [Download](https://github.com/Paulo-Borszcz/tics/releases)
* [Install](#installation)
* [Build from Source](#building-from-source)

## Features

* **Real-time ticket synchronization** with configurable polling interval
* **Desktop notifications** for newly created tickets
* **Quick response templates** for common replies
* **Auto-followup** sends a customizable acknowledgment card to assigned tickets automatically
* **Followup HTML editor** to personalize the auto-response message
* **GLPI session management** with automatic re-authentication on expiry
* **Split-pane interface** with ticket list and detail view side by side
* **Priority and status indicators** with color-coded badges and stripes
* **Open in browser** button to jump to the ticket in GLPI's web interface
* **First-run setup wizard** with connection validation

## System Requirements

* **Operating System:** Linux (64-bit)
* **Desktop:** GNOME or any GTK4-compatible environment
* **Libraries:** GTK4 4.x, libadwaita 1.6+

## Installation

### One-line install

Downloads the latest release binary, installs it to `~/.local/bin` and creates a desktop shortcut:

```
curl -fsSL https://raw.githubusercontent.com/Paulo-Borszcz/tics/main/install.sh | bash
```

To uninstall:

```
curl -fsSL https://raw.githubusercontent.com/Paulo-Borszcz/tics/main/install.sh | bash -s -- --uninstall
```

### Manual install

Download the binary from the [releases](https://github.com/Paulo-Borszcz/tics/releases) page, make it executable and place it somewhere in your `PATH`:

```
chmod +x tics-linux-amd64
mv tics-linux-amd64 ~/.local/bin/tics
```

## Building from Source

### Dependencies

Install the required development libraries for your distribution:

**Fedora / RHEL:**
```
sudo dnf install gtk4-devel libadwaita-devel gobject-introspection-devel gcc
```

**Ubuntu / Debian:**
```
sudo apt install libgtk-4-dev libadwaita-1-dev libgirepository1.0-dev gcc
```

**Arch Linux:**
```
sudo pacman -S gtk4 libadwaita gobject-introspection gcc
```

You will also need [Go](https://go.dev/dl/) 1.25 or later.

### Build and install

```
git clone https://github.com/Paulo-Borszcz/tics.git
cd tics
make build
make install
```

The first build takes several minutes due to CGo compilation of the GTK4 bindings. Subsequent builds are significantly faster.

## Configuration

On first launch, Tics presents a setup dialog where you enter:

* **GLPI URL** -- the base URL of your GLPI REST API (e.g. `https://glpi.example.com/apirest.php`)
* **User Token** -- your GLPI personal API token
* **App Token** -- optional application token if required by your GLPI instance

Configuration is stored in `~/.config/tics/config.json`. You can re-enter the setup from the settings menu at any time.

The sync interval defaults to 30 seconds and can be changed via the `TICS_SYNC_INTERVAL` environment variable (value in seconds).

## How It Works

Tics connects to the GLPI REST API and periodically fetches tickets that are either new or assigned to the authenticated user. When new tickets are detected, a desktop notification is shown. The auto-followup feature automatically sends an acknowledgment card to tickets in "processing (assigned)" status that haven't received one yet, so requesters know their ticket is being handled.

## License

This project is provided as-is for internal use.
