# Changelog

## [v0.5.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.5.3) - TBD
 - Fix issue with MQTT topic parsing
 - Add configurable message priority for Gotify
 - Add message template variable for event start time as Unix timestamp
 - Add custom User-Agent for outbound HTTP requests

## [v0.5.2](https://github.com/0x2142/frigate-notify/releases/tag/v0.5.2) - May 20 2025
 - Fix minor issue with additional default config values being overwritten if configured value is empty

## [v0.5.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.5.1) - May 13 2025
 - Fix issue with default startup check values being overwritten if configured value is empty

## [v0.5.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.5.0) - May 12 2025
 - Updated configuration management system
     - ⚠️ Environment variable keys are now separated by double-underscore. See [docs](https://frigate-notify.0x2142.com/latest/config/file) for examples
     - ⚠️ Please note that the `apprise-api` config is now `apprise_api`
     - Docker secrets are now supported
 - Add [option](https://frigate-notify.0x2142.com/latest/config/file/#general) to collect higher resolution snapshot from camera main stream
 - Add ability to send video clip instead of snapshot image via [Telegram](https://frigate-notify.0x2142.com/latest/config/file/#telegram)
 - Add ability to use Frigate [authenticated UI](https://frigate-notify.0x2142.com/latest/config/file/#server) on port 8971
 - Add ability to configure custom notification titles per alert provider profile, for those that support it
 - Add support for email threading via SMTP notifier

## [v0.4.4](https://github.com/0x2142/frigate-notify/releases/tag/v0.4.4) - Apr 08 2025
 - Add support for notifications via [Matrix](https://frigate-notify.0x2142.com/latest/config/file/#matrix)
 - Add support for `message_thread_id` to [Telegram](https://frigate-notify.0x2142.com/latest/config/file/#telegram) notifications
 - Add review/event link for Pushover action button
 - Add ability to disable Discord notifications being sent as embedded message
 - Add ability to wait for [license plate recognition](https://frigate-notify.0x2142.com/latest/config/file/#license-plate) data from Frigate to be included in notifications
 - Add ability to allow / block license plates from generating notification
 - Fix issue with Webhook default message template

## [v0.4.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.4.3) - Feb 21 2025
 - Add support for notifications via [Apprise API](https://frigate-notify.0x2142.com/latest/config/file/#apprise-api)
 - Add retry interval for collecting snapshots from Frigate
 - Updated notification templates to include sublabel information, if available
 - Fix app crash when no items in a review were eligible for notification

## [v0.4.2](https://github.com/0x2142/frigate-notify/releases/tag/v0.4.2) - Feb 10 2025
 - Fix issue with Signal notifications when no snapshot image is available
 - Fix issue with sending Signal notifications to groups or usernames

## [v0.4.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.4.1) - Feb 06 2025
 - Add support for notifications via [Signal](https://frigate-notify.0x2142.com/latest/config/file/#signal)
 - Add support for notifications via [Mattermost](https://frigate-notify.0x2142.com/latest/config/file/#mattermost)
 - Fix issue with alert-level filters where filters may prevent notifications on reviews with multiple object detections

## [v0.4.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.4.0) - Jan 27 2025
 - Support for notification based on Alerts & Detections via Frigate Reviews
     - ⚠️ **Note:** Reviews mode is now the default with this release
     - To continue using `events` mode, please see the [docs](https://frigate-notify.0x2142.com/latest/config/file/#app)
 - Add initial support for a [REST API](https://frigate-notify.0x2142.com/latest/api/)
     - Retrieve, edit, or reload configuration
     - Enable/disable notifications & send test notifications
     - App status/healthchecks
 - Added support for [multiple notification endpoints](https://frigate-notify.0x2142.com/latest/config/profilesandfilters/) per notification provider
     - Includes ability to filter which endpoint is used based on event conditions
     - For example: Two SMTP alert endpoints, each with different contacts. One group for high priority detections (person, car) & the other for anything else
 - Added support for specifying custom [Pushover](https://frigate-notify.0x2142.com/latest/config/file/#pushover) sounds
 - Added support for [re-checking](https://frigate-notify.0x2142.com/latest/config/file/#general) event details from Frigate after a delay
     - Useful if waiting for a 3rd-party app to update `sub_label`
 - Added support for processing audio events & optionally dropping [audio-only](https://frigate-notify.0x2142.com/latest/config/file/#general) events
 - Camera names are now transformed to title case (ex. "side_door" becomes "Side Door")
 - Fixed issue with retry logic used for webhook notifications

## [v0.3.5](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.5) - Oct 08 2024
 - Fixed issue where built-in alert templates were not being included in binary releases
 - Fixed issue where a notification may not be sent if previous event update from Frigate did not contain a snapshot
 - Fixed issue where Gotify snapshot was not using `public_url` if configured
 - Added `from` & `ignoressl` config items to `smtp` notifier
 - Added ability to send `webhook` notifications via HTTP GET requests
 - Added support for URL parameters with `webhook` notifications
 - Added option to only generate a [single notification](https://frigate-notify.0x2142.com/latest/config/file/#general) per Frigate event
 - Allow use of [template variables](https://frigate-notify.0x2142.com/latest/config/templates/#title-template) within alert `title` config
 - New options to set specific [log level](https://frigate-notify.0x2142.com/latest/config/options/) & additional `trace` level logs for troubleshooting
 - Minor enhancements to config file validation

## [v0.3.4](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.4) - Aug 15 2024
 - Fixed issue where `unzoned: drop` config would drop all notifications

## [v0.3.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.3) - Aug 14 2024
 - Additional filtering options for alerts:
     - [Sublabels](https://frigate-notify.0x2142.com/latest/config/file/#sublabels)
     - [Minimum label score](https://frigate-notify.0x2142.com/latest/config/file/#labels)
     - [Quiet hours](https://frigate-notify.0x2142.com/latest/config/file/#quiet-hours)
 - Ability to [modify snapshot](https://frigate-notify.0x2142.com/latest/config/file/#general) when retrieved from Frigate: include bounding box / timestamp or crop snapshot
 - Allow use of [template variables](https://frigate-notify.0x2142.com/latest/config/templates/#available-variables) within HTTP headers sent to notification providers
     - Ntfy action button can now be overrriden by defining a custom `X-Actions` header in the config file
     - Templates can also now include information from [environment variables](https://frigate-notify.0x2142.com/latest/config/templates/#environment-variables)

## [v0.3.2](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.2) - Jun 13 2024

 - Allow config to be provided via environment variables (Thanks [@ryan-willis](https://github.com/ryan-willis)!!)

## [v0.3.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.1) - Jun 11 2024

 - Fix misspelling for Ntfy
     - Currently config will accept either `nfty` or `ntfy`, but in a future release `nfty` will be removed
     - ⚠️ **Note**: If you are using `nfty`, please update your config to the correct spelling: `ntfy`
 - Fix issue where Ntfy action button would use local Frigate URL even if `public_url` was configured
 - Fix camera links for Frigate 0.14.0+
 - Add support for sending HTTP headers with Ntfy notifications
 - Added [config file](https://frigate-notify.0x2142.com/latest/config/file/#general) option to drop events that don't have a snapshot image
 - Move "No snapshot available" message into notification templates
 - Add support for disabling colored console logging via [config options](https://frigate-notify.0x2142.com/latest/config/options/)

## [v0.3.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.3.0) - Jun 03 2024

 - Reworked & reduced logging, added [debug option](https://frigate-notify.0x2142.com/latest/config/options/)
 - Add ability to define custom [notification templates](https://frigate-notify.0x2142.com/latest/config/templates/#alert-templates)
 - Add support for notifications via [webhook](https://frigate-notify.0x2142.com/latest/config/file/#webhook)
 - Add retry logic to startup check
     - App will now retry up to 5 times, checking once every 30 seconds
     - [Config](https://frigate-notify.0x2142.com/latest/config/file/#server) options now available to modify max attempts & retry interval
 - Enable concurrent notification processing
 - Fix issue where some event providers would fail if no snapshot was available via Frigate

## [v0.2.8](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.8) - May 15 2024

 - Add support for notifications via [Ntfy](https://frigate-notify.0x2142.com/latest/config/file/#ntfy)
 - Add ability to send additional HTTP [headers](https://frigate-notify.0x2142.com/latest/config/file/#frigate) to Frigate
 - Add new `public_url` config item for Frigate
     - This will be used in notification links & should be configured if Frigate is accessible via the internet
 - Add startup check to verify Frigate API is accessible
 - Rework event notifications to be built from templates
 - Fix default interval for querying evens via web API
 - Fix issue where label score is 0% via web API event query

## [v0.2.7](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.7) - May 06 2024

 - Allow changing default MQTT topic prefix via config
 - Added ability to filter notifications based on [labels](https://frigate-notify.0x2142.com/latest/config/file/#labels)
     - New option to allow only notifications with specified labels
     - New option to deny notifications based on labels

## [v0.2.6](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.6) - Apr 01 2024

 - Fixed issue with setting `unzoned: drop` under zone config, where alerts wouldn't be sent if event began outside a zone.
 - Fixed issue with correct timezone getting applied on container image
 - Added `timeformat` option for notifications, which allows custom date/time format
 - Minor updates to logging to ensure all event logs have an event ID attached

## [v0.2.5](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.5) - Mar 29 2024

 - Added support for alerts via [Pushover](https://frigate-notify.0x2142.com/latest/config/file/#pushover)
 - Reduced container image from ~1Gb to ~11Mb

## [v0.2.4](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.4) - Mar 28 2024

 - Added support for alerts via [Telegram](https://frigate-notify.0x2142.com/latest/config/file/#telegram)
 - Updated notification message to now include links to the camera & event clip

## [v0.2.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.3) - Mar 14 2024

 - Updated config management
 - Improved config file validation
 - Fixed issue where HTTP requests would fail if only IP & port were specified in config
 - Fixed crash if Frigate has no event info & app couldn't pull snapshots
 - Added ability to filter notifications based on [zones](https://frigate-notify.0x2142.com/latest/config/file/#zones)
     - New config option to only allow notifications based on zone
     - New config option to deny notifications based on zone
     - New config option to allow or drop events outside of a designated zone

## [v0.2.2](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.2) - Nov 29 2023

 - Fix SMTP issue where code was ignoring `tls: false` config flag
 - Fix issue where snapshots would only be sent to the first alerting method, if multiple were enabled
 - Update Debian image used for Docker image

## [v0.2.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.1) - Sep 20 2023

 - Added Docker release build for ARM

## [v0.2.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.0) - Sep 19 2023

 - Add support for excluding cameras from alerts
 - Added ability to set custom alert titles
 - Added ability to set custom MQTT client ID
 - Added ability to change MQTT port
 - Migrate Discord alerts to use Embeds
 - Improved SMTP & Gotify alerts to include a warning if a snapshot wasn't saved
 - Add GitHub actions for automated release builds
 - Added version number to startup logging

## [v0.1.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.1.1) - Apr 25 2023

 - Add support for MQTT authentication

## [v0.1.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.1.0) - Mar 17 2023

 - Initial release! 🎉
