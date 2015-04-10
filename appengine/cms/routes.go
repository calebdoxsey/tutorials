package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router

func serveTemplate(res http.ResponseWriter, name string, data interface{}) {
	tpl, err := template.ParseFiles("templates/" + name + ".gohtml")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	body := buf.String()
	tpl, err = template.ParseFiles("templates/layout.gohtml")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	buf.Reset()
	err = tpl.Execute(&buf, map[string]interface{}{
		"Body": body,
	})
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	res.Write(buf.Bytes())
}

func init() {
	router = httprouter.New()

	// API methods
	router.GET("/api/documents", serveDocumentsList)
	router.GET("/api/documents/:id", serveDocumentsGet)
	router.POST("/api/documents", serveDocumentsCreate)
	router.PUT("/api/documents/:id", serveDocumentsUpdate)
	router.DELETE("/api/documents/:id", serveDocumentsDelete)

	router.GET("/api/files/:id", serveFilesGet)
	router.POST("/api/files", serveFilesUpload)

	router.POST("/api/users", serveUsersCreate)
	router.POST("/api/users/login", serveUsersLogin)
	router.POST("/api/users/logout", serveUsersLogout)

	// HTML methods
	router.GET("/", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		session, _ := sessionStore.Get(req, "session")
		email, _ := session.Values["email"]

		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, `<!DOCTYPE html>
	<html>
	<head>
		<title>CMS Example</title>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.4/css/bootstrap.css" rel="stylesheet">
		<link href="/static/styles/main.css" rel="stylesheet">
	</head>
	<body>
		<script>var EMAIL = "`+fmt.Sprint(email)+`";</script>
		<script src="/static/scripts/hyperscript.js"></script>
		<script>h = hyperscript;</script>
		<script src="/static/scripts/views.js"></script>
		<script src="/static/scripts/main.js"></script>
	</body>
	</html>`)
	})

	http.Handle("/", router)
}
