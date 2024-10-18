package expressgo

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

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
	isNumber, _ := regexp.MatchString("[0-9]", name[0:1])
	return !isNumber
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
	isValidParamChar, _ := regexp.Compile("[A-Za-z0-9_]")
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
						return path, [][]string{}, errors.New("name of a path param is invalid")
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
							return path, [][]string{}, errors.New("name of a path param is invalid")
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
				return path, [][]string{}, errors.New("name of a path param is invalid")
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
			return path, [][]string{}, errors.New("name of a path param is invalid")
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

func (h *Handler) register(method string, path string, handler http.Handler) error {
	// apply config options
	if !h.app.allowHost && h.isHostIncluded(path) {
		return errors.New("path cannot contain host")
	}

	// parse params
	p, params, err := h.parseParams(path)
	if err != nil {
		return err
	}

	// apply config options
	if !h.app.caseSensitive {
		p = h.pathToLower(p)
	}
	if !h.app.coarse {
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
	callbacks := []Callback{}
	if uh, ok := handler.(*UserHandler); ok {
		callbacks = uh.callbacks
	}
	// if the route already exists, push the slice of callbacks to map and not register it to ServeMux
	if _, ok := h.app.callbacks[p]; ok {
		h.app.callbacks[p] = append(h.app.callbacks[p], callbacks)
		return nil
	}
	h.app.callbacks[p] = [][]Callback{callbacks}

	h.mux.Handle(p, handler)
	return nil
}

// For processing requests

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply config options
	if !h.app.caseSensitive {
		r.URL.Path = strings.ToLower(r.URL.Path)
	}

	h.mux.ServeHTTP(w, r)
}
