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

import "yunion.io/x/cloudmux/pkg/cloudprovider"

func init() {
	cmd := NewCommand("host")

	type HostListOptions struct {
		ListBaseOptions
		ZoneBaseOptions
	}

	ZoneR[HostListOptions](cmd).List("list", "List hosts", func(cli cloudprovider.ICloudZone, args *HostListOptions) (any, error) {
		hosts, err := cli.GetIHosts()
		if err != nil {
			return nil, err
		}
		objs := make([]interface{}, len(hosts))
		for i := range hosts {
			host := hosts[i]
			objs[i] = map[string]interface{}{
				"id":          host.GetId(),
				"global_id":   host.GetGlobalId(),
				"name":        host.GetName(),
				"is_emulated": host.IsEmulated(),
				"status":      host.GetStatus(),
				"host_status": host.GetHostStatus(),
			}
		}
		return objs, nil
	})
}
