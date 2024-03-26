# Config

This section will walk through the available configuration items.

Configuration snippets will be provided throughout this page. Feel free to copy the full [Sample Config](#sample-config) at the bottom of this page to customize.

## Frigate

### Server

- **server** (Required)
    - IP or hostname of the Frigate NVR
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates

```yaml title="Config File Snippet"
frigate:
  server: nvr.your.domain.local
  ignoressl: true
```

### WebAPI

!!! note
    Only one monitoring method can be configured, either `webapi` or `mqtt`. The other must be set to `enabled: false`.

- **enabled** (Optional - Default: `false`)
    - If set to `true`, Frigate events are collected by polling the web API
- **interval** (Optional - Default: `30`)
    - How frequently to check the Frigate web API for new events, in seconds 

```yaml title="Config File Snippet"
frigate:
  webapi:
    enabled: true
    interval: 60
```

### MQTT

!!! note
    Only one monitoring method can be configured, either `webapi` or `mqtt`. The other must be set to `enabled: false`.

- **enabled** (Optional - Default: `false`)
    - If set to `true`, Frigate events are collected via an MQTT broker
    - Note: This must be the same MQTT broker that Frigate is sending events to
- **server** (Required)
    - IP or hostname of the MQTT server
    - If MQTT monitoring is enabled, this field is required
- **port** (Optiona - Default: `1883`)
    - MQTT service port
- **clientid** (Optional - Default: `frigate-notify`)
    - Client ID of this app used when connecting to MQTT
    - Note: This must be a unique value & cannot be shared with other MQTT clients
- **username** (Optional)
    - MQTT username
    - If username & password are not set, then authentication is disabled
- **password** (Optional)
    - MQTT password
    - Required if `username` is set

```yaml title="Config File Snippet"
frigate:
  mqtt: 
    enabled: true
    server: mqtt.your.domain.local
    port: 1883
    clientid: frigate-notify
    username: mqtt-user
    password: mqtt-pass
```

### Cameras

- **exclude** (Optional)
    - If desired, provide a list of cameras to ignore
    - Any Frigate events on these cameras will not generate alerts
    - If left empty, this is disabled & all cameras will generate alerts

```yaml title="Config File Snippet"
frigate:
  cameras:
    exclude:
      - test_cam_01
      - test_cam_02
```

## Alerting

!!! note
    Any combination of alerting methods may be enabled, though you'll probably want to enable at least one! 😅

### General

- **title** (Optional - Default: `Frigate Alert`)
    - Title of alert messages that are generated (Email subject, etc)

```yaml title="Config File Snippet"
alerts:  
  general:
    title: Frigate Alert
```

### Zones

This config section allows control over whether to generate alerts on all zones, or only specific ones. By default, the app will generate notifications on **all** Frigate events, regardless of whether the event includes zone info.

- **unzoned** (Optional - Default: `allow`)
    - Controls alerts on events outside a zone
    - By default, events without a zone will generate alerts
    - Set to `drop` to prevent generating alerts from events without a zone
- **allow** (Optional)
    - Specify a list of zones to allow notifications
    - All other zones will be ignored
    - If `unzoned` is set to `allow`, notifications will still be sent on events without any zone info
- **block** (Optional)
    - Specify a list of zones to always ignore
    - This takes precedence over the `allow` list

```yaml title="Config File Snippet"
alerts:
  zones:
    unzoned: allow
    allow:
     - test_zone_01
    block:
     - test_zone_02
```

### Discord

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Discord webhooks
- **webhook** (Required)
    - Full URL of the desired Discord webhook to send alerts through
    - Required if this alerting method is enabled
    - Check [Discord's](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks) docs for how to create a webhook

```yaml title="Config File Snippet"
alerts:  
  discord:
    enabled: false
    webhook: https://<your-discord-webhook-here>
```

### Gotify

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Gotify
- **server** (Required)
    - IP or hostname of the target Gotify server
    - Required if this alerting method is enabled
- **token** (Required)
    - App token associated with this app in Gotify
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates

```yaml title="Config File Snippet"
alerts:  
  gotify:
    enabled: false
    server: gotify.your.domain.local
    token: ABCDEF
    ignoressl: true
```

### SMTP

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via SMTP
- **server** (Required)
    - IP or hostname of the target SMTP server
    - Required if this alerting method is enabled
- **port** (Required)
    - Port of the target SMTP server
    - Required if this alerting method is enabled
- **tls** (Optional - Default: `false`)
    - Set to `true` if SMTP TLS is required
- **user** (Optional)
    - Add SMTP username for authentication
    - If username & password are not set, then authentication is disabled
- **password** (Optional)
    - Password of SMTP user
    - Required if `user` is set
- **recipient** (Required)
    - Comma-separated list of email recipients
    - Required if this alerting method is enabled

```yaml title="Config File Snippet"
alerts:  
  smtp:
    enabled: false
    server: smtp.your.domain.local
    port: 587
    tls: true
    user: test_user@your.domain.local
    password: test_pass
    recipient: nvr_group@your.domain.local, someone_else@your.domain.local
```

### Telegram

In order to use Telegram for alerts, a bot token & chat ID are required.

To obtain a bot token, follow [this](https://core.telegram.org/bots#how-do-i-create-a-bot) doc to message @BotFather.

Once you have a bot token, make sure to initiate a chat message with your bot. Then visit the following URL:

- `https://api.telegram.org/bot<BOT_TOKEN>/getUpdates` 
- Replace `<BOT_TOKEN>` with the API token provided by @BotFather. 

Within the response, locate your message to the bot, then grab the ID under `message > chat > id`. An example response is shown below, where `999999999` is the ID we need to save:

```json
{
  "update_id": 1234567,
  "message": {
    "chat": {
      "id": 999999999,
      "first_name": "Test User",
      "username": "test-username1234",
      "type": "private"
    }
  }
}
```

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Telegram
- **chatid** (Required)
    - Chat ID for the alert recipient
    - Required if this alerting method is enabled
- **token** (Required)
    - Bot token generated from [@BotFather](https://core.telegram.org/bots#how-do-i-create-a-bot)
    - Required if this alerting method is enabled

```yaml title="Config File Snippet"
alerts:  
  telegram:
    enabled: true
    chatid: 123456789
    token: 987654321:ABCDEFGHIJKLMNOP
```


## Monitor

If enabled, this application will check in with tools like [HealthChecks](https://github.com/healthchecks/healthchecks) or [Uptime Kuma](https://github.com/louislam/uptime-kuma) on a regular interval for health / status monitoring.

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable health checks
- **url** (Required)
    - URL path for health check service
    - Required if monitoring is enabled
- **interval** (Required - Default: `60`)
    - Frequency of check-in events
    - Required if monitoring is enabled
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates

```yaml title="Config File Snippet"
monitor:
  enabled: false
  url: 
  interval: 
  ignoressl: 
```


---


## Sample Config { data-search-exclude }

A full config file template has been provided below:

```yaml
frigate:
  server: 
  ignoressl: 

  webapi:
    enabled: 
    interval: 
    
  mqtt: 
    enabled: 
    server: 
    port:
    clientid:
    username: 
    password: 
  
  cameras:
    exclude:
      - test_cam_01


alerts:  
  general:
    title:

  zones:
    unzoned: allow
    allow:
     - test_zone_01
    block:
     - test_zone_02

  discord:
    enabled: false
    webhook: 
  
  gotify:
    enabled: false
    server:
    token: 
    ignoressl:
  
  smtp:
    enabled: false
    server: 
    port: 
    tls: 
    user:
    password:
    recipient:


monitor:
  enabled: false
  url: 
  interval: 
  ignoressl: 
```