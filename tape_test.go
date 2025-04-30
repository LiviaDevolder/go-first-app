package main

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	// Arrange
	file, clean := createTmpFile(t, "12345")
	defer clean()

	tape := &tape{file}

	expect := "abc"
	tape.Write([]byte(expect))

	file.Seek(0, 0)
	newFileContent, _ := ioutil.ReadAll(file)

	// Act
	result := string(newFileContent)
	

	// Assert
	if result != expect {
		t.Errorf("result %s but expect %s", result, expect)
	}
}
