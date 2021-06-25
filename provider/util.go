package provider

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"oaf-server/codegen"
	"strconv"
	"strings"
)

const (
	GEOJSONContentType = "application/geo+json"
	JSONContentType    = "application/json"
	LDJSONContentType  = "application/ld+json"
	HTMLContentType    = "text/html"

	CapabilitesProvider = "CapabilitesProvider" // landingpage, collections, conformance | json, jsonld, html
	OASProvider         = "OASProvider"         // OAS spec | json, html
	DataProvider        = "DataProvider"        // GetFeature, GetFeatures | json, jsonld, html

)

type InvalidContentTypeError struct {
	content_type string
	err          string
}

func (e *InvalidContentTypeError) Error() string {
	return fmt.Sprintf("invalid request Content-Type %v", e.content_type)
}

type InvalidFormatError struct {
	format string
	err    string
}

func (e *InvalidFormatError) Error() string {
	return fmt.Sprintf("invalid request query format %v", e.format)
}

func GetContentTypeMap(provider string) map[string]string {
	var m = map[string]map[string]string{
		CapabilitesProvider: {
			"json":   JSONContentType,
			"jsonld": LDJSONContentType,
			"html":   HTMLContentType,
		},
		OASProvider: {
			"json": JSONContentType,
			"html": HTMLContentType,
		},
		DataProvider: {
			"json":   GEOJSONContentType,
			"jsonld": LDJSONContentType,
			"html":   HTMLContentType,
		},
	}
	return m[provider]
}

func GetContentFieldMap(provider string) map[string]string {
	var m = map[string]map[string]string{
		CapabilitesProvider: {
			JSONContentType:   "json",
			LDJSONContentType: "jsonld",
			HTMLContentType:   "html",
		},
		OASProvider: {
			JSONContentType: "json",
			HTMLContentType: "html",
		},
		DataProvider: {
			GEOJSONContentType: "json",
			LDJSONContentType:  "jsonld",
			HTMLContentType:    "html",
		},
	}
	return m[provider]

}

func GetFeatureContentFields() map[string]string {
	ct := make(map[string]string)

	ct[GEOJSONContentType] = "json"
	ct[LDJSONContentType] = "jsonld"
	ct[HTMLContentType] = "html"

	return ct
}

func GetContentType(r *http.Request, providerType string) (string, error) {

	ctMap := GetContentTypeMap(providerType)
	cfMap := GetContentFieldMap(providerType)

	// check query string
	queryFormat, ok := r.URL.Query()["f"]
	if ok && len(queryFormat) > 0 {
		resContentType, ok := ctMap[queryFormat[0]]
		if ok {
			return resContentType, nil
		} else {
			return "", &InvalidFormatError{format: queryFormat[0]}
		}
	}

	// otherwise use content-type set in request header
	reqContentType := r.Header.Get("Content-Type")
	r.Header.Set("Content-Type", reqContentType)

	// validate request Content-Type if set
	if reqContentType != "" {
		_, ok = cfMap[reqContentType]
		if !ok {
			return "", &InvalidContentTypeError{content_type: reqContentType}
		}
	}

	// if no Content-Type header set, default to text/html
	if reqContentType == "" {
		reqContentType = HTMLContentType
	}

	return reqContentType, nil
}

func GetRelationMap() map[string]string {
	ct := make(map[string]string)

	ct["alternate"] = "Alternative"
	ct["self"] = "This"

	return ct
}

func GetContentTypes() map[string]string {
	ct := make(map[string]string)

	ct["json"] = JSONContentType
	ct["jsonld"] = LDJSONContentType
	ct["html"] = HTMLContentType

	return ct
}

func ProcesLinksForParams(links []codegen.Link, queryParams url.Values) error {
	for l := range links {
		path, err := url.Parse(links[l].Href)
		if err != nil {
			return err
		}
		values := path.Query()

		for k, v := range queryParams {
			if k == "f" {
				continue
			}
			values.Add(k, v[0])
		}
		path.RawQuery = values.Encode()
		links[l].Href = path.String()
	}

	return nil

}

func CreateFeatureLinks(title, hrefPath, rel, ct string) ([]codegen.Link, error) {

	links := make([]codegen.Link, 0)

	href, err := ctLink(hrefPath, GetFeatureContentFields()[ct])
	if err != nil {
		return links, err
	}
	links = append(links, codegen.Link{Title: formatTitle(title, rel, GetFeatureContentFields()[ct]), Rel: rel, Href: href, Type: ct})

	if rel == "self" {
		rel = "alternate"
	}

	for k, sct := range GetContentTypeMap(DataProvider) {
		if ct == sct {
			continue
		}
		href, err := ctLink(hrefPath, k)
		if err != nil {
			return links, err
		}

		links = append(links, codegen.Link{Title: formatTitle(title, rel, k), Rel: rel, Href: href, Type: sct})

		if rel == "self" {
			rel = "alternate"
			links = append(links, codegen.Link{Title: formatTitle(title, rel, k), Rel: rel, Href: href, Type: sct})
		}
	}

	return links, nil
}

func GetApiLinks(hrefPath string) ([]codegen.Link, error) {

	links := make([]codegen.Link, 0)
	links = append(links, codegen.Link{Title: "Documentation of the API", Rel: "service-doc", Href: hrefPath + "?f=html", Type: "text/html"})
	links = append(links, codegen.Link{Title: "Definition of the API in OpenAPI 3.0", Rel: "service-desc", Href: hrefPath + "?f=json", Type: "application/vnd.oai.openapi+json;version=3.0"})

	return links, nil
}

func CreateLinks(providerName, providerType, hrefPath, rel, ct string) ([]codegen.Link, error) {

	links := make([]codegen.Link, 0)

	href, err := ctLink(hrefPath, GetContentFieldMap(providerType)[ct])
	if err != nil {
		return links, err
	}
	links = append(links, codegen.Link{Title: formatTitle(providerName, rel, GetContentFieldMap(providerType)[ct]), Rel: rel, Href: href, Type: ct})

	if rel == "self" {
		rel = "alternate"
	}

	//if rel != "self" {
	//	return links, nil
	//}

	for k, sct := range GetContentTypeMap(providerType) {
		if ct == sct {
			continue
		}
		href, err := ctLink(hrefPath, k)
		if err != nil {
			return links, err
		}
		links = append(links, codegen.Link{Title: formatTitle(providerName, rel, k), Rel: rel, Href: href, Type: sct})
		if rel == "self" {
			rel = "alternate"
			links = append(links, codegen.Link{Title: formatTitle(providerName, rel, k), Rel: rel, Href: href, Type: sct})
		}
	}

	return links, nil
}

func formatTitle(title, rel, format string) string {
	relation := rel
	if rel == "self" {
		relation = "this"
	}
	return strings.ToLower(fmt.Sprintf("%s %s in %s format", relation, title, format))
}

func ctLink(baselink, contentType string) (string, error) {

	u, err := url.Parse(baselink)
	if err != nil {
		log.Printf("Invalid link '%v', will return empty string.", baselink)
		return "", err
	}
	q := u.Query()

	var l string
	switch contentType {
	default:
		q["f"] = []string{contentType}
	}

	u.RawQuery = q.Encode()
	l = u.String()
	return l, nil
}

// copied,tweaked from https://github.com/go-spatial/jivan
func ConvertFeatureID(v interface{}) (interface{}, error) {
	switch aval := v.(type) {
	case float64:
		return uint64(aval), nil
	case int64:
		return uint64(aval), nil
	case uint64:
		return aval, nil
	case uint:
		return uint64(aval), nil
	case int8:
		return uint64(aval), nil
	case uint8:
		return uint64(aval), nil
	case uint16:
		return uint64(aval), nil
	case int32:
		return uint64(aval), nil
	case uint32:
		return uint64(aval), nil
	case []byte:
		return string(aval), nil
	case string:
		return aval, nil

	default:
		return 0, fmt.Errorf("cannot convert ID : %v", aval)
	}
}

func ParseLimit(limit string, defaultReturnLimit, maxReturnLimit uint64) uint64 {
	limitParam := defaultReturnLimit
	if limit != "" {
		newValue, err := strconv.ParseInt(limit, 10, 64)
		if err == nil && uint64(newValue) < maxReturnLimit {
			limitParam = uint64(newValue)
		} else {
			limitParam = maxReturnLimit
		}
	}
	return limitParam
}

func ParseUint(stringValue string, defaultValue uint64) uint64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseFloat64(stringValue string, defaultValue float64) float64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseBBox(stringValue string, defaultValue [4]float64) [4]float64 {
	if stringValue == "" {
		return defaultValue
	}
	bboxValues := strings.Split(stringValue, ",")
	if len(bboxValues) != 4 {
		return defaultValue
	}

	var value [4]float64
	for i, v := range bboxValues {
		value[i] = ParseFloat64(v, value[i])
	}

	return value
}

func UpperFirst(title string) string {
	return strings.Title(title)
}
