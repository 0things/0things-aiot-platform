# Logto

Logto runs with a dedicated PostgreSQL database and shares the `0things-net`
Docker network with the platform services. Its image is built from
`logto/Dockerfile`; the services are defined in the root Compose file.

Configure Logto in `logto/.env` and start the complete stack from the
repository root. For a new environment, copy `logto/.env.example` first:

```bash
cp logto/.env.example logto/.env
docker compose up -d
```

Logto is available at `http://localhost:3001` and the admin console at
`http://localhost:3002`.

Set the Logto database password and `DB_URL` in `logto/.env` before deployment.
Provider credentials are configured in Logto Admin and must not be committed.
