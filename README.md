<p align="center">
    <img src="/assets/gopher.png" width="230" />
    <h3 align="center">Walrus</h3>
    <p align="center">Fast, Secure and Reliable System Backup, Set up in Minutes.</p>
    <p align="center">
        <a href="https://travis-ci.com/Clivern/Walrus"><img src="https://travis-ci.com/Clivern/Walrus.svg?branch=master"></a>
        <a href="https://github.com/Clivern/Walrus/releases"><img src="https://img.shields.io/badge/Version-0.0.1-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Clivern/Walrus"><img src="https://goreportcard.com/badge/github.com/Clivern/Walrus?v=0.0.1"></a>
        <a href="https://github.com/Clivern/Walrus/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg"></a>
    </p>
</p>
<br/>
<p align="center">
    <img src="/assets/chart.png?v=0.0.1" width="60%" />
</p>


Walrus is a fast, secure and reliable backup system. With walrus, you can backup services like MySQL, PostgreSQL, Redis or a complete directory with a short interval and low overhead. It supports AWS S3, digitalocean spaces and any S3-compatible object storage service.


## Documentation

## Deployment

Download [the latest walrus binary](https://github.com/Clivern/Walrus/releases). Make it executable from everywhere.

```zsh
$ curl -sL https://github.com/Clivern/Walrus/releases/download/x.x.x/walrus_x.x.x_OS.tar.gz | tar xz
```


#### To run walrus as a tower:

Create the tower configs file `tower.config.yml` from `config.dist.yml`. Something like the following:

```yaml
# Tower configs
tower:
    # Env mode (dev or prod)
    mode: ${WALRUS_APP_MODE:-dev}
    # HTTP port
    port: ${WALRUS_API_PORT:-8080}
    # TLS configs
    tls:
        status: ${WALRUS_API_TLS_STATUS:-off}
        pemPath: ${WALRUS_API_TLS_PEMPATH:-cert/server.pem}
        keyPath: ${WALRUS_API_TLS_KEYPATH:-cert/server.key}

    # Message Broker Configs
    broker:
        # Broker driver (native)
        driver: ${WALRUS_BROKER_DRIVER:-native}
        # Native driver configs
        native:
            # Queue max capacity
            capacity: ${WALRUS_BROKER_NATIVE_CAPACITY:-5000}
            # Number of concurrent workers
            workers: ${WALRUS_BROKER_NATIVE_WORKERS:-4}

    # API Configs
    api:
        key: ${WALRUS_API_KEY:- }
        encryptionKey: ${WALRUS_ENCRYPTION_KEY:- }

    # Runtime, Requests/Response and Walrus Metrics
    metrics:
        prometheus:
            # Route for the metrics endpoint
            endpoint: ${WALRUS_METRICS_PROM_ENDPOINT:-/metrics}

    # Application Database
    database:
        # Database driver (sqlite3, mysql)
        driver: ${WALRUS_DATABASE_DRIVER:-sqlite3}
        # Hostname
        host: ${WALRUS_DATABASE_MYSQL_HOST:-localhost}
        # Port
        port: ${WALRUS_DATABASE_MYSQL_PORT:-3306}
        # Database
        name: ${WALRUS_DATABASE_MYSQL_DATABASE:-walrus.db}
        # Username
        username: ${WALRUS_DATABASE_MYSQL_USERNAME:-root}
        # Password
        password: ${WALRUS_DATABASE_MYSQL_PASSWORD:-root}

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


#### To run walrus as an agent:

Create the agent configs file `agent.config.yml` from `config.dist.yml`. Something like the following:

```yaml
# Agent configs
agent:
    # Env mode (dev or prod)
    mode: ${WALRUS_APP_MODE:-dev}
    # HTTP port
    port: ${WALRUS_API_PORT:-8081}
    # TLS configs
    tls:
        status: ${WALRUS_API_TLS_STATUS:-off}
        pemPath: ${WALRUS_API_TLS_PEMPATH:-cert/server.pem}
        keyPath: ${WALRUS_API_TLS_KEYPATH:-cert/server.key}

    # Message Broker Configs
    broker:
        # Broker driver (native)
        driver: ${WALRUS_BROKER_DRIVER:-native}
        # Native driver configs
        native:
            # Queue max capacity
            capacity: ${WALRUS_BROKER_NATIVE_CAPACITY:-5000}
            # Number of concurrent workers
            workers: ${WALRUS_BROKER_NATIVE_WORKERS:-4}

    # Tower Configs
    tower:
        url: ${WALRUS_TOWER_URL:-http://127.0.0.1:8080}

        # This must match the one defined in tower config file
        apiKey: ${WALRUS_TOWER_API_KEY:- }

        # This must match the one defined in tower config file
        encryptionKey: ${WALRUS_ENCRYPTION_KEY:- }

    # API Configs
    api:
        key: ${WALRUS_API_KEY:- }

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


#### To run the admin dashboard:

Clone the project:

```zsh
$ git clone https://github.com/Clivern/Walrus.git
```

Create the dashboard config file `web/.env` from `web/.env.dist`. Something like the following:

```
VUE_APP_TOWER_URL=http://localhost:8080
```

Then you can either build or run the dashboard

```zsh
$ make serve_ui

$ make build_ui
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
