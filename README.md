<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]

<br />
<div align="center">

<h3 align="center">Status Updates</h3>

  <p align="center">
    This bot is a simple way to keep your Discord community up to date with StatusPage updates.
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

Status Updates is a backend service that keeps your Discord community informed by syncing incident and status updates from Statuspage.io and Jira. Built with Go for reliability and extensibility, it powers the core API and background processing for the Tickets Bot ecosystem, enabling seamless integration of incident management, ticket tracking, and public status communication.

**Key Features:**

- **Incident & Status Management:** Track incidents, components, and status updates for your Discord community.
- **Jira Integration:** Sync incidents with Atlassian StatusPage for advanced incident management.
- **Extensible:** Modular internal structure for easy feature addition.

### Built With

- [Go (Golang)](https://go.dev/)

---

# Installation

## Using Docker


1. **Clone the repository:**
   ```sh
   git clone https://github.com/TicketsBot-cloud/status-updates.git
   cd status-updates
   ```
2. **Create a .env file by copying the provided [.env.example](./.env.example) file.**
  Fill out all Env Variables in the sections that are Required.
3. **Start the Bot**  
  Run the following command in the directory the Bot's files were copied to:
  `docker compose up -d`
4. **Set the Interaction Endpoint URL**
  Go to the [Discord Developer Portal](https://discord.com/developers/applications) and click on the Application you made for this bot.
  Set the Interaction Ednpoint URL to your Proxied domain with `/interactions` added behind it. (e.g. `https://{YOUR-CUSTOM-DOMAIN}/interactions`)
    * You will need to proxy the port you set as your SERVER_ADDR (e.g. 8080) to a Publicly Accessible URL.
    * Replace `{YOUR-CUSTOM-DOMAIN}` with your proxied URL (e.g. `gateway.example.com`)


## Manual


1. **Clone the repository:**
   ```sh
   git clone https://github.com/TicketsBot-cloud/status-updates.git
   cd status-updates
   ```
2. **Install dependencies:**
   ```sh
   go mod download
   ```
3. **Build the project:**
   ```sh
   go build -o status-updates ./cmd/status-updates
   ```

### Configuration

Configuration is managed via environment variables or a config file. See `internal/config/config.go` for all options.

**Example environment variables:**
```env
DISCORD_TOKEN=bot_token
DISCORD_PUBLIC_KEY=bot_public_key
DISCORD_GUILD_ID=guild_id
DISCORD_CHANNEL_ID=status_updates_channel_id
DISCORD_UPDATE_ROLE_ID=status_update_role_id

STATUSPAGE_API_KEY=statuspage_api_key
STATUSPAGE_PAGE_ID=statuspage_page_id

DATABASE_URI=postgres://postgres:postgres@localhost/postgres?sslmode=disable
```

### Usage

To run the service locally:
```sh
go run cmd/status-updates/main.go
```

The HTTP server will start and listen on the configured port (default: 8080).

The HTTP server is only used for the buttons to add the user to the role and thread for an incident.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or documentation improvements.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

---

## License

Distributed under the MIT License. See `LICENSE` for details.


<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/TicketsBot-cloud/status-updates.svg?style=for-the-badge
[contributors-url]: https://github.com/TicketsBot-cloud/status-updates/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/TicketsBot-cloud/status-updates.svg?style=for-the-badge
[forks-url]: https://github.com/TicketsBot-cloud/status-updates/network/members
[stars-shield]: https://img.shields.io/github/stars/TicketsBot-cloud/status-updates.svg?style=for-the-badge
[stars-url]: https://github.com/TicketsBot-cloud/status-updates/stargazers
[issues-shield]: https://img.shields.io/github/issues/TicketsBot-cloud/status-updates.svg?style=for-the-badge
[issues-url]: https://github.com/TicketsBot-cloud/status-updates/issues
[license-shield]: https://img.shields.io/github/license/TicketsBot-cloud/status-updates.svg?style=for-the-badge
[license-url]: https://github.com/TicketsBot-cloud/status-updates/blob/master/LICENSE.txt

[Golang]: https://img.shields.io/badge/Go-%2300ADD8?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
