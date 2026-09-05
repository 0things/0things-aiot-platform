## Context

See `proposal.md` for motivation. Currently, the frontend codebase has two differing patterns for API consumption:
1. **Modern standard (e.g., `events` module)**: Re-exports Orval-generated query parameter types (`GetDeviceEventsParams`), uses generated query key builders (`getGetDeviceEventsQueryKey`), and delegates query execution to generated hooks (`useGetDeviceEvents`) with `select: (res) => res?.data`.
2. **Legacy pattern (e.g., `devices`, `products`, `device-groups`, `scene-linkage`)**: Manually defines query parameter interfaces, duplicates `queryKey` construction logic, invokes raw HTTP functions inside manual `useQuery` blocks, and performs unsafe type assertions `as unknown as ...`.

## Goals / Non-Goals

**Goals:**
- Unify all feature API data-fetching wrappers onto the generated Orval TanStack Query hooks and types.
- Ensure 100% parameter type fidelity with backend Swagger definitions (e.g., numeric IDs vs string conversions).
- Preserve existing component public interfaces where possible while eliminating manual type assertions.
- Resolve naming collisions between query hooks and UI dialog context providers.

**Non-Goals:**
- Changing backend Swagger definitions or endpoints.
- Modifying UI component layouts or design.
- Hand-editing any generated files in `src/api/generated/`.

## Decisions

### 1. Feature Query Wrapper Pattern
Every feature `api/queries.ts` will follow the canonical structure established in `events/api/queries.ts`:

```typescript
import { getGetDevicesQueryKey, useGetDevices } from '@/api/generated'
import type {
  DeviceListDevicesResponse,
  GetDevicesParams,
} from '@/api/generated/model'

export type { GetDevicesParams }
export type DevicesResponse = DeviceListDevicesResponse

export const deviceKeys = {
  all: ['devices'] as const,
  lists: () => [...deviceKeys.all, 'list'] as const,
  list: (params?: GetDevicesParams) => getGetDevicesQueryKey(params),
  // Additional keys...
}

export function useDevices(params?: GetDevicesParams) {
  return useGetDevices(params, {
    query: {
      select: (res) => res?.data,
    },
  })
}
```

*Rationale*: Wrapping generated hooks inside feature `api/queries.ts` provides a stable import boundary for components, automatically unwraps `res?.data`, and maintains support for custom selection/options when necessary.

### 2. Parameter Typing for Number vs String (e.g., `productId`)
In Swagger, query parameters like `productId` are integers (`number`). Legacy hand-written types sometimes allowed `string`.
*Decision*: Strict typing with `GetDevicesParams` (`productId?: number`). Consumer components extracting route string params will parse `Number(productId)` when constructing query objects.

### 3. Disambiguating Provider Hooks
*Decision*: In `src/features/devices/components/devices-provider.tsx`, rename `useDevices` context hook to `useDevicesDialog` (or export both) and update dialog consumers, preventing namespace clash between data query `useDevices` and dialog state `useDevicesDialog`.

## Risks / Trade-offs

- **[Risk] Strict parameter type checks causing TypeScript compile errors in components**
  → *Mitigation*: Run `pnpm -C frontend format && pnpm -C frontend build` iteratively across all modified feature modules to ensure zero build errors.
- **[Risk] Cache invalidation key mismatches**
  → *Mitigation*: Ensure mutation hooks (`useCreateDevice`, `useUpdateProduct`, etc.) invalidate using the generated query key arrays (e.g. `getGetDevicesQueryKey()`).

