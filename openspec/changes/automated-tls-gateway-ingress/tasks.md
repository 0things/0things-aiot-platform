## 1. Chart networking

- [x] 1.1 Add one gateway API/MCP Traefik Ingress and `/api` StripPrefix middleware; verify rendered host, paths, backends, and middleware annotation with `helm template`.
- [x] 1.2 Add one optional production Certificate and shared TLS Ingress settings; verify production rendering contains one Certificate and all four hostnames use its Secret.
- [x] 1.3 Update local prod/test environment values without overwriting credentials; verify prod uses HTTPS/ClusterIP and test renders no public resources.
- [x] 1.4 Manage the Cloudflare ClusterIssuer and Token Secret in the 0things Chart using prod-only values; verify prod renders both and test renders neither.

## 2. Cluster prerequisites and frontend

- [x] 2.1 Keep only the cert-manager controller as a standalone k3s HelmChart; move ClusterIssuer and Token Secret deployment into the 0things release and document restricted Token handling.
- [x] 2.2 Align frontend production Logto build configuration with `auth.0thing.com` and API Resource Identifier `https://gateway.0thing.com/api`; verify Backend audience and frontend resource match.

## 3. Validation

- [x] 3.1 Update Helm README with domain routes, certificate operations and rollback; verify instructions match rendered manifests.
- [x] 3.2 Run `openspec validate --strict`, `helm lint`, and production/test `helm template` assertions; verify all checks pass.
