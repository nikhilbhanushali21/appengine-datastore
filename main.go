package main

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/", handleImages)
}

type Song struct {
	songName string
	URL      string
}

func handleImages(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		SongName := strings.Split(req.URL.Path, "/")[1]
		showSong(res, req, SongName)
		return
	}

	if req.Method == "POST" {
		saveImage(res, req)
		return
	}

	listSongs(res, req)
}

func listSongs(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	q := datastore.NewQuery("Song").Order("songName")

	html := "<h2>Welcome to your Database</h2>"

	iterator := q.Run(ctx)
	for {
		var entity Song
		_, err := iterator.Next(&entity)
		if err == datastore.Done {
			break
		} else if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		html += `
			<h3>` + entity.songName + `</h3>
			<br/>
			<img src='` + entity.URL + `' />
			<br/>
		`
	}

	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(res, `
		<!DOCTYPE html>
		<html>
			<body>

				<h1>Create Your Youtube Playlist</h1>

				<h2>Playlist</h2>

				<form method="POST">
					<table>
						<tr>
							<td><label for="Songname">Songname</label></td>
							<td><input type="text" name="Songname"></td>
						</tr>
						<tr>
							<td><label for="url">Song URL</label></td>
							<td><input type="text" name="url"></td>
						</tr>
						<tr>
							<td></td>
							<td><input type="submit"></td>
						</tr>
					</table>
				</form>

				<dl>
					`+html+`
				</dl>
			</body>
		</html>
	`)
}

func showSong(res http.ResponseWriter, req *http.Request, Songname string) {
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "url", Songname, 0, nil)
	var entity Song
	err := datastore.Get(ctx, key, &entity)
	if err == datastore.ErrNoSuchEntity {
		http.NotFound(res, req)
		return
	} else if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(res, `
		<h2>` + entity.songName + `</h2>
		<br/>
		<img src='` + entity.URL + `' />
		<br/>
	`)
}

func saveImage(res http.ResponseWriter, req *http.Request) {
	songname := req.FormValue("Songname")
	url := req.FormValue("url")
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Song", songname, 0, nil)
	entity := Song{
		songName: songname,
		URL: url,
	}

	_, err := datastore.Put(ctx, key, &entity)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	http.Redirect(res, req, "/", 302)
}