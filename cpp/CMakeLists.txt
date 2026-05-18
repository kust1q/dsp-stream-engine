cmake_minimum_required(VERSION 3.10)
project(dsp_stream_engine CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
# Important for static libraries linked into Go (shared or PIE executable)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

add_library(dsp STATIC dsp.cpp)
