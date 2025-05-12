# Running CronAI as a systemd Service

This document explains how to set up CronAI to run as a systemd service on Linux systems.

## Setup Steps

1. **Build and install the CronAI binary**

```bash
cd /path/to/cronai
go build -o cronai ./cmd/cronai
sudo cp cronai /usr/local/bin/cronai
```

2. **Create your configuration and prompt files**

```bash
mkdir -p /etc/cronai/cron_prompts
cp cronai.config.example /etc/cronai/cronai.config
cp -r cron_prompts/* /etc/cronai/cron_prompts/
```

3. **Set up your environment file**

```bash
cp .env.example /etc/cronai/.env
# Edit the .env file with your API keys and settings
sudo nano /etc/cronai/.env
```

4. **Create the systemd service file**

Copy the example service file and modify it for your system:

```bash
sudo cp cronai.service /etc/systemd/system/cronai.service
sudo nano /etc/systemd/system/cronai.service
```

Update the following fields in the service file:
- `User`: The user account that will run the service
- `WorkingDirectory`: The directory where your configuration is located (e.g., `/etc/cronai`)
- `ExecStart`: The path to the CronAI binary (e.g., `/usr/local/bin/cronai start --config /etc/cronai/cronai.config`)
- `EnvironmentFile`: The path to your .env file (e.g., `/etc/cronai/.env`)

5. **Enable and start the service**

```bash
sudo systemctl daemon-reload
sudo systemctl enable cronai
sudo systemctl start cronai
```

6. **Check the service status**

```bash
sudo systemctl status cronai
```

7. **View the logs**

```bash
sudo journalctl -u cronai -f
```

## Managing the Service

- **Restart the service**

```bash
sudo systemctl restart cronai
```

- **Stop the service**

```bash
sudo systemctl stop cronai
```

- **Disable the service** (prevents it from starting at boot)

```bash
sudo systemctl disable cronai
```