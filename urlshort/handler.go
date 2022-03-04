package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type YAMLUrlPath struct {
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

func ParseYAMLPaths(contents []byte) ([]YAMLUrlPath, error) {
	pathsToURLs := make([]YAMLUrlPath, 0)
	err := yaml.Unmarshal(contents, &pathsToURLs)
	return pathsToURLs, err
}

func ImportYAMLMappingIntoMap(contents []byte) (map[string]YAMLUrlPath, error) {
	urlPaths, err := ParseYAMLPaths(contents)
	if err != nil {
		return nil, err
	}

	yamlMappings := make(map[string]YAMLUrlPath)

	for _, value := range urlPaths {
		yamlMappings[value.Path] = value
	}

	return yamlMappings, nil
}
