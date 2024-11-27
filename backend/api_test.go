package backend

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"testing"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func setup() chan struct{} {
	gin.SetMode(gin.ReleaseMode)

	db, err := SetupDatabase(":memory:")
	if err != nil {
		log.Fatal(err.Error())
	}

	channel := make(chan struct{})
	go func() {
		defer CloseDatabase(db)
		SignalRoute(db, channel)
	}()

	return channel
}

func postRequest(path string, data string) (string, int, error) {
	buffer := bytes.NewBuffer([]byte(data))

	url := "http://localhost:8080/api/" + path
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return "", 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(body[:]), resp.StatusCode, nil
}

func TestRegisterRequirements(t *testing.T) {
	channel := setup()
	body, status, err := postRequest("register", `{"username":"","password":""}`)

	if err != nil {
		t.Fatal(err)
	}

	if status == 200 {
		t.Fatalf("False success at null uname/pass. %s", body)
	}

	body, status, err = postRequest("register", `{"username":"testing","password":""}`)

	if err != nil {
		t.Fatal(err)
	}

	if status == 200 {
		t.Fatalf("False success at null pass. %s", body)
	}

	body, status, err = postRequest("register", `{"username":"","password":"testing"}`)

	if err != nil {
		t.Fatal(err)
	}

	if status == 200 {
		t.Fatalf("False success at null user. %s", body)
	}

	body, status, err = postRequest("register", `{"username":"testing","password":"attempt_"}`)

	if err != nil {
		t.Fatal(err)
	}

	if status != 200 {
		t.Fatalf("False failure at proper register. %s", body)
	}

	channel <- struct{}{}
}

func TestLogin(t *testing.T) {
	channel := setup()

	body, status, err := postRequest("register", `{"username":"testing","password":"attempt_"}`)

	if err != nil {
		log.Fatal(err)
	}

	if status != 200 {
		log.Fatalf("Failed register. %s", body)
	}

	body, status, err = postRequest("login", `{"username":"testing","password":"badpass_"}`)

	if err != nil {
		log.Fatal(err)
	}

	if status == 200 {
		log.Fatalf("Wrong acceptance, different password. %s", body)
	}

	body, status, err = postRequest("login", `{"username":"baduser","password":"attempt"}`)

	if err != nil {
		log.Fatal(err)
	}

	if status == 200 {
		log.Fatalf("Wrong acceptance, different username. %s", body)
	}

	body, status, err = postRequest("login", `{"username":"testing","password":"attempt_"}`)

	if err != nil {
		log.Fatal(err)
	}

	if status != 200 {
		log.Fatalf("Failed proper login. %s", body)
	}

	channel <- struct{}{}
}
