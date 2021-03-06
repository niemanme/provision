package frontend

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/digitalrebar/digitalrebar/go/common/store"
	"github.com/digitalrebar/provision/backend"
)

func TestLeaseList(t *testing.T) {
	localDTI := testFrontend()

	localDTI.ListValue = nil
	req, _ := http.NewRequest("GET", "/api/v3/leases", nil)
	w := localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusOK)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")

	var list []backend.Lease
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 0 {
		t.Errorf("Response should be an empty list, but got: %d\n", len(list))
	}

	localDTI.ListValue = []store.KeySaver{&backend.Lease{Addr: net.ParseIP("1.1.1.1")}}
	req, _ = http.NewRequest("GET", "/api/v3/leases", nil)
	w = localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusOK)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")

	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Errorf("Response should be a list of 1, but got: %d\n", len(list))
	}
	if !list[0].Addr.Equal(net.ParseIP("1.1.1.1")) {
		t.Errorf("Response[0] is not named 1.1.1.1, %v\n", list[0].Addr)
	}
}

func TestLeasePost(t *testing.T) {
	localDTI := testFrontend()

	localDTI.CreateValue = nil
	localDTI.CreateError = &backend.Error{Code: 23, Messages: []string{"should not get this"}}
	req, _ := http.NewRequest("POST", "/api/v3/leases", nil)
	req.Header.Set("Content-Type", "text/html")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "Invalid content type: text/html")

	localDTI.CreateValue = nil
	localDTI.CreateError = &backend.Error{Code: 23, Messages: []string{"should not get this"}}
	req, _ = http.NewRequest("POST", "/api/v3/leases", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "EOF")

	localDTI.CreateValue = nil
	localDTI.CreateError = &backend.Error{Code: 23, Messages: []string{"should not get this"}}
	req, _ = http.NewRequest("POST", "/api/v3/leases", strings.NewReader("asgasgd"))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "invalid character 'a' looking for beginning of value")

	/* GREG: handle json failure? bad IP */

	localDTI.CreateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.CreateError = fmt.Errorf("this is a test: bad fred")
	v, _ := json.Marshal(localDTI.CreateValue)
	req, _ = http.NewRequest("POST", "/api/v3/leases", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "this is a test: bad fred")

	localDTI.CreateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.CreateError = nil
	v, _ = json.Marshal(localDTI.CreateValue)
	req, _ = http.NewRequest("POST", "/api/v3/leases", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusCreated)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")

	var be backend.Lease
	json.Unmarshal(w.Body.Bytes(), &be)
	if !be.Addr.Equal(net.ParseIP("1.1.1.1")) {
		t.Errorf("Returned Lease was not correct: %v %v\n", "1.1.1.1", be.Addr)
	}

	localDTI.CreateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.CreateError = &backend.Error{Code: 23, Messages: []string{"test one"}}
	v, _ = json.Marshal(localDTI.CreateValue)
	req, _ = http.NewRequest("POST", "/api/v3/leases", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w = localDTI.RunTest(req)
	localDTI.ValidateCode(t, 23)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	var berr backend.Error
	json.Unmarshal(w.Body.Bytes(), &berr)
	if berr.Messages[0] != "test one" {
		t.Errorf("Returned Error was not correct: %v %v\n", "test one", berr.Messages[0])
	}
}

func TestLeaseGet(t *testing.T) {
	localDTI := testFrontend()

	localDTI.GetValue = nil
	localDTI.GetBool = false
	req, _ := http.NewRequest("GET", "/api/v3/leases/1.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusNotFound)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "leases GET: 01010101: Not Found")

	localDTI.GetValue = nil
	localDTI.GetBool = false
	req, _ = http.NewRequest("GET", "/api/v3/leases/a.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "lease get: address not valid: a.1.1.1")

	localDTI.GetValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.GetBool = true
	req, _ = http.NewRequest("GET", "/api/v3/leases/1.1.1.1", nil)
	w := localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusOK)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	var be backend.Lease
	json.Unmarshal(w.Body.Bytes(), &be)
	if !be.Addr.Equal(net.ParseIP("1.1.1.1")) {
		t.Errorf("Returned Lease was not correct: %v %v\n", "1.1.1.1", be.Addr)
	}
}

/*
func TestLeasePatch(t *testing.T) {
	localDTI := testFrontend()

	localDTI.UpdateValue = nil
	localDTI.UpdateError = nil
	req, _ := http.NewRequest("PATCH", "/api/v3/leases/1.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusNotImplemented)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "lease patch: NOT IMPLEMENTED")
}
*/

func TestLeasePut(t *testing.T) {
	localDTI := testFrontend()

	localDTI.UpdateValue = nil
	localDTI.UpdateError = &backend.Error{Code: 23, Messages: []string{"should not get this one"}}
	req, _ := http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", nil)
	req.Header.Set("Content-Type", "text/html")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "Invalid content type: text/html")

	localDTI.UpdateValue = nil
	localDTI.UpdateError = &backend.Error{Code: 23, Messages: []string{"should not get this one"}}
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "EOF")

	localDTI.UpdateValue = nil
	localDTI.UpdateError = &backend.Error{Code: 23, Messages: []string{"should not get this one"}}
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader("asgasgd"))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "invalid character 'a' looking for beginning of value")

	localDTI.UpdateValue = nil
	localDTI.UpdateError = &backend.Error{Code: 23, Messages: []string{"should not get this one"}}
	req, _ = http.NewRequest("PUT", "/api/v3/leases/a.1.1.1", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "lease put: address not valid: a.1.1.1")

	/* GREG: handle json failure? hard to do - send a machine instead of lease */

	localDTI.UpdateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.UpdateError = fmt.Errorf("this is a test: bad 1.1.1.1")
	v, _ := json.Marshal(localDTI.UpdateValue)
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "this is a test: bad 1.1.1.1")

	localDTI.UpdateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.2"), Token: "kfred"}
	localDTI.UpdateError = nil
	v, _ = json.Marshal(localDTI.UpdateValue)
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "leases PUT: Key change from 01010101 to 01010102 not allowed")

	localDTI.UpdateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.UpdateError = &backend.Error{Code: 23, Type: "API_ERROR", Messages: []string{"test one"}}
	v, _ = json.Marshal(localDTI.UpdateValue)
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, 23)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "test one")

	localDTI.UpdateValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.UpdateError = nil
	v, _ = json.Marshal(localDTI.UpdateValue)
	req, _ = http.NewRequest("PUT", "/api/v3/leases/1.1.1.1", strings.NewReader(string(v)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusOK)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	var be backend.Lease
	json.Unmarshal(w.Body.Bytes(), &be)
	if be.Token != "kfred" {
		t.Errorf("Returned Lease was not correct: %v %v\n", "kfred", be.Token)
	}
}

func TestLeaseDelete(t *testing.T) {
	localDTI := testFrontend()

	localDTI.RemoveValue = nil
	localDTI.RemoveError = &backend.Error{Code: 23, Type: "API_ERROR", Messages: []string{"should get this one"}}
	req, _ := http.NewRequest("DELETE", "/api/v3/leases/1.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, 23)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "should get this one")

	localDTI.RemoveValue = nil
	localDTI.RemoveError = fmt.Errorf("this is a test: bad 1.1.1.1")
	req, _ = http.NewRequest("DELETE", "/api/v3/leases/1.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusNotFound)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "this is a test: bad 1.1.1.1")

	localDTI.RemoveValue = nil
	localDTI.RemoveError = fmt.Errorf("this is a test: bad 1.1.1.1")
	req, _ = http.NewRequest("DELETE", "/api/v3/leases/a.1.1.1", nil)
	localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusBadRequest)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	localDTI.ValidateError(t, "API_ERROR", "lease delete: address not valid: a.1.1.1")

	localDTI.RemoveValue = &backend.Lease{Addr: net.ParseIP("1.1.1.1"), Token: "kfred"}
	localDTI.RemoveError = nil
	req, _ = http.NewRequest("DELETE", "/api/v3/leases/1.1.1.1", nil)
	w := localDTI.RunTest(req)
	localDTI.ValidateCode(t, http.StatusOK)
	localDTI.ValidateContentType(t, "application/json; charset=utf-8")
	var be backend.Lease
	json.Unmarshal(w.Body.Bytes(), &be)
	if be.Token != "kfred" {
		t.Errorf("Returned Lease was not correct: %v %v\n", "kfred", be.Token)
	}
}
