package gonfiguration

func SetDefault(key string, val interface{}) {
	gonfig.setDefault(key, val)
}

func SetDefaults(defaults map[string]interface{}) {
	gonfig.setDefaults(defaults)
}

func GetDefaults() map[string]interface{} {
	return gonfig.getDefaults()
}

func (g *gonfiguration) setDefault(key string, val interface{}) {
	g.Lock()
	defer g.Unlock()

	g.defaults[key] = val
}

func (g *gonfiguration) getDefault(key string) interface{} {
	g.RLock()
	defer g.RUnlock()

	if _, ok := g.defaults[key]; !ok {
		return nil
	}

	return g.defaults[key]
}

func (g *gonfiguration) setDefaults(defaults map[string]interface{}) {
	for key, val := range defaults {
		g.setDefault(key, val)
	}
}

func (g *gonfiguration) getDefaults() map[string]interface{} {
	g.RLock()
	defer g.RUnlock()

	defaultsCopy := make(map[string]interface{}, len(g.defaults))
	for key, val := range g.defaults {
		defaultsCopy[key] = val
	}

	return defaultsCopy
}
