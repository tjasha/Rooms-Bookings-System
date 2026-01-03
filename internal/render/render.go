package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewTemplates(a *config.AppConfig) {
	app = a
}

// this will be data added to every template
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	//pop string in the session until there is different data - this are populated when session doesn't have data
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// td = template data
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		//get template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		//log.Println("Could not get template from template cache ")
		return errors.New("Could not get template from template cache ")
	}
	buf := new(bytes.Buffer) //used for better error checking

	////adding data to every template - currently doesn't actually hold data
	//templatedata = AddDefaultData(templatedata, r)
	//
	//// we need to send some data here to not have nil
	//err := template.Execute(buf, templatedata)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{} //it's creating empty map

	// i want to add automatically all available templates that exist
	// they should be added in order

	// i want to first add all *page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) //we just look for all files with this pattern
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *page.tmpl that we found before
	for _, page := range pages { //page is full path to the template
		name := filepath.Base(page)                                     //returns last element of the path = name of the file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) //(ts = template set) we  parse this file (page) and store in the template (name)
		if err != nil {
			return myCache, err
		}

		//now we look for all layouts - we use the same syntax as for the pages
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		//checking how many elements we have
		if len(matches) > 0 {
			// check if any of the pages needs layout inside of them to be rendered. if yes, it adds it to the ts
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts //adding template set to cache map
	}

	return myCache, nil
}
