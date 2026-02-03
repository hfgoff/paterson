#!/usr/bin/env bash
set -e

SERVICE_NAME=bus
INSTALL_DIR=/opt/bus
SERVICE_FILE=/etc/systemd/system/${SERVICE_NAME}.service
RUN_USER=hunter

echo "Installing ${SERVICE_NAME} service..."

# Ensure running as root
if [[ $EUID -ne 0 ]]; then
  echo "Please run as root: sudo ./svc.sh"
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
Description=Bus Service
After=network.target

[Service]
Type=simple
User=${RUN_USER}
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/run.sh
KillMode=control-group
Restart=always
RestartSec=10
Nice=10
CPUQuota=50%
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "Starting service..."
systemctl daemon-reload
systemctl enable ${SERVICE_NAME}
systemctl restart ${SERVICE_NAME}

echo "Done!"
echo "Logs: journalctl -u ${SERVICE_NAME} -f"
echo "Ex: journalctl -u bus --since \"5 minutes ago\""
