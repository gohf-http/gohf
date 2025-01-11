package gohf

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// net/http pattern syntax: [METHOD] [HOST]/[PATH]
// METHOD, HOST and PATH are all optional; that is, the string can be "/".

type pattern struct {
	method string
	host   string
	path   string
}

func parsePattern(s string) (pattern, error) {
	var pat pattern

	rest := strings.TrimSpace(s)
	if len(rest) == 0 {
		return pat, errors.New("empty pattern")
	}

	var method, host, path string
	if i := strings.IndexAny(rest, " \t"); i >= 0 {
		method, rest = rest[:i], strings.TrimLeft(rest[i+1:], " \t")
	}

	if i := strings.IndexByte(rest, '/'); i >= 0 {
		host, path = rest[:i], rest[i+1:]
	} else {
		return pat, fmt.Errorf("invalid pattern: host/path missing \"/\". received: \"%s\"", s)
	}

	pat.method = method
	pat.host = host
	pat.path = path
	return pat, nil
}

func (pat pattern) String() string {
	var s string
	if pat.method != "" {
		s = pat.method + " "
	}
	s += pat.host + "/"
	s += strings.TrimLeft(pat.path, "/")
	return s
}

func mergePattern(p1, p2 pattern) (pattern, error) {
	var pat pattern
	var method, host, path string

	if p1.method != "" && p2.method != "" && p1.method != p2.method {
		return pat, fmt.Errorf("method conflict: \"%s\" and \"%s\"", p1.method, p2.method)
	}
	if p1.host != "" && p2.host != "" && p1.host != p2.host {
		return pat, fmt.Errorf("host conflict: \"%s\" and \"%s\"", p1.host, p2.host)
	}

	if p1.method != "" {
		method = p1.method
	} else {
		method = p2.method
	}

	if p1.host != "" {
		host = p1.host
	} else {
		host = p2.host
	}

	path, err := url.JoinPath(p1.path, p2.path)
	if err != nil {
		return pat, fmt.Errorf("url.JoinPath failed: %w", err)
	}

	pat.method = method
	pat.host = host
	pat.path = path
	return pat, nil
}
