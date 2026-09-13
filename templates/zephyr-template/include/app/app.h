#pragma once

/**
 * @brief Initialize all system services and run the application lifecycle.
 *
 * Must be called from the main/boot thread only. Initializes core subsystems,
 * loads persisted configuration, starts sensor collectors, and connects networking.
 *
 * @return 0 on success, negative errno code on failure.
 */
int app_run(void);

