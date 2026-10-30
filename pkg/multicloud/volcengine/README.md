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
