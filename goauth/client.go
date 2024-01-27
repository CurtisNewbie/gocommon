package goauth

import (
	"errors"
	"strings"

	"github.com/curtisnewbie/miso/miso"
	"github.com/sirupsen/logrus"
)

const (
	// Property Key for enabling GoAuth Client, by default it's true
	//
	// goauth-client-go doesn't use it internally, it's only useful for the Callers
	PropGoAuthClientEnabled = "goauth.client.enabled"

	// event bus name for adding paths
	addPathEventBus = "goauth.add-path"

	// event bus name for adding resources
	addResourceEventBus = "goauth.add-resource"

	ScopeProtected string = "PROTECTED"
	ScopePublic    string = "PUBLIC"
)

func init() {
	miso.SetDefProp(PropGoAuthClientEnabled, true)
}

type PathDoc struct {
	Desc string
	Type string
	Code string
}

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" `
}

type RoleInfoResp struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type CreatePathReq struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	Group   string `json:"group"`
	Desc    string `json:"desc"`
	ResCode string `json:"resCode"`
	Method  string `json:"method"`
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
}

type TestResAccessResp struct {
	Valid bool `json:"valid"`
}

type AddResourceReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Test whether this role has access to the url
func TestResourceAccess(rail miso.Rail, req TestResAccessReq) (*TestResAccessResp, error) {
	tr := miso.NewDynTClient(rail, "/remote/path/resource/access-test", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return nil, tr.Err
	}

	if err := tr.Require2xx(); err != nil {
		return nil, err
	}

	r, e := miso.ReadGnResp[*TestResAccessResp](tr)
	if e != nil {
		return nil, e
	}

	if r.Error {
		return nil, r.Err()
	}

	if r.Data == nil {
		return nil, errors.New("data is nil, unable to retrieve TestResAccessResp")
	}

	return r.Data, nil
}

// Create resource
func AddResource(rail miso.Rail, req AddResourceReq) error {
	tr := miso.NewDynTClient(rail, "/remote/resource/add", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return tr.Err
	}

	if err := tr.Require2xx(); err != nil {
		return err
	}

	r, e := miso.ReadGnResp[any](tr)
	if e != nil {
		return e
	}

	if r.Error {
		return r.Err()
	}

	logrus.Infof("Reported resource, Name: %s, Code: %s", req.Name, req.Code)
	return nil
}

// Report path
func AddPath(rail miso.Rail, req CreatePathReq) error {
	tr := miso.NewDynTClient(rail, "/remote/path/add", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return tr.Err
	}

	if err := tr.Require2xx(); err != nil {
		return err
	}

	r, e := miso.ReadGnResp[any](tr)
	if e != nil {
		return e
	}

	if r.Error {
		return r.Err()
	}

	return nil
}

// Retrieve role information
func GetRoleInfo(rail miso.Rail, req RoleInfoReq) (*RoleInfoResp, error) {
	tr := miso.NewDynTClient(rail, "/remote/role/info", "goauth").
		EnableTracing().
		PostJson(req)

	if tr.Err != nil {
		return nil, tr.Err
	}

	if err := tr.Require2xx(); err != nil {
		return nil, err
	}

	r, e := miso.ReadGnResp[*RoleInfoResp](tr)
	if e != nil {
		return nil, e
	}

	if r.Error {
		return nil, r.Err()
	}

	if r.Data == nil {
		return nil, errors.New("data is nil, unable to retrieve RoleInfoResp")
	}

	return r.Data, nil
}

// Check whether goauth client is enabled
//
//	"goauth.miso.enabled"
func IsEnabled() bool {
	return miso.GetPropBool(PropGoAuthClientEnabled)
}

// Report path asynchronously
func AddPathAsync(rail miso.Rail, req CreatePathReq) error {
	return miso.PubEventBus(rail, req, addPathEventBus)
}

// Report resource asynchronously
func AddResourceAsync(rail miso.Rail, req AddResourceReq) error {
	return miso.PubEventBus(rail, req, addResourceEventBus)
}

// Register a hook to report paths and resources to GoAuth on server bootstrapped
//
// This method checks if the goauth client is enabled, nothing will happen if the client is disabled.
func ReportOnBoostrapped(rail miso.Rail, res []AddResourceReq) {
	if !IsEnabled() {
		rail.Debug("GoAuth client disabled, will not report resources")
		return
	}

	miso.NewEventBus(addResourceEventBus)
	miso.NewEventBus(addPathEventBus)

	miso.PostServerBootstrapped(func(rail miso.Rail) error {

		app := miso.GetPropStr(miso.PropAppName)
		for _, res := range res {
			if res.Code == "" || res.Name == "" {
				continue
			}

			// report resource synchronously
			if e := AddResource(rail, AddResourceReq(res)); e != nil {
				return e
			}
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

			// report the path asynchronously
			if err := AddPathAsync(rail, r); err != nil {
				return err
			}
		}
		return nil
	})
}
