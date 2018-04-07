package route

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)


func Parameters(req *http.Request) Params {
	if p, _ := req.Body.(*parameters); p != nil {
		return p.params
	}
	return nil
}


func Pattern(req *http.Request) string {
	if p, _ := req.Body.(*parameters); p != nil {
		return p.pattern
	}
	return req.URL.Path // if matched will be same as url path
}


func Recycle(req *http.Request) {
	if p, _ := req.Body.(*parameters); p != nil {
		p.reset(req)
	}
}


type Params []struct{ Key, Value string }

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}


func (ps *Params) push(key, val string) {
	n := len(*ps)
	*ps = (*ps)[:n+1]
	(*ps)[n].Key, (*ps)[n].Value = key, val
}


type Router interface {
	http.Handler

	Route(*http.Request) http.Handler
}

type RouterFunc func(*http.Request) http.Handler

func (f RouterFunc) Route(req *http.Request) http.Handler {
	h := f(req)
	if r, ok := h.(Router); ok {
		return r.Route(req)
	}
	return h
}

func (f RouterFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h := f(req); h != nil {
		h.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}


func Chain(routes ...Router) Router {
	return RouterFunc(func(req *http.Request) http.Handler {
		for _, router := range routes {
			if handler := router.Route(req); handler != nil {
				return handler
			}
		}
		return nil
	})
}


func New(path string, handler interface{}) Router {
	p := "/" + strings.TrimLeft(path, "/")

	var h http.Handler = nil
	switch t := handler.(type) {
	case http.HandlerFunc:
		h = t
	case func(http.ResponseWriter, *http.Request):
		h = http.HandlerFunc(t)
	case nil:
		panic("given handler cannot be: nil")
	default:
		panic(fmt.Sprintf("not a handler given: %T - %+v", t, t))
	}

	// maybe static route
	if strings.IndexAny(p, ":*") == -1 {
		return RouterFunc(func(req *http.Request) http.Handler {
			if p == req.URL.Path {
				return h
			}
			return nil
		})
	}

	// prepare and validate pattern segments to match
	segments := strings.Split(strings.Trim(p, "/"), "/")
	for i, seg := range segments {
		segments[i] = "/" + seg
		if pos := strings.IndexAny(seg, ":*"); pos == -1 {
			continue
		} else if pos != 0 {
			panic("special param matching signs, must follow after slash: " + p)
		} else if len(seg)-1 == pos {
			panic("param must be named after sign: " + p)
		} else if seg[0] == '*' && i+1 != len(segments) {
			panic("match all, must be the last segment in pattern: " + p)
		} else if strings.IndexAny(seg[1:], ":*") != -1 {
			panic("only one param per segment: " + p)
		}
	}
	ts := p[len(p)-1] == '/' // whether we need to match trailing slash

	// pool for parameters
	num := strings.Count(p, ":") + strings.Count(p, "*")
	pool := sync.Pool{}
	pool.New = func() interface{} {
		return &parameters{params: make(Params, 0, num), pool: &pool, pattern: p}
	}

	// extend handler in order to salvage parameters
	handle := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
		if p, _ := req.Body.(*parameters); p != nil {
			p.reset(req)
		}
	})

	// dynamic route matcher
	return RouterFunc(func(req *http.Request) http.Handler {
		ps := pool.Get().(*parameters)
		if match(segments, req.URL.Path, &ps.params, ts) {
			ps.ReadCloser = req.Body
			req.Body = ps
			return handle
		}
		ps.params = ps.params[0:0]
		pool.Put(ps)
		return nil
	})
}

// matches pattern segments to an url and pushes named parameters to ps
func match(segments []string, url string, ps *Params, ts bool) bool {
	for _, segment := range segments {
		switch {
		case len(url) == 0 || url[0] != '/':
			return false
		case segment[1] == ':' && len(url) > 1:
			end := 1
			for end < len(url) && url[end] != '/' {
				end++
			}
			ps.push(segment[2:], url[1:end])
			url = url[end:]
		case segment[1] == '*':
			ps.push(segment[2:], url)
			return true
		case len(url) < len(segment) || url[:len(segment)] != segment:
			return false
		default:
			url = url[len(segment):]
		}
	}
	return (!ts && url == "") || (ts && url == "/") // match trailing slash
}

type parameters struct {
	io.ReadCloser
	params  Params
	pattern string
	pool    *sync.Pool
}

func (p *parameters) reset(req *http.Request) {
	req.Body = p.ReadCloser
	p.params = p.params[0:0]
	p.pool.Put(p)
}
