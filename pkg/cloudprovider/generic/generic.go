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

package generic

import (
	"yunion.io/x/pkg/errors"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func GetResourceByFuncs[T cloudprovider.ICloudResource](objs []T, idFs []func(T) string, id string) (T, error) {
	for i := range objs {
		obj := objs[i]
		for _, f := range idFs {
			if f(obj) == id {
				return obj, nil
			}
		}
	}
	return *new(T), cloudprovider.ErrNotFound
}

func GetResourceId[T cloudprovider.ICloudResource](obj T) string {
	return obj.GetId()
}

func GetResourceName[T cloudprovider.ICloudResource](obj T) string {
	return obj.GetName()
}

func GetResourceGlobalId[T cloudprovider.ICloudResource](obj T) string {
	return obj.GetGlobalId()
}

func GetResourceById[T cloudprovider.ICloudResource](objs []T, id string) (T, error) {
	return GetResourceByFuncs(objs, [](func(T) string){GetResourceId[T]}, id)
}

func GetResourceByName[T cloudprovider.ICloudResource](objs []T, id string) (T, error) {
	return GetResourceByFuncs(objs, [](func(T) string){GetResourceName[T]}, id)
}

func GetResourceByGlobalId[T cloudprovider.ICloudResource](objs []T, id string) (T, error) {
	return GetResourceByFuncs(objs, [](func(T) string){GetResourceGlobalId[T]}, id)
}

func GetResourceByIdOrName[T cloudprovider.ICloudResource](objs []T, idOrName string) (T, error) {
	return GetResourceByFuncs(objs, [](func(T) string){
		GetResourceId[T],
		GetResourceName[T],
		GetResourceGlobalId[T],
	}, idOrName)
}

func Iter[T cloudprovider.ICloudResource](objs []T, f func(T) error, continueOnErr bool) error {
	var errs []error
	for _, obj := range objs {
		if err := f(obj); err != nil {
			if !continueOnErr {
				return errors.Wrapf(err, "resource %q", obj.GetGlobalId())
			} else {
				errs = append(errs, err)
			}
		}
	}
	return errors.NewAggregate(errs)
}
