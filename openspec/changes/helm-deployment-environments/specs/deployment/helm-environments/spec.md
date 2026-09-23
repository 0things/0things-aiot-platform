## Purpose

Provide a reproducible Helm deployment contract for the complete 0things platform, with independently configurable production and test environments and no undocumented deployment prerequisites.

## ADDED Requirements

### Requirement: Complete platform installation
The Helm chart SHALL install the 0things application workloads and the infrastructure services required by the enabled deployment configuration.

#### Scenario: Install a complete environment
- **WHEN** an operator installs the chart with a supported environment values file
- **THEN** Kubernetes resources are rendered for backend, MCP server, data engine, HTTP and MQTT transports, AI copilot, frontend, Logto, PostgreSQL, Redis, NATS, EMQX, and TDengine

#### Scenario: Disable an optional component
- **WHEN** an operator disables a component through its `enabled` value
- **THEN** the chart omits that component's workload and service resources without producing invalid references

### Requirement: Environment-specific values
The chart SHALL provide shared defaults and separate production and test values overlays, with each overlay defining its own endpoints, image settings, credentials, replica counts, resources, and persistence settings where they differ.

#### Scenario: Render test resources
- **WHEN** the chart is rendered with `values-test.yaml`
- **THEN** the rendered resources use the test environment configuration and do not require production values

#### Scenario: Render production resources
- **WHEN** the chart is rendered with `values-prod.yaml`
- **THEN** the rendered resources use the production environment configuration and do not inherit test credentials or endpoints

### Requirement: Values-driven credentials
The chart SHALL expose required database, middleware, application, and registry credential inputs through values and SHALL render the corresponding Kubernetes Secrets when secret creation is enabled.

#### Scenario: Create deployment secrets
- **WHEN** complete credential values are supplied and secret creation is enabled
- **THEN** the chart renders every Secret referenced by an enabled workload

#### Scenario: Reject missing required credentials
- **WHEN** an enabled workload lacks a required credential value or existing Secret reference
- **THEN** chart validation fails with an actionable error before resources are installed

### Requirement: Runtime service configuration
The chart SHALL render runtime configuration using Kubernetes service discovery names and SHALL mount configuration at the paths expected by each workload.

#### Scenario: Connect internal services
- **WHEN** workloads start from a rendered release
- **THEN** their PostgreSQL, Redis, NATS, EMQX, TDengine, Logto, MCP, and backend addresses resolve to services in the release namespace

#### Scenario: Roll out a configuration change
- **WHEN** a rendered ConfigMap or Secret changes during an upgrade
- **THEN** affected application workloads receive a pod-template change that triggers a rolling update

### Requirement: Persistent state
The chart SHALL define configurable persistent storage for stateful infrastructure and SHALL keep production and test data isolated by Kubernetes namespace and release.

#### Scenario: Preserve state across pod replacement
- **WHEN** a PostgreSQL, Redis, NATS, EMQX, or TDengine pod is replaced
- **THEN** its configured persistent volume remains associated with that environment

### Requirement: Kubernetes and Helm conventions
Rendered resources SHALL use stable selectors, recommended application labels, configurable resource requests and limits, appropriate health checks, and valid Kubernetes API versions.

#### Scenario: Validate the chart
- **WHEN** maintainers run chart linting, environment rendering, and Kubernetes schema validation for production and test overlays
- **THEN** both environments pass without template, values schema, or resource schema errors

### Requirement: Operational documentation
The chart SHALL document installation, upgrade, validation, rollback, credential configuration, image pull configuration, and environment selection.

#### Scenario: Deploy from documented commands
- **WHEN** an operator follows the documented production or test installation command
- **THEN** the selected environment can be installed without discovering additional undocumented Kubernetes resources or secrets

### Requirement: Helm-only repository deployment
The repository SHALL expose the Helm chart as its only checked-in full-platform deployment method and SHALL NOT include Docker Compose manifests or Compose-dependent launch scripts.

#### Scenario: Inspect deployment entry points
- **WHEN** a maintainer inspects the repository's deployment manifests, build paths, and operational documentation
- **THEN** the supported full-platform deployment path uses `deploy/helm/0things` without requiring Docker Compose files or commands
