package bodyparser

import (
	_ "embed"
	"encoding/json"
	"strings"
)

type Mime struct {
	Extensions []string `json:"extensions"`
}

// mimeType -> [extention1, extention2]
var mimeMap map[string][]string

//go:embed mime-db.json
var mimeDb []byte

// gather mime extention names
func init() {
	mimeMap = map[string][]string{}
	mimes := map[string]json.RawMessage{}
	json.Unmarshal(mimeDb, &mimes)

	for m, mRaw := range mimes {
		mObj := Mime{}
		json.Unmarshal(mRaw, &mObj)
		mExts := mObj.Extensions

		if len(mExts) > 0 {
			mimeMap[m] = mExts
		}
	}
}

func normalize(t string) string {
	if t == "urlencoded" {
		return "application/x-www-form-urlencoded"
	}
	if t == "multipart" {
		return "multipart/*"
	}
	if t[0] == '+' {
		// "+json" -> "*/*+json" expando
		return "*/*" + t
	}

	if !strings.Contains(t, "/") {
		exts := strings.Split(t, ".")
		ext := strings.ToLower(exts[len(exts)-1])

		for m, mExts := range mimeMap {
			for _, mExt := range mExts {
				if ext == mExt {
					return m
				}
			}
		}
	}

	return t
}

func mimeMatch(actual string, expected string) bool {
	actualParts := strings.Split(actual, "/")
	expectedParts := strings.Split(expected, "/")

	if len(actualParts) != 2 || len(expectedParts) != 2 {
		return false
	}

	if expectedParts[0] != "*" && expectedParts[0] != actualParts[0] {
		return false
	}

	if len(expectedParts[1]) >= 2 && expectedParts[1][0:2] == "*+" {
		return (len(expectedParts[1]) <= len(actualParts[1])+1) && (expectedParts[1][1:] == actualParts[1][len(actualParts[1])-len(expectedParts[1])+1:])
	}

	if expectedParts[1] != "*" && expectedParts[1] != actualParts[1] {
		return false
	}

	return true
}
