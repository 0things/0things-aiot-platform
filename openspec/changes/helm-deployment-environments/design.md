## Context

See `proposal.md` for motivation and `specs/deployment/helm-environments/spec.md` for observable behavior. The current deployment is a single Docker Compose topology whose application processes load one YAML file selected by `APP_CONF`. Those configuration files mix public settings and credentials, and the frontend Nginx configuration assumes fixed Compose DNS names.

The Helm implementation must remain self-contained, render without access to a cluster, and support two releases in separate namespaces without name collisions. It must also preserve the user's requirement that all required configuration keys are discoverable in values files and that production and test credentials are independently defined.

## Goals / Non-Goals

**Goals:**

- Produce one installable application chart under `deploy/helm/0things`.
- Keep the template tree organized by application, infrastructure, configuration, and networking modules.
- Make both environment overlays independently renderable and lintable.
- Generate every required application and infrastructure Secret from values by default.
- Use fixed component names and rendered internal service addresses; production and test releases use separate namespaces.
- Make Helm the repository's only full-platform deployment method.

**Non-Goals:**

- Provision the Kubernetes cluster, ingress controller, DNS records, TLS certificates, or cloud storage classes.
- Change application APIs or add generic environment-variable overrides to the Go services.
- Implement high-availability database clustering or managed-cloud database integration.
- Manage or execute application database migrations; required schemas are prepared outside this chart.

## Decisions

### Use one self-contained application chart

Infrastructure resources will be maintained as templates in this chart rather than hidden behind third-party subchart value schemas. This keeps all deployment inputs visible in `values.yaml`, avoids dependency download requirements, and makes the complete test and production manifests reviewable with `helm template`.

Alternative considered: upstream PostgreSQL, Redis, NATS, EMQX, and TDengine subcharts. This reduces infrastructure template ownership but introduces several independent values contracts, release-name-dependent service discovery, and external chart availability. Those trade-offs conflict with the requested simple, explicit deployment package.

### Organize templates by deployment concern

The chart will use four template modules:

```text
templates/
├── application/
├── infrastructure/
├── configuration/
└── networking/
```

Each application remains a file inside the shared application module. `application/transport-mqtt.yaml` contains the MQTT transport workload. Common labels, names, image rendering, storage, and checksums live in `_helpers.tpl`.

Alternative considered: one directory per project. It was rejected because it creates many shallow directories and obscures cross-service deployment concerns.

### Use base values plus explicit environment overlays

`values.yaml` defines the complete shape and safe defaults. `values-test.yaml` and `values-prod.yaml` override environment name, public URLs, image tags, credentials, replicas, resources, and storage. Operators select one overlay with `-f`; test and production releases use separate namespaces.

Credentials remain explicit in the environment values as requested. Templates convert them into Kubernetes Secrets. Each secret group also supports `create: false` and `existingSecret` so a deployment can later adopt an external secret controller without changing workloads.

### Render complete application configuration into Secrets

Because the Go configuration loaders only read `APP_CONF` and do not apply environment-variable overrides, Helm will render each complete runtime YAML file into a Secret and mount it at `/data/app/config/helm.yml`. Workloads set `APP_CONF=config/helm.yml`. Non-sensitive frontend Nginx configuration is rendered into a ConfigMap.

Pod templates carry checksums of the configuration and secret templates so Helm upgrades roll pods when values change.

Alternative considered: `envFrom` only. It was rejected because the current services would ignore those variables for nested Viper keys.

### Use fixed Kubernetes names within isolated namespaces

Resources use fixed component names such as `postgres` and `backend`. Production and test are installed in separate namespaces, so their names do not collide. Runtime configuration and frontend Nginx configuration use those same names for service discovery.

### Manage infrastructure as single-instance StatefulSets

PostgreSQL, Logto PostgreSQL, Redis, NATS JetStream, EMQX, and TDengine use StatefulSets with per-release volume claim templates. Storage class, access mode, and size are configurable. This is appropriate for the current single-server deployment target but is not presented as a highly available production database architecture.

### Expose one HTTP entry point through the frontend

The frontend Nginx configuration proxies `/api`, `/copilot`, `/mcp`, and `/swagger` to fixed component Services in the release namespace. The chart exposes frontend Ingress configuration and separate optional Logto public/admin Ingress rules. MQTT exposure remains configurable through the EMQX Service type.

### Validate before installation

Templates use `required` for credentials whose necessity depends on component enablement. Verification consists of `helm lint`, production and test `helm template`, OpenSpec validation, and Kubernetes schema validation when `kubeconform` is available.

### Remove Docker Compose deployment entry points

The root and service-local Compose manifests and Compose-dependent start/stop scripts are removed. Helm is the only checked-in full-platform deployment path. Local development may run application processes directly against operator-provided dependencies or a test Helm release.

The frontend image still needs a standalone Nginx configuration at build time, so that file moves from the Compose deployment directory to `frontend/nginx.conf` and the Dockerfile references its new application-owned location.

## Risks / Trade-offs

- [Plaintext credentials may be committed in environment values] → Use placeholders in repository defaults, document private override files and `existingSecret`, and fail production rendering when required values are empty.
- [Single-instance infrastructure is not highly available] → Make persistence and resources explicit and document external/managed infrastructure as future work.
- [StatefulSet credential changes do not update an already initialized database] → Document that Secret rotation requires a database-level password change before updating the values file.
- [Ingress hosts and TLS depend on cluster configuration] → Keep ingress optional and expose class name, hosts, annotations, and TLS secret names in each environment overlay.

## Migration Plan

1. Render and validate the test overlay locally.
2. Publish all configured application images.
3. Prepare the required application database schema outside the chart, then install the test release in its own namespace and verify application readiness, persistence, and MQTT routing.
4. Back up production data and render the production overlay for review.
5. Prepare any required production schema changes, then install or upgrade production with `--atomic --wait`.
6. On failure, inspect pod events and logs; Helm automatically rolls back when `--atomic` is used.
7. Remove the obsolete Docker Compose manifests and launch scripts after Helm validation, then verify that active build and deployment documentation contains no Compose dependency.
