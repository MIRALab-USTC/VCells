
#ifndef VCELLSCPP_TOOLS_H
#define VCELLSCPP_TOOLS_H

#include <string>
using namespace std;

// Generate filenames of Init Segmented Image and Generated SuperPixel Image.
// The passed-in source image filename is supposed to have *file extension*,
// such as .png, .jpg, etc.
//
// Params:
//     strSrcImg:     passed-in filename of source image
//     strInitSegImg: passed-out filename of init segmented image
//     strDstImg:     passed-out filename of generated superpixel image
//
// Return:
//     0 if success
//     non-zero if failure
//
// Example:
//      strSrcImg:     demo.png
//      strInitSegImg: demo_init_seg.png
//      strDstImg:     demo_superpixel.png

int GenerateImgFilenames(const string& strSrcImg, string& strInitSegImg, string& strDstImg);


#endif //VCELLSCPP_TOOLS_H
