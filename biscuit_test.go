package biscuit

import (
	"fmt"
	"strings"
)

func ExampleReadCookies() {

	example := `# Netscape HTTP Cookie File
# http://curl.haxx.se/docs/http-cookies.html
# This file was generated by libcurl! Edit at your own risk.

.theguardian.com	TRUE	/	FALSE	1504398562	GU_mvt_id	975374
www.theguardian.com	FALSE	/	FALSE	0	GU_geo_continent	OC
#HttpOnly_.stackoverflow.com	TRUE	/	FALSE	2682374400	prov	7341b18d-42b0-5071-d792-f57c4c74188f
`

	cookies, err := ReadCookies(strings.NewReader(example))
	if err != nil {
		panic(err)
	}
	for _, c := range cookies {
		fmt.Println(c)
	}

	// Output:
	// GU_mvt_id=975374; Path=/; Domain=theguardian.com; Expires=Sun, 03 Sep 2017 00:29:22 GMT
	// GU_geo_continent=OC; Path=/; Domain=www.theguardian.com
	// prov=7341b18d-42b0-5071-d792-f57c4c74188f; Path=/; Domain=stackoverflow.com; Expires=Fri, 01 Jan 2055 00:00:00 GMT; HttpOnly
	//
}
