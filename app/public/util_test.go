package public

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestUploadFile(t *testing.T) {
	fileName := "09a137d7-14cd-41ab-bf1f-a6e43bdf3f6a.jpg"
	uploadFileName := "testing-" + fileName
	host := "lab.huayuan-iot.com"
	author := "admin"
	project := "video"
	token := "b4da0ed0-d1b7-11e8-b75e-435a751a1801"
	user := "admin"
	file := "/tmp/jpg/" + fileName
	dlUrl, err := UploadFile(
		file,
		uploadFileName, host, author, project, token, user,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dlUrl)
}

func TestOnvifHTTPCapture(t *testing.T) {
	data, err := OnvifHTTPCapture("192.168.1.64", "onvif", "hyiot123")
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("a.jpg", data, 0666); err != nil {
		t.Fatal(err)
	}
}
