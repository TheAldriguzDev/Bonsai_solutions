# Bonsai Solutions

**Modular retail monitoring and web automation toolkit | 2020–2023**

Bonsai Solutions is a Go-based toolkit for monitoring product listings, new arrivals, and restocks across multiple online retailers. It combines retailer-specific modules with reusable networking utilities, proxy support, and experimental anti-bot challenge handling.

The broader Bonsai Solutions project was developed and commercialized as a subscription-based retail monitoring service, reaching **300+ paying customers**. This public repository showcases the codebase and technical approach; it is not a complete deployment of the original commercial service.

> **Project status:** Historical/portfolio project. Retailer endpoints, challenge mechanisms, and third-party integrations may have changed since active development. The code is provided as a technical reference and is not guaranteed to work against current websites.

## Highlights

- **Multi-site monitoring:** Dedicated modules for retailers such as Foot Locker, Adidas, Zalando, JD Sports, Size?, Footpatrol, ASOS, and SkateDeluxe.
- **Product and restock workflows:** Separate monitoring paths for new product listings and inventory changes where implemented.
- **HTTP/TLS networking:** Customizable request handling built with Go HTTP libraries and TLS-client tooling.
- **Proxy support:** Proxy configurations for retailer-specific monitoring workflows.
- **Anti-bot research:** A DataDome-related solver module exploring challenge and cookie-handling mechanisms.
- **Utility API:** A lightweight Fiber service exposing retailer-specific product-link redirects.

## Tech stack

| Area | Technology |
| --- | --- |
| Core language | Go (module targets Go 1.20) |
| HTTP API | Fiber v2 |
| Network requests | `fhttp`, `tls-client`, `uTLS` |
| HTML parsing | `goquery` |
| Concurrency | Go's native concurrency primitives |
| Interface | Command-line menu |

The libraries above are listed in [`go.mod`](go.mod). Individual modules may use different subsets of them.

## Repository structure

```text
Bonsai_solutions/
├── api/                 # Fiber API and URL redirect helpers
├── modules/             # Retailer-specific monitors
├── monitorData/         # Monitoring-related data
├── packages/            # Shared utilities (including logging)
├── solvers/
│   └── datadome/        # DataDome-related research components
├── main.go              # CLI entry point and module selection
├── go.mod
└── README.md
```

## Getting started

### Prerequisites

- Go 1.20 or later
- An internet connection for dependency installation and requests
- Proxy configuration where required by an individual module

### Install and run

```bash
git clone https://github.com/TheAldriguzDev/Bonsai_solutions.git
cd Bonsai_solutions
go mod download
go run .
```

The application opens a command-line menu. Select a retailer module to launch its workflow.

The original project expects a `proxies/` directory containing proxy-list files:

```bash
mkdir -p proxies
```

**Configuration note:** File names, proxy formats, and other requirements can vary by module. Consult the relevant module source before running it. No universal example configuration is provided because the repository does not define a single shared configuration format.

### Build a binary

```bash
go build -o bonsai .
./bonsai
```

On Windows, use an appropriate `.exe` output name when building.

## Utility API

The CLI also includes an option to start a small **Fiber** HTTP server on port **8080**. It provides the following endpoint:

```text
GET /redirect?site=<retailer>&sku=<product-id>
```

This endpoint constructs and redirects to a retailer-specific product URL for supported retailer identifiers. Refer to [`api/api.go`](api/api.go) for the implemented mappings. It is a utility endpoint, not a full public monitoring API.

## Technical focus

Bonsai Solutions was an exercise in building and operating a specialized web-monitoring product. The main engineering areas included:

1. **Modular integrations:** Isolating retailer-specific parsing and monitoring logic from shared utilities.
2. **Network-level behavior:** Working with HTTP clients, TLS behavior, and proxy-aware request flows.
3. **Monitoring workflows:** Handling different retail events, including new additions and restocks.
4. **Anti-bot analysis:** Studying challenge and session-cookie behavior in protected retail environments.
5. **Product delivery:** Packaging the broader monitoring service for subscription customers.

## Limitations and responsible use

- Retail websites, APIs, and defensive mechanisms change frequently; some modules may require updates.
- The repository does not provide a production deployment or a maintained configuration guide for every retailer.
- Use this code only on systems you own or are authorized to access, and respect applicable laws, website terms, privacy requirements, and rate limits.
- Bonsai Solutions is an independent project and is not affiliated with the retailers or security vendors mentioned here.

## Author

Developed by [TheAldriguzDev](https://github.com/TheAldriguzDev).

For questions about the historical project, open a GitHub issue.

## Webhook configuration

Discord webhook URLs are read from the environment variables listed in `.env.example`. Set them before starting the application. `.env` is ignored by Git but is not loaded automatically. Do not commit real webhook URLs.
