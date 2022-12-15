package opencv_ifc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
)

var client *http.Client

var openCVServerAddr = "http://localhost:8090"

func CreateClient() {
	//setup a mocked http client.
	println("")
	println("Setup HTTP client")
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			panic(err)
		}
		println(fmt.Sprintf("%s", b))
	}))
	defer ts.Close()
	client = ts.Client()
}

func SetServerAddress(serverAddress string) {
	openCVServerAddr = serverAddress
}

func SendImageToImageServer(img *image.Image) string {
	//println("Encoding new frame")
	// Convert image to jpg and obtain the bytes
	var imageBuf bytes.Buffer
	_ = jpeg.Encode(&imageBuf, *img, nil)

	// Prepare the reader instances to encode
	values := map[string]io.Reader{
		"file": bytes.NewReader(imageBuf.Bytes()),
	}

	// Upload and get back the json response
	resp, err := Upload(client, openCVServerAddr, values)
	if err != nil {
		println("Response error!")
		return ""
	}

	// Return json string
	//println("Response received: " + resp)
	return resp
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (response string, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return "", err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return "", err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return "", err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	//println("Encoded data")
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		println("Error when performing HTTP request")
		return "", err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	//println("POSTing...")
	res, err := client.Do(req)
	if err != nil {
		println(err.Error())
		return "", err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		println(fmt.Errorf("bad status: %s", res.Status))
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", err
}
