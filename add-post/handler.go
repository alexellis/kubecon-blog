package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"

	// "github.com/go-git/go-git"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		input = body
	}

	log.Println(string(input), r.Header, r.Method, r.URL.Path, r.Form)
	if r.Method == http.MethodPost {
		parts, err := url.ParseQuery(string(input))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unable to read input from form: %s", err.Error())))
			return
		}

		title := parts.Get("postTitle")
		bodyVal := parts.Get("postBody")
		err = clonePush(title, bodyVal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Title: %q\n", title)
		log.Printf("Body: %q\n", bodyVal)

		//push to git here
		w.Write([]byte("Thank you, your post has been submitted for publishing."))
		return
	}

	res, err := ioutil.ReadFile("./templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Unable to locate file: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func clonePush(title string, body string) error {
	tmpPath := path.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		return err
	}

	tokenBytes, err := ioutil.ReadFile("/var/openfaas/secrets/github-token")
	if err != nil {
		return err
	}

	token := strings.TrimSpace(string(tokenBytes))
	r, err := git.PlainClone(tmpPath, false, &git.CloneOptions{
		Auth: &ghttp.BasicAuth{
			Username: "alexellis",
			Password: token,
		},
		URL: os.Getenv("GITHUB_REPO"),
	})

	if err != nil {
		return errors.Wrapf(err, "unable to clone: %s", os.Getenv("GITHUB_REPO"))
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	titleSafe := strings.ReplaceAll(title, " ''@#$<>", "-")
	titleSafe = strings.ReplaceAll(title, " ", "-")
	filename := filepath.Join(tmpPath, "blog/content/posts/"+titleSafe+".md")
	postContent := fmt.Sprintf(`---
title: "%s"
date: 2020-10-20T20:21:04+01:00
draft: false
---
%s`, title, body)

	err = ioutil.WriteFile(filename, []byte(postContent), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = w.Add("blog/content/posts/" + titleSafe + ".md")
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Println(status)

	commit, err := w.Commit("Blog post by add-post\n\nSigned-off-by: add-post <add-post@openfaas.com>", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "add-post",
			Email: "add-post@openfaas.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	fmt.Println(commit)

	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}
	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		Auth: &ghttp.BasicAuth{
			Username: "alexellis",
			Password: token,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
