# cayman - System Monitoring Dashboard

A real-time system monitoring dashboard built with Go and SvelteKit. cayman provides live metrics for CPU usage, system load, memory, and systemd unit status through a modern web interface with real-time updates via Server-Sent Events (SSE).

## Features

- **Real-time System Monitoring**: Live updates for CPU usage, system load averages, and memory information
- **Systemd Integration**: Monitor active and failed systemd units
- **Modern Web Interface**: Built with SvelteKit, TypeScript, and Tailwind CSS
- **Server-Sent Events**: Real-time data streaming without polling
- **Responsive Design**: Works on desktop and mobile devices
- **Embedded Frontend**: Single binary deployment with embedded static assets

## Tech Stack

### Backend
- **Go 1.24.5**: Core backend language
- **Echo v4**: HTTP web framework
- **go-sse**: Server-Sent Events implementation
- **gopsutil**: System and process utilities
- **go-systemd**: Systemd integration

### Frontend
- **SvelteKit**: Modern frontend framework
- **TypeScript**: Type-safe JavaScript
- **Tailwind CSS**: Utility-first CSS framework
- **Vite**: Build tool and development server
- **DaisyUI**: Tailwind CSS component library

## Project Structure

```
cayman/
├── go.mod                  # Go module definition
├── Taskfile.yml           # Task runner configuration
├── bin/                   # Build output
├── cmd/                   # Application entry points
│   └── cayman/           # Main application
│       └── main.go       # Application entry point
├── frontend/              # SvelteKit frontend
│   ├── src/
│   │   ├── routes/        # SvelteKit routes
│   │   └── lib/           # Shared components and types
│   ├── build/             # Production build output
│   └── package.json       # Node.js dependencies
└── internal/              # Internal Go packages
    ├── data/              # Data collection modules
    │   ├── hardware/      # CPU information
    │   └── systemd/       # Systemd integration
    └── modules/           # Application modules
        └── dashboard/     # Dashboard monitoring module
```

## Prerequisites

- **Go 1.24.5 or later**
- **Node.js 18+ and npm**
- **Task** (optional, for using Taskfile commands)
- **Linux system** with systemd (for systemd monitoring features)

## Development Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd cayman
```

### 2. Install Dependencies

**Backend Dependencies:**
```bash
go mod download
```

**Frontend Dependencies:**
```bash
cd frontend
npm install
cd ..
```

### 3. Development Mode

The project uses [Task](https://taskfile.dev/) for build automation. Install it for the best development experience:

```bash
# Install Task (if not already installed)
go install github.com/go-task/task/v3/cmd/task@latest
```

**Run in Development Mode:**
```bash
# Start both frontend and backend in development mode
task default
```

This will:
- Start the frontend development server with hot reload
- Start the backend with hot reload using Air

**Manual Development Setup:**

If you prefer not to use Task:

```bash
# Terminal 1: Start frontend development server
cd frontend
npm run dev

# Terminal 2: Start backend with hot reload
go install github.com/air-verse/air@latest
air
```

### 4. Available Commands

```bash
# Build production version
task build

# Run production build
task run

# Build only frontend
task build:frontend

# Development mode (frontend only)
task dev:frontend

# Development mode (backend only)
task dev:go
```

**Note**: During development with `task default` or `task dev:go`, the backend runs with the default settings (0.0.0.0:8080). To test different address/port configurations, build the binary and run it manually with the desired flags.

## Building for Production

### Build Single Binary

```bash
task build
```

This creates a single binary at `bin/cayman` with the frontend assets embedded.

### Run Production Build

```bash
task run
# or
./bin/cayman
```

The application will be available at `http://localhost:8080` by default.

## CLI Usage

The application supports the following command-line flags:

```bash
./bin/cayman [flags]
```

### Available Flags

- `--addr string`: Listen address (default: "0.0.0.0")
- `--port string`: Listen port (default: "8080")
- `--help`: Show help message

### Examples

```bash
# Run with default settings (0.0.0.0:8080)
./bin/cayman

# Run on a different port
./bin/cayman --port 3000

# Run on localhost only
./bin/cayman --addr 127.0.0.1

# Run on a specific address and port
./bin/cayman --addr 192.168.1.100 --port 9000

# IPv6 address
./bin/cayman --addr "::1" --port 8080

# Show help
./bin/cayman --help
```

## API Endpoints

### REST API
- `GET /api/host/current` - Get current system state
- `GET /api/stop` - Gracefully stop the server

### Server-Sent Events
- `GET /api/host/events` - Real-time system metrics stream

Events published:
- `cpu` - CPU usage percentage
- `load` - System load averages
- `mem` - Memory information

## Configuration

The application supports configuration through command-line flags:

### Network Configuration
- **Listen Address**: Configurable via `--addr` flag (default: "0.0.0.0")
- **Listen Port**: Configurable via `--port` flag (default: "8080")

### System Configuration
- **Update Interval**: 3 seconds for metrics (hardcoded)
- **SSE Replay Window**: 5 minutes (hardcoded)
- **Shutdown Timeout**: 5 seconds for graceful shutdown (hardcoded)

## Development Guidelines

### Backend Development

- Follow Go conventions and use `gofmt`
- Add new monitoring modules in `internal/modules/`
- Data collection logic goes in `internal/data/`
- Use structured logging with `slog`

### Frontend Development

- Use TypeScript for all new code
- Follow SvelteKit conventions
- Components go in `src/lib/components/`
- Types are defined in `src/lib/types.ts`
- Use Tailwind CSS for styling

### Adding New Metrics

1. **Backend**: Add data collection in `internal/data/`
2. **Module**: Extend the host module in `internal/modules/host/`
3. **Types**: Update types in `internal/modules/host/types.go`
4. **Frontend**: Add UI components and update types in `frontend/src/lib/types.ts`

## Monitoring Features

### System Metrics
- **CPU Usage**: Real-time CPU utilization percentage
- **Load Averages**: 1, 5, and 15-minute load averages
- **Memory**: Total, available, and used memory
- **Host Information**: Hostname, FQDN, OS details

### Systemd Integration
- **Active Units**: Count of currently active systemd units
- **Failed Units**: Count of failed systemd units (highlighted in red)

## Browser Support

- Modern browsers with Server-Sent Events support
- Mobile responsive design
- Progressive enhancement for older browsers

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Check what's using port 8080 (or your chosen port)
lsof -i :8080
# Kill the process or use a different port
./bin/cayman --port 3000
```

**Bind address issues:**
```bash
# If you can't bind to 0.0.0.0, try localhost only
./bin/cayman --addr 127.0.0.1

# For IPv6 systems
./bin/cayman --addr "::1"
```

**Permission denied for systemd:**
```bash
# Ensure your user has access to systemd
systemctl --user status
```

**Frontend build fails:**
```bash
# Clear node_modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**Hot reload not working:**
```bash
# Ensure Air is installed and configured
go install github.com/air-verse/air@latest
```

## Performance Notes

- The application polls system metrics every 3 seconds
- SSE connections are automatically managed and cleaned up
- Frontend uses efficient state management with Svelte stores
- Production build includes asset optimization and compression

## Security Considerations

- The application currently allows CORS from all origins (development only)
- No authentication is implemented
- **Network Binding**: Default binding to 0.0.0.0 exposes the service to all network interfaces
  - Use `--addr 127.0.0.1` to restrict to localhost only
  - Use `--addr <specific-ip>` to bind to a specific interface
- Intended for internal/local network use
- Consider adding authentication for production deployments
