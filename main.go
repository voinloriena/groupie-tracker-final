package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type errorData struct {
	Num  int
	Text string
}

type artistData struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Relation     string              `json:"relations"`
	Concerts     map[string][]string `json:"datesLocations"`
}

type relation struct {
	ID       int                 `json:"id"`
	Concerts map[string][]string `json:"datesLocations"`
}

var allData []artistData

func main() {

	fmt.Println("...Uno momento por favor ¯\\_(ツ)_/¯ ... Give me just a minute to gather data...")
	FileServer := http.FileServer(http.Dir("docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", FileServer))

	allData = gatherDataUp("https://groupietrackers.herokuapp.com/api/artists")
	if allData == nil {
		fmt.Println("Failed to gather Data from API")
		os.Exit(1)
	}

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/response", response)
	http.HandleFunc("/search", search)

	fmt.Println()
	fmt.Println("Thanks, man (ಥ﹏ಥ) Now Server is listening to port #8080   ᕦ(ò_óˇ)ᕤ")
	http.ListenAndServe(":8080", nil)
}

func mainPage(res http.ResponseWriter, req *http.Request) {
	temp, er := template.ParseFiles("docs/htmlTemplates/index.html")
	if er != nil {
		err(res, req, http.StatusInternalServerError)
		return
	}
	if req.URL.Path != "/" {
		err(res, req, http.StatusNotFound)
		return
	}
	f := filter(req)
	result := []artistData{}
	for _, v := range allData {
		if hasString(f, v.Name) {
			result = append(result, v)
		}
	}
	temp.Execute(res, result)
}

func response(res http.ResponseWriter, req *http.Request) {
	temp, er := template.ParseFiles("docs/htmlTemplates/response.html")
	if er != nil {
		log.Fatal(er)
		err(res, req, http.StatusInternalServerError)
		return
	}
	name := req.FormValue("name")
	for _, v := range allData {
		if v.Name == name {
			fmt.Println(v)
			temp.Execute(res, v)
			break
		}
	}
	return
}

func search(res http.ResponseWriter, req *http.Request) {
	temp, e1 := template.ParseFiles("docs/htmlTemplates/search.html")
	if e1 != nil {
		err(res, req, http.StatusInternalServerError)
		return
	}
	temp.Execute(res, allData)
}

func err(res http.ResponseWriter, req *http.Request, err int) {
	temp, er := template.ParseFiles("docs/htmlTemplates/error.html")
	if er != nil {
		log.Fatal(er)
		return
	}
	res.WriteHeader(err)
	errData := errorData{Num: err}
	if err == 404 {
		errData.Text = "Page Not Found"
	} else if err == 400 {
		errData.Text = "Bad Request"
	} else if err == 500 {
		errData.Text = "Internal Server Error"
	}
	fmt.Println(errData)
	temp.Execute(res, errData)
}

func gatherDataUp(link string) []artistData {
	data1 := getData(link)
	Artists := []artistData{}
	e := json.Unmarshal(data1, &Artists)
	if e != nil {
		log.Fatal(e)
		return nil
	}
	for i := 0; i < len(Artists); i++ {
		r := relation{}
		json.Unmarshal(getData(Artists[i].Relation), &r)
		Artists[i].Concerts = r.Concerts
	}
	return Artists
}

func getData(link string) []byte {
	data1, e1 := http.Get(link)
	if e1 != nil {
		log.Fatal(e1)
		return nil
	}
	data2, e2 := ioutil.ReadAll(data1.Body)
	if e2 != nil {
		log.Fatal(e2)
		return nil
	}
	return data2
}

func filter(req *http.Request) []string {
	req.ParseForm()

	fromCD, e1 := strconv.Atoi(req.FormValue("fromCD"))
	toCD, e2 := strconv.Atoi(req.FormValue("toCD"))
	if e1 != nil {
		fromCD = 1957
	}
	if e2 != nil {
		toCD = 2020
	}
	fmt.Println("fromCD:", fromCD)
	fmt.Println("toCD:", toCD)

	f1 := filterCD(fromCD, toCD)
	fmt.Println("f1:", f1)

	memNum := req.Form["member"]
	f2 := filterMem(memNum)
	fmt.Println("memNum:", memNum)
	if len(memNum) == 0 {
		f2 = f1
	}
	fmt.Println("f2:", f2)

	fromFA := req.FormValue("fromFA")
	toFA := req.FormValue("toFA")
	if fromFA == "" {
		fromFA = "1957-01-01"
	}
	if toFA == "" {
		toFA = "2025-01-01"
	}
	f3 := filterFA(fromFA, toFA)
	fmt.Println("f3:", f3)

	loc := req.FormValue("location")
	fmt.Println("loc:", loc)
	f4 := filterLoc(loc)
	if loc == "" {
		f4 = f1
	}
	fmt.Println("f4:", f4)

	result := []string{}
	for _, v := range allData {
		if hasString(f1, v.Name) && hasString(f2, v.Name) && hasString(f3, v.Name) && hasString(f4, v.Name) {
			result = append(result, v.Name)
		}
	}
	return result
}

func filterCD(min, max int) []string {
	if max < min {
		min, max = max, min
	}
	a := []string{}
	for _, v := range allData {
		if v.CreationDate >= min && v.CreationDate <= max {
			a = append(a, v.Name)
		}
	}
	return a
}

func filterMem(a []string) []string {
	b := []int{}
	for _, v := range a {
		i, _ := strconv.Atoi(v)
		b = append(b, i)
	}
	c := []string{}
	for _, v := range b {
		for _, g := range allData {
			if len(g.Members) == v {
				c = append(c, g.Name)
			}
		}
	}
	return c
}

func filterFA(from, to string) []string {
	a := []string{}
	date1, e1 := time.Parse("2006-01-02", from)
	date2, e2 := time.Parse("2006-01-02", to)
	fmt.Println("fromFA:", date1)
	fmt.Println("toFA:", date2)
	if e1 != nil || e2 != nil {
		fmt.Println("Time parse problem")
	}
	if date2.Before(date1) {
		date1, date2 = date2, date1
	}
	for _, v := range allData {
		FAdate, e3 := time.Parse("02-01-2006", v.FirstAlbum)
		if e3 != nil {
			fmt.Println("FADate parse problem")
		}
		if date2.After(FAdate) && date1.Before(FAdate) {
			fmt.Println(date1, "<", FAdate, "<", date2)
			a = append(a, v.Name)
		} else {
			fmt.Println("ELSE:")
			fmt.Println(date1, "<", FAdate, "<", date2)
		}
	}
	return a
}

func filterLoc(a string) []string {
	b := []string{}
	for _, artist := range allData {
		for key := range artist.Concerts {
			if searchWord(key, a) {
				b = append(b, artist.Name)
			}
		}
	}
	return b
}

func searchWord(s string, toFind string) bool {
	for i := 0; i <= len(s)-len(toFind); i++ {
		if toFind == s[i:i+len(toFind)] {
			return true
		}
	}
	return false
}

func hasString(ar []string, a string) bool {
	for _, v := range ar {
		if v == a {
			return true
		}
	}
	return false
}
