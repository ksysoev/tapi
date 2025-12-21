# Example OpenAPI Specifications for Testing

This directory contains example OpenAPI specifications for testing TAPI.

## Included Examples

### 1. Pet Store API (`example-petstore.yaml`)
A classic Pet Store API demonstrating:
- Multiple HTTP methods (GET, POST, DELETE)
- Path parameters
- Query parameters
- Request bodies
- Various response codes

## Testing with Examples

### Validate the spec:
```bash
tapi validate -f example-petstore.yaml
```

### Explore interactively:
```bash
tapi explore -f example-petstore.yaml
```

### Try these endpoints:
1. GET /pet/{petId} - Requires path parameter `petId`
2. GET /pet/findByStatus - Optional query parameter `status`
3. POST /pet - Requires JSON request body
4. DELETE /pet/{petId} - Delete operation

## Remote Specs

You can also explore remote OpenAPI specifications:

```bash
# Official Swagger Petstore
tapi explore -u https://petstore3.swagger.io/api/v3/openapi.json

# GitHub API (if available)
tapi explore -u https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json
```

## Creating Your Own Specs

To test with your own API:

1. Export your OpenAPI spec (YAML or JSON)
2. Place it in this directory or anywhere on your system
3. Run: `tapi explore -f your-spec.yaml`

## Tips

- Use `?` to see keyboard shortcuts
- Use `j/k` for Vim-style navigation
- Press `e` on an endpoint to execute a request
- Press `Esc` to go back
- Press `q` to quit
