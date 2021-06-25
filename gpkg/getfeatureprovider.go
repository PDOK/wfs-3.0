package gpkg

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

type GetFeatureProvider struct {
	data        *Feature
	srsid       string
	contenttype string
}

func (gp *GeoPackageProvider) NewGetFeatureProvider(r *http.Request) (codegen.Provider, error) {

	collectionId, featureId, _ := codegen.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := gp.GeoPackage.DefaultBBox

	p := &GetFeatureProvider{srsid: fmt.Sprintf("EPSG:%d", gp.GeoPackage.Srid)}

	path := r.URL.Path
	ct := r.Header.Get("Content-Type")

	if ct == provider.JSONContentType {
		ct = provider.GEOJSONContentType
	}

	p.contenttype = ct

	for _, cn := range gp.GeoPackage.Collections {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := gp.GeoPackage.GetFeatures(r.Context(), gp.GeoPackage.DB, cn, collectionId, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) == 1 {

			feature := fcGeoJSON.Features[0]

			hrefBase := fmt.Sprintf("%s%s", gp.CommonProvider.ServiceEndpoint, path) // /collections
			links, _ := provider.CreateFeatureLinks("feature", hrefBase, "self", ct)
			feature.Links = links

			p.data = feature
		}

		break
	}

	return p, nil
}

func (gfp *GetFeatureProvider) Provide() (interface{}, error) {
	return gfp.data, nil
}

func (gfp *GetFeatureProvider) ContentType() string {
	return gfp.contenttype
}

func (gfp *GetFeatureProvider) String() string {
	return "getfeature"
}

func (gfp *GetFeatureProvider) SrsId() string {
	return gfp.srsid
}
