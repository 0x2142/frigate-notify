# Frigate-notify

## About

This project is designed to generate notifications based on [Frigate](https://github.com/blakeblackshear/frigate) NVR events.

Currently Frigate only supports notifications through Home Assistant, which I don't use. So I set out to build a simple notification app that would work with a standalone Frigate server.

## Features

**Event Polling**
- MQTT
- Direct via Frigate API

**Notification Methods**
- Discord
- Gotify
- SMTP

**Other**
- Aliveness monitor via HTTP GET (for use with tools like [Uptime Kuma](https://github.com/louislam/uptime-kuma))

## Setup 

### Configuration

An example config file has been provided (`example-config.yml`). Please copy this & edit to fit your needs.

The example config file contains additional details about how to configure the app. A few important notes:
- Only one event polling method can be configured - Either MQTT or API.
- Regardless of polling method, app needs access to Frigate API to pull screenshots
- Any combination of notification services can be enabled

By default, this app will look for a configuration file in the local directory named `config.yml`. 
- This can be changed via `-c` flag (example: `./frigate-notify -c newconfig.yml`)
- or specifying a location via environment variable: `FN_CONFIGFILE`


### Run

Run locally with 
```bash
$ ./frigate-notify
```

Or as a container:
```bash
$ docker run -v /path/to/config.yml:/app/config.yml ghcr.io/0x2142/frigate-notify:latest
```
## Future

Just a quick list of things I may add later:
- Ability to specify which cameras to alert on
- Dampening time between notifications
- Additional alerting methods


## Screenshots

**Discord**

![Discord](/screenshots/discord.png)

**Gotify**

![Gotify](/screenshots/gotify.png)