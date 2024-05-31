# Frigate-notify

<div align="center">

[![Static Badge](https://img.shields.io/badge/Project-Documentation-blue)](https://frigate-notify.0x2142.com) [![GitHub Repo stars](https://img.shields.io/github/stars/0x2142/frigate-notify)]() [![GitHub release (with filter)](https://img.shields.io/github/v/release/0x2142/frigate-notify)](https://github.com/0x2142/frigate-notify/releases) [![Static Badge](https://img.shields.io/badge/Docker-latest-blue)](https://github.com/0x2142/frigate-notify/pkgs/container/frigate-notify)

</div>

## About

This project is designed to generate event notifications from a standalone [Frigate](https://github.com/blakeblackshear/frigate) NVR instance.

Currently Frigate only supports notifications through Home Assistant, which I'm not using right now. So I set out to build a simple notification app that would work with a standalone Frigate server.

## Features

**Event Polling**
- MQTT
- Direct via Frigate API

**Notification Methods**
- Discord
- Gotify
- SMTP
- Telegram
- Pushover
- Nfty
- Generic Webhook

**Other**
- Aliveness monitor via HTTP GET (for use with tools like [HealthChecks](https://github.com/healthchecks/healthchecks) or [Uptime Kuma](https://github.com/louislam/uptime-kuma))


## Installation

See [Install / Setup](https://frigate-notify.0x2142.com/install/).

## Configuration

An example config file has been provided (`example-config.yml`). Please copy this & edit to fit your needs.

The sample config contains inline descriptions for each field. For additional details on the available config options & the defaults, see the [Config](https://frigate-notify.0x2142.com/config/) docs.

## Screenshots

**Discord**

![Discord](/screenshots/discord.png)

**Gotify**

![Gotify](/screenshots/gotify.png)

**SMTP**

![SMTP](/screenshots/smtp.png)

**Telegram**

![Telegram](/screenshots/telegram.png)

## Changelog

Additional details on each release are in the [Changelog](https://frigate-notify.0x2142.com/changelog/)

## Support

<a href="https://www.buymeacoffee.com/0x2142" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-blue.png" alt="Buy Me A Coffee" height="41" width="174"></a>
