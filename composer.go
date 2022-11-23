package main

import (
	"log"
)

type schema struct {
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Version          string                 `json:"version,omitempty"`
	Type             string                 `json:"type,omitempty"`
	Keywords         []string               `json:"keywords,omitempty"`
	Homepage         string                 `json:"homepage,omitempty"`
	Time             string                 `json:"time,omitempty"`
	License          string                 `json:"license,omitempty"`
	Authors          map[string]interface{} `json:"authors,omitempty"`
	Support          map[string]string      `json:"support,omitempty"`
	Require          map[string]string      `json:"require,omitempty"`
	RequireDev       map[string]string      `json:"require-dev,omitempty"`
	Conflict         map[string]string      `json:"conflict,omitempty"`
	Replace          map[string]string      `json:"replace,omitempty"`
	Provide          map[string]string      `json:"provide,omitempty"`
	Suggest          map[string]string      `json:"suggest,omitempty"`
	Autoload         map[string]interface{} `json:"autoload,omitempty"`
	AutoloadDev      map[string]interface{} `json:"autoload-dev,omitempty"`
	TargetDir        string                 `json:"target-dir,omitempty"`
	MinimumStability string                 `json:"minimum-stability,omitempty"`
	Repositories     []struct {
		Type    string `json:"type,omitempty"`
		Url     string `json:"url,omitempty"`
		Package *struct {
			Name    string `json:"name,omitempty"`
			Version string `json:"version,omitempty"`
			Dist    struct {
				Url  string `json:"url,omitempty"`
				Type string `json:"type,omitempty"`
			} `json:"dist,omitempty"`
			Source struct {
				Url       string `json:"url,omitempty"`
				Type      string `json:"type,omitempty"`
				Reference string `json:"reference,omitempty"`
			} `json:"source,omitempty"`
		} `json:"package,omitempty"`
		PackagistOrg bool `json:"packagist.org,omitempty"`
	} `json:"repositories,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
	Archive *struct {
		Exclude []string `json:"exclude,omitempty"`
	} `json:"archive,omitempty"`
	PreferStable bool                   `json:"prefer-stable,omitempty"`
	Scripts      map[string]interface{} `json:"scripts,omitempty"`
	Extra        interface{}            `json:"extra,omitempty"`
	Bin          []string               `json:"bin,omitempty"`
}

func (s *schema) setPhpVersion(version string) {
	s.Require["php"] = version
	if s.Config != nil {
		platform, isMap := s.Config["platform"].(map[string]interface{})
		if isMap {
			platform["php"] = version
			s.Config["platform"] = platform
		}
	}
}

func (s *schema) setRequireDepVersion(dep, version string) {
	_, ok := s.Require[dep]
	if !ok {
		log.Printf("Dependency %s not found in require\n", dep)
		return
	}

	s.Require[dep] = version
}

func (s *schema) setRequireDevDepVersion(dep, version string) {
	_, ok := s.RequireDev[dep]
	if !ok {
		log.Printf("Dependency %s not found in require-dev\n", dep)
		return
	}

	s.RequireDev[dep] = version
}
