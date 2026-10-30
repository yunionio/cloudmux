# Volcengine PostgreSQL inventory

The Volcengine provider supports **read-only RDS for PostgreSQL discovery**.
MySQL and SQL Server are not included in the RDS capability of this adapter.

The adapter uses the `rds_postgresql` service, API version `2022-01-01`, with
signed JSON POST requests over verified HTTPS. The cloud account needs access
to these read operations:

- `DescribeRegions`: discover the regions that offer PostgreSQL.
- `DescribeDBInstances`: enumerate every page of instances in each region.
- `DescribeDBInstanceDetail`: read instance configuration and endpoints.

Only the above actions are accepted by the PostgreSQL request helper. The
capability is advertised as `rds.readonly`, using the existing framework's
read-only capability convention. Resource creation, deletion, reboot, renewal,
resizing, public-access changes, and tag modifications are not implemented.

## Inventory fields

Instance ID and name, status, PostgreSQL version, primary-node CPU and memory,
allocated storage, SKU, storage type, availability zones, VPC/subnet, primary
read/write endpoints, project, tags, creation/expiry time, and billing mode.

The primary-node memory reported in GB is converted to MB. High-availability
replicas are not added to the primary instance's sizing. The adapter does not
query database contents or export passwords. Parameters, database accounts,
databases, backups, performance metrics and maintenance windows are not part
of the initial implementation.

Listing fails instead of returning an empty or partial inventory on denied
requests, malformed responses, interrupted pagination, repeated instance IDs,
or invalid detail identities. An unsupported region returns `ErrNotSupported`.
The cloud platform project name is returned as the external project identifier,
consistent with the existing Volcengine resource adapters.

## Validation

```sh
go test -mod=vendor ./pkg/multicloud/volcengine ./pkg/multicloud/volcengine/provider
```

Tests use synthetic IDs and a mocked HTTP transport. They cover pagination,
stable IDs on repeated discovery, metadata/unit conversion, primary endpoints,
permission and partial-response failures, unsupported regions, rejection of
write actions, and TLS certificate verification. Live integration also needs
comparison of the discovered IDs with the cloud console or API and two
successful Cloudpods syncs without duplicates.

API references:

- https://docs.volcengine.com/docs/RDS_for_PG/describeregions
- https://docs.volcengine.com/docs/RDS_for_PG/describedbinstances
- https://docs.volcengine.com/docs/RDS_for_PG/describedbinstancedetail


## Additional read-only inventory

Redis (`2020-12-07`), CLB (`2020-04-01`) and VKE (`2022-05-12`) are
supported through explicit Describe/List allowlists, verified HTTPS, account
proxy settings and complete-list validation. Schemas were checked against
`volcengine/volcengine-go-sdk` at `7a94574ceb30be91c364dfdac484e901438c0743`.

- Redis: instance configuration, MiB capacity, primary private/public endpoints,
  VPC/subnet, global zone ID, tags and project mapping. No credentials, accounts,
  data, backups, parameters or performance measurements are read.
- CLB: instances, listeners, referenced server groups and backend ECS/IP targets.
  TCP/UDP listeners correctly have no HTTP forwarding rules. HTTP forwarding
  rules, certificate/ACL inventory and security-group/EIP relationship discovery
  are not implemented and return unsupported errors rather than fake empty lists.
  Unknown health/performance values must not be interpreted as measured zeros.
- VKE: clusters, node pools and nodes, including cluster-to-pool/node ownership.
  Kubeconfig retrieval and Kubernetes workload access are deliberately unsupported.

Writes return unsupported errors. Pagination, permission failures, changed totals,
missing arrays, duplicate IDs and mismatched resource/parent identities fail the
relevant discovery without returning partial inventory. Unimplemented child
resource discovery may be reported by Cloudpods as unsupported; this is not a
claim of full provider product parity.

MySQL and SQL Server use read-only RDS discovery (`2022-01-01`). The RDS
reconciliation result is atomic across PostgreSQL, MySQL and SQL Server: any
permission, detail, or pagination failure fails the complete regional discovery,
preserving previously synchronized engines. Regions explicitly absent from an
engine region catalog are skipped; denied access is never interpreted as absence.

SLS and model usage/billing require new resource types or pages and are out of
scope. This extension only adds adapters for existing Cloudpods resource types.
