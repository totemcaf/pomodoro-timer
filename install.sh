#!/bin/bash

# Pomodoro Timer Installation Script
# This script installs the Pomodoro Timer application and creates desktop integration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="pomodoro"
BINARY_SOURCE="bin/pomodoro"
INSTALL_DIR="/usr/local/bin"
DESKTOP_FILE="pomodoro.desktop"
APPLICATIONS_DIR="/usr/share/applications"
LOCAL_APPLICATIONS_DIR="$HOME/.local/share/applications"
ICON_DIR="/usr/share/pixmaps"
LOCAL_ICON_DIR="$HOME/.local/share/icons"

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_requirements() {
    print_info "Checking requirements..."
    
    # Check if binary exists
    if [[ ! -f "$BINARY_SOURCE" ]]; then
        print_error "Binary not found at $BINARY_SOURCE"
        print_info "Please build the application first with: make build"
        exit 1
    fi
    
    # Check if desktop file exists
    if [[ ! -f "$DESKTOP_FILE" ]]; then
        print_error "Desktop file not found at $DESKTOP_FILE"
        exit 1
    fi
    
    print_success "All requirements met"
}

install_system_wide() {
    print_info "Installing system-wide (requires sudo)..."
    
    # Install binary
    print_info "Installing binary to $INSTALL_DIR..."
    sudo cp "$BINARY_SOURCE" "$INSTALL_DIR/$APP_NAME"
    sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    print_success "Binary installed to $INSTALL_DIR/$APP_NAME"
    
    # Install desktop file
    print_info "Installing desktop file to $APPLICATIONS_DIR..."
    sudo cp "$DESKTOP_FILE" "$APPLICATIONS_DIR/"
    sudo chmod 644 "$APPLICATIONS_DIR/$DESKTOP_FILE"
    print_success "Desktop file installed to $APPLICATIONS_DIR/$DESKTOP_FILE"
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        print_info "Updating desktop database..."
        sudo update-desktop-database "$APPLICATIONS_DIR"
        print_success "Desktop database updated"
    fi
}

install_user_local() {
    print_info "Installing for current user only..."
    
    # Create directories if they don't exist
    mkdir -p "$HOME/.local/bin"
    mkdir -p "$LOCAL_APPLICATIONS_DIR"
    
    # Install binary
    print_info "Installing binary to $HOME/.local/bin..."
    cp "$BINARY_SOURCE" "$HOME/.local/bin/$APP_NAME"
    chmod +x "$HOME/.local/bin/$APP_NAME"
    print_success "Binary installed to $HOME/.local/bin/$APP_NAME"
    
    # Update desktop file to use local path
    print_info "Creating user desktop file..."
    sed "s|/usr/local/bin/pomodoro|$HOME/.local/bin/pomodoro|g" "$DESKTOP_FILE" > "$LOCAL_APPLICATIONS_DIR/$DESKTOP_FILE"
    chmod 644 "$LOCAL_APPLICATIONS_DIR/$DESKTOP_FILE"
    print_success "Desktop file installed to $LOCAL_APPLICATIONS_DIR/$DESKTOP_FILE"
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        print_info "Updating user desktop database..."
        update-desktop-database "$LOCAL_APPLICATIONS_DIR"
        print_success "Desktop database updated"
    fi
    
    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        print_warning "$HOME/.local/bin is not in your PATH"
        print_info "Add the following line to your ~/.bashrc or ~/.zshrc:"
        echo "export PATH=\"\$HOME/.local/bin:\$PATH\""
        print_info "Then run: source ~/.bashrc (or source ~/.zshrc)"
    fi
}

uninstall_system() {
    print_info "Uninstalling system-wide installation..."
    
    # Remove binary
    if [[ -f "$INSTALL_DIR/$APP_NAME" ]]; then
        sudo rm "$INSTALL_DIR/$APP_NAME"
        print_success "Binary removed from $INSTALL_DIR"
    fi
    
    # Remove desktop file
    if [[ -f "$APPLICATIONS_DIR/$DESKTOP_FILE" ]]; then
        sudo rm "$APPLICATIONS_DIR/$DESKTOP_FILE"
        print_success "Desktop file removed from $APPLICATIONS_DIR"
    fi
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        sudo update-desktop-database "$APPLICATIONS_DIR"
        print_success "Desktop database updated"
    fi
}

uninstall_user() {
    print_info "Uninstalling user installation..."
    
    # Remove binary
    if [[ -f "$HOME/.local/bin/$APP_NAME" ]]; then
        rm "$HOME/.local/bin/$APP_NAME"
        print_success "Binary removed from $HOME/.local/bin"
    fi
    
    # Remove desktop file
    if [[ -f "$LOCAL_APPLICATIONS_DIR/$DESKTOP_FILE" ]]; then
        rm "$LOCAL_APPLICATIONS_DIR/$DESKTOP_FILE"
        print_success "Desktop file removed from $LOCAL_APPLICATIONS_DIR"
    fi
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$LOCAL_APPLICATIONS_DIR"
        print_success "Desktop database updated"
    fi
}

show_help() {
    echo "Pomodoro Timer Installation Script"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  install-system    Install system-wide (requires sudo)"
    echo "  install-user      Install for current user only"
    echo "  uninstall-system  Remove system-wide installation"
    echo "  uninstall-user    Remove user installation"
    echo "  help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 install-user       # Install for current user"
    echo "  $0 install-system     # Install system-wide"
    echo "  $0 uninstall-user     # Remove user installation"
}

# Main script
case "${1:-help}" in
    "install-system")
        check_requirements
        install_system_wide
        print_success "Installation complete! You can now find 'Pomodoro Timer' in your applications menu."
        ;;
    "install-user")
        check_requirements
        install_user_local
        print_success "Installation complete! You can now find 'Pomodoro Timer' in your applications menu."
        ;;
    "uninstall-system")
        uninstall_system
        print_success "System-wide uninstallation complete."
        ;;
    "uninstall-user")
        uninstall_user
        print_success "User uninstallation complete."
        ;;
    "help"|*)
        show_help
        ;;
esac
