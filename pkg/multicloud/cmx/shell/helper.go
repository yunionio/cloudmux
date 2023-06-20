// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/printutils"
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/cloudprovider/generic"
)

var (
	CommandTable    = &shellutils.CommandTable
	R               = shellutils.R
	PrintList       = printutils.PrintInterfaceList
	PrintObject     = printutils.PrintInterfaceObject
	PrintGetterList = printutils.PrintGetterList
)

type ICloudProvider interface {
	GetProvider() cloudprovider.ICloudProvider
	GetDefaultRegionId() string
}

type cloudProvider struct {
	provider  cloudprovider.ICloudProvider
	globalOpt *GlobalOptions
}

func NewCloudProvider(opt *GlobalOptions) (ICloudProvider, error) {
	if len(opt.AccessKey) == 0 {
		return nil, errors.Errorf("Missing accessKey")
	}

	if len(opt.Secret) == 0 {
		return nil, errors.Errorf("Missing secret")
	}

	cfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := cfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	factory, err := cloudprovider.GetProviderFactory(opt.Provider)
	if err != nil {
		return nil, errors.Wrapf(err, "GetProviderFactory")
	}

	p, err := factory.GetProvider(cloudprovider.ProviderConfig{
		URL:           opt.CloudEnv,
		Account:       opt.AccessKey,
		Secret:        opt.Secret,
		ProxyFunc:     proxyFunc,
		DefaultRegion: opt.Region,
		Debug:         opt.Debug,
	})

	if err != nil {
		return nil, errors.Wrap(err, "GetProvider")
	}

	var dr cloudprovider.ICloudRegion = nil

	if opt.Region != "" {
		dr, _ = generic.GetResourceByIdOrName(p.GetIRegions(), opt.Region)
		if dr == nil {
			rs := p.GetIRegions()
			regions := make([]string, len(rs))
			for i := range rs {
				regions[i] = rs[i].GetId()
			}
			return nil, errors.Errorf("Not found region by %q of %q, choose one of:\n%s", opt.Region, opt.Provider, jsonutils.Marshal(regions).PrettyString())
		}
	}

	return &cloudProvider{
		provider:  p,
		globalOpt: opt,
	}, nil
}

func (p *cloudProvider) GetProvider() cloudprovider.ICloudProvider {
	return p.provider
}

func (p *cloudProvider) GetDefaultRegionId() string {
	return p.globalOpt.Region
}

func getResources[T cloudprovider.ICloudResource](
	resType string,
	nf func() ([]T, error),
	specifyIdOrName string,
	mustMatch bool) ([]T, error) {
	objs, err := nf()
	if err != nil {
		return nil, errors.Wrapf(err, "get %q resources", resType)
	}

	if specifyIdOrName == "" {
		if mustMatch {
			return nil, errors.Errorf("Use '--%s' to specify %s", resType, resType)
		} else {
			return objs, nil
		}
	}
	obj, err := generic.GetResourceByIdOrName(objs, specifyIdOrName)
	if err != nil {
		if errors.Cause(err) == cloudprovider.ErrNotFound {
			return nil, errors.Wrapf(err, "%s %s not found", resType, specifyIdOrName)
		}
	}
	return []T{obj}, nil
}

func iterResources[T cloudprovider.ICloudResource](
	resType string,
	nf func() ([]T, error),
	specifyIdOrName string,
	mustMatch bool,
	ef func(T) error,
) error {
	op, err := generic.NewOperator(func() ([]T, error) {
		return getResources(resType, nf, specifyIdOrName, mustMatch)
	})
	if err != nil {
		return err
	}
	return op.Iter(func(t T) error {
		log.Infof("With %s %q", resType, t.GetGlobalId())
		return ef(t)
	}, false)
}

// func regionR[OPT any](opts OPT, command string, desc string, requireDefaultRegion bool, cb func(cloudprovider.ICloudRegion, OPT) error) {
// 	R(opts, command, desc, func(cli ICloudProvider, args OPT) error {
// 		return iterResources(
// 			"region",
// 			func() ([]cloudprovider.ICloudRegion, error) {
// 				return cli.GetProvider().GetIRegions(), nil
// 			},
// 			cli.GetDefaultRegionId(),
// 			requireDefaultRegion,
// 			func(region cloudprovider.ICloudRegion) error {
// 				if err := cb(region, args); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		)
// 	})
// }
//
// func RegionR[T any](opts T, command string, desc string, cb func(cloudprovider.ICloudRegion, T) error) {
// 	regionR(opts, command, desc, true, cb)
// }
//
// func RegionRList[T any](opts T, command string, desc string, cb func(cloudprovider.ICloudRegion, T) error) {
// 	regionR(opts, command, desc, false, cb)
// }
//
// func zoneR[T IZoneBaseOptions](opts T, command string, desc string, requireZone bool, cb func(cloudprovider.ICloudZone, T) error) {
// 	RegionRList(opts, command, desc, func(region cloudprovider.ICloudRegion, opts T) error {
// 		return iterResources(
// 			"zone",
// 			region.GetIZones,
// 			opts.GetZoneId(),
// 			requireZone,
// 			func(zone cloudprovider.ICloudZone) error {
// 				if err := cb(zone, opts); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		)
// 	})
// }
//
// func ZoneRList[T IZoneBaseOptions](opts T, command string, desc string, cb func(cloudprovider.ICloudZone, T) error) {
// 	zoneR(opts, command, desc, false, cb)
// }
//
// func ZoneR[T IZoneBaseOptions](opts T, command string, desc string, cb func(cloudprovider.ICloudZone, T) error) {
// 	zoneR(opts, command, desc, true, cb)
// }
//
// func hostR[T IHostBaseOptions](opts T, command string, desc string, requireHost bool, cb func(cloudprovider.ICloudHost, T) error) {
// 	ZoneRList(opts, command, desc, func(zone cloudprovider.ICloudZone, opts T) error {
// 		return iterResources(
// 			"host",
// 			zone.GetIHosts,
// 			opts.GetHostId(),
// 			requireHost,
// 			func(host cloudprovider.ICloudHost) error {
// 				if err := cb(host, opts); err != nil {
// 					return err
// 				}
// 				return nil
// 			},
// 		)
// 	})
// }
//
// func HostRList[T IHostBaseOptions](opts T, command string, desc string, cb func(cloudprovider.ICloudHost, T) error) {
// 	hostR(opts, command, desc, false, cb)
// }
//
// func HostR[T IHostBaseOptions](opts T, command string, desc string, cb func(cloudprovider.ICloudHost, T) error) {
// 	hostR(opts, command, desc, true, cb)
// }
