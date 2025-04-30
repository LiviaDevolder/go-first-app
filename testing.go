package poker

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func VerifyPlayerWin(t *testing.T, storage *PlayerFileStorageSystem, winner string) {
	t.Helper()

	if len(storage.league) != 1 {
		t.Fatal("was waiting for a win call but got nothing")
	}

	if storage.league[0].Name != winner {
		t.Errorf("did not record the correct value, result %s, expect %s", storage.league[0].Name, winner)
	}
}

func VerifyRequestBody(t *testing.T, result, expect string) {
	t.Helper()
	if result != expect {
		t.Errorf("result %s, expect %s", result, expect)
	}
}

func VerifyStatusCode(t *testing.T, result, expect int) {
	t.Helper()
	if result != expect {
		t.Errorf("wrong status code, result %d, expected %d", result, expect)
	}
}

func GetLeagueAnswer(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Cant parse %s player server answer: %v", body, err)
	}

	return
}

func VerifyLeague(t *testing.T, result, expected []Player) {
	t.Helper()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("obtido %v esperado %v", result, expected)
	}
}

func NewLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func VerifyContentType(t *testing.T, response *httptest.ResponseRecorder, expect string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != expect {
		t.Errorf("wrong response type, expect %s and got %s", expect, response.Result().Header)
	}
}

func VerifyScore(t *testing.T, result, expect int) {
	t.Helper()

	if result != expect {
		t.Errorf("result %d expect %d", result, expect)
	}
}

func CreateTmpFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("can not write temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func VerifyNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expect no error but got %v", err)
	}
}
