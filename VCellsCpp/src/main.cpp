#include <iostream>
#include <string>
#include <opencv2/core/core.hpp>
#include <opencv2/imgcodecs.hpp>
#include <opencv2/highgui/highgui.hpp>
#include "tools.h"
#include "VCells.h"

using namespace std;
using namespace cv;

int main(int argc, char* argv[])
{
    if (argc != 2)
    {
        printf("Usage: ./VCells [src_img]\n");
        return -1;
    }

    string strSrcImg = argv[1];
    string strInitSegImg = "";
    string strDstImg = "";

    int iRet = GenerateImgFilenames(strSrcImg, strInitSegImg, strDstImg);
    if (iRet != 0)
    {
        printf("GenerateImgFilenames() failed, src image path: %s\n", strSrcImg.c_str());
        return iRet;
    }

    image = imread(strSrcImg, CV_LOAD_IMAGE_COLOR);   // Read the file
    if(!image.data)
    {
        printf("failed to open or find the image: %s\n", strSrcImg.c_str());
        return -1;
    }

    // image sizes
    width = image.size().width;
    height = image.size().height;


    /************************************************************************************/
    /************                     initialization                         ************/
    /************************************************************************************/
    struct pixel* pixelArray = (struct pixel*) malloc(sizeof(pixel)* height * width);
    struct centroid* generators = (struct centroid*) malloc(NUM_CLUSTER * sizeof(struct centroid));

    printf("Initializing the segments ...\n");
    initializePixel(pixelArray);
    initializeGenerators(generators, pixelArray);
    classicCVT(pixelArray,generators);
    cvtEnd = clock();
    cpu_time_used_initialize = ((double) (cvtEnd - start)) / CLOCKS_PER_SEC;
    printf("Initialization is done.\n\n");
    drawSketch(strInitSegImg, pixelArray, &borderColor);

    /************************************************************************************/
    /************                     VCells by EWCVT                        ************/
    /************************************************************************************/
    printf("VCells is running ...\n");
    EWCVT(generators,pixelArray);
    printf("Superpixels are ready ...\n\n");
    ewcvtEnd = clock();
    cpu_time_used_EWCVT = ((double) (ewcvtEnd - cvtEnd)) / CLOCKS_PER_SEC;

    /************************************************************************************/
    /************                     output the results                     ************/
    /************************************************************************************/
    printf("Time for initialization: %f\n", cpu_time_used_initialize);
    printf("Time used for VCells: %f\n", cpu_time_used_EWCVT);

    drawSketch(strDstImg, pixelArray, &borderColor);

    free(pixelArray);
    free(generators);

    return 0;
}