package router

import (
	"sync/atomic"

	"github.com/connyay/litellm-go/internal/provider"
)

type providerPool struct {
	pvs  []provider.Provider
	next uint64
}

func (p *providerPool) pick() provider.Provider {
	idx := atomic.AddUint64(&p.next, 1)
	return p.pvs[(int(idx)-1)%len(p.pvs)]
}

// Router maps model_name -> provider pool

type Router struct {
	m map[string]*providerPool
}

func New() *Router { return &Router{m: make(map[string]*providerPool)} }

func (r *Router) Register(modelName string, p provider.Provider) {
	pool, ok := r.m[modelName]
	if !ok {
		pool = &providerPool{}
		r.m[modelName] = pool
	}
	pool.pvs = append(pool.pvs, p)
}

func (r *Router) Get(modelName string) (provider.Provider, bool) {
	pool, ok := r.m[modelName]
	if !ok || len(pool.pvs) == 0 {
		return nil, false
	}
	return pool.pick(), true
}

// Len returns number of models registered.
func (r *Router) Len() int { return len(r.m) }
