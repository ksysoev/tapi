package tui

import "github.com/ksysoev/tapi/pkg/openapi"

func (m Model) getCurrentPath() *openapi.Path {
	if m.selectedEndpoint < 0 || m.selectedEndpoint >= len(m.endpointsList) {
		return nil
	}

	idx := 0
	for i, path := range m.spec.Paths {
		for range path.Operations {
			if idx == m.selectedEndpoint {
				return &m.spec.Paths[i]
			}
			idx++
		}
	}
	return nil
}

func (m Model) getCurrentOperation() *openapi.Operation {
	if m.selectedEndpoint < 0 || m.selectedEndpoint >= len(m.endpointsList) {
		return nil
	}

	idx := 0
	for _, path := range m.spec.Paths {
		for j := range path.Operations {
			if idx == m.selectedEndpoint {
				return &path.Operations[j]
			}
			idx++
		}
	}
	return nil
}
