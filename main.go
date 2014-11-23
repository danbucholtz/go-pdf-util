package main


import (
    "code.google.com/p/gofpdf"
    "fmt"
    "image"
    "os"
    "path/filepath"

    _ "image/jpeg"
    _ "image/png"
)

func pixelsToMM(pixelValue float64, dpi int) float64{
    millimetersPerInch := 25.4

    return pixelValue * (millimetersPerInch/ float64(dpi))
}

func main() {
    var args = os.Args[1:]

    numArgs := len(args)

    pageWidth := 8.5
    pageHeight := 11.0
    dpi := 200

    if numArgs == 0 {
    	fmt.Println("Invalid usage: Missing arguments")
    } else{
        // initialize the PDF
        pdf := gofpdf.New("P", "mm", "A4", "")
    	
        for i:=0; i < numArgs; i++ {
    		fmt.Println("Opening File: " + args[i]);
    		
    		reader, err := os.Open(args[i]);
    		
    		if err == nil {
    			// go ahead and read in the config data for the image
    			
    			im, _, err := image.DecodeConfig(reader);
    			
    			if err != nil {
    				fmt.Println("Failed to decode image config: " + err.Error())
    				continue
    			}

                if err != nil {
                    fmt.Println("Failed to decode image: " + err.Error())
                    continue
                }
                

                // w00t, we read in the image.  Close the reader at some point in the future since we're done with it
                defer reader.Close()

    			fmt.Printf("Page #%d ... START\n", i + 1)

                // use floats instead of ints
                height := float64(im.Height)
                width := float64(im.Width)

                fmt.Printf("Image Height: %f   Image Width: %f\n", height, width)

                totalHorizontalPixels := float64(dpi) * pageWidth
                totalVerticalPixels := float64(dpi) * pageHeight

                fmt.Printf("totalHorizontalPixels: %f   totalVerticalPixels: %f\n", totalHorizontalPixels, totalVerticalPixels)

                isPortraitImage := height >= width

                // resize the image if it is bigger than the page at a given dpi

                if isPortraitImage {
                    fmt.Printf("Image is Portrait\n")
                    if height > totalVerticalPixels {

                        scaleFactor := totalVerticalPixels/height

                        height = height * scaleFactor
                        width = width * scaleFactor

                        fmt.Printf("Scaling necessary...   scaleFactor: %f   height: %f   width %f\n", scaleFactor, height, width)
                    }
                } else{
                    fmt.Printf("Image is Landscape\n")
                    if width > totalHorizontalPixels {

                        scaleFactor := totalHorizontalPixels/width

                        height = height * scaleFactor
                        width = width * scaleFactor

                        fmt.Printf("Scaling necessary...   scaleFactor: %f   height: %f   width %f\n", scaleFactor, height, width)
                    }
                }
                
                // figure out margins
                widthDifferenceInPixels := totalHorizontalPixels - width
                heightDifferenceInPixels := totalVerticalPixels - height

                fmt.Printf("widthDifferenceInPixels: %f   heightDifferenceInPixels: %f\n", widthDifferenceInPixels, heightDifferenceInPixels )

                xMargin := 0.0
                yMargin := 0.0

                if widthDifferenceInPixels > 0 {
                    xMargin = widthDifferenceInPixels / 2.0
                }

                if heightDifferenceInPixels > 0 {
                    yMargin = heightDifferenceInPixels / 2.0
                }

                fmt.Printf("xMargin: %f   yMargin: %f\n", xMargin, yMargin )

                // convert values to mm, since that's what gofpdf is using

                xMarginMM := pixelsToMM(xMargin, dpi)
                yMarginMM := pixelsToMM(yMargin, dpi)
                widthMM := pixelsToMM(width, dpi)
                heightMM := pixelsToMM(height, dpi)

                fmt.Printf("Converting values from pixels to MM \n")

                fmt.Printf("xMarginMM: %f   yMarginMM: %f   widthMM: %f   heightMM: %f   \n", xMarginMM, yMarginMM, widthMM, heightMM )

                fmt.Printf("Adding Image to Page ...\n")

                pdf.AddPage()
                pdf.Image(args[i], xMarginMM, yMarginMM, widthMM, heightMM, false, "", 0, "")

                fmt.Printf("Adding Image to Page ... DONE\n")
                
    		} else{
    			// abort
    			fmt.Println("Failed to open file: " + err.Error());
    		}
    	}
        fileStr := filepath.Join(".", "pdf/output.pdf")
        pdfErr := pdf.OutputFileAndClose(fileStr)
        if pdfErr == nil {
                fmt.Println("Successfully generated pdf/output.pdf")
        } else {
                fmt.Println("Error occurred outputting PDF: " + pdfErr.Error())
        }
    }
}