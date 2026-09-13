cmake_minimum_required(VERSION 3.20.0)
set(APP_ROOT ${CMAKE_CURRENT_LIST_DIR}/../..)
set(KCONFIG_ROOT ${APP_ROOT}/Kconfig)
find_package(Zephyr REQUIRED HINTS $ENV{ZEPHYR_BASE})
project(iot_core_tests C)
target_include_directories(app PRIVATE ${APP_ROOT}/include)
zephyr_linker_sources(SECTIONS ${APP_ROOT}/collectors.ld)
target_sources(app PRIVATE src/main.c
 ${APP_ROOT}/src/core/app_event.c ${APP_ROOT}/src/core/time_manager.c
 ${APP_ROOT}/src/core/extensions.c
 ${APP_ROOT}/src/collector/collector_manager.c
 ${APP_ROOT}/src/collector/temperature_collector.c
 ${APP_ROOT}/src/collector/imu_collector.c
 ${APP_ROOT}/src/collector/battery_collector.c
 ${APP_ROOT}/src/config/config_manager.c ${APP_ROOT}/src/config/remote_config.c
 ${APP_ROOT}/src/telemetry/telemetry_store.c ${APP_ROOT}/src/telemetry/telemetry_serializer.c)
