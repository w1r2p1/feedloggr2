package feedloggr2

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"html/template"

	rss "github.com/jteeuwen/go-pkg-rss"
)

func UpdateFeeds(c *Config) error {
	db, e := OpenSqliteDB(c.Database)
	if e != nil {
		return e
	}

	if c.Verbose {
		fmt.Println("Updating feeds...")
	}

	var feeds []*Feed
	for _, f := range c.Feeds {
		r := rss.NewWithHandlers(5, false, db, db)
		e = r.Fetch(f.Url, nil)
		if e != nil {
			fmt.Printf("Error connecting to %s: %s\n", f.Title, e)
			continue
		}

		items, e := db.GetItems(f.Url)
		if e != nil {
			fmt.Printf("Error updating %s: %s\n", f.Title, e)
			continue
		}
		if c.Verbose {
			fmt.Printf("Got %d new items from: %s\n", len(items), f.Title)
		}

		feeds = append(feeds, &Feed{
			Title: f.Title,
			Url:   f.Url,
			Items: items,
		})
	}

	if c.Verbose {
		fmt.Println("Generating page...")
	}
	funcmap := template.FuncMap{
		"date_link": func(h int, t time.Time) string {
			d := t.Add(time.Hour * time.Duration(h)).Format("2006-01-02")
			return fmt.Sprintf("%s.html", d)
		},
	}
	t := template.Must(template.New("TemplateName").Funcs(funcmap).Parse(HTML_BODY))
	d := struct {
		Date  time.Time
		Feeds []*Feed
	}{
		Date:  Now(),
		Feeds: feeds,
	}
	file := fmt.Sprintf("%s.html", Now().Format("2006-01-02"))
	path := filepath.Join(c.OutputPath, file)
	f, e := os.Create(path)
	if e != nil {
		panic(e) // TODO
	}
	e = t.Execute(f, d)
	if e != nil {
		panic(e) // TODO
	}

	if c.Verbose {
		fmt.Println("Updating symlink...")
	}
	path = filepath.Join(c.OutputPath, "index.html")
	e = os.Remove(path)
	if e != nil {
		perr, ok := e.(*os.PathError)
		// Ignore any "no such file" errors
		// It works correctly, but "|| perr.Err == "is logically wrong. Bug?
		if ok == false || perr.Err == fmt.Errorf("no such file or directory") {
			panic(e) // TODO
		}
	}
	e = os.Symlink(file, path)
	if e != nil {
		panic(e) // TODO
	}

	if c.Verbose {
		fmt.Println("Done.")
	}
	return nil
}