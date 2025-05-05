# Config File

The following section details options available via the `config.yml` file. Config snippets are provided for each section, however it is recommended to start with a blank copy of the full [sample config](./sample.md).

Config may also be provided via environment variables. Frigate-notify will load environment variables prefixed with `FN_`. Environment variables follow the same structure as the config file below, with heirarchy separated by two underscores (`__`). For example, setting the Frigate server address would be `FN_FRIGATE__SERVER`, or enabling Discord alerts would use `FN_ALERTS__DISCORD__ENABLED`. If multiple alert profiles are configured, they are specified by index starting at zero - for example: `FN_ALERTS__DISCORD__0__ENABLED`, `FN_ALERTS__DISCORD__1__ENABLED`, etc. Environment variable keys are shown alongside config file keys below.

Config is loaded & merged in the following order: Config file > Environment Variables > Docker secrets . This allows most configuration to be provided via config file, but provide secrets via environment variables or docker secrets.

## App

- **mode** (Optional - Default: `reviews`)
    - Env: `FN_APP__MODE`
    - Specify whether to notifications are based on Frigate [Events or Reviews](https://docs.frigate.video/configuration/review/#review-items-vs-events)
    - `events` will notify on every detected object from Frigate
    - `reviews` will only notify on Frigate **Alerts** (Requires Frigate 0.14+)
        - When in `reviews` mode, toggle `notify_detections` under the `alerts` [config section](#alerts) to also notify on **Detections**
        - See also [Alerts vs Detections](https://docs.frigate.video/configuration/review/#alerts-and-detections)
- **API**
    - **enabled** (Optional - Default: `false`)
        - Env: `FN_APP__API__ENABLED`
        - Set to `true` to enable Frigate-Notify's REST API server
    - **port** (Optional - Default: `8000`)
        - Env: `FN_APP__API__PORT`
        - Change default port for API server

```yaml title="Config File Snippet"
app:
  mode: events
  api:
    enabled: true
    port: 8000
```

## Frigate

### Server

- **server** (Required)
    - Env: `FN_FRIGATE__SERVER`
    - IP, hostname, or URL of the Frigate NVR
    - If IP or hostname specified, app will prepend `http://`
    - If Frigate is not behind a reverse proxy, append port number if necessary
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_FRIGATE__IGNORESSL`
    - Set to `true` to allow self-signed certificates for `server`
- **public_url** (Optional)
    - Env: `FN_FRIGATE__PUBLIC_URL`
    - Should be set if Frigate is available via an external, public URL
    - This value is used for the links used in notifications
    - Format should be full URL (example: `https://nvr.your.public.domain.tld`)
- **username** (Optional)
    - Frigate username to log in with, if using authenticated UI on port 8971
    - If username is configured, password must also be configured
    - Recommended to create a unique user for frigate-notify with **viewer** role
- **password** (Optional)
    - Frigate password to log in with, if using authenticated UI on port 8971
    - If password is configured, username must also be configured
- **headers** (Optional)
    - Env: `FN_FRIGATE__HEADERS`
    - Send additional HTTP headers to Frigate
    - Useful for things like authentication
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **startup_check** (Optional)
    - On startup, frigate-notify will attempt to reach the configured Frigate NVR to validate connectivity
    - These options allow customization of the max attempts & retry interval
    - **attempts** (Optional - Default: `5`)
        - Env: `FN_FRIGATE__STARTUP_CHECK__ATTEMPTS`
        - Max number of attempts to reach Frigate server
    - **interval** (Optional - Default: `30`)
        - Env: `FN_FRIGATE__STARTUP_CHECK__INTERVAL`
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
    - Env: `FN_FRIGATE__WEBAPI__ENABLED`
    - If set to `true`, Frigate events are collected by polling the web API
- **interval** (Optional - Default: `30`)
    - Env: `FN_FRIGATE__WEBAPI__INTERVAL`
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
    - Env: `FN_FRIGATE__MQTT__ENABLED`
    - If set to `true`, Frigate events are collected via an MQTT broker
    - Note: This must be the same MQTT broker that Frigate is sending events to
- **server** (Required)
    - Env: `FN_FRIGATE__MQTT__SERVER`
    - IP or hostname of the MQTT server
    - If MQTT monitoring is enabled, this field is required
- **port** (Optiona - Default: `1883`)
    - Env: `FN_FRIGATE__MQTT__PORT`
    - MQTT service port
- **clientid** (Optional - Default: `frigate-notify`)
    - Env: `FN_FRIGATE__MQTT__CLIENTID`
    - Client ID of this app used when connecting to MQTT
    - Note: This must be a unique value & cannot be shared with other MQTT clients
- **username** (Optional)
    - Env: `FN_FRIGATE__MQTT__USERNAME`
    - MQTT username
    - If username & password are not set, then authentication is disabled
- **password** (Optional)
    - Env: `FN_FRIGATE__MQTT__PASSWORD`
    - MQTT password
    - Required if `username` is set
- **topic_prefix** (Optional - Default: `frigate`)
    - Env: `FN_FRIGATE__MQTT__TOPIC_PREFIX`
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
    - Env: `FN_FRIGATE__CAMERAS__EXCLUDE`
    - If desired, provide a list of cameras to ignore
    - Any Frigate events on these cameras will not generate alerts
    - If left empty, this is disabled & all cameras will generate alerts
    - If configuring via environment variable, separate camera names by semicolon

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

All alert providers (Discord, Gotify, etc) also support optional filters & the ability to configure multiple profiles per provider. Please see [Alert Profiles & Filters](https://frigate-notify.0x2142.com/latest/config/profilesandfilters/) for more information.

### General

- **title** (Optional - Default: `Frigate Alert`)
    - Env: `FN_ALERTS__GENERAL__TITLE`
    - Title of alert messages that are generated (Email subject, etc)
    - Title value can utilize [template variables](./templates.md#available-variables)
- **timeformat** (Optional - Default: `2006-01-02 15:04:05 -0700 MST`)
    - Env: `FN_ALERTS__GENERAL__TIMEFORMAT`
    - Optionally set a custom date/time format for notifications
    - This utilizes Golang's [reference time](https://go.dev/src/time/format.go) for formatting
    - See [this](https://www.geeksforgeeks.org/time-formatting-in-golang) guide for help
    - Example below uses RFC1123 format
- **nosnap** (Optional - Default: `allow`)
    - Env: `FN_ALERTS__GENERAL__NOSNAP`
    - Specify what to do with events that have no snapshot image
    - By default, these events will be sent & notification message will say "No snapshot available"
    - Set to `drop` to silently drop these events & not send notifications
- **snap_bbox** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__SNAP_BBOX`
    - Includes object bounding box on snapshot when retrieved from Frigate
    - Note: Per [Frigate docs](https://docs.frigate.video/integrations/api/#get-apieventsidsnapshotjpg), only applied when event is in progress
- **snap_timestamp** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__SNAP_TIMESTAMP`
    - Includes timestamp on snapshot when retrieved from Frigate
    - Note: Per [Frigate docs](https://docs.frigate.video/integrations/api/#get-apieventsidsnapshotjpg), only applied when event is in progress
- **snap_crop** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__SNAP_CROP`
    - Crops snapshot when retrieved from Frigate
    - Note: Per [Frigate docs](https://docs.frigate.video/integrations/api/#get-apieventsidsnapshotjpg), only applied when event is in progress
- **snap_hires** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__SNAP_HIRES`
    - By default, snapshots are collected from Frigate detect stream which may be lower resolution
    - Set this to `true` to collect snapshot from camera recording stream
    - **Note**: If enabled, the above settings for `snap_bbox`, `snap_timestamp`, and `snap_crop` settings have no effect
    - **Note**: Snapshot generated via this method is based on the event start time provided by Frigate
        - This means that the snapshot collected may differ from the snapshot Frigate choses to use for the event
        - This may also mean that, depending on the timing of the detection, this snapshot *may* not include the detected object
        - If it is a priority to ensure that snapshots always include the detected object, then it is recommended to leave this option disabled
- **max_snap_retry** (Optional - Default: `10`)
    - Env: `FN_ALERTS__GENERAL__MAX_SNAP_RETRY`
    - Max number of retry attempts when waiting for snapshot to become available
    - Retries are every 2 seconds
    - Default is 10, which means waiting up to 20 seconds for snapshot
    - Note: Does not apply if event received from Frigate contains `has_snapshot: false`
- **notify_once** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__NOTIFY_ONCE`
    - By default, each Frigate event may generate several notifications as the object changes zones, etc
    - Set this to `true` to only notify once per event
- **notify_detections** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GENERAL__NOTIFY_DETECTIONS`
    - Only used when app `mode` is `reviews`
    - By default, notifications will only be sent on Frigate alerts
    - Set to `true` to also enable on detections
- **recheck_delay** (Optional - Default: `0`)
    - Env: `FN_ALERTS__GENERAL__RECHECK_DELAY`
    - Optionally re-check event details from Frigate before sending notifications
    - Delay period in seconds
    - If set to `0`, events are sent immediately upon receipt from Frigate
    - This setting can be useful if needing to wait for a 3rd-party app to set sub_labels
- **audio_only** (Optional - Default: `allow`)
    - Env: `FN_ALERTS__GENERAL__AUDIO_ONLY`
    - Specify what to do with events that only contain audio detection
    - By default, these events will generate notifications
    - Set to `drop` to silently drop these events & not send notifications

```yaml title="Config File Snippet"
alerts:
  general:
    title: Frigate Alert
    timeformat: Mon, 02 Jan 2006 15:04:05 MST
    nosnap:
    snap_bbox:
    snap_timestamp:
    snap_crop:
    snap_hires:
    max_snap_retry:
    notify_once:
    notify_detections:
    audio_only:
```

### Quiet Hours

Define a quiet period & supress alerts during this time.

- **start** (Optional)
    - Env: `FN_ALERTS__QUIET__START`
    - When quiet period begins, in 24-hour format
    - Required if `end` is configured
- **end** (Optional)
    - Env: `FN_ALERTS__QUIET__END`
    - When quiet period ends, in 24-hour format
    - Required if `start` is configured

```yaml title="Config File Snippet"
alerts:
  quiet:
    start: 08:00
    end: 17:00
```

### Zones

This config section allows control over whether to generate alerts on all zones, or only specific ones. By default, the app will generate notifications on **all** Frigate events, regardless of whether the event includes zone info.

??? note "A note about how this works"
    
    **With MQTT**, Frigate will send a `new` event when a detection starts. Subsequent changes, like the detected object transitioning from one zone to another, will trigger `update` events. These `update` events will contain a list of current zone(s) that the object is in, as well as a list of all zones that the object has entered during the event.

    In order to reduce the number of notifications generated, this app will only alert on the *first time* the detected object enters a zone.

    For example, let's say you have a camera in your front yard with zones for sidewalk, driveway, and lawn - but only allow notifications for driveway and lawn. During an event someone was detected originally on the sidewalk, then driveway, lawn, and back to driveway. In this case, you should only receive two notifications. Once for the first time the person entered the driveway zone, and a second when they entered the lawn zone. 

    **With Web API event query**, we only query the event from Frigate one time. So currently, only one alert would be sent depending on the detected zones at the time the web API was queried for new events.

- **unzoned** (Optional - Default: `allow`)
    - Env: `FN_ALERTS__ZONES__UNZONED`
    - Controls alerts on events outside a zone
    - By default, events without a zone will generate alerts
    - Set to `drop` to prevent generating alerts from events without a zone
- **allow** (Optional)
    - Env: `FN_ALERTS__ZONES__ALLOW`
    - Specify a list of zones to allow notifications
    - All other zones will be ignored
    - If `unzoned` is set to `allow`, notifications will still be sent on events without any zone info
    - If configuring via environment variable, separate zone names by semicolon
- **block** (Optional)
    - Env: `FN_ALERTS__ZONES__BLOCK`
    - Specify a list of zones to always ignore
    - This takes precedence over the `allow` list
    - If configuring via environment variable, separate zone names by semicolon

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
    - Env: `FN_ALERTS__LABELS__MIN_SCORE`
    - Filter by minimum label score, based on Frigate `top_score` value
    - Scores are a percent accuracy of object identification (0-100)
    - For example, to filter objects under 80% accuracy, set `min_score: 80`
    - By default, any score above 0 will generate an alert
- **allow** (Optional)
    - Env: `FN_ALERTS__LABELS__ALLOW`
    - Specify a list of labels to allow notifications
    - If set, all other labels will be ignored
    - If not set, all labels will generate notifications
    - If configuring via environment variable, separate label names by semicolon
- **block** (Optional)
    - Env: `FN_ALERTS__LABELS__BLOCK`
    - Specify a list of labels to always ignore
    - This takes precedence over the `allow` list
    - If configuring via environment variable, separate label names by semicolon

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

### Sublabels

Filter by sublabels, just like normal [labels](#labels).

- **allow** (Optional)
    - Env: `FN_ALERTS__SUBLABELS__ALLOW`
    - Specify a list of sublabels to allow notifications
    - If set, all other sublabels will be ignored
    - If not set, all sublabels will generate notifications
    - If configuring via environment variable, separate sublabel names by semicolon
- **block** (Optional)
    - Env: `FN_ALERTS__SUBLABELS__BLOCK`
    - Specify a list of sublabels to always ignore
    - This takes precedence over the `allow` list
    - If configuring via environment variable, separate sublabel names by semicolon

```yaml title="Config File Snippet"
alerts:
  sublabels:
    allow:
     - ABCD
     - EFGH
    block:
     - XYZ
```

### License Plate

Include license plate recognition data in notifications, if enabled in Frigate.

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__LICENSE_PLATE__ENABLED`
    - Specify whether to wait for license plate recognition data when Frigate detects a car & license plate
    - This will re-check the Frigate for license plate information every 2 seconds with a 10 second maximum
- **allow** (Optional)
    - Env: `FN_ALERTS__LICENSE_PLATE__ALLOW`
    - Specify a list of license plates to allow notifications
    - If set, all other license plates will be ignored
    - If not set, all license plates will generate notifications
    - If configuring via environment variable, separate license plates by semicolon
- **block** (Optional)
    - Env: `FN_ALERTS__LICENSE_PLATE__BLOCK`
    - Specify a list of license plates to always ignore
    - This takes precedence over the `allow` list
    - If configuring via environment variable, separate license plates by semicolon

```yaml title="Config File Snippet"
alerts:
  license_plate:
    enabled: true
    allow:
     - ABCD
     - EFGH
    block:
     - XYZ
```

### Apprise API

!!!important
    Notifications via Apprise require an external service: [https://github.com/caronc/apprise-api](https://github.com/caronc/apprise-api)

    Please follow the instructions on the [apprise api](https://github.com/caronc/apprise-api) repo for set up & configuration. This service exposes a REST API that Frigate-Notify uses to forward notifications to various notification providers.

    ⚠️ Not all Apprise notification providers support sending attachments. See the full list [here](https://github.com/caronc/apprise/wiki).

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__APPRISE_API__ENABLED`
    - Set to `true` to enable alerting via Apprise API
- **server** (Required)
    - Env: `FN_ALERTS__APPRISE_API__SERVER`
    - Full URL of the desired [apprise api](https://github.com/caronc/apprise-api) container
    - Required if this alerting method is enabled
- **token** (Required - Unless `urls` is used)
    - Env: `FN_ALERTS__APPRISE_API__TOKEN`
    - Config token in apprise api
    - Required if this alerting method is enabled
- **urls** (Required - Unless `token` is used)
    - Env: `FN_ALERTS__APPRISE_API__URLS`
    - Destination Apprise URLs to forward notifications to
    - See supported providers & example URL formats [here](https://github.com/caronc/apprise?tab=readme-ov-file#supported-notifications)
    - If configuring via environment variable, separate URLs by semicolon
- **tags** (Optional - Required if `token` is used)
    - Env: `FN_ALERTS__APPRISE_API__TAGS`
    - If using a config token, specify target tags to notify
    - If configuring via environment variable, separate tags by semicolon
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__APPRISE_API__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **template** (Optional)
    - Env: `FN_ALERTS__APPRISE_API__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:
  apprise_api:
    enabled: false
    server:
    token:
    urls:
      - ntfy://xxxxxxxx/frigate
      - discord://xxxxxxxxxxx
    tags:
      - ntfy
    ignoressl: true
    template:
```

### Discord

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__DISCORD__ENABLED`
    - Set to `true` to enable alerting via Discord webhooks
- **webhook** (Required)
    - Env: `FN_ALERTS__DISCORD__WEBHOOK`
    - Full URL of the desired Discord webhook to send alerts through
    - Required if this alerting method is enabled
    - Check [Discord's](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks) docs for how to create a webhook
- **disable_embed** (Optional)
    - Env: `FN_ALERTS__DISCORD__DISABLE_EMBED`
    - By default, notifications are sent as Discord embedded message
    - Set to `true` to disable this
- **template** (Optional)
    - Env: `FN_ALERTS__DISCORD__TEMPLATE`
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
    - Env: `FN_ALERTS__GOTIFY__ENABLED`
    - Set to `true` to enable alerting via Gotify
- **server** (Required)
    - Env: `FN_ALERTS__GOTIFY__SERVER`
    - IP or hostname of the target Gotify server
    - Required if this alerting method is enabled
- **token** (Required)
    - Env: `FN_ALERTS__GOTIFY__TOKEN`
    - App token associated with this app in Gotify
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__GOTIFY__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **template** (Optional)
    - Env: `FN_ALERTS__GOTIFY__TEMPLATE`
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

### Matrix

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__MATRIX__ENABLED`
    - Set to `true` to enable alerting via Matrix webhooks
- **server** (Required)
    - Env: `FN_ALERTS__MATRIX__SERVER`
    - Full URL of the desired Matrix homeserver
    - Required if this alerting method is enabled
- **username** (Required)
    - Env: `FN_ALERTS__MATRIX__USERNAME`
    - Username of Matrix user
- **password** (Required)
    - Env: `FN_ALERTS__MATRIX__PASSWORD`
    - Password for Matrix user
- **roomid** (Required)
    - Env: `FN_ALERTS__MATRIX__ROOMID`
    - Target Room ID to send notifications
    - Format: `"!<roomid>:<matrixhomeserver>"`
        - Note: This **must** be wrapped in quotes
    - Notification user must be invited to this room
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__MATRIX__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **template** (Optional)
    - Env: `FN_ALERTS__MATRIX__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
  matrix:
    enabled: false
    server: https://matrix.your.domain.tld
    username: someuser
    password: somepass
    roomid: "!abcd1234:matrix.your.domain.tld"
    ignoressl: true
    template:
```

### Mattermost

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__MATTERMOST__ENABLED`
    - Set to `true` to enable alerting via Mattermost webhooks
- **webhook** (Required)
    - Env: `FN_ALERTS__MATTERMOST__WEBHOOK`
    - Full URL of the desired Mattermost webhook to send alerts through
    - Required if this alerting method is enabled
    - Check [Mattermost's](https://developers.mattermost.com/integrate/webhooks/incoming/) docs for how to create a webhook
- **channel** (Optional)
    - Env: `FN_ALERTS__MATTERMOST__CHANNEL`
    - Override destination channel to post messages, if allowed by Mattermost config
- **username** (Optional)
    - Env: `FN_ALERTS__MATTERMOST__USERNAME`
    - Override username to post messages as, if allowed by Mattermost config
- **priority** (Optional - Default: `standard`)
    - Env: `FN_ALERTS__MATTERMOST__PRIORITY`
    - Set message priority
    - Options: `standard`, `important`, `urgent`
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__MATTERMOST__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **headers** (Optional)
    - Env: `FN_ALERTS__MATTERMOST__HEADERS`
    - Send additional HTTP headers with Mattermost webhook
    - Header values can utilize [template variables](./templates.md#available-variables)
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **template** (Optional)
    - Env: `FN_ALERTS__MATTERMOST__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
  mattermost:
    enabled: false
    webhook: https://mattermost.your.domain.tld
    channel: frigate-notifications
    username: frigate-notify
    priority: standard
    ignoressl: true
    headers:
    template:
```

### Ntfy

!!!note
    If you're self-hosting Ntfy, you'll need to ensure support for [attachments](https://docs.ntfy.sh/config/#attachments) is enabled.

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__NTFY__ENABLED`
    - Set to `true` to enable alerting via Ntfy
- **server** (Required)
    - Env: `FN_ALERTS__NTFY__SERVER`
    - Full URL of the desired Ntfy server
    - Required if this alerting method is enabled
- **topic** (Required)
    - Env: `FN_ALERTS__NTFY__TOPIC`
    - Destination topic that will receive alert notifications
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__NTFY__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **headers** (Optional)
    - Env: `FN_ALERTS__NTFY__HEADERS`
    - Send additional HTTP headers to Ntfy server
    - Header values can utilize [template variables](./templates.md#available-variables)
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
    - **Note:** Notifications via Ntfy are sent with a default action button that links to the event clip. This can be overridden by defining a custom `X-Action` header here
- **template** (Optional)
    - Env: `FN_ALERTS__NTFY__TEMPLATE`
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

### Pushover

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__PUSHOVER__ENABLED`
    - Set to `true` to enable alerting via Pushover
- **token** (Required)
    - Env: `FN_ALERTS__PUSHOVER__TOKEN`
    - Pushover application API token
    - Required if this alerting method is enabled
- **userkey** (Required)
    - Env: `FN_ALERTS__PUSHOVER__USERKEY`
    - Recipient user or group key from Pushover dashboard
    - Required if this alerting method is enabled
- **devices** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__DEVICES`
    - Optionally specify list of devices to send notifications to
    - If left empty, all devices will receive the notification
- **sound** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__SOUND`
    - Specify custom sound for notifications from this app
    - For available values, see the [Pushover Docs](https://pushover.net/api#sounds)
- **priority** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__PRIORITY`
    - Optionally set message priority
    - Valid priorities are -2, -1, 0, 1, 2
- **retry** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__RETRY`
    - Message retry in seconds until message is acknowledged
    - If `priority` is set to 2, this is required
    - Minimum value is 30 seconds
- **expire** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__EXPIRE`
    - Expiration timer for message retry
    - If `priority` is set to 2, this is required
- **ttl** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__TTL`
    - Optionally set lifetime of message, in seconds
    - If set, message notifications are deleted from devices after this time
- **template** (Optional)
    - Env: `FN_ALERTS__PUSHOVER__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
  pushover:
    enabled: true
    token: aaaaaaaaaaaaaaaaaaaaaa
    userkey: bbbbbbbbbbbbbbbbbbbbbb
    devices: device1,device2
    sound:
    priority: 0
    retry:
    expire:
    ttl:
    template:
```

### Signal

!!!important
    Signal notifications rely on an external service to handle communication to Signal: [https://github.com/bbernhard/signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api)

    Please follow the instructions on the [signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api) repo for set up & configuration. This service exposes a REST API that Frigate-Notify uses to forward notifications to Signal.

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__SIGNAL__ENABLED`
    - Set to `true` to enable alerting via Ntfy
- **server** (Required)
    - Env: `FN_ALERTS__SIGNAL__SERVER`
    - Full URL of the desired [signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api) container
    - Required if this alerting method is enabled
- **account** (Required)
    - Env: `FN_ALERTS__SIGNAL__ACCOUNT`
    - Signal account used to send notifications
    - This is the full phone number including country code (ex. `+12223334444`)
    - Required if this alerting method is enabled
- **recipients** (Required)
    - Env: `FN_ALERTS__SIGNAL__RECIPIENTS`
    - One or more Signal recipients that will receive notifications
    - This is the full phone number including country code (ex. `+12223334444`)
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__SIGNAL__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **template** (Optional)
    - Env: `FN_ALERTS__SIGNAL__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:
  signal:
    enabled: false
    server: https://signal-cli-rest-api.your.domain.tld
    account: +12223334444
    recipients:
     - +15556667777
    ignoressl: true
    template:
```

### SMTP

- **enabled** (Optional - Default: `false`)
    - Env: `FN_ALERTS__SMTP__ENABLED`
    - Set to `true` to enable alerting via SMTP
- **server** (Required)
    - Env: `FN_ALERTS__SMTP__SERVER`
    - IP or hostname of the target SMTP server
    - Required if this alerting method is enabled
- **port** (Required)
    - Env: `FN_ALERTS__SMTP__PORT`
    - Port of the target SMTP server
    - Required if this alerting method is enabled
- **tls** (Optional - Default: `false`)
    - Env: `FN_ALERTS__SMTP__TLS`
    - Set to `true` if SMTP TLS is required
- **user** (Optional)
    - Env: `FN_ALERTS__SMTP__USER`
    - Add SMTP username for authentication
    - If username & password are not set, then authentication is disabled
- **password** (Optional)
    - Env: `FN_ALERTS__SMTP__PASSWORD`
    - Password of SMTP user
    - Required if `user` is set
- **from** (Optional)
    - Env: `FN_ALERTS__SMTP__FROM`
    - Set sender address for outgoing messages
    - If left blank but authentication is configured, then `user` will be used
- **recipient** (Required)
    - Env: `FN_ALERTS__SMTP__RECIPIENT`
    - Comma-separated list of email recipients
    - Required if this alerting method is enabled
- **template** (Optional)
    - Env: `FN_ALERTS__SMTP__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__SMTP__IGNORESSL`
    - Set to `true` to allow self-signed certificates

```yaml title="Config File Snippet"
alerts:  
  smtp:
    enabled: false
    server: smtp.your.domain.tld
    port: 587
    tls: true
    from: test_user@your.domain.tld
    user: test_user@your.domain.tld
    password: test_pass
    recipient: nvr_group@your.domain.tld, someone_else@your.domain.tld
    template:
    ignoressl:
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
    - Env: `FN_ALERTS__TELEGRAM__ENABLED`
    - Set to `true` to enable alerting via Telegram
- **chatid** (Required)
    - Env: `FN_ALERTS__TELEGRAM__CHATID`
    - Chat ID for the alert recipient
    - Required if this alerting method is enabled
- **token** (Required)
    - Env: `FN_ALERTS__TELEGRAM__TOKEN`
    - Bot token generated from [@BotFather](https://core.telegram.org/bots#how-do-i-create-a-bot)
    - Required if this alerting method is enabled
- **send_clip** (Optional - Default: `false`)
    - Env: `FN_ALERTS__TELEGRAM__SEND_CLIP`
    - Optionally send event video clip instead of snapshot image
    - **Note**: Clips may take a short while to become available. Use [`max_snap_retry`](#general) to control how long frigate-notify will wait for the clip to become available
- **message_thread_id** (Optional)
    - Env: `FN_ALERTS__TELEGRAM__MESSAGE_THREAD_ID`
    - Optionally send notification to a message thread by ID
- **template** (Optional)
    - Env: `FN_ALERTS__TELEGRAM__TEMPLATE`
    - Optionally specify a custom notification template
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)

```yaml title="Config File Snippet"
alerts:  
  telegram:
    enabled: true
    chatid: 123456789
    message_thread_id: 100
    send_clip:
    token: 987654321:ABCDEFGHIJKLMNOP
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
    - Env: `FN_ALERTS__WEBHOOK__ENABLED`
    - Set to `true` to enable alerting via webhook
- **server** (Required)
    - Env: `FN_ALERTS__WEBHOOK__SERVER`
    - Full URL of the desired webhook server
    - Required if this alerting method is enabled
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_ALERTS__WEBHOOK__IGNORESSL`
    - Set to `true` to allow self-signed certificates
- **method** (Optional - Default: `POST`)
    - Env: `FN_ALERTS__WEBHOOK__METHOD`
    - Set HTTP method for webhook notifications
    - Supports `GET` and `POST`
- **params** (Optional)
    - Env: `FN_ALERTS__WEBHOOK__PARAMS`
    - Set optional HTTP params that will be appended to URL
    - Params can utilize [template variables](./templates.md#available-variables)
    - Format: `param: value`
    - Example: `token: abcd1234`
- **headers** (Optional)
    - Env: `FN_ALERTS__WEBHOOK__HEADERS`
    - Send additional HTTP headers to webhook receiver
    - Header values can utilize [template variables](./templates.md#available-variables)
    - Header format: `Header: Value`
    - Example: `Authorization: Basic abcd1234`
- **template** (Optional)
    - Env: `FN_ALERTS__WEBHOOK__TEMPLATE`
    - Optionally specify a custom notification template
    - Only applies when `method` is `POST`
    - For more information on template syntax, see [Alert Templates](./templates.md#alert-templates)
    - Note: Webhook templates **must** be valid JSON

```yaml title="Config File Snippet"
  webhook:
    enabled: false
    server: 
    ignoressl:
    method:
    params:
    headers:
    template:
```

## Monitor

If enabled, this application will check in with tools like [HealthChecks](https://github.com/healthchecks/healthchecks) or [Uptime Kuma](https://github.com/louislam/uptime-kuma) on a regular interval for health / status monitoring.

- **enabled** (Optional - Default: `false`)
    - Env: `FN_MONITOR__ENABLED`
    - Set to `true` to enable health checks
- **url** (Required)
    - Env: `FN_MONITOR__URL`
    - URL path for health check service
    - Required if monitoring is enabled
- **interval** (Required - Default: `60`)
    - Env: `FN_MONITOR__INTERVAL`
    - Frequency of check-in events
    - Required if monitoring is enabled
- **ignoressl** (Optional - Default: `false`)
    - Env: `FN_MONITOR__IGNORESSL`
    - Set to `true` to allow self-signed certificates

```yaml title="Config File Snippet"
monitor:
  enabled: false
  url: 
  interval: 
  ignoressl: 
```
