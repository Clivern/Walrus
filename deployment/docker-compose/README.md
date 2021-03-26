## Deploy with Docker

Install docker & docker-compose

```bash
$ apt-get update
$ sudo apt install docker.io
$ sudo systemctl enable docker
$ sudo apt install docker-compose
```

Clone Walrus for docker-compose.yml file

```bash
$ git clone https://github.com/Clivern/Walrus.git walrus
$ cd walrus/deployment/docker-compose
```

Feel free to update walrus tower port, `api.key` and `api.encryptionKey`. make sure you also use these values in walrus agent config file since agents require the `tower URL`, tower `API key` and `encryptionKey` to be able to reach and communicate with walrus tower.

Run the tower and etcd. It is also recommended to run etcd anywhere where data loss is mitigated.

```bash
$ docker-compose up -d
```

Now tower should be running. User your server public IP and tower port configured before to open the dashboard and setup the admin account.

```bash
# To get the public IP
$ curl https://ipinfo.io/ip
```

In the host where backups have to take place, download walrus binary.

```bash
$ curl -sL https://github.com/Clivern/Walrus/releases/download/v0.1.8/walrus_0.1.8_Linux_x86_64.tar.gz | tar xz
```

Create agent config file. Don't forget to replace `agent.tower` configs with the `tower URL`, `apiKey` and `encryptionKey`, you can get these values from tower configs you created earlier.

```yaml
# Agent configs
agent:
    # Env mode (dev or prod)
    mode: ${WALRUS_APP_MODE:-prod}
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

Then run the host agent.

```
$ walrus agent -c /path/to/agent.config.yml
```

If everything is right, you should be able to see the host shown in the tower dashboard with one active agent. You can create backup crons under that host and update s3 configs in `settings` tab.

