package codegen

// this code is generated by go generate
// DO NOT EDIT BY HAND!

import (
	"net/http"
	"strings"
)

const (
	JSONContentType = "application/json"
	HTMLContentType = "text/html"
)

// These are the MIME types that the handlers support.
var SupportedContentTypes []string = []string{JSONContentType, HTMLContentType}

type Provider interface {
	String() string
	Provide() (interface{}, error)
}

type Providers interface {
	Init() error

	/*
	   The landing page provides links to the API definition, the conformance
	   statements and to the feature collections in this dataset.
	*/
	NewGetLandingPageProvider(r *http.Request) (Provider, error)

	/*

	 */
	NewGetCollectionsProvider(r *http.Request) (Provider, error)

	/*

	 */
	NewDescribeCollectionProvider(r *http.Request) (Provider, error)

	/*
	   Fetch features of the feature collection with id `collectionId`.

	   Every feature in a dataset belongs to a collection. A dataset may
	   consist of multiple feature collections. A feature collection is often a
	   collection of features of a similar type, based on a common schema.

	   Use content negotiation to request HTML or GeoJSON.
	*/
	NewGetFeaturesProvider(r *http.Request) (Provider, error)

	/*
	   Fetch the feature with id `featureId` in the feature collection
	   with id `collectionId`.

	   Use content negotiation to request HTML or GeoJSON.
	*/
	NewGetFeatureProvider(r *http.Request) (Provider, error)

	/*
	   A list of all conformance classes specified in a standard that the
	   server conforms to.
	*/
	NewGetConformanceDeclarationProvider(r *http.Request) (Provider, error)
}

// generate convenient functions to retrieve path params

// GetLandingPage
// no parameters present

// GetCollections
// no parameters present

// DescribeCollection
func ParametersForDescribeCollection(r *http.Request) (collectionId string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	collectionId = pathSplit[2]
	return
}

// GetFeatures
func ParametersForGetFeatures(r *http.Request) (collectionId string, limit string, bbox string, datetime string, offset string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	collectionId = pathSplit[2]
	limitArray, ok := r.URL.Query()["limit"]
	if ok {
		limit = limitArray[0]
	}

	bboxArray, ok := r.URL.Query()["bbox"]
	if ok {
		bbox = bboxArray[0]
	}

	datetimeArray, ok := r.URL.Query()["datetime"]
	if ok {
		datetime = datetimeArray[0]
	}

	offsetArray, ok := r.URL.Query()["offset"]
	if ok {
		offset = offsetArray[0]
	}

	return
}

// GetFeature
func ParametersForGetFeature(r *http.Request) (collectionId string, featureId string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	collectionId = pathSplit[2]
	featureId = pathSplit[4]
	return
}

// GetConformanceDeclaration
// no parameters present

