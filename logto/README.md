# Logto

Logto runs with a dedicated PostgreSQL database and shares the `0things-net`
Docker network with the platform services. Its image is built from
`logto/Dockerfile`; the services are defined in the root Compose file.

Copy the root environment template and start the complete stack from the
repository root:

```bash
cp .env.example .env
docker compose up -d
```

Logto is available at `http://localhost:3001` and the admin console at
`http://localhost:3002`.

Set `LOGTO_POSTGRES_PASSWORD` in the root `.env` before deployment. Provider
credentials are configured in Logto Admin and must not be committed.
