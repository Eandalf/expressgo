package cors

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Eandalf/expressgo"
)

type CorsConfig struct {
	Origin               any
	originString         string
	originSlice          []string
	originRegExp         *regexp.Regexp
	originBool           bool
	Methods              any
	methodString         string
	methodSlice          []string
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
		originString:         "*",
		originSlice:          []string{},
		methodString:         "GET,HEAD,PUT,PATCH,POST,DELETE",
		methodSlice:          []string{},
		AllowedHeaders:       []string{},
		ExposedHeaders:       []string{},
		PreflightContinue:    false,
		OptionsSuccessStatus: 204,
	}

	// merge configs
	if len(corsConfig) > 0 {
		userConfig := corsConfig[0]

		if o, ok := userConfig.Origin.(string); ok && o != "" {
			config.originString = o
		} else if os, ok := userConfig.Origin.([]string); ok && len(os) > 0 {
			config.originSlice = os
		} else if or, ok := userConfig.Origin.(*regexp.Regexp); ok && or != nil {
			config.originRegExp = or
		} else if ob, ok := userConfig.Origin.(bool); ok && ob {
			config.originBool = ob
		}

		if ms, ok := userConfig.Methods.([]string); ok && len(ms) > 0 {
			config.methodSlice = ms
		} else if m, ok := userConfig.Methods.(string); ok && m != "" {
			config.methodString = m
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
		if len(config.originSlice) > 0 {
			origin := req.Get("Origin")
			if origin != "" {
				for _, o := range config.originSlice {
					if o == origin {
						res.Set("Access-Control-Allow-Origin", origin)
						res.Append("Vary", "Origin")
						break
					}
				}
			}
		} else if config.originRegExp != nil {
			origin := req.Get("Origin")
			if origin != "" && config.originRegExp.MatchString(origin) {
				res.Set("Access-Control-Allow-Origin", origin)
				res.Append("Vary", "Origin")
			}
		} else if config.originBool {
			origin := req.Get("Origin")
			if origin != "" {
				res.Set("Access-Control-Allow-Origin", origin)
				res.Append("Vary", "Origin")
			}
		} else {
			res.Set("Access-Control-Allow-Origin", config.originString)
			if config.originString != "*" {
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
			if len(config.methodSlice) > 0 {
				m := strings.Join(config.methodSlice, ",")
				res.Set("Access-Control-Allow-Methods", m)
			} else if config.methodString != "" {
				res.Set("Access-Control-Allow-Methods", config.methodString)
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
