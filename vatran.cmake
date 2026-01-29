# CAPI source files
set(KATRAN_CAPI_SOURCES
    src/katran_capi.cpp
)

# CAPI header files
set(KATRAN_CAPI_HEADERS
    include/katran_capi.h
    include/katran_capi_types.h
)

# Static library
add_library(katran_capi_static STATIC ${KATRAN_CAPI_SOURCES})

target_include_directories(katran_capi_static
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
    PRIVATE
        # Include path to katran headers (from submodule)
        ${CMAKE_SOURCE_DIR}/external/katran
)

# Find libbpf in deps
find_library(LIBBPF_LIBRARY
    NAMES bpf
    PATHS ${CMAKE_CURRENT_SOURCE_DIR}/_build/deps/lib
          ${CMAKE_CURRENT_SOURCE_DIR}/_build/deps/usr/lib64
    NO_DEFAULT_PATH
)

if(NOT LIBBPF_LIBRARY)
    message(STATUS "libbpf not found in deps, using system libbpf")
    set(LIBBPF_LIBRARY bpf)
endif()

target_link_libraries(katran_capi_static
    PUBLIC
        -Wl,--start-group
        katranlb
        Folly::folly
        bpfadapter
        chhelpers
        murmur3
        iphelpers
        pcapwriter
        katransimulator
        ${LIBBPF_LIBRARY}
        elf
        z
        mnl
        ${GFLAGS}
        ${PTHREAD}
        -Wl,--end-group
)

set_target_properties(katran_capi_static PROPERTIES
    CXX_STANDARD 20
    CXX_STANDARD_REQUIRED ON
)

# Smoke test
add_executable(katran_capi_smoke src/katran_capi_smoke.c)
target_include_directories(katran_capi_smoke
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)
target_link_libraries(katran_capi_smoke katran_capi_static)
set_target_properties(katran_capi_smoke PROPERTIES LINKER_LANGUAGE CXX)

# Installation
install(FILES ${KATRAN_CAPI_HEADERS} DESTINATION include/katran/capi)
install(TARGETS katran_capi_static
    EXPORT katran-capi-exports
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
)

install(EXPORT katran-capi-exports
    FILE katran-capi-config.cmake
    DESTINATION lib/cmake/katran
)
