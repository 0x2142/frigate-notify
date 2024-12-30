# Alert Profiles & Filters

While the config examples & guides mostly demonstrate creating a single alert profile per provider, it is possible to define multiple. For example, you could create multiple Discord profiles - each sending notifications to a different Discord channel based on configured filters.

## Alert Profiles

In order to configure multiple profiles for an alert provider, we just create a YAML list with each item being an instance of the alert provider config.

For example, if we configured a single profile each for Discord & Ntfy - we may have a configuration that looks like the following:

```yaml title="Config File Snippet"
alerts:  
  discord:
    enabled: true
    webhook: https://discord.com/api/webhooks/<webhook-one>
  ntfy:
    enabled: true
    server: https://ntfy.your.domain.tld
    topic: frigate
    ignoressl: true
```

Now if we wanted to create multiple profiles for Discord, our configuration would look like the following:

```yaml title="Config File Snippet"
alerts:  
  discord:
    - enabled: true
      webhook: https://discord.com/api/webhooks/<webhook-one>
    - enabled: true
      webhook: https://discord.com/api/webhooks/<webhook-two>
  ntfy:
    enabled: true
    server: https://ntfy.your.domain.tld
    topic: frigate
    ignoressl: true
```

When an alert is sent via frigate-notify, each profile is processed in the order of configuration. So the Discord profile with `<webhook-one>` would be processed first, then the profile with `<webhook-two>`.

In addition, each profile is assigned an internal `id` by frigate-notify, which you may see in the logs & API. These IDs are assigned in incremental order of configuration, starting at `0`. In the example above, Discord profile `<webhook-one>` would be ID `0`, and `<webhook-two>` would be ID `1`.

## Alert Filters

All alert profiles can also used based on conditional filters. This can be helpful if you have multiple profiles for a single notification provider, or if you're using multiple different providers and want to control which notifications are sent to each.

When configuring filters for an alert provider, you are specifying the criteria to match for a notification to be sent via this provider. You may configure any or all of the available filters - however please note that if multiple filter criteria are applied, then the event must match at least one item in **all** filters to be sent.

For example, if we configure Ntfy to only send alerts on zone `driveway`, then only events with that zone will be sent via Ntfy. If we also configured that same profile to filter on labels `person` and `car`, then an event would need to match both the `driveway` zone **AND** one of the two labels (`person` **OR** `car) to be sent.

!!!info
    Global event filters are processed **before** alert-level filters.

    For example, if a configuration is set to block zone `driveway`, but an alert-level filter is configured to send notifications on zone `driveway` - then the event would be dropped by the global setting.

The following filters are currently available for alerts:

- **zones** - List of one or more zones
- **labels** - List of one or more labels
- **sublabels** - List of one or more sublabels
- **cameras** - List of one or more cameras
- **quiet** - Start/Stop times for quiet hours (see [here](https://frigate-notify.0x2142.com/latest/config/file/#quiet-hours) for more information on how to configure this)

Example below uses Ntfy to demonstrate configuring filters - but this works with any alert provider:

```yaml title="Config File Snippet"
alerts: 
  ntfy:
    enabled: true
    server: https://ntfy.your.domain.tld
    topic: frigate
    filters:
        zones:
          - driveway
        labels:
          - car
          - person
        sublabels:
        cameras:
          - camera_one
          - camera_two
          - camera_three
        quiet:
          start: 09:00
          end: 18:00
```
