package server

// this code is generated by go generate
// DO NOT EDIT BY HAND!
import (
	"net/http"
	"regexp"
	"sort"
	pc "wfs3_server/provider_common"
)

type routes []*route

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	Routes routes
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.Routes = append(h.Routes, &route{pattern, handler})
	sort.Sort(sort.Reverse(h.Routes))
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.Routes = append(h.Routes, &route{pattern, http.HandlerFunc(handler)})
	sort.Sort(sort.Reverse(h.Routes))
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.Routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func (s routes) Len() int {
	return len(s)
}
func (s routes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s routes) Less(i, j int) bool {
	return len(s[i].pattern.String()) < len(s[j].pattern.String())
}

func (s *Server) Router() *RegexpHandler {
	router := &RegexpHandler{}
	router.HandleFunc(regexp.MustCompile("/api"), s.HandleForProvider(pc.NewGetApiProvider(s.ServiceSpecPath)))
	// path: /
	router.HandleFunc(regexp.MustCompile("/"), s.HandleForProvider(s.Providers.NewGetLandingPageProvider))
	// path: /collections
	router.HandleFunc(regexp.MustCompile("/collections"), s.HandleForProvider(s.Providers.NewDescribeCollectionsProvider))
	// path: /collections/{collectionId}
	router.HandleFunc(regexp.MustCompile("/collections/.*"), s.HandleForProvider(s.Providers.NewDescribeCollectionProvider))
	// path: /collections/{collectionId}/items
	router.HandleFunc(regexp.MustCompile("/collections/.*/items"), s.HandleForProvider(s.Providers.NewGetFeaturesProvider))
	// path: /collections/{collectionId}/items/{featureId}
	router.HandleFunc(regexp.MustCompile("/collections/.*/items/.*"), s.HandleForProvider(s.Providers.NewGetFeatureProvider))
	// path: /conformance
	router.HandleFunc(regexp.MustCompile("/conformance"), s.HandleForProvider(s.Providers.NewGetRequirementsClassesProvider))
	return router
}
