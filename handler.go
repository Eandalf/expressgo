package expressgo

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

var isValidParamChar *regexp.Regexp
var isNumber *regexp.Regexp

func init() {
	isValidParamChar = regexp.MustCompile("[A-Za-z0-9_]")
	isNumber = regexp.MustCompile("[0-9]")
}

type Handler struct {
	mux *http.ServeMux
	app *App
}

// For path registration

func (h *Handler) isHostIncluded(path string) bool {
	return path[0] != '/'
}

func (h *Handler) makePrecise(path string) string {
	// remove the trailing "/"
	return strings.TrimSuffix(path, "/") + "/{$}"
}

func (h *Handler) pathToLower(path string) string {
	isParam := false
	output := ""

	for _, char := range path {
		if char == '{' {
			isParam = true
		}

		if isParam {
			output += string(char)
		} else {
			output += strings.ToLower(string(char))
		}

		if char == '}' {
			isParam = false
		}
	}

	return output
}

// param name should not start with a number
func (h *Handler) isValidParamName(name string) bool {
	return !isNumber.MatchString(name[0:1])
}

// Parse the params in the current param zone from :name to {name}, and merge them if any separator is found.
//
// Return the following in order, paramName, paramsInZone.
//
// hyphen and dot are translated to 0H and 0D, which is impossible for variables to start with 0, thus, prevents the collisions
func (h *Handler) parseParamZone(currentParams []string, currentSeparators []rune) (string, []string) {
	// separators
	const (
		hyphen = '-'
		dot    = '.'
	)

	paramName := ""
	paramsInZone := []string{}

	for i, param := range currentParams {
		paramName += param
		paramsInZone = append(paramsInZone, param)

		if i < len(currentSeparators) {
			switch currentSeparators[i] {
			case hyphen:
				paramName += "0H"
				paramsInZone = append(paramsInZone, "0H")
			case dot:
				paramName += "0D"
				paramsInZone = append(paramsInZone, "0D")
			}
		}
	}

	if paramName != "" {
		paramName = "{" + paramName + "}"
	}

	return paramName, paramsInZone
}

// Parse the params in path from :name to {name}.
//
// Return the following in order, parsedPath, params, error.
//
// Separators: hyphen(-): "0H", dot(.): "0D".
//
// parsePath: the parsed path
//
// params: a list of params divided by param zones, for example, /:one-:two/:three -> [["one", "0H", "two"], ["three"]]
//
// error: any illegal variable naming or param format found
func (h *Handler) parseParams(path string) (string, [][]string, error) {
	// separators
	const (
		hyphen = '-'
		dot    = '.'
	)

	parsedPath := ""
	// all params
	params := [][]string{}
	// params within a param zone (/: ... /)
	currentParams := []string{}
	// separators within a param zone (/: ... /)
	currentSeparators := []rune{}
	currentParam := ""
	isInParamZone := false
	isAfterSeparator := false

	for pos, char := range path {
		// if being inside the zone of a param string, collect the chars
		if isInParamZone {
			// if leaving the "/:[A-Za-z0-9_]{1,}" interval (seeing another '/')
			if char == '/' {
				if currentParam != "" {
					// invalid variable name found
					if !h.isValidParamName(currentParam) {
						return path, [][]string{}, errors.New("name of a path param is invalid, " + currentParam + " is found")
					}

					currentParams = append(currentParams, currentParam)
					currentParam = ""
				}

				if len(currentParams) > 0 {
					// parse params to the form of "{paramName}"
					paramName, paramsInZone := h.parseParamZone(currentParams, currentSeparators)
					parsedPath = strings.TrimSuffix(parsedPath, ":") + paramName
					params = append(params, paramsInZone)
					currentParams = []string{}
					currentSeparators = []rune{}
				}

				// leave param zone
				isInParamZone = false
				isAfterSeparator = false
				parsedPath += string(char)
				continue
			}

			if isAfterSeparator {
				if char == ':' {
					isAfterSeparator = false
					continue
				}

				return path, [][]string{}, errors.New("path param should start with a colon (:), instead " + string(char) + " is found")
			} else {
				// if seeing separators
				if char == hyphen || char == dot {
					if currentParam != "" {
						// invalid variable name found
						if !h.isValidParamName(currentParam) {
							return path, [][]string{}, errors.New("name of a path param is invalid, " + currentParam + " is found")
						}

						currentParams = append(currentParams, currentParam)
						currentParam = ""
					}

					currentSeparators = append(currentSeparators, char)
					isAfterSeparator = true
					continue
				}
			}

			matched := isValidParamChar.MatchString(string(char))
			// invalid variable name found
			if !matched {
				return path, [][]string{}, errors.New("name of a path param is invalid, " + string(char) + " should not be a part of a param name")
			}

			currentParam += string(char)
		} else {
			// enter the param zone when "/:" is found (for next char)
			if char == ':' && pos-1 >= 0 && path[pos-1:pos] == "/" {
				isInParamZone = true
			}

			// append current char to parsedPath if it is not in the param zone
			parsedPath += string(char)
		}
	}

	// if any remaining currentParam found
	if currentParam != "" {
		// invalid variable name found
		if !h.isValidParamName(currentParam) {
			return path, [][]string{}, errors.New("name of a path param is invalid, " + currentParam + " is found")
		}

		currentParams = append(currentParams, currentParam)
		currentParam = ""
	}

	if len(currentParams) > 0 {
		// parse params to the form of "{paramName}"
		paramName, paramsInZone := h.parseParamZone(currentParams, currentSeparators)
		parsedPath = strings.TrimSuffix(parsedPath, ":") + paramName
		params = append(params, paramsInZone)
	}

	return parsedPath, params, nil
}

func (h *Handler) register(method string, path string, handler *UserHandler) error {
	// apply config options
	if !h.app.config.allowHost && h.isHostIncluded(path) {
		return errors.New("path cannot contain host")
	}

	// parse params
	p, params, err := h.parseParams(path)
	if err != nil {
		return err
	}

	// apply config options
	if !h.app.config.caseSensitive {
		p = h.pathToLower(p)
	}
	if !h.app.config.coarse {
		p = h.makePrecise(p)
	}

	// add method if provided, e.g., get + path => "GET path"
	if method != "" {
		p = method + " " + p
	}

	// register params
	h.app.params[p] = params

	// register callbacks
	// register the slice of callbacks with the route formed by the method and the path
	callbacks := handler.callbacks
	// if the route already exists, push the slice of callbacks to map and not register it to ServeMux
	if _, ok := h.app.callbacks[p]; ok {
		h.app.callbacks[p] = append(h.app.callbacks[p], callbacks)
		return nil
	}
	// if global middlewares exist
	if len(*h.app.globalCallbacks) > 0 {
		// register existing global middlewares first for first-seen routes
		h.app.callbacks[p] = append(*h.app.globalCallbacks, callbacks)
		h.mux.Handle(p, &UserHandler{app: h.app, callbacks: (*h.app.globalCallbacks)[0]})
	} else {
		h.app.callbacks[p] = [][]Callback{callbacks}
		h.mux.Handle(p, handler)
	}

	return nil
}

// For processing requests

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply config options

	// for precise path matching with 301 instead of 308 returned
	//
	// > 1. Requests from http clients to POST paths need to have the path *very* precise. For example, `app.Post("/test/body/base", ...)` would need the path to be set to `/test/body/base/` in client requests.
	// > 2. This is caused by the default behavior of **ExpressGo** to make path precise and the default redirect http status code (301) used by **net/http**.
	// > 3. While making the path precise, **ExpressGo** actually forces each path to have a trailing slash (/).
	// > 4. While an http client sends a request to the originally designated path (`/path`), **net/http** would send a redirect with status code 301 to point to `/path/`.
	// > 5. This would cause the client to drop the request body and resend the request through GET method as per status code 301 indicated.
	// > 6. Related issue: [golang/go#60769](https://github.com/golang/go/issues/60769)
	//
	// This check should be removed once the issue above being implemented.
	if !h.app.config.coarse {
		path := r.URL.Path
		lastChar := path[len(path)-1:]
		if lastChar != "/" {
			r.URL.Path = path + "/"
		}
	}

	// for case-insensitive path matching
	if !h.app.config.caseSensitive {
		r.URL.Path = strings.ToLower(r.URL.Path)
	}

	h.mux.ServeHTTP(w, r)
}
