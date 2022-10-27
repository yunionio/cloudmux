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

type Operator[T cloudprovider.ICloudResource] struct {
	objects []T
}

func NewOperator[T cloudprovider.ICloudResource](nf func() ([]T, error)) (*Operator[T], error) {
	objs, err := nf()
	if err != nil {
		return nil, errors.Wrap(err, "new resources")
	}

	return &Operator[T]{
		objects: objs,
	}, nil
}

func (o Operator[T]) Iter(f func(T) error, continueOnErr bool) error {
	return Iter(o.objects, f, continueOnErr)
}
