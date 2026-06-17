
## Discord
---

- DISCORD_TOKEN
- DISCORD_PUBLIC_KEY
- DISCORD_GUILD_ID
- DISCORD_CHANNEL_ID
- DISCORD_UPDATE_ROLE_ID
- DISCORD_SHOULD_CROSSPOST <sup>default: `true`</sup>


## Statuspage
---

- STATUSPAGE_API_KEY
- STATUSPAGE_PAGE_ID 
- STATUSPAGE_URL <sup>default: `status.ticketsbot.cloud`</sup>


## Gateway
---

- SERVER_ADDR <sup>default: `8080`</sup>


## Database
---

**Note:** The default values below are only used when using the provided `docker-compose.yaml` file.
- DATABASE_URI <sup>default: `postgres://postgres:${DATABASE_PASSWORD:-null}@postgres-statusbot:5432/postgres?sslmode=disable`</sup>
- DATABASE_PASSWORD <sup>default: `null`</sup>



# Optional Variables:

## Daemon
---

- DAEMON_ENABLED <sup>default: `false`</sup>
- DAEMON_FREQUENCY <sup>default: `30s`</sup>
- DAEMON_EXECUTION_TIMEOUT <sup>default: `30m`</sup>

## Debug
---

- JSON_LOGS <sup>default: `false`</sup>
- LOG_LEVEL <sup>default: `info`</sup>
