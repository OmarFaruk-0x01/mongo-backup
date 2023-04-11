package utils

import (
	"errors"
	"fmt"
	"net/url"
)

func IsEmpty(a string) bool {
	return a == ""
}

func GetURLOrigin(_url string) (origin *string, err error) {
	pUrl, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}

	if IsEmpty(pUrl.Host) || IsEmpty(pUrl.Scheme) || IsEmpty(pUrl.User.String()) {
		return nil, errors.New("invalid url")
	}

	_origin := fmt.Sprintf("%s://%s@%s", pUrl.Scheme, pUrl.User.String(), pUrl.Host)
	return &_origin, nil
}
