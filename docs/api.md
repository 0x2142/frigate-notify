# API

Frigate-Notify has a REST API available for managing configuration, notifications, and retrieving app status.

## Setup

The API is disabled by default. To enable, set `enabled: true` under **app > api** in the app configuration:

```yaml title="Config File Snippet"
app:
  api:
    enabled: true
    port: 8000
```

For more details, see the [config](./config/file.md#app) docs.

## Usage

Once the API is enabled, interactive documentation can be found at `:8000/docs`.

API endpoints will be under the following prefix: `:8000/api/v1`.

OpenAPI spec available at `:8000/openapi.json` or `:8000/openapi.yaml`.

## Supported Operations

### Config

 - (GET / PUT) `/api/v1/config`
     - Retrieve or set application configuration

### Control

 - (GET / POST) `/api/v1/notif_state`
     - Retrieve or set global notification state
     - Can be used to dynamically silence all notifications from Frigate-Notify

 - (POST) `/api/v1/notif_test`
     - Trigger test notification via all configured notification providers
     - Can be used to test alert filters, templates, or alert provider configuration
     - Note: This works by manually querying Frigate for the last received event & sending a notification based on that event
         - Notifications sent via this method will bypass global filters, but are **still processed** by alert-level filters

 - (POST) `/api/v1/reload`
     - Trigger reload of configuration & restart of application

### Status

 - (GET) `/api/v1/status`
     - Retrieve application status
     - Includes:
         - App health status & reachability of Frigate API and MQTT broker (if used)
         - Stats on last Frigate event & last notification sent
         - Stats on alerts sent/failed & errors for each notification provider

 - (GET) `/api/v1/version`
     - Retrieve application version

 - (GET) `/api/v1/healthz`
     - Retrieve application health status
     - Returns `ok` if app is healthy

 - (GET) `/api/v1/readyz`
     - Retrieve application ready status
     - Returns `ok` if app is ready
