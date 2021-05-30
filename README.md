<p align="center">
    <img src="https://raw.githubusercontent.com/Clivern/Walrus/main/assets/gopher.png" width="230" />
    <h3 align="center">Walrus</h3>
    <p align="center">Fast, Secure and Reliable System Backup, Set up in Minutes.</p>
    <p align="center">
        <a href="https://github.com/Clivern/Walrus/actions"><img src="https://github.com/Clivern/Walrus/workflows/Build/badge.svg"></a>
        <a href="https://github.com/Clivern/Walrus/actions"><img src="https://github.com/Clivern/Walrus/workflows/Release/badge.svg"></a>
        <a href="https://github.com/Clivern/Walrus/releases"><img src="https://img.shields.io/badge/Version-1.1.0-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Clivern/Walrus"><img src="https://goreportcard.com/badge/github.com/Clivern/Walrus?v=1.1.0"></a>
        <a href="https://hub.docker.com/r/clivern/walrus"><img src="https://img.shields.io/badge/Docker-Latest-green"></a>
        <a href="https://github.com/Clivern/Walrus/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg"></a>
    </p>
</p>
<br/>
<p align="center">
    <img src="https://raw.githubusercontent.com/Clivern/Walrus/main/assets/chart.png?v=1.1.0" width="80%" />
</p>
<p align="center">
    <h3 align="center">Dashboard Screenshots</h3>
    <p align="center">
        <img src="https://raw.githubusercontent.com/Clivern/Walrus/main/assets/screenshot_02.png?v=1.1.0" width="90%" />
        <img src="https://raw.githubusercontent.com/Clivern/Walrus/main/assets/screenshot_03.png?v=1.1.0" width="90%" />
    </p>
</p>

Walrus is a fast, secure and reliable backup system suitable for modern infrastructure. With walrus, you can backup services like SQLite, MySQL, PostgreSQL, Redis, etcd or a complete directory with a short interval and low overhead. It supports AWS S3, digitalocean spaces and any S3-compatible object storage service.


## Documentation

## Deployment

Download [the latest walrus binary](https://github.com/Clivern/Walrus/releases). Make it executable from everywhere.

```zsh
$ curl -sL https://github.com/Clivern/Walrus/releases/download/vx.x.x/walrus_x.x.x_OS.tar.gz | tar xz
```

Then install etcd cluster or single node, please refer to etcd docs or bin directory inside this repository.


#### Run Walrus Tower:

Create the tower configs file `tower.config.yml` from `config.dist.yml`. Something like the following:

Please make sure to update the `apiKey` and `encryptionKey` to a different random values.

```yaml
# Tower configs
tower:
    # Env mode (dev or prod)
    mode: ${WALRUS_APP_MODE:-dev}
    # HTTP port
    port: ${WALRUS_API_PORT:-8000}
    # URL
    url: ${WALRUS_API_URL:-http://127.0.0.1:8000}
    # TLS configs
    tls:
        status: ${WALRUS_API_TLS_STATUS:-off}
        pemPath: ${WALRUS_API_TLS_PEMPATH:-cert/server.pem}
        keyPath: ${WALRUS_API_TLS_KEYPATH:-cert/server.key}

    # API Configs
    api:
        key: ${WALRUS_API_KEY:-6c68b836-6f8e-465e-b59f-89c1db53afca}
        encryptionKey: ${WALRUS_ENCRYPTION_KEY:-B?E(H+Mb}

    # Async Workers
    workers:
        # Queue max capacity
        buffer: ${WALRUS_WORKERS_CHAN_CAPACITY:-5000}
        # Number of concurrent workers
        count: ${WALRUS_WORKERS_COUNT:-4}

    # Runtime, Requests/Response and Walrus Metrics
    metrics:
        prometheus:
            # Route for the metrics endpoint
            endpoint: ${WALRUS_METRICS_PROM_ENDPOINT:-/metrics}

    # Application Database
    database:
        # database driver
        driver: ${WALRUS_DB_DRIVER:-etcd}

        etcd:
            # etcd database name or prefix
            databaseName: ${WALRUS_DB_ETCD_DB:-walrus}
            # etcd username
            username: ${WALRUS_DB_ETCD_USERNAME:- }
            # etcd password
            password: ${WALRUS_DB_ETCD_PASSWORD:- }
            # etcd endpoints
            endpoints: ${WALRUS_DB_ETCD_ENDPOINTS:-http://127.0.0.1:2379}
            # Timeout in seconds
            timeout: 30

    # Log configs
    log:
        # Log level, it can be debug, info, warn, error, panic, fatal
        level: ${WALRUS_LOG_LEVEL:-info}
        # output can be stdout or abs path to log file /var/logs/walrus.log
        output: ${WALRUS_LOG_OUTPUT:-stdout}
        # Format can be json
        format: ${WALRUS_LOG_FORMAT:-json}
```

The run the `tower` with `systemd`

```
walrus tower -c /path/to/tower.config.yml
```


#### Run Walrus Agent:

Create the agent configs file `agent.config.yml` from `config.dist.yml`. Something like the following:

```yaml
# Agent configs
agent:
    # Env mode (dev or prod)
    mode: ${WALRUS_APP_MODE:-dev}
    # HTTP port
    port: ${WALRUS_API_PORT:-8001}
    # URL
    url: ${WALRUS_API_URL:-http://127.0.0.1:8001}
    # TLS configs
    tls:
        status: ${WALRUS_API_TLS_STATUS:-off}
        pemPath: ${WALRUS_API_TLS_PEMPATH:-cert/server.pem}
        keyPath: ${WALRUS_API_TLS_KEYPATH:-cert/server.key}

    # API Configs
    api:
        key: ${WALRUS_API_KEY:-56e1a911-cc64-44af-9c5d-8c7e72ec96a1}

    # Async Workers
    workers:
        # Queue max capacity
        buffer: ${WALRUS_WORKERS_CHAN_CAPACITY:-5000}
        # Number of concurrent workers
        count: ${WALRUS_WORKERS_COUNT:-4}

    # Tower Configs
    tower:
        url: ${WALRUS_TOWER_URL:-http://127.0.0.1:8000}
        # This must match the one defined in tower config file
        apiKey: ${WALRUS_TOWER_API_KEY:-6c68b836-6f8e-465e-b59f-89c1db53afca}
        # This must match the one defined in tower config file
        encryptionKey: ${WALRUS_ENCRYPTION_KEY:-B?E(H+Mb}
        # Time interval between agent ping checks
        pingInterval: ${WALRUS_CHECK_INTERVAL:-60}

    # Backup settings
    backup:
        tmpDir: ${WALRUS_BACKUP_TMP_DIR:-/tmp}

    # Log configs
    log:
        # Log level, it can be debug, info, warn, error, panic, fatal
        level: ${WALRUS_LOG_LEVEL:-info}
        # output can be stdout or abs path to log file /var/logs/walrus.log
        output: ${WALRUS_LOG_OUTPUT:-stdout}
        # Format can be json
        format: ${WALRUS_LOG_FORMAT:-json}
```

The run the `agent` with `systemd`

```
walrus agent -c /path/to/agent.config.yml
```

Now you can open the walrus tower dashboard `http://127.0.0.1:8000` and start the setup.


#### To run the Admin Dashboard (Development Only):

Clone the project or your own fork:

```zsh
$ git clone https://github.com/Clivern/Walrus.git
```

Create the dashboard config file `web/.env` from `web/.env.dist`. Something like the following:

```
VUE_APP_TOWER_URL=http://localhost:8080
```

Then you can either build or run the dashboard

```zsh
# Install npm packages
$ cd web
$ npm install
$ npm install -g npx

# Add tower url to frontend
$ echo "VUE_APP_TOWER_URL=http://127.0.0.1:8000" > .env

$ cd ..

# Validate js code format
$ make check_ui_format

# Format UI
$ make format_ui

# Run Vuejs app
$ make serve_ui

# Build Vuejs app
$ make build_ui

# Any changes to the dashboard, must be reflected to cmd/pkged.go
# You can use these commands to do so
$ go get github.com/markbates/pkger/cmd/pkger
$ make package
```

## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Walrus is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/walrus/releases) for changelogs for each release version of Walrus. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/walrus/issues


## Security Issues

If you discover a security vulnerability within Walrus, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2020, Clivern. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Walrus** is authored and maintained by [@clivern](http://github.com/clivern).
