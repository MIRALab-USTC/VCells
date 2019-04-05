# VCells

*VCells* is the implementation of VCells algorithms described in [VCells: Simple and Efficient Superpixels Using Edge-Weighted Centroidal Voronoi Tessellations, TPAMI 34(6), 1241-1247" by Jie Wang and Xiaoqiang Wang](http://www-personal.umich.edu/~jwangumi/publications/VCells.pdf).

## Description
This project has 2 components:

```
VCells
├── VCellsCpp
├── VCellsServer
```
### VCellsCpp
It adopts the code of main algorithm in original implementation: [VCells.cpp](http://staff.ustc.edu.cn/~jwangx/software.html). We've replaced the image reading and writing part with OpenCV's component, so now it can support most image format. Besides, it now supports command line to specify which image to operate.

### VCellsServer
In order to provide a friendly user interface, we build a simple web application, where Golang is used to build a simple server and jQuery and Bootstrap are used to design the interface.


## Build
### VCellsCpp
We use CMake in our project, so you have to generate makefile with CMake and then build target with makefile.

#### OpenCV
We use [OpenCV](https://opencv.org/) in our project. So, you should build and install OpenCV on your machine.
You can also download a release of VCellsCpp on our [Release Page](https://github.com/MIRALab-USTC/VCells/releases).

#### Generate makefile
```
cd VCellsCpp
mkdir cmake-build-debug && cd cmake-build-debug
cmake ..
```
You will see a makefile generated in directory `cmake-build-debug`.

#### Make
```
make
```
This will generate an executable file.

### VCellsServer

```
cd VCellsServer
go build main.go file_upload.go
```

And you will get an executable file: `main`

## Usage
You can use the VCellsCpp part independently as a command line program. Or, you can run a simple server on your machine and interact with it from web browser.

### VCellsCpp
VCellsCpp serves as a command line program. Its usage is as follows:

```
./VCellsCpp [srcImgPath]
```

It will generate a `initial segmented image` and a `superpixel image` in the same directory of source image.


### VCellsServer
Before you run the server, you should copy the executable `VCellsCpp` to the directory `VCellsServer/upload_files/`.

You can specify the port that the server listens to and the log file path. The default port will be set to `8080`,
and the default log file path will be set to `./vcells.log`.

```
./main -port [port] -log [log]
```

## License
[MIT license](https://opensource.org/licenses/MIT)
