# Install / Setup 

Both methods below will require a [configuration](config.md) file to be provided. A template can be found [here](config.md#sample-config)

## Docker (Recommended)

A Docker image is built with every release. Those images can be found [here](https://github.com/0x2142/frigate-notify/pkgs/container/frigate-notify).

The app can be run as a container with the bundled [docker-compose.yml](https://github.com/0x2142/frigate-notify/blob/main/docker-compose.yml) file:

```yaml
version: "3.0"
services:
  frigate-notify:
    image: ghcr.io/0x2142/frigate-notify:latest
    environment:
      - TZ=Etc/UTC
    volumes:
      - /path/to/config:/app/config
    restart: unless-stopped

```

Update `volumes` in the compose file to the location of the `config.yml` file. By default the app will check the `/app/config` directory for this file.

The file name & path can also be modified by specifying a path via the `FN_CONFIGFILE` environment variable.

Or optionally, run without docker compose:

```bash
docker run -v /path/to/config.yml:/app/config.yml ghcr.io/0x2142/frigate-notify:latest
```

## Local

Binary executables are provided with every [release](https://github.com/0x2142/frigate-notify/releases).

Download the latest version & run with:

```bash
./frigate-notify
```

By default, the app will look for a config file in the same directory named `config.yml`.

You can manually specify a config file name or location with the `-c` flag:

```bash
./frigate-notify -c <new-config-file-location>
```