package cors

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Eandalf/expressgo"
)

type CorsConfig struct {
	Origin               string
	Origins              []string
	OriginRegExp         *regexp.Regexp
	OriginBool           bool
	Methods              string
	MethodSlice          []string
	AllowedHeaders       []string
	ExposedHeaders       []string
	Credentials          bool
	MaxAge               int
	PreflightContinue    bool
	OptionsSuccessStatus int
}

func Use(corsConfig ...CorsConfig) expressgo.Callback {
	// the default config
	config := CorsConfig{
		Origin:               "*",
		Origins:              []string{},
		Methods:              "GET,HEAD,PUT,PATCH,POST,DELETE",
		MethodSlice:          []string{},
		AllowedHeaders:       []string{},
		ExposedHeaders:       []string{},
		PreflightContinue:    false,
		OptionsSuccessStatus: 204,
	}

	// merge configs
	if len(corsConfig) > 0 {
		userConfig := corsConfig[0]

		if userConfig.Origin != "" {
			config.Origin = userConfig.Origin
		}
		if len(userConfig.Origins) > 0 {
			config.Origins = userConfig.Origins
			config.Origin = ""
		}
		if userConfig.OriginRegExp != nil {
			config.OriginRegExp = userConfig.OriginRegExp
			config.Origin = ""
		}
		if userConfig.OriginBool {
			config.OriginBool = userConfig.OriginBool
			config.Origin = ""
		}

		if userConfig.Methods != "" {
			config.Methods = userConfig.Methods
		}
		if len(userConfig.MethodSlice) > 0 {
			config.MethodSlice = userConfig.MethodSlice
		}

		if len(userConfig.AllowedHeaders) > 0 {
			config.AllowedHeaders = userConfig.AllowedHeaders
		}
		if len(userConfig.ExposedHeaders) > 0 {
			config.ExposedHeaders = userConfig.ExposedHeaders
		}

		if userConfig.Credentials {
			config.Credentials = userConfig.Credentials
		}
		if userConfig.MaxAge != 0 {
			config.MaxAge = userConfig.MaxAge
		}
		if config.PreflightContinue != userConfig.PreflightContinue {
			config.PreflightContinue = userConfig.PreflightContinue
		}
		if userConfig.OptionsSuccessStatus != 0 {
			config.OptionsSuccessStatus = userConfig.OptionsSuccessStatus
		}
	}

	// create CORS middleware
	cors := func(req *expressgo.Request, res *expressgo.Response, next *expressgo.Next) {
		// set Access-Control-Allow-Origin
		if config.Origin == "*" {
			res.Set("Access-Control-Allow-Origin", config.Origin)
		} else if config.Origin != "" {
			res.Set("Access-Control-Allow-Origin", config.Origin)
			res.Append("Vary", "Origin")
		} else if len(config.Origins) > 0 {
			origin := req.Get("Origin")
			if origin != "" {
				for _, o := range config.Origins {
					if o == origin {
						res.Set("Access-Control-Allow-Origin", origin)
						res.Append("Vary", "Origin")
						break
					}
				}
			}
		} else if config.OriginRegExp != nil {
			origin := req.Get("Origin")
			if origin != "" && config.OriginRegExp.MatchString(origin) {
				res.Set("Access-Control-Allow-Origin", origin)
				res.Append("Vary", "Origin")
			}
		} else if config.OriginBool {
			origin := req.Get("Origin")
			if origin != "" {
				res.Set("Access-Control-Allow-Origin", origin)
				res.Append("Vary", "Origin")
			}
		}

		// set Access-Control-Allow-Credentials
		if config.Credentials {
			res.Set("Access-Control-Allow-Credentials", "true")
		}

		// set Access-Control-Expose-Headers
		if len(config.ExposedHeaders) > 0 {
			h := strings.Join(config.ExposedHeaders, ",")
			if h != "" {
				res.Set("Access-Control-Expose-Headers", h)
			}
		}

		if req.Native.Method == http.MethodOptions {
			// for preflights

			// set Access-Control-Allow-Methods
			if len(config.MethodSlice) > 0 {
				m := strings.Join(config.MethodSlice, ",")
				res.Set("Access-Control-Allow-Methods", m)
			} else if config.Methods != "" {
				res.Set("Access-Control-Allow-Methods", config.Methods)
			}

			// set Access-Control-Allow-Headers
			if len(config.AllowedHeaders) > 0 {
				h := strings.Join(config.AllowedHeaders, ",")
				if h != "" {
					res.Set("Access-Control-Allow-Headers", h)
				}
			} else {
				h := req.Get("Access-Control-Request-Headers")
				if h != "" {
					res.Set("Access-Control-Allow-Headers", h)
					res.Append("Vary", "Access-Control-Request-Headers")
				}
			}

			// set Access-Control-Max-Age
			if config.MaxAge > 0 {
				t := strconv.Itoa(config.MaxAge)
				res.Set("Access-Control-Max-Age", t)
			}

			if config.PreflightContinue {
				next.Next = true
				next.Route = true
			} else {
				// Support some old clients.
				// Safari and potentially other browsers need content-length 0.
				// Set status code 204 to prevent them from hanging for waiting a body.
				res.Status(config.OptionsSuccessStatus)
				res.Set("Content-Length", "0")
				res.End()
			}
		} else {
			next.Next = true
			next.Route = true
		}
	}

	return cors
}
