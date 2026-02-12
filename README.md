# Paterson e-paper Bus Display

First Time Setup:

1. Download Go

```shell
chmod +x && ./install_go_pi.sh
```

2. Let non-root control the pi gpio

```shell
sudo usermod -aG gpio hunter
```

3. Start it

```shell
sudo ./svc.sh
```

To make changes / edits:

1. `sudo ./svc.sh`

## Debugging

1. ssh into the machine (usually a pi) running the service.

2. `journalctl -u bus`
