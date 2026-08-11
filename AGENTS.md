# Repository Guidelines

## Project Structure & Module Organization

This repository contains three deployable parts:

- `aiot-frontend/`: Vite + React + TypeScript administration UI. Feature code lives in `src/features/`, shared UI in `src/components/`, API clients in `src/api/`, and translations in `public/locales/{zh,en}/`.
- `aiot-backend/`: Go HTTP service. Keep application code under `cmd/`, `internal/`, and reusable packages under `pkg/`; server tests are in `test/server/`.
- `telemetry-service/`: Go/Kratos telemetry service. Protobuf definitions are in `api/` and `internal/`; generated code stays beside its source.

Do not edit generated frontend clients in `aiot-frontend/src/api/generated/` by hand; regenerate them from the OpenAPI source when contracts change.

## Build, Test, and Development Commands

Run commands from the relevant package directory:

```bash
cd aiot-frontend && pnpm dev          # local Vite server
cd aiot-frontend && pnpm build        # TypeScript check and production build
cd aiot-frontend && pnpm lint         # ESLint
cd aiot-frontend && pnpm format:check # Prettier validation
cd aiot-backend && make test          # Go server tests with coverage report
cd aiot-backend && make build         # build bin/server
cd telemetry-service && make build    # build telemetry binaries
```

Use `make api`, `make config`, or `make generate` in `telemetry-service/` after changing protobuf or generated Go code.

## Coding Style & Naming Conventions

Follow existing formatting: TypeScript uses two spaces, single quotes, and semicolon-free code; Go uses `gofmt`. Name React components in `PascalCase`, hooks as `use-*.ts`, feature directories in kebab-case, and Go files with lowercase snake case. Keep UI strings in both locale files rather than hard-coding language-specific copy. Reuse shadcn components and Lucide icons before adding dependencies.

## Testing Guidelines

Run the narrowest relevant check before handoff, then the package build. Add Go tests as `*_test.go` alongside the existing `test/server` organization. The frontend currently has no dedicated unit-test script; verify changed flows in the local browser and run `pnpm build`, `pnpm lint`, and `pnpm format:check` when practical.

## Commit & Pull Request Guidelines

Use concise Conventional Commit subjects, such as `fix: unify list page scrolling` or `chore: update project configuration`. Keep each commit scoped to one concern. PRs should state the affected service or route, describe behavior changes, link the issue when available, and include before/after screenshots for visible frontend work. Never commit credentials, local tool directories, build output, or generated artifacts unless regeneration is the intended change.
