# orange-is-stupid

orange is stupid because for no reason they completely disallow having a static ip.

So I'm using cloudflare API and raspberry pi to update my subdomain to my ip automatically.

---

Building binary for raspberry pi on win machine _(if you use linux, you can figure this out :D)_.

```pwsh
.\build.bat
```

---

Remember to supply 2 env variables:

```bash
CF_API_TOKEN="Tv4..."
URLS="subdomain.domain.tld|true"
```

The `URLS` can be multiple separated by comma (`,`) and the pipe (`|`) that you see here, is determining whether the url should use cloudflare proxy.

In order to have a valid cloudflare api token (`CF_API_TOKEN`):

1. Log in to cloudflare and proceed to `https://dash.cloudflare.com/`
2. From the menu of left expand `Manage account`
3. Go to `Account API token`
4. Click `Create Token` button
5. Next to `Edit zone DNS` click `Use template` button
6. In the `Zone Resources` either select specific zone, or change the `Specific zone` to all.
7. Click `Continue to summary`
8. Click `Create Token`
9. The API token will be under the `Copy this token to access the Cloudflare API. For security this will not be shown again. Learn more` text.

## scheduler example

You can run the binary on schedule with the following setup, it will be preserved even after raspberry pi reboot:

> `/etc/systemd/system/cfrunner.service`

```yaml
[Unit]
Description=Cloudflare Runner (daily job)
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
WorkingDirectory=/home/pi/cloudflare-runner
User=tomek
Group=tomek
ExecStart=/home/pi/cloudflare-runner/cfrunner
# don't forget to put /home/pi/cloudflare-runner/.env file
# Log to the journal
StandardOutput=journal
StandardError=journal
# Be polite to the system a bit
Nice=5
# Reasonable file limit
LimitNOFILE=65536
# Hardening (avoid anything that would block writing in /home/tpiomek)
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full

[Install]
WantedBy=multi-user.target
```

> `/etc/systemd/system/cfrunner.timer`

```yaml
[Unit]
Description=Run Cloudflare Runner every hour

[Timer]
OnCalendar=hourly
# every 3:00 AM = OnCalendar=*-*-* 03:00:00
# If the machine was off at 03:00, run it shortly after boot:
Persistent=true
Unit=cfrunner.service

[Install]
WantedBy=timers.target
```

In order to run the above, you need to kickstart the timer once:

```bash
sudo systemctl daemon-reload # reloading the systemd service
sudo systemctl enable --now cfrunner.timer # enable the timer
```

Check the schedule:

```bash
systemctl list-timers cfrunner.timer
```

In order to manually start the service _(for testing and not having to wait for the timer)_:

```bash
sudo systemctl start cfrunner.service
```
