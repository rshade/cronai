# Running CronAI as a systemd Service

This document explains how to set up CronAI to run as a systemd service on Linux systems.

## Setup Steps

1. **Build and install the CronAI binary**

```bash
cd /path/to/cronai
go build -o cronai ./cmd/cronai
sudo cp cronai /usr/local/bin/cronai
```text

2. **Create your configuration and prompt files**

```bash
mkdir -p /etc/cronai/cron_prompts
cp cronai.config.example /etc/cronai/cronai.config
cp -r cron_prompts/* /etc/cronai/cron_prompts/
```text

3. **Set up your environment file**

```bash
cp .env.example /etc/cronai/.env
# Edit the .env file with your API keys and settings
sudo nano /etc/cronai/.env
```text

4. **Create the systemd service file**

Copy the example service file and modify it for your system:

```bash
sudo cp cronai.service /etc/systemd/system/cronai.service
sudo nano /etc/systemd/system/cronai.service
```text

Update the following fields in the service file:

- `User`: The user account that will run the service
- `WorkingDirectory`: The directory where your configuration is located (e.g., `/etc/cronai`)
- `ExecStart`: The path to the CronAI binary (e.g., `/usr/local/bin/cronai start --config /etc/cronai/cronai.config`)
  - Since v0.0.2: You can also specify the operation mode with `--mode cron` (default)
  - Future modes (bot, queue) will be available in upcoming releases
- `EnvironmentFile`: The path to your .env file (e.g., `/etc/cronai/.env`)

5. **Enable and start the service**

```bash
sudo systemctl daemon-reload
sudo systemctl enable cronai
sudo systemctl start cronai
```text

6. **Check the service status**

```bash
sudo systemctl status cronai
```text

7. **View the logs**

```bash
sudo journalctl -u cronai -f
```text

## Managing the Service

- **Restart the service**

```bash
sudo systemctl restart cronai
```text

- **Stop the service**

```bash
sudo systemctl stop cronai
```text

- **Disable the service** (prevents it from starting at boot)

```bash
sudo systemctl disable cronai
```text
