cmake_minimum_required(VERSION 3.20.0)
find_package(Zephyr REQUIRED HINTS $ENV{ZEPHYR_BASE})
project(iot_scaffold LANGUAGES C)
set_property(TARGET app PROPERTY C_STANDARD 11)
target_include_directories(app PRIVATE include)
zephyr_linker_sources(SECTIONS collectors.ld)
file(GLOB collector_sources CONFIGURE_DEPENDS src/collector/*_collector.c)
target_sources(app PRIVATE src/main.c src/core/app.c src/core/app_event.c
  src/core/device_identity.c src/core/time_manager.c src/core/health_manager.c
  src/core/shell.c src/core/extensions.c src/collector/collector_manager.c
  src/telemetry/telemetry_store.c src/telemetry/telemetry_serializer.c
  src/telemetry/telemetry_uploader.c src/config/config_manager.c
  src/config/remote_config.c ${collector_sources})
target_sources_ifdef(CONFIG_APP_WIFI app PRIVATE src/network/wifi_manager.c)
target_sources_ifdef(CONFIG_APP_MQTT app PRIVATE src/network/mqtt_manager.c)
if(NOT CONFIG_APP_MQTT)
  target_sources(app PRIVATE src/network/mqtt_disabled.c)
endif()
if(CONFIG_APP_TLS)
  if(NOT APP_CA_CERT OR NOT EXISTS "${APP_CA_CERT}")
    message(FATAL_ERROR "TLS requires -DAPP_CA_CERT=/absolute/path/to/broker-ca.pem")
  endif()
  generate_inc_file_for_target(app ${APP_CA_CERT} ${ZEPHYR_BINARY_DIR}/include/generated/app_ca.inc)
endif()
