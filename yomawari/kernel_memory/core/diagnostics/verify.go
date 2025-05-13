package diagnostics

import (
	"errors"
	"net/url"
	"strings"
)

func ValidateUrl(urlStr string, requireHttps bool, allowReservedIp bool, allowQuery bool) error {
	if urlStr == "" {
		return errors.New("the URL is empty")
	}

	if requireHttps && strings.HasPrefix(urlStr, "http://") {
		return errors.New("the URL `" + urlStr + "` is not safe, it must start with https://")
	}
	if requireHttps && !strings.HasPrefix(urlStr, "https://") {
		return errors.New("the URL `" + urlStr + "` is incomplete, enter a valid URL starting with 'https://'")
	}

	parsedUrl, err := url.Parse(urlStr)
	if err != nil || parsedUrl.Host == "" {
		return errors.New("the URL `" + urlStr + "` is not valid")
	}

	if requireHttps && parsedUrl.Scheme != "https" {
		return errors.New("the URL `" + urlStr + "` is not safe, it must start with https://")
	}

	isReservedIp := func(host string) bool {
		return strings.HasPrefix(host, "0.") ||
			strings.HasPrefix(host, "10.") ||
			strings.HasPrefix(host, "127.") ||
			strings.HasPrefix(host, "169.254.") ||
			strings.HasPrefix(host, "192.0.0.") ||
			strings.HasPrefix(host, "192.88.99.") ||
			strings.HasPrefix(host, "192.168.") ||
			host == "255.255.255.255"
	}

	if !allowReservedIp && (parsedUrl.Host == "localhost" || isReservedIp(parsedUrl.Host)) {
		return errors.New("the URL `" + urlStr + "` is not safe, it cannot point to a reserved network address")
	}

	if !allowQuery && parsedUrl.RawQuery != "" {
		return errors.New("the URL `" + urlStr + "` is not valid, it cannot contain query parameters")
	}

	if parsedUrl.Fragment != "" {
		return errors.New("the URL `" + urlStr + "` is not valid, it cannot contain URL fragments")
	}

	return nil
}
