# TAPI - Terminal API Explorer

A beautiful terminal-based OpenAPI specification explorer and API testing tool built with Go and Bubbletea.

![TAPI Demo](https://img.shields.io/badge/TUI-Bubbletea-purple) ![Go Version](https://img.shields.io/badge/Go-1.24+-blue) ![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 🎨 **Beautiful TUI** - Built with Bubbletea and Lipgloss for an elegant terminal interface
- 📖 **OpenAPI Support** - Load and explore local and remote OpenAPI 3.x specifications
- 🚀 **API Testing** - Make API requests directly from the TUI (Swagger-like experience)
- ⌨️ **Vim Keybindings** - Navigate efficiently with j/k/h/l and other Vim shortcuts
- 🔍 **Browse Endpoints** - Quickly find and explore API operations
- 📝 **Request Builder** - Fill required fields and customize requests interactively
- 🎯 **Multiple Views** - Endpoints list, operation details, request builder, and response viewer
- ✅ **Validation** - Validate OpenAPI specifications for correctness

## Installation

### From Source

```bash
git clone https://github.com/ksysoev/tapi
cd tapi
make install
```

Or build manually:

```bash
go build -o tapi ./cmd/tapi
```

## Quick Start

### 1. Validate an OpenAPI spec:

```bash
tapi validate -f ./example-petstore.yaml
```

### 2. Explore the example Pet Store API:

```bash
tapi explore -f ./example-petstore.yaml
```

### 3. Explore a remote OpenAPI spec:

```bash
tapi explore -u https://petstore3.swagger.io/api/v3/openapi.json
```

## Usage

### Commands

- `tapi explore -f <file>` - Explore a local OpenAPI specification
- `tapi explore -u <url>` - Explore a remote OpenAPI specification  
- `tapi validate -f <file>` - Validate an OpenAPI specification
- `tapi --help` - Show help information

### TUI Navigation

#### Endpoints List View
- **j/k or ↓/↑** - Navigate through endpoints
- **g/G** - Jump to top/bottom
- **Enter or l** - View endpoint details
- **?** - Toggle help
- **q** - Quit

#### Operation Details View
- **j/k** - Scroll up/down
- **d/u** - Half-page scroll down/up
- **e or Enter** - Open request builder
- **h** - Go back to endpoints list
- **Esc** - Return to main view

#### Request Builder View
- **Tab or j** - Next input field
- **Shift+Tab or k** - Previous input field
- **Ctrl+S or Alt+Enter** - Send request
- **h** - Go back
- **Esc** - Cancel

#### Response View
- **j/k** - Scroll through response
- **d/u** - Half-page scroll
- **h** - Go back to request builder
- **Esc** - Return to endpoints

### Example Workflow

1. Start TAPI with your OpenAPI spec
2. Browse the list of endpoints using `j/k`
3. Press `Enter` to view endpoint details
4. Press `e` to build a request
5. Fill in required parameters using `Tab` to navigate
6. Press `Ctrl+S` to send the request
7. View the response with formatted output

## Features in Detail

### OpenAPI Support
- OpenAPI 3.x (JSON and YAML formats)
- Local file loading
- Remote URL fetching
- Comprehensive validation
- Support for:
  - Path parameters
  - Query parameters
  - Request bodies
  - Multiple response codes
  - Multiple content types

### Styling
- Color-coded HTTP methods:
  - **GET** - Green
  - **POST** - Blue  
  - **PUT** - Yellow
  - **PATCH** - Orange
  - **DELETE** - Red
- Focused/unfocused input states
- Beautiful borders and panels
- Clear visual hierarchy

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for architecture and contribution guidelines.

```bash
# Run tests
make test

# Format code
make fmt

# Lint
make lint

# Build
make build

# Clean
make clean
```

## Examples

See [EXAMPLES.md](EXAMPLES.md) for example OpenAPI specifications and usage patterns.

## Project Structure

```
tapi/
├── cmd/tapi/              # Main application entry
├── pkg/
│   ├── cmd/              # CLI commands (Cobra)
│   ├── openapi/          # OpenAPI parsing
│   ├── tui/              # TUI components (Bubbletea)
│   └── request/          # HTTP client
├── internal/styles/      # UI styling (Lipgloss)
└── example-petstore.yaml # Sample OpenAPI spec
```

## Requirements

- Go 1.24+
- Terminal with ANSI color support

## Roadmap

- [ ] Search/filter endpoints
- [ ] Request history
- [ ] Authentication support (Bearer, API keys, OAuth)
- [ ] Environment variables
- [ ] Save/load request collections
- [ ] Export to cURL/Postman
- [ ] JSON/XML syntax highlighting
- [ ] WebSocket support

## License

MIT

## Acknowledgments

Built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [kin-openapi](https://github.com/getkin/kin-openapi) - OpenAPI parsing

---

Made with ❤️ for API developers who love the terminal
