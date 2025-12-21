package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

type Spec struct {
	Title       string
	Version     string
	Description string
	Servers     []Server
	Paths       []Path
	raw         *openapi3.T
}

type Server struct {
	URL         string
	Description string
}

type Path struct {
	Path       string
	Operations []Operation
}

type Operation struct {
	Method      string
	Summary     string
	Description string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[string]Response
	Tags        []string
}

type Parameter struct {
	Name        string
	In          string
	Description string
	Required    bool
	Schema      *Schema
}

type RequestBody struct {
	Description string
	Required    bool
	Content     map[string]MediaType
}

type MediaType struct {
	Schema *Schema
}

type Response struct {
	Description string
	Content     map[string]MediaType
}

type Schema struct {
	Type       string
	Format     string
	Properties map[string]*Schema
	Required   []string
	Example    interface{}
}

func LoadFromFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return parseSpec(data)
}

func LoadFromURL(url string) (*Spec, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseSpec(data)
}

func parseSpec(data []byte) (*Spec, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	var doc *openapi3.T
	var err error

	if json.Valid(data) {
		doc, err = loader.LoadFromData(data)
	} else {
		doc, err = loader.LoadFromData(data)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	return convertSpec(doc), nil
}

func convertSpec(doc *openapi3.T) *Spec {
	spec := &Spec{
		Title:       doc.Info.Title,
		Version:     doc.Info.Version,
		Description: doc.Info.Description,
		raw:         doc,
	}

	for _, server := range doc.Servers {
		spec.Servers = append(spec.Servers, Server{
			URL:         server.URL,
			Description: server.Description,
		})
	}

	for path, pathItem := range doc.Paths.Map() {
		p := Path{Path: path}

		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}

			operation := Operation{
				Method:      method,
				Summary:     op.Summary,
				Description: op.Description,
				OperationID: op.OperationID,
				Tags:        op.Tags,
			}

			for _, param := range op.Parameters {
				if param.Value != nil {
					operation.Parameters = append(operation.Parameters, Parameter{
						Name:        param.Value.Name,
						In:          param.Value.In,
						Description: param.Value.Description,
						Required:    param.Value.Required,
						Schema:      convertSchema(param.Value.Schema),
					})
				}
			}

			if op.RequestBody != nil && op.RequestBody.Value != nil {
				rb := &RequestBody{
					Description: op.RequestBody.Value.Description,
					Required:    op.RequestBody.Value.Required,
					Content:     make(map[string]MediaType),
				}
				for contentType, mediaType := range op.RequestBody.Value.Content {
					rb.Content[contentType] = MediaType{
						Schema: convertSchema(mediaType.Schema),
					}
				}
				operation.RequestBody = rb
			}

			operation.Responses = make(map[string]Response)
			if op.Responses != nil {
				for status, resp := range op.Responses.Map() {
					if resp.Value != nil {
						r := Response{
							Description: *resp.Value.Description,
							Content:     make(map[string]MediaType),
						}
						for contentType, mediaType := range resp.Value.Content {
							r.Content[contentType] = MediaType{
								Schema: convertSchema(mediaType.Schema),
							}
						}
						operation.Responses[status] = r
					}
				}
			}

			p.Operations = append(p.Operations, operation)
		}

		if len(p.Operations) > 0 {
			spec.Paths = append(spec.Paths, p)
		}
	}

	return spec
}

func convertSchema(schemaRef *openapi3.SchemaRef) *Schema {
	if schemaRef == nil || schemaRef.Value == nil {
		return nil
	}

	s := schemaRef.Value
	schemaType := ""
	if s.Type != nil && len(*s.Type) > 0 {
		schemaType = (*s.Type)[0]
	}
	
	schema := &Schema{
		Type:       schemaType,
		Format:     s.Format,
		Properties: make(map[string]*Schema),
		Required:   s.Required,
		Example:    s.Example,
	}

	for name, propRef := range s.Properties {
		schema.Properties[name] = convertSchema(propRef)
	}

	return schema
}
