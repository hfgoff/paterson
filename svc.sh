#!/usr/bin/env bash
set -e

SERVICE_NAME=bus
INSTALL_DIR=/opt/bus
SERVICE_FILE=/etc/systemd/system/${SERVICE_NAME}.service
RUN_USER=pi

echo "Installing ${SERVICE_NAME} service..."

# Ensure running as root
if [[ $EUID -ne 0 ]]; then
  echo "Please run as root: sudo ./install.sh"
  exit 1
fi

echo "Creating install directory..."
mkdir -p "$INSTALL_DIR"

echo "Copying project files..."
rsync -a --delete ./ "$INSTALL_DIR/"

echo "Fixing permissions..."
chown -R ${RUN_USER}:${RUN_USER} "$INSTALL_DIR"
chmod +x "$INSTALL_DIR/svc.sh"

echo "Writing systemd service..."
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Bus + E-Paper Service
After=network.target

[Service]
Type=simple
User=${RUN_USER}
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/svc.sh
Restart=always
RestartSec=3
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd..."
systemctl daemon-reexec
systemctl daemon-reload

echo "Enabling service..."
systemctl enable ${SERVICE_NAME}

echo "Starting service..."
systemctl restart ${SERVICE_NAME}

echo "Done!"
echo "Logs: journalctl -u ${SERVICE_NAME} -f"
