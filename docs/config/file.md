# Config File

The following section details options available via the `config.yml` file. Config snippets are provided for each section, however it is recommended to start with a blank copy of the full [sample config](./sample.md).

Config may also be provided via environment variables. Frigate-notify will load environment variables prefixed with `FN_`. Environment variables follow the same structure as the config file below, with heirarchy separated by an underscore (`_`). For example, setting the Frigate server address would be `FN_FRIGATE_SERVER`, or enabling Discord alerts would use `FN_ALERTS_DISCORD_ENABLED`.

## Frigate

### Server

- **server** (Required)
    - IP, hostname, or URL of the Frigate NVR
    - If IP or hostname specified, app will prepend `http://`
    - If Frigate is not behind a reverse proxy, append port number if necessary
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates for `server`
- **public_url** (Optional)
    - Should be set if Frigate is available via an external, public URL
    - This value is used for the links used in notifications
    - Format should be full URL (example: `https://nvr.your.public.domain.tld`)
- **headers** (Optional)
    - Send additional HTTP headers to Frigate
    - Useful for things like authentication
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **startup_check** (Optional)
    - On startup, frigate-notify will attempt to reach the configured Frigate NVR to validate connectivity
    - These options allow customization of the max attempts & retry interval
    - **attempts** (Optional - Default: `5`)
        - Max number of attempts to reach Frigate server
    - **interval** (Optional - Default: `30`)
        - Interval between retries, in seconds

```yaml title="Config File Snippet"
frigate:
  server: nvr.your.domain.tld
  ignoressl: true
  public_url: https://nvr.your.public.domain.tld
  headers:
    - Authorization: Basic abcd1234
  startup_check:
    attempts: 5
    interval: 30
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
- **topic_prefix** (Optional - Default: `frigate`)
    - Optionally change MQTT topic prefix
    - This should match the topic prefix used by Frigate

```yaml title="Config File Snippet"
frigate:
  mqtt: 
    enabled: true
    server: mqtt.your.domain.tld
    port: 1883
    clientid: frigate-notify
    username: mqtt-user
    password: mqtt-pass
    topic_prefix: frigate
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

## Alerts

!!! note
    Any combination of alerting methods may be enabled, though you'll probably want to enable at least one! 😅

### General

- **title** (Optional - Default: `Frigate Alert`)
    - Title of alert messages that are generated (Email subject, etc)
- **timeformat** (Optional - Default: `2006-01-02 15:04:05 -0700 MST`)
    - Optionally set a custom date/time format for notifications
    - This utilizes Golang's [reference time](https://go.dev/src/time/format.go) for formatting
    - See [this](https://www.geeksforgeeks.org/time-formatting-in-golang) guide for help
    - Example below uses RFC1123 format
- **nosnap** (Optional - Default: `allow`)
    - Specify what to do with events that have no snapshot image
    - By default, these events will be sent & notification message will say "No snapshot available"
    - Set to `drop` to silently drop these events & not send notifications

```yaml title="Config File Snippet"
alerts:  
  general:
    title: Frigate Alert
    timeformat: Mon, 02 Jan 2006 15:04:05 MST
    nosnap: 
```

### Zones

This config section allows control over whether to generate alerts on all zones, or only specific ones. By default, the app will generate notifications on **all** Frigate events, regardless of whether the event includes zone info.

??? note "A note about how this works"
    
    **With MQTT**, Frigate will send a `new` event when a detection starts. Subsequent changes, like the detected object transitioning from one zone to another, will trigger `update` events. These `update` events will contain a list of current zone(s) that the object is in, as well as a list of all zones that the object has entered during the event.

    In order to reduce the number of notifications generated, this app will only alert on the *first time* the detected object enters a zone.

    For example, let's say you have a camera in your front yard with zones for sidewalk, driveway, and lawn - but only allow notifications for driveway and lawn. During an event someone was detected originally on the sidewalk, then driveway, lawn, and back to driveway. In this case, you should only receive two notifications. Once for the first time the person entered the driveway zone, and a second when they entered the lawn zone. 

    **With Web API event query**, we only query the event from Frigate one time. So currently, only one alert would be sent depending on the detected zones at the time the web API was queried for new events.

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

### Labels

Similar to [zones](#zones), notifications can be filtered based on labels. By default, the app will generate notifications regardless of any labels received from Frigate. Using this config section, certain labels can be blocked from sending notifications - or an allowlist can be provided to only generate alerts from specified labels.

- **min_score** (Optional - Default: `0`)
    - Filter by minimum label score, based on Frigate `top_score` value
    - Scores are a percent accuracy of object identification (0-100)
    - For example, to filter objects under 80% accuracy, set `min_score: 80`
    - By default, any score above 0 will generate an alert
- **allow** (Optional)
    - Specify a list of labels to allow notifications
    - If set, all other labels will be ignored
    - If not set, all labels will generate notifications
- **block** (Optional)
    - Specify a list of labels to always ignore
    - This takes precedence over the `allow` list

```yaml title="Config File Snippet"
alerts:
  labels:
    min_score: 80
    allow:
     - person
     - dog
    block:
     - bird
```

### Subabels

Filter by sublabels, just like normal [labels](#labels).

- **allow** (Optional)
    - Specify a list of sublabels to allow notifications
    - If set, all other sublabels will be ignored
    - If not set, all sublabels will generate notifications
- **block** (Optional)
    - Specify a list of sublabels to always ignore
    - This takes precedence over the `allow` list

```yaml title="Config File Snippet"
alerts:
  sublabels:
    allow:
     - ABCD
     - EFGH
    block:
     - XYZ
```

### Discord

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Discord webhooks
- **webhook** (Required)
    - Full URL of the desired Discord webhook to send alerts through
    - Required if this alerting method is enabled
    - Check [Discord's](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks) docs for how to create a webhook
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:  
  discord:
    enabled: false
    webhook: https://<your-discord-webhook-here>
    template:
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
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:  
  gotify:
    enabled: false
    server: gotify.your.domain.tld
    token: ABCDEF
    ignoressl: true
    template:
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
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:  
  smtp:
    enabled: false
    server: smtp.your.domain.tld
    port: 587
    tls: true
    user: test_user@your.domain.tld
    password: test_pass
    recipient: nvr_group@your.domain.tld, someone_else@your.domain.tld
    template:
```

### Telegram

!!! note
    There is an [issue](https://github.com/0x2142/frigate-notify/issues/54#issuecomment-2148564526) with Telegram alerts if you use URL-embedded credentials for your Frigate links, for example: `https://user:pass@frigate.domain.tld`

    Telegram appears to incorrectly process these URLs, which will cause the camera & clip links  to become unclickable within Telegram.

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
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:  
  telegram:
    enabled: true
    chatid: 123456789
    token: 987654321:ABCDEFGHIJKLMNOP
    template:
```

### Pushover

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Pushover
- **token** (Required)
    - Pushover application API token
    - Required if this alerting method is enabled
- **userkey** (Required)
    - Recipient user or group key from Pushover dashboard
    - Required if this alerting method is enabled
- **devices** (Optional)
    - Optionally specify list of devices to send notifications to
    - If left empty, all devices will receive the notification
- **priority** (Optional)
    - Optionally set message priority
    - Valid priorities are -2, -1, 0, 1, 2
- **retry** (Optional)
    - Message retry in seconds until message is acknowledged
    - If `priority` is set to 2, this is required
    - Minimum value is 30 seconds
- **expire** (Optional)
    - Expiration timer for message retry
    - If `priority` is set to 2, this is required
- **ttl** (Optional)
    - Optionally set lifetime of message, in seconds
    - If set, message notifications are deleted from devices after this time
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
  pushover:
    enabled: true
    token: aaaaaaaaaaaaaaaaaaaaaa
    userkey: bbbbbbbbbbbbbbbbbbbbbb
    devices: device1,device2
    priority: 0
    retry:
    expire:
    ttl:
    template:
```

### Ntfy

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via Ntfy
- **server** (Required)
    - Full URL of the desired Ntfy server
    - Required if this alerting method is enabled
- **topic** (Required)
    - Destination topic that will receive alert notifications
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates
- **headers** (Optional)
    - Send additional HTTP headers to Ntfy server
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts: 
  ntfy:
    enabled: true
    server: https://ntfy.your.domain.tld
    topic: frigate
    ignoressl: true
    headers:
    template:
```

### Webhook

!!! note
    Webhook alerts are JSON only, and do not contain an image from the event.

```json title="Default webhook message"
{
    "time": "",
    "id": "",
    "camera": "",
    "label": "", 
    "score": "",
    "current_zones": "",
    "entered_zones": "",
    "has_clip": "",
    "has_snapshot": "",
    "links": {
         "camera": "",
         "clip": "",
         "snapshot": "",
    },
}
```

- **enabled** (Optional - Default: `false`)
    - Set to `true` to enable alerting via webhook
- **server** (Required)
    - Full URL of the desired webhook server
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Set to `true` to allow self-signed certificates
- **headers** (Optional)
    - Send additional HTTP headers to webhook receiver
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **template** (Optional)
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)
    - Note: Webhook templates **must** be valid JSON

```yaml title="Config File Snippet"
  webhook:
    enabled: false
    server: 
    ignoressl:
    headers:
    template:
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
