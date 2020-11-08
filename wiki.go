package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "html/template"
    "regexp"
    "errors"
)

type Page struct {
    Title string
    Body []byte
}

// global variable for storing templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// global variable to store valid paths to the webserver
// mustCompile will parse and compile the regular expression returning regexp.Regexp and error values (Compile will panic)
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// method function to allow the saving of pages
func (p *Page) save() error {
    filename := p.Title + ".txt"
    file := ioutil.WriteFile(filename, p.Body, 0600)
    return file
}

// method function to allow the loading of pages
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    // NO LONGER NESSESARY DUE TO TEMPLATE CACHING
    // t, err := template.ParseFiles(tmpl + ".html")
    // if err != nil {
    //     // throw a 500 error code
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }
    // err = t.Execute(w, p)
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// function to validate and get the title of the page
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil
}

// basic handler to hold the root directory of the webapp
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Basic handler on the path %s", r.URL.Path[1:len(r.URL.Path)])
}

// handles the page views
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        // adds a location header and a 302 status code to the http response
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

// hander for editing pages using html form and template feature
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

// handler for saving the pages edits upon form submission
func saveHandler (w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// wrapper function to take a function of type w, r, title and return an http.HandlerFunc
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func main() {
    // p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
    // p1.save()
    // p2, _ := loadPage("TestPage")
    // fmt.Println(string(p2.Body))

    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
