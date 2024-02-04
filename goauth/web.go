package goauth

import (
	"strings"
	"sync"

	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	ScopeProtected string = "PROTECTED"
	ScopePublic    string = "PUBLIC"
)

var (
	loadResourcePathOnce sync.Once
	loadedResources      = []AddResourceReq{}
	loadedPaths          = []CreatePathReq{}
)

type CreatePathReq struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	Group   string `json:"group"`
	Desc    string `json:"desc"`
	ResCode string `json:"resCode"`
	Method  string `json:"method"`
}

type AddResourceReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type QueryResourcePathReq struct {
	Resources []AddResourceReq
	Paths     []CreatePathReq
}

// Create endpoint to expose resources and endpoint paths to be collected by goauth.
func ReportOnBoostrapped(resources []AddResourceReq) {

	miso.PreServerBootstrap(func(rail miso.Rail) error {

		// resources and paths are polled by goauth
		miso.Get("/auth/resource", func(c *gin.Context, rail miso.Rail) (any, error) {

			// resources and paths are lazily loaded
			loadResourcePathOnce.Do(func() {
				app := miso.GetPropStr(miso.PropAppName)
				for _, res := range resources {
					if res.Code == "" || res.Name == "" {
						continue
					}
					loadedResources = append(loadedResources, res)
				}

				routes := miso.GetHttpRoutes()
				for _, route := range routes {
					if route.Url == "" {
						continue
					}
					var routeType = ScopeProtected
					if route.Scope == miso.ScopePublic {
						routeType = ScopePublic
					}

					url := route.Url
					if !strings.HasPrefix(url, "/") {
						url = "/" + url
					}

					r := CreatePathReq{
						Method:  route.Method,
						Group:   app,
						Url:     "/" + app + url,
						Type:    routeType,
						Desc:    route.Desc,
						ResCode: route.Resource,
					}
					loadedPaths = append(loadedPaths, r)
				}
			})

			return QueryResourcePathReq{
				Resources: loadedResources,
				Paths:     loadedPaths,
			}, nil
		}).Build()
		return nil
	})
}
