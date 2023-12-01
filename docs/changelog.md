# Changelog

## [v0.2.2](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.2)
 - Fix SMTP issue where code was ignoring `tls: false` config flag
 - Fix issue where snapshots would only be sent to the first alerting method, if multiple were enabled 
 - Update Debian image used for Docker image 

## [v0.2.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.1)
 - Added Docker release build for ARM

## [v0.2.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.2.0)
 - Add support for excluding cameras from alerts
 - Added ability to set custom alert titles
 - Added ability to set custom MQTT client ID
 - Added ability to change MQTT port
 - Migrate Discord alerts to use Embeds
 - Improved SMTP & Gotify alerts to include a warning if a snapshot wasn't saved
 - Add GitHub actions for automated release builds
 - Added version number to startup logging

## [v0.1.1](https://github.com/0x2142/frigate-notify/releases/tag/v0.1.1)
 - Add support for MQTT authentication

## [v0.1.0](https://github.com/0x2142/frigate-notify/releases/tag/v0.1.0)
 - Initial release! 🎉