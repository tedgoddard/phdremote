// Copyright 2014 Ted Goddard. All rights reserved.
// Use of this source code is governed the Apache 2.0
// license that can be found in the LICENSE file.

package fits

import "bytes"
import "fmt"
import "io"
import "bufio"
import "math"
import "image"
import "image/color"
import "image/png"
import "encoding/binary"
import "unsafe"


/*
#cgo LDFLAGS: -lcfitsio

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "fitsio.h"

double* fit2png(char *fileName, int* width, int* height, int* len)  {

    fitsfile *fptr;   // FITS file pointer, defined in fitsio.h
    int status = 0;   // CFITSIO status value MUST be initialized to zero!
    int bitpix, naxis, ii, anynul;
    long naxes[2] = {1,1}, fpixel[2] = {1,1};
    double *pixels;
    double *pixelss;
    char format[20], hdformat[20];

    //fileName format
    //image.fits                    -  the whole image
    //image.fits[100:110,400:410]   -  a section
    //table.fits[2][bin (x,y) = 32] -  the pixels in
    //    an image constructed from a 2D histogram of X and Y
    //    columns in a table with a binning factor = 32

    if (!fits_open_file(&fptr, fileName, READONLY, &status))  {
        if (!fits_get_img_param(fptr, 2, &bitpix, &naxis, naxes, &status))  {
          if (naxis > 2 || naxis == 0)
             printf("Error: only 1D or 2D images are supported\n");
          else
          {
            // get memory for 1 row
            pixels = (double *) malloc(naxes[0] * sizeof(double));
            // get memory for complete image
           *len = naxes[0] * naxes[1] * sizeof(double);
           pixelss = (double *) malloc(*len);

            if (pixels == NULL) {
                printf("Memory allocation error\n");
                return(NULL);
            }

            *width = naxes[0];
            *height = naxes[1];
            // loop over all the rows in the image, top to bottom
            for (fpixel[1] = naxes[1]; fpixel[1] >= 1; fpixel[1]--)  {

               if (fits_read_pix(fptr, TDOUBLE, fpixel, naxes[0], NULL,
                    pixels, NULL, &status) )  {
                    if (status) fits_report_error(stderr, status);
                    printf("    at %d %d\n", fpixel[0], naxes[0]);
                }

                if (fits_read_pix(fptr, TDOUBLE, fpixel, naxes[0], NULL,
                    &pixelss[*width * (naxes[1] - fpixel[1])], NULL, &status) ) {
                    if (status) fits_report_error(stderr, status);
                    printf("    at %d %d\n", fpixel[0], naxes[0]);
                }

            }
            free(pixels);
          }
        }
        fits_close_file(fptr, &status);
    } 

    if (status) fits_report_error(stderr, status); // print any error message

    return(pixelss);
}
*/
import "C"


type GrayFloat64 struct  {
    //maximum pixel value to allow scaling to 16 bits
    MaxPixel float64
    Pix []float64
    Stride int
    Rect image.Rectangle
}

func (gray64 GrayFloat64) ColorModel() color.Model  {
    return color.Gray16Model
}

func (gray64 GrayFloat64) Bounds() image.Rectangle  {
    return gray64.Rect
}

func (gray64 GrayFloat64) At(x, y int) color.Color  {
    offset := (y - gray64.Rect.Min.Y) * gray64.Stride +
        (x - gray64.Rect.Min.X)
    pixelFloat := gray64.Pix[offset]
    grayValue := color.Gray16{
            uint16(math.MaxUint16 * pixelFloat / gray64.MaxPixel)}
    return grayValue
}

func Convert(fileName string, imageWriter io.Writer)  {

    var widthC C.int
    var heightC C.int
    var lenC C.int
    greyCBytes := C.fit2png(C.CString(fileName), &widthC, &heightC, &lenC)
    greyBytes := C.GoBytes(unsafe.Pointer(greyCBytes), lenC)

    width := int(widthC)
    height := int(heightC)
    len := int(lenC)

    maxPixel := 0.0
    var pixel float64
    floatLen := len / 8
    buf := bytes.NewReader(greyBytes)
    var greyFloats = make([]float64, floatLen)
    for i := 0; i < floatLen; i++ {
        err := binary.Read(buf, binary.LittleEndian, &pixel)
        if err != nil {
            fmt.Println("binary.Read failed:", err)
        }
        greyFloats[i] = pixel
        if (pixel > maxPixel)  {
            maxPixel = pixel
        }
    }

    greyImage := GrayFloat64{
        MaxPixel:maxPixel,
        Pix:greyFloats,
        Stride:width,
        Rect:image.Rect(0, 0, width, height),
    }

    C.free(unsafe.Pointer(greyCBytes))

    bufWriter := bufio.NewWriter(imageWriter)
    png.Encode(bufWriter, greyImage)
    bufWriter.Flush()
}
