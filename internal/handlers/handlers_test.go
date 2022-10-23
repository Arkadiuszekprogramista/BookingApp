package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
}

var theTestForPost = []struct{
	name string
	url string
	method string
	expectedStatusCode int
	postData []postData
}{
	{"is form valid", "/search-availability", "POST", http.StatusOK, []postData{
		{key:"start", value: "2020-01-01"},
		{key:"end", value: "2020-01-12"},
		},
	},
}


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
		RoomID: "1 ",
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
		body := reqBody.urlValues(e.reqBody.StartDate, e.reqBody.EndDate, e.reqBody.FirstName, e.reqBody.LastName, e.reqBody.Email,  e.reqBody.Phone, e.reqBody.RoomID)

		req, _ := http.NewRequest(e.method, e.url, strings.NewReader(body.Encode()))
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
	body := reqBody.urlValues("2050-01-01", "2050-01-02","Johny","Smith","email@email.com", "123131 31313  133", "2")

	req, _ = http.NewRequest("POST","/make-reservation", strings.NewReader(body.Encode()))
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
	body = reqBody.urlValues("2050-01-01", "2050-01-02","Johny","Smith","email@email.com", "123131 31313  133", "1000")

	req, _ = http.NewRequest("POST","/make-reservation", strings.NewReader(body.Encode()))
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

func TestRepository_PostAvailabilty(t *testing.T) {
	//room are not available
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-02")

	req, _ := http.NewRequest("POST","/search-availability", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no room availabe gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	postedData = url.Values{}
	postedData.Add("start","2040-01-01")
	postedData.Add("end", "2040-01-02")

	req, _ = http.NewRequest("POST","/search-availability", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong  status code: got %d, wnated %d", rr.Code, http.StatusOK)
	}

	//empty post body
	postedData = url.Values{}
	postedData.Add("start","2040-01-01")
	postedData.Add("end", "2040-01-02")

	req, _ = http.NewRequest("POST","/search-availability", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//star date is invalid
	postedData = url.Values{}
	postedData.Add("start","invalid")
	postedData.Add("end", "2040-01-02")

	req, _ = http.NewRequest("POST","/search-availability", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wnated %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//end date is invalid
	postedData = url.Values{}
	postedData.Add("start","2040-01-01")
	postedData.Add("end", "invalid")

	req, _ = http.NewRequest("POST","/search-availability", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wnated %d", rr.Code, http.StatusTemporaryRedirect)
	}
	
	//database query fails
	postedData = url.Values{}
	postedData.Add("start","2060-01-01")
	postedData.Add("end", "2060-01-02")

	req, _ = http.NewRequest("POST","/search-availability", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wnated %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	//rooms are not available
	//create request body
	postData := url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-02")
	postData.Add("room_id", "3")

	//create requestt
	req, _ := http.NewRequest("POST","/search-availability-json", strings.NewReader(postData.Encode()))

	// get contex whit session
	ctx := getCtx(req)

	req = req.WithContext(ctx)

	//set the request header 
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//get response recorder
	rr := httptest.NewRecorder()

	//make handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	//make request to out handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}


	if j.Ok {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	//rooms are available
	//create request body
	postData = url.Values{}
	postData.Add("start", "2050-01-01")
	postData.Add("end", "2050-01-02")
	postData.Add("room_id", "3")

	//create requestt
	req, _ = http.NewRequest("POST","/search-availability-json", strings.NewReader(postData.Encode()))

	// get contex whit session
	ctx = getCtx(req)

	req = req.WithContext(ctx)

	//set the request header 
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//get response recorder
	rr = httptest.NewRecorder()

	//make handlerfunc
	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	//make request to out handler
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if j.Ok {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	//no request body
	req, _ = http.NewRequest("POST","/search-availability-json", nil)
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
	if j.Ok || j.Messeage != "internal server error" {
		t.Error("Got availability when request body was empty")
	}

	//database error
	postData = url.Values{}
	postData.Add("start", "2060-01-01")
	postData.Add("end", "2060-01-02")
	postData.Add("room_id", "1")

	req, _ = http.NewRequest("POST","/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}
	if j.Ok != false && j.Messeage != "error querying database" {
		t.Error("Got availability when simulating database error")
	}


}

func getCtx(req *http.Request) context.Context{
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil{
		log.Println(err)
	}
	return ctx
}

func (b *reqBody) urlValues(startDate, endDate, firstName, lastName, email, phone, roomID string) url.Values {

	postData := url.Values{}
	postData.Add("start_date", startDate)
	postData.Add("end_data", endDate)
	postData.Add("first_name", firstName)
	postData.Add("last_name", lastName)
	postData.Add("email", email)
	postData.Add("phone", phone)
	postData.Add("room_id", roomID)

	return postData
}
