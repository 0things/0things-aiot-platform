## Why

The repository currently relies on Docker Compose and lacks a Kubernetes deployment package, making production and test deployments difficult to reproduce and manage consistently. A standards-compliant Helm chart is needed to deploy the complete platform with explicit environment configuration and without undocumented prerequisite secrets.

## What Changes

- Add an application Helm chart under `deploy/helm/0things` for the platform workloads and required infrastructure.
- Organize templates by deployment module: application, infrastructure, configuration, and networking, instead of creating a directory per project.
- Provide shared defaults in `values.yaml` and environment overrides in `values-test.yaml` and `values-prod.yaml`.
- Define database, middleware, image registry, and application credentials through environment values and render the required Kubernetes Secrets automatically.
- Deploy and configure backend, MCP server, data engine, MQTT transport, AI copilot, frontend, Logto, PostgreSQL, Redis, NATS, EMQX, and TDengine.
- Apply standard Helm and Kubernetes metadata, selectors, resource configuration, persistence, probes, ingress configuration, configuration checksums, and values validation.
- Document repeatable test and production installation, upgrade, validation, and rollback commands.
- Remove repository Docker Compose manifests and Compose-dependent launch scripts so Helm is the only checked-in full-platform deployment method.

## Capabilities

### New Capabilities

- `deployment/helm-environments`: Defines complete Helm-based platform deployment and independent production/test environment configuration.

### Modified Capabilities

None.

## Impact

- Adds deployment artifacts under `deploy/helm/0things`.
- Removes the root and service-local Docker Compose manifests and their launch scripts, and updates deployment documentation to use Helm.
- Introduces Kubernetes resources with explicitly versioned container images for platform infrastructure.
- Requires separately published images for each runnable workload.
- Replaces hard-coded Docker service addresses in mounted runtime configuration with Helm-rendered Kubernetes service addresses.
- Does not change public HTTP APIs or application domain behavior.
