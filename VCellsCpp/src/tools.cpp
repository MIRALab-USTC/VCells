
#include "tools.h"

int GenerateImgFilenames(const string& strSrcImg, string& strInitSegImg, string& strDstImg)
{
    for (unsigned int i=strSrcImg.length()-1; i >1; --i)
    {
        if (strSrcImg[i] == '.')
        {
            strInitSegImg = strSrcImg;
            strDstImg = strSrcImg;
            strInitSegImg.insert(i, "_init_seg");
            strDstImg.insert(i, "_superpixel");
            return 0;
        }
    }

    strInitSegImg = "";
    strDstImg = "";

    printf("invalid image path: %s\n", strSrcImg.c_str());
    return -1;
}