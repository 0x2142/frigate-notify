# Changelog

## [v0.2.3](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.3) - Upcoming Release

 - Updated config management
 - Improved config file validation
 - Fixed issue where HTTP requests would fail if only IP & port were specified in config
 - Fixed crash if Frigate has no event info & app couldn't pull snapshots

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