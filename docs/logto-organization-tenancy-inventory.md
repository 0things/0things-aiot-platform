# Logto organization tenancy inventory

The application no longer stores local users, organizations, or organization memberships. Logto is the identity and organization system of record.

| Layer | Field or source | Representation | Rule |
| --- | --- | --- | --- |
| Logto token | `sub` | string | Backend request user identifier |
| Logto token | `organization_id` | string | Required organization scope |
| Backend model | `organization_id` on products, devices, device groups, OTA packages, message parsers, and rule chains | string | Persisted Logto Organization ID |
| Request context | `tenant.WithOrganization` | string | Repository scope source; empty means unauthenticated/invalid request |
| Repository scope | organization predicates | string | Every tenant-owned read/write is filtered by request organization |
| Frontend | Logto access-token claims | string | Current organization display comes from the token; no local organization cache or switch API |

The following tables are intentionally not created by backend migration:

- `users`
- `organizations`
- `organization_users`

The backend does not translate organization IDs to integers and does not fall back to a default organization. Organization membership and organization-token issuance are managed by Logto.
