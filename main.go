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

func main() {
    var args = os.Args[1:]

    numArgs := len(args)

    pageWidth := 8.5
    pageHeight := 11.0

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

                dpi := 96.0

                fmt.Printf("Image Height: %f   Image Width: %f\n", height, width)

                totalHorizontalPixels := dpi * pageWidth
                totalVerticalPixels := dpi * pageHeight

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

                // okay, cool, we have our margins and the proper image dimensions, so now try to render the pdf page
                _ = xMargin
                _ = yMargin

                fmt.Printf("xMargin: %f   yMargin: %f\n", xMargin, yMargin )

                fmt.Printf("Adding Image to Page ...\n")

                pdf.AddPage()
                pdf.Image(args[i], 0, 0, 200, 0, false, "", 0, "")

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