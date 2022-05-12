package config

import (
	"strconv"
	"sync"
)

type MemConfig struct {
	config map[string]string
	sync.RWMutex
}

var memConfig *MemConfig

func NewMemConfig() *MemConfig {
	if memConfig == nil {
		memConfig = &MemConfig{
			config: make(map[string]string),
		}
	}
	return memConfig
}

func (p *MemConfig) Set(k, v string) error {
	p.Lock()
	defer p.Unlock()

	p.config[k] = v
	return nil
}

func (p *MemConfig) Get(k string) string {
	p.RLock()
	defer p.RUnlock()

	if v, ok := p.config[k]; ok {
		return v
	}
	return ""
}

func (p *MemConfig) GetDef(k, def string) string {
	p.RLock()
	defer p.RUnlock()

	if v, ok := p.config[k]; ok {
		if v == "" {
			return def
		}
		return v
	}
	return def
}

func (p *MemConfig) GetIntDef(k string, idef int) int {
	p.RLock()
	defer p.RUnlock()

	if v, ok := p.config[k]; ok {
		if v == "" {
			return idef
		}
		ival, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return ival
	}
	return idef
}

func (p *MemConfig) GetAll() map[string]string {
	p.RLock()
	defer p.RUnlock()

	buf := make(map[string]string)
	for k, v := range p.config {
		buf[k] = v
	}
	return buf
}
