# Changelog

## [v0.2.8](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.8) - Upcoming Release

- Add support for notifications via [Nfty](https://frigate-notify.0x2142.com/config/#nfty)

## [v0.2.7](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.7) - May 06 2024

- Allow changing default MQTT topic prefix via config
- Added ability to filter notifications based on [labels](https://frigate-notify.0x2142.com/config/#labels)
    - New option to allow only notifications with specified labels
    - New option to deny notifications based on labels

## [v0.2.6](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.6) - Apr 01 2024

 - Fixed issue with setting `unzoned: drop` under zone config, where alerts wouldn't be sent if event began outside a zone.
 - Fixed issue with correct timezone getting applied on container image
 - Added `timeformat` option for notifications, which allows custom date/time format
 - Minor updates to logging to ensure all event logs have an event ID attached

## [v0.2.5](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.5) - Mar 29 2024

 - Added support for alerts via [Pushover](https://frigate-notify.0x2142.com/config/#pushover)
 - Reduced container image from ~1Gb to ~11Mb

## [v0.2.4](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.4) - Mar 28 2024

 - Added support for alerts via [Telegram](https://frigate-notify.0x2142.com/config/#telegram)
 - Updated notification message to now include links to the camera & event clip

## [v0.2.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.3) - Mar 14 2024

 - Updated config management
 - Improved config file validation
 - Fixed issue where HTTP requests would fail if only IP & port were specified in config
 - Fixed crash if Frigate has no event info & app couldn't pull snapshots
 - Added ability to filter notifications based on [zones](https://frigate-notify.0x2142.com/config/#zones)
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