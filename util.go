package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func parseUser(r *http.Request, required map[string]struct{}) (User, error) {
	var uc User
	if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		fs := r.PostForm.Get("fs")
		group := r.PostForm.Get("groupname")
		uc = User{username, password, fs, group}
	} else if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&uc)
	} else {
		return User{}, errors.New("could not parse user (invalid format)")
	}

	// Validate user
	if _, ok := required["username"]; ok && uc.Username == "" {
		return User{}, errors.New("could not parse user. (no username supplied)")
	}
	if _, ok := required["password"]; ok && uc.Password == "" {
		return User{}, errors.New("could not parse user. (no password supplied)")
	}
	if _, ok := required["fs"]; ok && uc.Fs == "" {
		return User{}, errors.New("could not parse user. (no fs supplied)")
	}
	if _, ok := required["group"]; ok && uc.Group == "" {
		return User{}, errors.New("could not parse user. (no group supplied)")
	}

	return uc, nil
}

func parseGroup(r *http.Request) (string, error) {
	var name string
	if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		r.ParseForm()
		name = r.PostForm.Get("groupname")
	} else if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var uc Group
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&uc)
		name = uc.Name
	} else {
		return "", errors.New("could not parse group (invalid format)")
	}

	if name == "" {
		return "", errors.New("could not parse user (no groupname supplied)")
	}
	return name, nil
}

// Decodes hex-encoded SHA512 Hash to Base64 encoding with `{SHA512}` prefix
func ldapEncodePassword(password string) ([]string, error) {
	src := make([]byte, hex.DecodedLen(len(password)))

	_, err := hex.Decode(src, []byte(password))
	if err != nil {
		return []string{}, err
	}

	encoded := base64.StdEncoding.EncodeToString(src)

	return []string{"{SHA512}" + encoded}, nil
}
