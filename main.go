package main

import "fmt"
import "os"
import "net/http"
import "html/template"
import "path/filepath"

// new iter stuff
import "maps"  // maps.Keys()
import "slices" // slices.Collect()
// why in two different packages?

var index = template.Must(template.New("page").Parse(`
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<ul>
			{{range .Repos}}
				<li><a href="{{.}}">{{.}}</a></li>
			{{end}}
		</ul>
	</body>
</html>
	`))

var repoPage = template.Must(template.New("repopage").Parse(`
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<ul>
			{{range .Files}}
				<li>{{.}}</li>
			{{end}}
		</ul>
		<a href="/">back</a>
	</body>
</html>
	`))

func ScanForGitRepos(path string) (map[string]repo, error) {
	libgit_init()
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	repos := make(map[string]repo)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		repo, err := repository_open(filepath.Join(path, e.Name()))
		if err != nil {
			continue
		}

		repos[e.Name()] = repo
	}

	return repos, nil
}

func main() {
	repos, err := ScanForGitRepos("..")
	if err != nil {
		fmt.Println("Could not load repos")
		os.Exit(1)
	}

	// index repo chooser
	repoNames := slices.Collect(maps.Keys(repos))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := struct {
			Title string
			Repos []string
		}{
			Title: "repos",
			Repos: repoNames,
		}

		if err := index.Execute(w, data); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})

	// Interesting. Prebaking all the paths, yo!
	//  Although I could just use r.PathValue("id") with "/{id}" path.
	for name, repo := range repos {
		http.HandleFunc(filepath.Join("/", name), func(w http.ResponseWriter, r *http.Request) {
			filenames, err := list_all(repo)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Server error", http.StatusInternalServerError)
			}

			
			data := struct {
				Title string
				Files []string
			}{
				Title: name,
				Files: filenames,
			}

			if err := repoPage.Execute(w, data); err != nil {
				fmt.Println(err)
				http.Error(w, "Server error", http.StatusInternalServerError)
			}
		})
	}

	addr := ":8080"
	fmt.Println("Listenin' on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Server error:", err)
	}

	// fmt.Println("asd asd cock");
	// libgit_init()

	// repo, err := repository_open(".")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)  // maybe copy the exit status of the libgit error?
	// }

	// names, err := list_all(repo)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println(names)
}
