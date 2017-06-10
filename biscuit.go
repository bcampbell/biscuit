// Package biscuit contains functions for reading http cookies from files.
package biscuit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ReadCookies reads a set of cookies from a Netscape-format cookies.txt file
// format spec, such as it is: http://www.cookiecentral.com/faq/#3.5
func ReadCookies(r io.Reader) ([]*http.Cookie, error) {
	cookies := []*http.Cookie{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// skip blank lines and comments
		// curl supports a "#HttpOnly_" prefix on first field for httponly cookies. Worth supporting, I think.
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			if !strings.HasPrefix(trimmed, "#HttpOnly_") {
				continue
			}
		}
		c, err := parseCookie(line)
		if err != nil {
			return nil, err
		}
		cookies = append(cookies, c)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cookies, nil
}

func toBool(s string) bool {
	// should we be a bit more accepting here?
	return (strings.TrimSpace(s) == "TRUE")
}

func parseCookie(line string) (*http.Cookie, error) {
	f := strings.SplitN(line, "\t", 7)
	if len(f) < 7 {
		return nil, errors.New("syntax error")
	}

	expires := time.Time{}
	if len(strings.TrimSpace(f[4])) > 0 {
		i, err := strconv.ParseInt(f[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad expires '%s'", f[4])
		}
		expires = time.Unix(i, 0)
	}

	httponly := false
	domain := strings.TrimSpace(f[0])
	if strings.HasPrefix(domain, "#HttpOnly_") {
		// convention used by curl (and others?)
		httponly = true
		domain = domain[10:]
	}

	// TODO: if f[1] is TRUE, ensure domain begins with a "."?
	c := &http.Cookie{
		Name:       f[5],
		Value:      f[6],
		Path:       f[2],
		Domain:     domain,
		Expires:    expires,
		RawExpires: f[4],
		MaxAge:     0,
		Secure:     toBool(f[3]),
		HttpOnly:   httponly,
		Raw:        "",
		Unparsed:   []string{},
	}
	return c, nil
}
