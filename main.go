package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type githubContent struct {
	organization string
}

// die prints error then exit.
func die(err interface{}) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	var remote string
	flag.StringVar(&remote, "remote", "origin", "remote to fork")
	flag.Parse()

	user := os.Getenv("GITHUB_USER")
	if user == "" {
		die("GITHUB_USER not specified")
	}
	token := os.Getenv("GITHUB_AUTH")
	if user == "" {
		die("GITHUB_AUTH not specified")
	}

	var url string
	cmd := exec.Command("git", "remote", "get-url", remote)
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprint(os.Stderr, string(b))
		die("could not get url of remote")
	}
	url = string(b)
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "https://") {
		die(fmt.Sprintf("only https protocol is supported, got: %s", url))
	}

	addr := strings.TrimPrefix(url, "https://")
	hostPath := strings.SplitN(addr, "/", 2)
	if len(hostPath) != 2 {
		die("invalid form of address:" + addr)
	}
	host := hostPath[0]
	if strings.Contains(host, "@") {
		// value before @ is a token
		host = strings.SplitN(host, "@", 2)[1]
	}
	path := hostPath[1]
	switch host {
	case "github.com":
		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			die("invalid repository path:" + path)
		}
		org := paths[0]
		repo := paths[1]

		forkApiAddr := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", org, repo)
		content, err := json.Marshal(githubContent{organization: org})
		if err != nil {
			die(fmt.Errorf("unable to marshal githubContent: %w", err))
		}
		req, err := http.NewRequest("POST", forkApiAddr, bytes.NewBuffer(content))
		if err != nil {
			die(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", "token "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			die(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			die(fmt.Errorf("could not read response body: %w", err))
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			die(fmt.Sprintf("bad response status: %d\n%s", resp.StatusCode, string(body)))
		}
		// successfully forked, or it has existed already.

		cmd = exec.Command("git", "remote", "add", user, "https://"+token+"@"+host+"/"+user+"/"+repo)
		b, err := cmd.CombinedOutput()
		if err != nil {
			die(fmt.Sprintf("%s\n%s", string(b), err))
		}
	default:
		die(fmt.Sprintf("unsupported host: %s", host))
	}
}
