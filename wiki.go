package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "html/template"
)

type Page struct {
    Title string
    Body []byte
}

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

// basic handler to hold the root directory of the webapp
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Basic handler on the path %s", r.URL.Path[1:len(r.URL.Path)-1])
}

// handles the views
func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", title, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("edit.html")
    t.Execute(w, p)
}

func saveHandler (w http.ResponseWriter, r *http.Request) {

}

func main() {
    p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))

    http.HandleFunc("/", handler)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
