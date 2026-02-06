#!/usr/bin/env bash
set -euo pipefail

REPO="Paulo-Borszcz/tics"
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
APP_ID="com.github.pauloborszcz.Tics"

# --- Cores ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

# Cores da barra de progresso
PROGRESS_COLOR_BAR="\e[38;5;114m"
PROGRESS_COLOR_TEXT="\e[38;5;255m"
PROGRESS_COLOR_INFO="\e[38;5;245m"
PROGRESS_COLOR_RESET="\e[0m"

PROGRESS_CHAR_FILLED="█"
PROGRESS_CHAR_EMPTY="░"
PROGRESS_BAR_WIDTH=30

info()  { echo -e "${GREEN}[+]${NC} $*"; }
warn()  { echo -e "${YELLOW}[!]${NC} $*"; }
error() { echo -e "${RED}[x]${NC} $*"; exit 1; }

# --- format_bytes: converte bytes para formato legivel ---
format_bytes() {
    local bytes=$1
    if [[ "$bytes" -ge 1073741824 ]]; then
        local int=$((bytes / 1073741824))
        local dec=$(( (bytes % 1073741824) * 10 / 1073741824 ))
        printf "%d.%d GB" "$int" "$dec"
    elif [[ "$bytes" -ge 1048576 ]]; then
        local int=$((bytes / 1048576))
        local dec=$(( (bytes % 1048576) * 10 / 1048576 ))
        printf "%d.%d MB" "$int" "$dec"
    elif [[ "$bytes" -ge 1024 ]]; then
        local int=$((bytes / 1024))
        local dec=$(( (bytes % 1024) * 10 / 1024 ))
        printf "%d.%d KB" "$int" "$dec"
    else
        printf "%d B" "$bytes"
    fi
}

# --- Obter tamanho remoto do arquivo ---
get_remote_file_size() {
    local url="$1"
    curl -sI -L "$url" 2>/dev/null | grep -i content-length | tail -1 | awk '{print $2}' | tr -d '\r'
}

# --- Desenhar barra de progresso ---
draw_progress_bar() {
    local percent=$1
    local downloaded=$2
    local total=$3
    local speed=$4
    local name=$5

    local filled=$(( percent * PROGRESS_BAR_WIDTH / 100 ))
    local empty=$(( PROGRESS_BAR_WIDTH - filled ))

    local bar=""
    for ((i=0; i<filled; i++)); do
        bar+="$PROGRESS_CHAR_FILLED"
    done
    for ((i=0; i<empty; i++)); do
        bar+="$PROGRESS_CHAR_EMPTY"
    done

    local downloaded_fmt
    downloaded_fmt=$(format_bytes "$downloaded")
    local total_fmt
    total_fmt=$(format_bytes "$total")
    local speed_fmt="${speed}/s"

    printf "\r${PROGRESS_COLOR_TEXT}  %-12s ${PROGRESS_COLOR_BAR}%s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_TEXT}%3d%%${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_INFO}%8s/%s %10s${PROGRESS_COLOR_RESET}" \
        "$name" "$bar" "$percent" "$downloaded_fmt" "$total_fmt" "$speed_fmt"
}

# --- Spinner para quando tamanho e desconhecido ---
draw_spinner_bar() {
    local frame=$1
    local name=$2
    local spinner_chars=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
    local spinner="${spinner_chars[$((frame % ${#spinner_chars[@]}))]}"

    printf "\r${PROGRESS_COLOR_BAR}  %s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_TEXT}%-12s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_INFO}Baixando...${PROGRESS_COLOR_RESET}    " \
        "$spinner" "$name"
}

# --- Download com barra de progresso visual ---
download_with_progress() {
    local url="$1"
    local destino="$2"
    local nome="${3:-$(basename "$destino")}"

    if [[ ${#nome} -gt 12 ]]; then
        nome="${nome:0:11}…"
    fi

    local tamanho_total
    tamanho_total=$(get_remote_file_size "$url")

    echo ""

    if [[ -n "$tamanho_total" ]] && [[ "$tamanho_total" -gt 0 ]] 2>/dev/null; then
        local start_time
        start_time=$(date +%s)

        curl -fsSL "$url" -o "$destino" 2>/dev/null &
        local curl_pid=$!

        while kill -0 "$curl_pid" 2>/dev/null; do
            if [[ -f "$destino" ]]; then
                local current_size
                current_size=$(stat -c%s "$destino" 2>/dev/null || stat -f%z "$destino" 2>/dev/null || echo 0)

                local percent=0
                if [[ "$tamanho_total" -gt 0 ]] && [[ "$current_size" -gt 0 ]]; then
                    percent=$((current_size * 100 / tamanho_total))
                    [[ "$percent" -gt 100 ]] && percent=100
                fi

                local current_time
                current_time=$(date +%s)
                local elapsed=$((current_time - start_time))
                local speed="--"
                if [[ "$elapsed" -gt 0 ]] && [[ "$current_size" -gt 0 ]]; then
                    local bytes_per_sec=$((current_size / elapsed))
                    speed=$(format_bytes "$bytes_per_sec")
                fi

                draw_progress_bar "$percent" "$current_size" "$tamanho_total" "$speed" "$nome"
            fi
            sleep 0.1
        done

        wait "$curl_pid" && local curl_status=0 || local curl_status=$?

        if [[ $curl_status -eq 0 ]] && [[ -f "$destino" ]] && [[ -s "$destino" ]]; then
            local final_size
            final_size=$(stat -c%s "$destino" 2>/dev/null || stat -f%z "$destino" 2>/dev/null || echo "$tamanho_total")
            local end_time
            end_time=$(date +%s)
            local total_time=$((end_time - start_time))
            [[ "$total_time" -eq 0 ]] && total_time=1
            local avg_speed
            avg_speed=$(format_bytes $((final_size / total_time)))
            draw_progress_bar 100 "$final_size" "$tamanho_total" "$avg_speed" "$nome"
            echo ""
            return 0
        else
            printf "\r${PROGRESS_COLOR_BAR}  ✗${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_TEXT}%-12s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_INFO}Falhou${PROGRESS_COLOR_RESET}                              \n" "$nome"
            rm -f "$destino"
            return 1
        fi
    else
        local temp_file="${destino}.tmp"
        local frame=0
        local curl_pid

        curl -fsSL "$url" -o "$temp_file" 2>/dev/null &
        curl_pid=$!

        while kill -0 "$curl_pid" 2>/dev/null; do
            draw_spinner_bar $frame "$nome"
            frame=$((frame + 1))
            sleep 0.1
        done

        wait "$curl_pid" && local curl_status=0 || local curl_status=$?

        if [[ $curl_status -eq 0 ]] && [[ -f "$temp_file" ]]; then
            local final_size
            final_size=$(stat -c%s "$temp_file" 2>/dev/null || stat -f%z "$temp_file" 2>/dev/null || echo "?")
            printf "\r${PROGRESS_COLOR_BAR}  ✓${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_TEXT}%-12s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_INFO}Concluido (%s)${PROGRESS_COLOR_RESET}       \n" \
                "$nome" "$(format_bytes "$final_size")"
            mv "$temp_file" "$destino"
            return 0
        else
            printf "\r${PROGRESS_COLOR_BAR}  ✗${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_TEXT}%-12s${PROGRESS_COLOR_RESET} ${PROGRESS_COLOR_INFO}Falhou${PROGRESS_COLOR_RESET}                    \n" "$nome"
            rm -f "$temp_file"
            return 1
        fi
    fi
}

# --- Detect architecture ---
detect_arch() {
    local arch
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *) error "Arquitetura nao suportada: $arch" ;;
    esac
}

# --- Detect OS ---
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)  echo "linux" ;;
        *) error "SO nao suportado: $os (apenas Linux)" ;;
    esac
}

# --- Check dependencies ---
check_deps() {
    for cmd in curl; do
        command -v "$cmd" &>/dev/null || error "'$cmd' nao encontrado. Instale e tente novamente."
    done

    if ! pkg-config --exists gtk4 2>/dev/null; then
        warn "GTK4 nao detectado. O Tics precisa do GTK4 para funcionar."
        echo "  Fedora/RHEL:   sudo dnf install gtk4-devel"
        echo "  Ubuntu/Debian: sudo apt install libgtk-4-dev"
        echo "  Arch:          sudo pacman -S gtk4"
        echo ""
        read -rp "Continuar mesmo assim? [s/N] " answer
        if [[ ! "${answer,,}" =~ ^(s|sim|y|yes)$ ]]; then
            exit 1
        fi
    fi
}

# --- Get latest release tag ---
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
        | grep '"tag_name"' | head -1 | cut -d'"' -f4)

    if [ -z "$version" ]; then
        error "Nao foi possivel obter a versao mais recente. Verifique sua conexao ou se existem releases em github.com/${REPO}/releases"
    fi
    echo "$version"
}

# --- Download and install binary ---
install_binary() {
    local version="$1" os="$2" arch="$3"
    local filename="tics-${os}-${arch}"
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"

    mkdir -p "$INSTALL_DIR"

    if ! download_with_progress "$url" "${INSTALL_DIR}/tics" "tics"; then
        error "Falha ao baixar. Verifique se o release existe em github.com/${REPO}/releases"
    fi

    chmod +x "${INSTALL_DIR}/tics"
    info "Binario instalado em ${INSTALL_DIR}/tics"
}

# --- Install .desktop file ---
install_desktop() {
    mkdir -p "$DESKTOP_DIR"
    cat > "${DESKTOP_DIR}/${APP_ID}.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=Tics
Comment=GLPI Ticket Manager
Exec=${INSTALL_DIR}/tics
Icon=preferences-system-notifications
Categories=Utility;GTK;
Keywords=glpi;tickets;helpdesk;
StartupNotify=true
Terminal=false
EOF

    if command -v update-desktop-database &>/dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi

    info "Atalho criado no menu de aplicativos"
}

# --- Check PATH ---
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        warn "${INSTALL_DIR} nao esta no seu PATH."
        echo "    Adicione ao seu ~/.bashrc ou ~/.zshrc:"
        echo ""
        echo "      export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
    fi
}

# --- Uninstall ---
uninstall() {
    info "Desinstalando Tics..."

    rm -f "${INSTALL_DIR}/tics"
    rm -f "${DESKTOP_DIR}/${APP_ID}.desktop"

    if command -v update-desktop-database &>/dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi

    info "Tics desinstalado com sucesso."
    echo "    Configuracoes em ~/.config/tics/ foram mantidas."
    echo "    Para remover tudo: rm -rf ~/.config/tics"
}

# --- Main ---
main() {
    echo -e "${BOLD}"
    echo "  _____ _"
    echo " |_   _(_) ___ ___"
    echo "   | | | |/ __/ __|"
    echo "   | | | | (__\__ \\"
    echo "   |_| |_|\___|___/"
    echo -e "${NC}"
    echo "  GLPI Ticket Manager"
    echo ""

    if [ "${1:-}" = "--uninstall" ] || [ "${1:-}" = "uninstall" ]; then
        uninstall
        exit 0
    fi

    check_deps

    local os arch version
    os=$(detect_os)
    arch=$(detect_arch)
    version=$(get_latest_version)

    info "Versao: ${version} | SO: ${os} | Arch: ${arch}"

    install_binary "$version" "$os" "$arch"
    install_desktop

    echo ""
    check_path

    info "Instalacao concluida!"
    echo ""
    echo "  Abra o Tics pelo menu de apps ou execute 'tics'"
    echo "  Na primeira execucao, configure a URL e token do GLPI."
    echo ""
    echo "  Para desinstalar: $0 --uninstall"
}

main "$@"
