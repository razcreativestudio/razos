/**
 * RAZOS Internal Reverse Proxy Engine
 * Updated: 2026-10-08
 * Deskripsi: Core reverse proxy dinamis ringan yang memetakan domain hostname pengguna
 * langsung ke port internal container target (port 8100-8200) tanpa restart server.
 */

package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// ProxyRouter mengelola mapping routing hostname domain ke port backend lokal
type ProxyRouter struct {
	mu     sync.RWMutex
	routes map[string]int // Domain -> Host Port
}

var globalRouter = &ProxyRouter{
	routes: make(map[string]int),
}

// GetRouter mengembalikan instance router global
func GetRouter() *ProxyRouter {
	return globalRouter
}

// AddRoute mendaftarkan rute domain ke port target
func (r *ProxyRouter) AddRoute(domain string, targetPort int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[domain] = targetPort
}

// RemoveRoute menghapus rute domain dari tabel routing
func (r *ProxyRouter) RemoveRoute(domain string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.routes, domain)
}

// GetTargetPort mengambil port tujuan untuk domain tertentu
func (r *ProxyRouter) GetTargetPort(domain string) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	port, ok := r.routes[domain]
	return port, ok
}

// ProxyHandler melayani reverse proxy HTTP request ke port lokal yang dituju
func (r *ProxyRouter) ProxyHandler(w http.ResponseWriter, req *http.Request) bool {
	host := req.Host
	port, found := r.GetTargetPort(host)
	if !found {
		return false
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		http.Error(w, "Bad Gateway Target", http.StatusBadGateway)
		return true
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	proxy.ServeHTTP(w, req)
	return true
}
