#pragma once
/* Boot-thread only. Initialize all services, then start acquisition before networking.
 * Returns zero on success or negative errno; no network availability required. */
int app_run(void);
