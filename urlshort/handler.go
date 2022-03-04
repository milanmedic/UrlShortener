package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type YAMLUrlMapping struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path
		if dest, ok := pathsToUrls[reqPath]; ok {
			http.Redirect(rw, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(rw, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		yamlMappings, err := ImportYAMLMappingIntoMap(yml)

		if err != nil {
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}

		reqPath := r.URL.Path
		if dest, ok := yamlMappings[reqPath]; ok {
			http.Redirect(w, r, dest.Url, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func ImportYAMLMappingIntoMap(contents []byte) (map[string]YAMLUrlMapping, error) {
	yamlMappings := make(map[string]YAMLUrlMapping)
	err := yaml.Unmarshal(contents, &yamlMappings)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return yamlMappings, nil
}
