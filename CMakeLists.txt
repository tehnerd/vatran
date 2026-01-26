# Copyright (C) 2018-present, Facebook, Inc.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 2 of the License.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License along
# with this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

cmake_minimum_required(VERSION 3.5)

# C API source files
set(KATRAN_CAPI_SOURCES
    src/katran_capi.cpp
)

# C API header files
set(KATRAN_CAPI_HEADERS
    include/katran_capi.h
    include/katran_capi_types.h
)

# Shared library
#add_library(katran_capi SHARED ${KATRAN_CAPI_SOURCES})

# Static library
add_library(katran_capi_static STATIC ${KATRAN_CAPI_SOURCES})

# Configure both targets
foreach(target katran_capi_static)
    target_include_directories(${target}
        PUBLIC
            $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
            $<INSTALL_INTERFACE:include>
        PRIVATE
            ${KATRAN_INCLUDE_DIR}
            ${FOLLY_INCLUDE_DIR}
    )

    target_link_libraries(${target}
        PRIVATE
            katranlb
            Folly::folly
    )

    set_target_properties(${target} PROPERTIES
        CXX_STANDARD 20
        CXX_STANDARD_REQUIRED ON
    )
endforeach()

# Shared library versioning
#set_target_properties(katran_capi PROPERTIES
#    VERSION 1.0.0
#    SOVERSION 1
#)

# Installation rules
install(
    FILES ${KATRAN_CAPI_HEADERS}
    DESTINATION include/katran/capi
)

install(
    TARGETS katran_capi_static
    EXPORT katran-capi-exports
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
    RUNTIME DESTINATION bin
)

install(
    EXPORT katran-capi-exports
    FILE KatranCapiTargets.cmake
    NAMESPACE katran::
    DESTINATION lib/cmake/katran
)
