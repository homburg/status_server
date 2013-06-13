package main

import (
	"bytes"
	"fmt"
	"github.com/eknkc/amber"
	"homburg/status_server/res"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

const listenAddr = ":8086"

// All access control handled by nginx
func accessControl(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// Get remote address regardless of proxy
// NB. X-Forwarded-For might be a comma-separated list (chain) of ip addresses
func getRemoteAddr(r *http.Request) string {
	headers := r.Header
	forward := headers.Get("X-Forwarded-For")
	if "" != forward {
		return forward
	}
	return r.RemoteAddr
}

// Convert newlines to <br>
func newlineToHtmlBreak(s string) string {
	return strings.Replace(s, "\n", "<br>", -1)
}

// Run command and return html
func commandToHtml(cmds []string) (string, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.Output()
	if nil != err {
		return "", err
	}

	outStr := strings.TrimRight(string(out), "\n")
	return outStr, nil
}

type templateData struct {
	Hostname  string
	GoVersion string
	Script    template.JS
}

var dropboxCommandMatch *regexp.Regexp

func getTemplate() *template.Template {
	// template.Must(tmpl.Parse(status_server.ServerTemplate))
	// Amber template compiler
	compiler := amber.New()
	compiler.Options.PrettyPrint = false
	compiler.Options.LineNumbers = false

	// err := compiler.Parse(status_server.ServerTemplateAmber)
	err := compiler.Parse(status_server.ServerTemplateAmber)

	if nil != err {
		log.Println(err)
	}

	return template.Must(compiler.Compile())
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// Dropbox handler dependencies
	dropboxCommandMatch = regexp.MustCompile("/dropbox/(.*)")
	dropboxAllowedCommands := []string{"status", "help", "start"}

	hostname, _ := os.Hostname()
	tData := templateData{hostname, runtime.Version(), status_server.ServerTemplateScript}

	log.Println("Started")

	// Handle dropbox addresses
	http.HandleFunc("/dropbox/", func(w http.ResponseWriter, r *http.Request) {

		if !accessControl(w, r) {
			return
		}

		path, _ := url.QueryUnescape(r.URL.String())
		matches := dropboxCommandMatch.FindStringSubmatch(path)

		if len(matches) == 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		args := strings.SplitN(matches[1], "/", 2)

		command := args[0]
		if len(args) == 0 {
			command = "help"
		} else {
			allowedCommand := false
			for _, str := range dropboxAllowedCommands {
				if str == command {
					allowedCommand = true
					break
				}
			}

			if !allowedCommand {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		cmdWithArgs := []string{"sudo", "-u", "thomas", "dropbox"}

		for _, a := range args {
			cmdWithArgs = append(cmdWithArgs, a)
		}

		outStr, _ := commandToHtml(cmdWithArgs)
		fmt.Fprintln(w, outStr)
	})

	// landscape-sysinfo
	http.HandleFunc("/landscape/sysinfo", func(w http.ResponseWriter, r *http.Request) {
		if !accessControl(w, r) {
			return
		}

		outStr, _ := commandToHtml([]string{"landscape-sysinfo"})
		fmt.Fprintln(w, outStr)
	})

	// dstat
	http.HandleFunc("/dstat", func(w http.ResponseWriter, r *http.Request) {
		if !accessControl(w, r) {
			return
		}

		outStr, _ := commandToHtml([]string{"dstat", "1", "7"})
		fmt.Fprintln(w, outStr)
	})

	// post actions
	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		if !accessControl(w, r) {
			return
		}

		if r.Method == "POST" {
			action := r.FormValue("action")

			if action == "server-sickbeard-restart" {
				cmd := exec.Command("sudo", "-u", "root", "/home/thomas/bin/service_sickbeard_restart.sh")
				out, err := cmd.Output()
				if nil != err {
					fmt.Fprint(w, err.Error())
				} else {
					fmt.Fprint(w, string(out))
				}
			}
		}
	})

	// Setup template for main layout
	var buf bytes.Buffer

	tmpl := getTemplate()
	err := tmpl.Execute(&buf, tData)
	html := buf.String()

	if nil != err {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if "/" != r.URL.String() {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if !accessControl(w, r) {
			return
		}

		fmt.Fprint(w, string(html))
	})

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
