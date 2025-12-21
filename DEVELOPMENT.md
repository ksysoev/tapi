# TAPI Development Guide

## Project Structure

```
tapi/
├── cmd/
│   └── tapi/           # Main application entry point
│       └── main.go
├── pkg/
│   ├── cmd/            # Cobra CLI commands
│   │   ├── root.go
│   │   ├── explore.go
│   │   └── validate.go
│   ├── openapi/        # OpenAPI spec parsing
│   │   └── loader.go
│   ├── tui/            # Bubbletea TUI components
│   │   └── model.go
│   └── request/        # HTTP client for API requests
│       └── client.go
├── internal/
│   └── styles/         # Lipgloss styling
│       └── styles.go
└── example-petstore.yaml  # Sample OpenAPI spec
```

## Key Features Implemented

### 1. CLI Framework (Cobra)
- `explore` command: Launch interactive TUI
- `validate` command: Validate OpenAPI specs
- Support for both local files and remote URLs

### 2. OpenAPI Support
- Parse OpenAPI 3.x specifications (JSON/YAML)
- Load from local files or remote URLs
- Extract endpoints, operations, parameters, request/response schemas
- Validation using kin-openapi library

### 3. TUI (Bubbletea)
Beautiful terminal interface with multiple views:
- **Endpoints List**: Browse all available API endpoints
- **Operation Details**: View detailed info about selected endpoint
- **Request Builder**: Fill parameters and build requests
- **Response Viewer**: Display API responses with syntax

### 4. Vim Keybindings
Full Vim-style navigation:
- `j/k`: Move down/up
- `h/l`: Navigate between panels
- `g/G`: Jump to top/bottom
- `d/u`: Half-page scroll
- `/`: Search (future enhancement)

### 5. API Testing
- Interactive request builder
- Fill required parameters
- Add request body (JSON)
- Send HTTP requests
- Display formatted responses

## Architecture

### Model-Update-View Pattern (Bubbletea)

The TUI follows Elm Architecture:

```go
type Model struct {
    spec              *openapi.Spec
    currentView       view
    selectedEndpoint  int
    viewport          viewport.Model
    inputs            []textinput.Model
}

func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

### Views

1. **viewEndpoints**: List all API endpoints
2. **viewOperationDetails**: Show operation details with scrolling
3. **viewRequestBuilder**: Build and send requests
4. **viewResponse**: Display API responses

### Styling

Centralized styling using Lipgloss:
- Color-coded HTTP methods (GET=green, POST=blue, DELETE=red)
- Focused/unfocused input states
- Panel borders and highlighting
- Status indicators

## Usage Examples

### Explore Local Spec
```bash
tapi explore -f openapi.yaml
```

### Explore Remote Spec
```bash
tapi explore -u https://petstore3.swagger.io/api/v3/openapi.json
```

### Validate Spec
```bash
tapi validate -f openapi.yaml
```

## Testing

Run tests:
```bash
make test
```

## Future Enhancements

- [ ] Search/filter endpoints
- [ ] Request history
- [ ] Save/load request collections
- [ ] Authentication support (API keys, OAuth, Bearer tokens)
- [ ] Environment variables
- [ ] Export requests to cURL/Postman
- [ ] Response syntax highlighting (JSON/XML)
- [ ] Multiple server selection
- [ ] Request/response body editor with syntax highlighting
- [ ] WebSocket support
- [ ] GraphQL support

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Add tests
5. Run `make fmt lint test`
6. Submit a pull request
