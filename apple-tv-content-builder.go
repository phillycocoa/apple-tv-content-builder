package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
)

// Response models the response we get when asking for /me/videos
type Response struct {
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	VideoRecords []VideoRecord `json:"data"`
}

// CountOfVideoRecords returns the number of video records in the response
func (response Response) CountOfVideoRecords() int {
	return len(response.VideoRecords)
}

// VideoRecord models the videos of our video collection in Response
type VideoRecord struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URI         string         `json:"uri"`
	Pictures    PictureCatalog `json:"pictures"`
	VideoFiles  []VideoFile    `json:"files"`
}

// PictureCatalog models the pictures for a video record
type PictureCatalog struct {
	URI   string        `json:"uri"`
	Sizes []PictureSize `json:"sizes"`
}

// ThumbnailSize return our prefered PictureSize when needing a thumbnail
func (pc PictureCatalog) ThumbnailSize() PictureSize {
	for _, element := range pc.Sizes {
		if element.Height == 360 {
			return element
		}
	}
	return pc.Sizes[0]
}

// PictureSize models the multiple size offerings for a picture
type PictureSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

// HDFile returns the file of files we pref when needing to link to the HD version of this video
func (vr VideoRecord) HDFile() VideoFile {
	for _, element := range vr.VideoFiles {
		if element.Height == 1080 {
			return element
		}
	}
	return vr.VideoFiles[0]
}

// VideoFile models the perminent file urls for the video
type VideoFile struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

// Basic error checking function: check for error, if present panic
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// Build a http request with auth headers
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.vimeo.com/me/videos", nil)
	req.Header.Set("Authorization", "bearer 4a19abe051f1f17f6dc5b5204f2d461c")
	res, err := client.Do(req)
	checkError(err)

	// parse the response body into a byte slice
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	checkError(err)

	// unmarshal the response bytes from JSON into a struct
	var r Response
	err = json.Unmarshal(bodyBytes, &r)
	checkError(err)

	// FIXME: Need to find a nicer way than hardcoding this path
	PROJECTFOLDER := "/Users/zorn/go/src/github.com/phillycocoa/apple-tv-content-builder/"

	// grab out template file
	filepath := PROJECTFOLDER + "main-template.tmpl"
	contents, err := ioutil.ReadFile(filepath)
	checkError(err)

	// build a new template object with the contents of our file
	t := template.New("main template")
	t, err = t.Parse(string(contents[:]))
	checkError(err)

	// process the template given the JSON struct
	var templateOutput bytes.Buffer
	err = t.Execute(&templateOutput, r)
	checkError(err)

	// take the result of the template processing and write it to a file
	filepathOut := PROJECTFOLDER + "out/PCHTemplate.xml.js"
	err = ioutil.WriteFile(filepathOut, templateOutput.Bytes(), 0644)
	checkError(err)

}

// Previous code we used to build a JSON response
// func sampleJSONResponse() []byte {
// 	v1 := VideoRecord{"123"}
// 	v2 := VideoRecord{"456"}
// 	r := Response{11, 1, []videoRecord{v1, v2}}
// 	b, _ := json.Marshal(r)
// 	return b
// }
