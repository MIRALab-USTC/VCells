cmake_minimum_required(VERSION 3.9)
project(VCellsCpp)

set(CMAKE_CXX_STANDARD 11)

include_directories(include src lib/opencv/include)

add_library(opencv STATIC IMPORTED include/tools.h)

find_package( OpenCV REQUIRED PATHS ${CMAKE_SOURCE_DIR}/lib/opencv/bin)
set(OpenCV_LIBS "opencv_core" "opencv_imgcodecs")

set(INC include/VCells.h include/tools.h)
set(SRC src/main.cpp src/tools.cpp)
set(EXECUTABLE_OUTPUT_PATH bin)

add_executable(VCellsCpp ${INC} ${SRC})
target_link_libraries(VCellsCpp ${OpenCV_LIBS})