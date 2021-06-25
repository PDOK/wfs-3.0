package provider

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
)

type GetLandingPageProvider struct {
	Links       []codegen.Link `json:"links,omitempty"`
	contenttype string
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	*Service
}

func NewGetLandingPageProvider(serviceConfig Service) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {
		p := &GetLandingPageProvider{}
		reqContentType, err := GetContentType(r, p.ProviderType())

		if err != nil {
			return nil, err
		}

		p.contenttype = reqContentType
		p.Service = &serviceConfig
		links, _ := CreateLinks("landing page", CapabilitesProvider, serviceConfig.Url, "self", reqContentType)
		apiLink, _ := CreateLinks("openapi3 specification", OASProvider, fmt.Sprintf("%s/api", serviceConfig.Url), "service", reqContentType)                   // /api, "service", ct)
		conformanceLink, _ := CreateLinks("capabilities", CapabilitesProvider, fmt.Sprintf("%s/conformance", serviceConfig.Url), "conformance", reqContentType) // /conformance, "conformance", ct)
		dataLink, _ := CreateLinks("collections", CapabilitesProvider, fmt.Sprintf("%s/collections", serviceConfig.Url), "data", reqContentType)                // /collections, "collections", ct)
		p.Links = append(p.Links, links...)
		p.Links = append(p.Links, apiLink...)
		p.Links = append(p.Links, conformanceLink...)
		p.Links = append(p.Links, dataLink...)

		if p.contenttype == "application/json" {
			p.Title = p.Service.Name
			p.Description = p.Service.Description
			p.Service = nil
		} else if p.contenttype == "application/ld+json" {
			p.Links = nil
		}
		return p, nil
	}
}

func (glp *GetLandingPageProvider) Provide() (interface{}, error) {
	return glp, nil
}

func (glp *GetLandingPageProvider) ContentType() string {
	return glp.contenttype
}

func (glp *GetLandingPageProvider) String() string {
	return "landingpage"
}

func (glp *GetLandingPageProvider) SrsId() string {
	return "n.a"
}

func (glp *GetLandingPageProvider) ProviderType() string {
	return CapabilitesProvider
}
