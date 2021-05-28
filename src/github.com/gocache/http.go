package gocache

import (
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_gogocache/"

type HTTPPool struct {
	self string
	basePath string
}



func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{self: self,basePath: defaultBasePath}
}

func (pool *HTTPPool) Log(format string,v ...interface{}){
	log.Printf("[Server %s] %s",pool.self,pool.basePath)
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter,r *http.Request){
	if !strings.HasPrefix(r.URL.Path,pool.basePath){
		panic("HTTPPool serving unexpected path: "+r.URL.Path)
	}
	pool.Log("%s %s",r.Method,r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(pool.basePath):],"/",2)
	if len(parts) != 2{
		http.Error(w,"bad request",http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil{
		http.Error(w, "no such group: "+groupName,http.StatusNotFound)
		return
	}
	view ,err := group.Get(key)
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/octet-stream")
	w.Write(view.ByteSlice())
}
