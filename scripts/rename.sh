#!/bin/bash

target_dir=$1

# find $target_dir | grep '.go$' | xargs -I{} sed -r -i "" "s|yunion.io/x/onecloud/pkg/cloudprovider|yunion.io/x/cloudmux/pkg/cloudprovider|g" {}

    # s|yunion.io/x/onecloud/pkg/apis|yunion.io/x/cloudmux/pkg/apis|g; \

find $target_dir | grep '.go$' | xargs -I{} sed -r -i "" "s|multicloud.ApsaraTags|ApsaraTags|g; \
    s|yunion.io/x/onecloud/pkg/cloudprovider|yunion.io/x/cloudmux/pkg/cloudprovider|g; \
    s|yunion.io/x/onecloud/pkg/multicloud|yunion.io/x/cloudmux/pkg/multicloud|g; \
    s|multicloud.CloudpodsTags|CloudpodsTags|g; \
    s|multicloud.BingoTags|BingoTags|g; \
    s|multicloud.AwsTags|AwsTags|g; \
    s|multicloud.AzureTags|AzureTags|g; \
    s|multicloud.EcloudTags|EcloudTags|g; \
    s|multicloud.AliyunTags|AliyunTags|g; \
    s|multicloud.HuaweiTags|HuaweiTags|g; \
    s|HuaweiTags|huawei.HuaweiTags|g; \
    s|multicloud.InCloudSphereTags|InCloudSphereTags|g; \
    s|multicloud.JdcloudTags|JdcloudTags|g; \
    s|multicloud.ProxmoxTags|ProxmoxTags|g; \
    s|multicloud.UcloudTags|UcloudTags|g; \
    s|multicloud.GoogleTags|GoogleTags|g; \
    s|multicloud.ZStackTags|ZStackTags|g; \
    s|multicloud.QcloudTags|QcloudTags|g; \
    s|multicloud.RemoteFileTags|RemoteFileTags|g; \
    s|multicloud.OpenStackTags|OpenStackTags|g; \
    s|multicloud.CtyunTags|CtyunTags|g" {}
