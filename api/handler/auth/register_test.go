package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	_, err := http.NewRequest("POST", "localhost:8080/auth/register", nil)
	if err != nil {
		t.Errorf("failed %v", err)
	}

	rec := httptest.NewRecorder()

	t.Error(rec)

	// b, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Errorf("failed %v", err)
	// }

	// fmt.Println(b)

}
