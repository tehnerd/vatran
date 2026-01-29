cmake_minimum_required(VERSION 3.9)
project(vatran VERSION 1.0.0 LANGUAGES C CXX)

# Option to build katran from submodule or use pre-built
option(BUILD_KATRAN_FROM_SOURCE "Build katran from submodule (requires folly, glog, etc)" OFF)

if(BUILD_KATRAN_FROM_SOURCE)
    # Set katran-specific variables before adding subdirectory
    set(KATRAN_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/katran")

    # Build katran from submodule
    # This brings in: katranlb, bpfadapter, chhelpers, iphelpers, pcapwriter, katransimulator
    add_subdirectory(external/katran)
else()
    # Use pre-built katran library
    # Expects KATRAN_INSTALL_DIR to point to katran installation
    # e.g., cmake -DKATRAN_INSTALL_DIR=/path/to/katran/install ..
    if(NOT DEFINED KATRAN_INSTALL_DIR)
        # Default: assume katran is built via build_katran.sh in _build/build directory
        set(KATRAN_INSTALL_DIR "${CMAKE_CURRENT_SOURCE_DIR}/_build/build" CACHE PATH "Path to katran installation")
    endif()

    set(KATRAN_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/katran")

    # Find required dependencies
    find_package(folly CONFIG REQUIRED)
    find_package(Glog REQUIRED)
    find_package(Gflags REQUIRED)

    # Import katran libraries
    foreach(lib katranlb bpfadapter chhelpers iphelpers pcapwriter katransimulator murmur3)
        add_library(${lib} STATIC IMPORTED)
        set_target_properties(${lib} PROPERTIES
            IMPORTED_LOCATION "${KATRAN_INSTALL_DIR}/katran/lib/lib${lib}.a"
        )
    endforeach()
endif()

# Include Vatran build rules
include(vatran.cmake)
