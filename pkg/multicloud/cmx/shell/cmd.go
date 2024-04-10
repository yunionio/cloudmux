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
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type Command struct {
	prefix string
}

func NewCommand(prefix string) *Command {
	return &Command{
		prefix: prefix,
	}
}

// NewCO returns a Command with Option
func NewCO[OPT any](cmd *Command) *CommandOption[OPT] {
	return &CommandOption[OPT]{
		Command: cmd,
		opt:     new(OPT),
	}
}

type CommandOption[OPT any] struct {
	*Command
	opt           *OPT
	isList        bool
	useGetterList bool
	requireRegion bool
	requireZone   bool
	requireHost   bool
}

func (co *CommandOption[OPT]) RawRun(suffix string, desc string, cb func(ICloudProvider, *OPT) error) {
	R(co.opt, fmt.Sprintf("%s-%s", co.prefix, suffix), desc, cb)
}

func (co *CommandOption[OPT]) Run(suffix string, desc string, cb func(ICloudProvider, *OPT) (any, error)) {
	co.RawRun(suffix, desc, func(cli ICloudProvider, opt *OPT) error {
		data, err := cb(cli, opt)
		if err != nil {
			return err
		}
		return co.processData(data)
	})
}

func (co *CommandOption[OPT]) processData(data any) error {
	if data == nil {
		return nil
	}
	if co.isList {
		lo, ok := interface{}(co.opt).(IListOption)
		if ok {
			cols := []string{}
			if !lo.IsDetails() {
				cols = lo.GetColumns()
			}
			PrintList(data, 0, lo.GetOffset(), lo.GetLimit(), cols)
			return nil
		} else {
			if co.useGetterList {
				PrintGetterList(data, nil)
				return nil
			}
			PrintList(data, 0, 0, 0, nil)
			return nil
		}
	}
	PrintObject(data)
	return nil
}

func (co *CommandOption[OPT]) UseList() *CommandOption[OPT] {
	co.isList = true
	return co
}

func (co *CommandOption[OPT]) UseGetterList() *CommandOption[OPT] {
	co.UseList().useGetterList = true
	return co
}

func (co *CommandOption[OPT]) RunByProvider(suffix string, desc string, cb func(cloudprovider.ICloudProvider, *OPT) (any, error)) {
	co.Run(suffix, desc, func(cli ICloudProvider, opt *OPT) (any, error) {
		return cb(cli.GetProvider(), opt)
	})
}

func (co *CommandOption[OPT]) RequireRegion() *CommandOption[OPT] {
	co.requireRegion = true
	return co
}

func (co *CommandOption[OPT]) RunByRegion(
	suffix string, desc string,
	cb func(cli cloudprovider.ICloudRegion, args *OPT) (any, error),
) {
	co.RawRun(suffix, desc, func(cli ICloudProvider, o *OPT) error {
		return iterResources(
			"region",
			func() ([]cloudprovider.ICloudRegion, error) {
				return cli.GetProvider().GetIRegions()
			},
			cli.GetDefaultRegionId(),
			co.requireRegion,
			func(region cloudprovider.ICloudRegion) error {
				data, err := cb(region, o)
				if err != nil {
					return err
				}
				return co.processData(data)
			},
		)
	})
}

func (co *CommandOption[OPT]) RequireZone() *CommandOption[OPT] {
	co.requireZone = true
	return co
}

func (co *CommandOption[OPT]) RunByZone(
	suffix string, desc string,
	cb func(cloudprovider.ICloudZone, *OPT) (any, error),
) {
	co.RunByRegion(suffix, desc, func(ir cloudprovider.ICloudRegion, o *OPT) (any, error) {
		return nil, iterResources(
			"zone",
			ir.GetIZones,
			interface{}(o).(IZoneBaseOptions).GetZoneId(),
			co.requireZone,
			func(zone cloudprovider.ICloudZone) error {
				data, err := cb(zone, o)
				if err != nil {
					return err
				}
				return co.processData(data)
			},
		)
	})
}

func (co *CommandOption[OPT]) RequireHost() *CommandOption[OPT] {
	co.requireHost = true
	return co
}

func (co *CommandOption[OPT]) RunByHost(
	suffix string, desc string,
	cb func(cloudprovider.ICloudHost, *OPT) (any, error),
) {
	co.RunByZone(suffix, desc, func(iz cloudprovider.ICloudZone, o *OPT) (any, error) {
		return nil, iterResources(
			"host",
			iz.GetIHosts,
			interface{}(o).(IHostBaseOptions).GetHostId(),
			co.requireHost,
			func(host cloudprovider.ICloudHost) error {
				data, err := cb(host, o)
				if err != nil {
					return err
				}
				return co.processData(data)
			},
		)
	})
}

type IRunner[O any, C any] interface {
	RequireRegion() IRunner[O, C]
	RequireZone() IRunner[O, C]
	RequireHost() IRunner[O, C]

	List(suffix, desc string, cb func(cli C, args *O) (any, error))
	GetterList(suffix, desc string, cb func(cli C, args *O) (any, error))
	Run(suffix, desc string, cb func(cli C, args *O) (any, error))
}

type sBaseRunner[O any, C any] struct {
	cmd  *Command
	runF func(*CommandOption[O]) func(string, string, func(C, *O) (any, error))
	coFs []func(*CommandOption[O]) *CommandOption[O]
}

func newBaseRunner[O any, C any](
	cmd *Command,
	runF func(*CommandOption[O]) func(string, string, func(C, *O) (any, error)),
) IRunner[O, C] {
	return &sBaseRunner[O, C]{
		cmd:  cmd,
		runF: runF,
		coFs: make([]func(*CommandOption[O]) *CommandOption[O], 0),
	}
}

func (r *sBaseRunner[O, C]) newCO() *CommandOption[O] {
	co := NewCO[O](r.cmd)
	for _, f := range r.coFs {
		co = f(co)
	}
	return co
}

func (r *sBaseRunner[O, C]) addCOF(f func(*CommandOption[O]) *CommandOption[O]) *sBaseRunner[O, C] {
	r.coFs = append(r.coFs, f)
	return r
}

func (r *sBaseRunner[O, C]) useList() *sBaseRunner[O, C] {
	r.addCOF(func(co *CommandOption[O]) *CommandOption[O] {
		co.UseList()
		return co
	})
	return r
}

func (r *sBaseRunner[O, C]) useGetterList() *sBaseRunner[O, C] {
	r.addCOF(func(co *CommandOption[O]) *CommandOption[O] {
		co.UseGetterList()
		return co
	})
	return r
}

func (r *sBaseRunner[O, C]) RequireRegion() IRunner[O, C] {
	r.addCOF(func(co *CommandOption[O]) *CommandOption[O] {
		co.RequireRegion()
		return co
	})
	return r
}

func (r *sBaseRunner[O, C]) RequireZone() IRunner[O, C] {
	r.addCOF(func(co *CommandOption[O]) *CommandOption[O] {
		co.RequireZone()
		return co
	})
	return r
}

func (r *sBaseRunner[O, C]) RequireHost() IRunner[O, C] {
	r.addCOF(func(co *CommandOption[O]) *CommandOption[O] {
		co.RequireHost()
		return co
	})
	return r
}

func (r *sBaseRunner[O, C]) list(suffix, desc string, cb func(cli C, args *O) (any, error)) {
	r.Run(suffix, desc, cb)
}

func (r *sBaseRunner[O, C]) List(suffix, desc string, cb func(cli C, args *O) (any, error)) {
	r.useList()
	r.list(suffix, desc, cb)
}

func (r *sBaseRunner[O, C]) GetterList(suffix, desc string, cb func(cli C, args *O) (any, error)) {
	r.useGetterList()
	r.list(suffix, desc, cb)
}

func (r *sBaseRunner[O, C]) Run(suffix, desc string, cb func(cli C, args *O) (any, error)) {
	r.runF(r.newCO())(suffix, desc, cb)
}

type sProviderRunner[O any] struct {
	IRunner[O, cloudprovider.ICloudProvider]
}

func ProviderR[O any](cmd *Command) IRunner[O, cloudprovider.ICloudProvider] {
	return &sProviderRunner[O]{
		IRunner: newBaseRunner(
			cmd,
			func(co *CommandOption[O]) func(string, string, func(cloudprovider.ICloudProvider, *O) (any, error)) {
				return co.RunByProvider
			}),
	}
}

type sRegionRunner[O any] struct {
	IRunner[O, cloudprovider.ICloudRegion]
}

func RegionR[O any](cmd *Command) IRunner[O, cloudprovider.ICloudRegion] {
	return &sRegionRunner[O]{
		IRunner: newBaseRunner(
			cmd,
			func(co *CommandOption[O]) func(string, string, func(cloudprovider.ICloudRegion, *O) (any, error)) {
				return co.RunByRegion
			}),
	}
}

type sZoneRunner[O any] struct {
	IRunner[O, cloudprovider.ICloudZone]
}

func ZoneR[O any](cmd *Command) IRunner[O, cloudprovider.ICloudZone] {
	return &sZoneRunner[O]{
		IRunner: newBaseRunner(
			cmd,
			func(co *CommandOption[O]) func(string, string, func(cloudprovider.ICloudZone, *O) (any, error)) {
				return co.RunByZone
			}),
	}
}

type sHostRunner[O IHostBaseOptions] struct {
	IRunner[O, cloudprovider.ICloudHost]
}

func HostR[O IHostBaseOptions](cmd *Command) IRunner[O, cloudprovider.ICloudHost] {
	return &sHostRunner[O]{
		IRunner: newBaseRunner(
			cmd,
			func(co *CommandOption[O]) func(string, string, func(cloudprovider.ICloudHost, *O) (any, error)) {
				return co.RunByHost
			},
		),
	}
}

type SEmptyOptionRunner[C any] struct {
	runner IRunner[EmptyOption, C]
}

func EmptyOptionRunner[C any](
	prefix string,
	rf func(*Command) IRunner[EmptyOption, C],
) *SEmptyOptionRunner[C] {
	cmd := NewCommand(prefix)
	r := rf(cmd)
	return &SEmptyOptionRunner[C]{
		runner: r,
	}
}

func (er *SEmptyOptionRunner[C]) Run(suffix, desc string, cb func(cli C) (any, error)) {
	er.runner.Run(suffix, desc, func(cli C, args *EmptyOption) (any, error) {
		return cb(cli)
	})
}

func (er *SEmptyOptionRunner[C]) List(suffix, desc string, cb func(cli C) (any, error)) {
	er.runner.List(suffix, desc, func(cli C, args *EmptyOption) (any, error) {
		return cb(cli)
	})
}

func (er *SEmptyOptionRunner[C]) GetterList(suffix, desc string, cb func(cli C) (any, error)) {
	er.runner.GetterList(suffix, desc, func(cli C, args *EmptyOption) (any, error) {
		return cb(cli)
	})
}

func EmptyOptionProviderR(prefix string) *SEmptyOptionRunner[cloudprovider.ICloudProvider] {
	return EmptyOptionRunner(prefix, ProviderR[EmptyOption])
}

func EmptyOptionRegionR(prefix string) *SEmptyOptionRunner[cloudprovider.ICloudRegion] {
	return EmptyOptionRunner(prefix, RegionR[EmptyOption])
}

func EmptyOptionZoneR(prefix string) *SEmptyOptionRunner[cloudprovider.ICloudZone] {
	return EmptyOptionRunner(prefix, ZoneR[EmptyOption])
}
