package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arkadiuszekprogramista/bookingapp/internal/models"
)

type reqBody struct {
	StartDate string
	EndDate string
	FirstName string
	LastName string
	Email string
	Phone string
	//RoomId is  not int because its for parms in url addres 
	RoomID string 
}

type postData struct {
	key string
	value  string
}

var theTest = []struct {
	name string
	url string
	method string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals", "/generals-quarters", "GET", http.StatusOK},
	{"majors", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// {"post-search","/search-availability","POST",[]postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-12"},
	// }, http.StatusOK},
	// {"post-search-json","/search-availability-json","POST",[]postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-12"},
	// }, http.StatusOK},
	// {"post-make-reservation","/make-reservation","POST",[]postData{
	// 	{key: "first_name", value: "Arka"},
	// 	{key: "last_name", value: "Las"},
	// 	{key: "email", value: "email@email.com"},
	// 	{key: "phone", value: "123-123-123"},
	// },http.StatusOK},
}

var TestForPostMethods = []struct {
	name string
	url string
	method string
	expectedStatusCode int
	reqBody reqBody
}{
	{"checking if all data are valid", "/make-reservation", "POST", http.StatusSeeOther, reqBody{
		StartDate: "2050-01-01",
		EndDate: "2050-01-02",
		FirstName: "John",
		LastName: "Smith",
		Email: "john@email.com",
		Phone: "233 111 333",
		RoomID: "3",
		},
	},
	{"invalid start date", "/make-reservation", "POST", http.StatusTemporaryRedirect, reqBody{
		StartDate: "invalid",
		EndDate: "2050-01-02",
		FirstName: "John",
		LastName: "Smith",
		Email: "john@email.com",
		Phone: "233 111 333",
		RoomID: "3",
		},
	},
	{"invalid end date", "/make-reservation", "POST", http.StatusTemporaryRedirect, reqBody{
		StartDate: "2050-01-01",
		EndDate: "invalid",
		FirstName: "John",
		LastName: "Smith",
		Email: "john@email.com",
		Phone: "233 111 333",
		RoomID: "3",
		},
	},
	{"invalid room id", "/make-reservation", "POST", http.StatusTemporaryRedirect, reqBody{
		StartDate: "invalid",
		EndDate: "2050-01-02",
		FirstName: "John",
		LastName: "Smith",
		Email: "john@email.com",
		Phone: "233 111 333",
		RoomID: "invalid",
		},
	},

}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {
		resp, err :=  ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}


func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 3,
		Room: models.Room{
			ID: 3,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "./make-reservation", nil)
	ctx := getCtx(req)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everething)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test whit non-existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 100

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_PostReservation(t *testing.T) {

	for _, e := range TestForPostMethods {

		var reqBody reqBody
		body := reqBody.getUrlWihtParams(e.reqBody.StartDate, e.reqBody.EndDate, e.reqBody.FirstName, e.reqBody.Email, e.reqBody.LastName, e.reqBody.Phone, e.reqBody.RoomID)

		req, _ := http.NewRequest(e.method,e.url, strings.NewReader(body))
		ctx := getCtx(req)
		req = req.WithContext(ctx)
	
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	
		rr := httptest.NewRecorder()
	
		handler := http.HandlerFunc(Repo.PostReservation)
		
		handler.ServeHTTP(rr, req)
	
		if rr.Code != e.expectedStatusCode {
			t.Errorf("PostReservation handler test %s returned wrong response code: got %d, wanted %d",e.name, rr.Code, e.expectedStatusCode)
		}
	}

	//test for missing body
	req, _ := http.NewRequest("POST","/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test inserction reservation
	var reqBody reqBody
	body := reqBody.getUrlWihtParams("2050-01-01", "2050-01-02","Johny","Smith","email@email.com", "123131 31313  133", "2")

	req, _ = http.NewRequest("POST","/make-reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for test insercion reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test inserction room restriction
	body = reqBody.getUrlWihtParams("2050-01-01", "2050-01-02","Johny","Smith","email@email.com", "123131 31313  133", "1000")

	req, _ = http.NewRequest("POST","/make-reservation", strings.NewReader(body))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for insertion room restriction: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context{
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil{
		log.Println(err)
	}
	return ctx
}

func (b *reqBody) getUrlWihtParams(startDate, endDate, firstName, lastName, email, phone, roomID string) string {

	body := reqBody{
		StartDate: startDate,
		EndDate: endDate,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Phone: phone,
		RoomID: roomID,
	}
	return fmt.Sprintf("start_date=%s&end_date=%s&first_name=%s&last_name=%s&email=%s&phone=%s&room_id=%s", body.StartDate, body.EndDate, body.FirstName, body.LastName, body.Email, body.Phone, body.RoomID)
}