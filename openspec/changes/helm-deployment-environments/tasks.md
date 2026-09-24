## 1. Chart Foundation

- [x] 1.1 Create the Helm chart metadata, module directories, common helpers, and `.helmignore`; verify `helm show chart deploy/helm/0things` succeeds
- [x] 1.2 Define complete shared, test, and production values including independent credentials, images, resources, persistence, services, and ingress; verify `helm show values` and YAML parsing succeed
- [x] 1.3 Validate required credentials in Helm templates with actionable `required` errors; verify chart values render through Helm

## 2. Configuration and Infrastructure

- [x] 2.1 Render application runtime files, frontend Nginx configuration, registry authentication, and component credentials from values; verify both overlays produce every referenced ConfigMap and Secret
- [x] 2.2 Implement PostgreSQL and Logto PostgreSQL StatefulSets with Services, health checks, and persistent volume claims; verify both overlays render valid database resources
- [x] 2.3 Implement Redis, NATS, EMQX, and TDengine StatefulSets with configurable persistence and service ports; verify enabled and disabled component renders contain no dangling resources

## 3. Application Workloads

- [x] 3.1 Implement backend, MCP server, data engine, transport, AI copilot, frontend, and Logto Deployments with standard labels, resources, probes, configuration mounts, and checksums; verify every enabled workload renders with matching selectors
- [x] 3.2 Implement application, infrastructure, and externally configurable Services; verify all rendered runtime DNS names correspond to rendered Services

## 4. External Access and Operations

- [x] 4.1 Implement optional frontend and Logto Ingress resources and verify disabled ingress renders no Ingress while configured test and production ingress renders valid routes
- [x] 4.2 Add `NOTES.txt` and chart documentation for install, upgrade, validation, rollback, images, credentials, persistence, and environment selection; verify every documented path and command matches the chart
- [x] 4.3 Remove Docker Compose manifests and launch scripts, move shared build configuration to its owning module, and update active deployment documentation to use Helm

## 5. Validation

- [x] 5.1 Run `helm lint` and `helm template` for base, test, and production values and fix all errors
- [x] 5.2 Run Kubernetes schema validation with `kubeconform` when available, validate the OpenSpec change strictly, and confirm the implementation only touches the intended deployment and change files
