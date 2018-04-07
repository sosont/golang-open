package main

import "net/http"

// ResponseWriteReader for middleware
type ResponseWriteReader interface {
	StatusCode() int
	ContentLength() int
	http.ResponseWriter
}

// WrapResponseWriter implement ResponseWriteReader interface
type WrapResponseWriter struct {
	status int
	length int
	http.ResponseWriter
}

// 创建 wrapResponseWriter
func NewWrapResponseWriter(w http.ResponseWriter) *WrapResponseWriter {
	wr := new(WrapResponseWriter)
	wr.ResponseWriter = w
	wr.status = 200
	return wr
}

// 写头
func (p *WrapResponseWriter) WriteHeader(status int) {
	p.status = status
	p.ResponseWriter.WriteHeader(status)
}

func (p *WrapResponseWriter) Write(b []byte) (int, error) {
	n, err := p.ResponseWriter.Write(b)
	p.length += n
	return n, err
}

// 返回状态码
func (p *WrapResponseWriter) StatusCode() int {
	return p.status
}

// 获取内容大小
func (p *WrapResponseWriter) ContentLength() int {
	return p.length
}

// Middleware方法类型
type MiddlewareFunc func(ResponseWriteReader, *http.Request, func())

// Middleware
type MiddlewareServe struct {
	middlewares []MiddlewareFunc
	Handler     http.Handler
}

// ServeHTTP 
func (p *MiddlewareServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := 0
	wr := NewWrapResponseWriter(w)
	var next func()
	next = func() {
		if i < len(p.middlewares) {
			i++
			p.middlewares[i-1](wr, r, next)
		} else if p.Handler != nil {
			p.Handler.ServeHTTP(wr, r)
		}
	}
	next()
}

// 遍历MiddewareFunc
func (p *MiddlewareServe) Use(funcs ...MiddlewareFunc) {
	for _, f := range funcs {
		p.middlewares = append(p.middlewares, f)
	}
}