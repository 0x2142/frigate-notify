# Install / Setup

Configuration can be provided via config file or environment variables, check [here](config/file.md) for details. A configuration file template can be found [here](config/sample.md)

## Docker (Recommended)

A Docker image is built with every release. Those images can be found [here](https://github.com/0x2142/frigate-notify/pkgs/container/frigate-notify).

Each image is assigned a version tag upon release. Using the `latest` tag is recommended & will ensure you pull down the most recent version of the app. If needed, prior releases can be used by specifying the individual release tag, like `v0.3.0`.

The app can be run as a container with the bundled [docker-compose.yml](https://github.com/0x2142/frigate-notify/blob/main/docker-compose.yml) file:

```yaml
services:
  frigate-notify:
    image: ghcr.io/0x2142/frigate-notify:latest
    # Uncomment below if REST API server is enabled
    # ports:
    #   - "8000:8000"
    environment:
      - TZ=Etc/UTC
    volumes:
      - /path/to/config.yml:/app/config.yml
    restart: unless-stopped
```

Update `volumes` in the compose file to the location of the `config.yml` file. By default the app will check the `/app/config` directory for this file.

The file name & path can also be modified by specifying a path via the `FN_CONFIGFILE` environment variable:

```yaml
services:
  frigate-notify:
    image: ghcr.io/0x2142/frigate-notify:latest
    environment:
      FN_CONFIGFILE: /app/config/config.yml
      TZ: Etc/UTC
    volumes:
      - /path/to/config_directory:/app/config
    restart: unless-stopped
```

Optionally, you can run without compose:

```bash
docker run -v /path/to/config.yml:/app/config.yml ghcr.io/0x2142/frigate-notify:latest
```

!!!note
    Sensitive information (like passwords) can be configured via environment variables or docker secrets.

    Environment variables are listed in the [config file reference](config/file.md).

    For docker secrets, configure your secret - then point the appropriate environment variable to your secret file. For example, if you had a secret named `mqtt_password` and wanted to provide that to frigate-notify:

    ```yaml
    services:
      frigate-notify:
        image: ghcr.io/0x2142/frigate-notify:latest
        environment:
          FN_FRIGATE__MQTT__PASSWORD: /run/secrets/mqtt_password
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
