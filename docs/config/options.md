# Config Options

The following options are available as command line flags or environment variables:

| Flag        | Environment Variable | Description                                                                                              |
|-------------|----------------------|----------------------------------------------------------------------------------------------------------|
| -c          | FN_CONFIGFILE        | Specify alternate config file location                                                                   |
| -debug      | FN_DEBUG             | Set to `true` to enable debug logging (Overrides -loglevel)                                              |
| -loglevel   | FN_LOGLEVEL          | Specify desired log level: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` (Default: `info`) |
| -jsonlog    | FN_JSONLOG           | Set to `true` to enable logging in JSON                                                                  |
| -nocolor    | FN_NOCOLOR           | Set to `true` to disable color for console logging                                                       |
