package util

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

var pearHeaders = []string{
	"authorization",
	"X-Pear-Token",
	"X-Pear-Inner-Auth",
}

func FilterHeader(key string) bool {
	for _, v := range pearHeaders {
		if v == key {
			return true
		}
	}

	for _, v := range hopHeaders {
		if v == key {
			return true
		}
	}

	return false
}
