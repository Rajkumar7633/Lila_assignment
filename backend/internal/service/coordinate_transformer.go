package service

import (
	"fmt"
	"sync"

	"github.com/lila/player-viz-backend/internal/domain"
)

const imageSize = 1024

type mapDefinition struct {
	id       string
	label    string
	scale    float64
	originX  float64
	originZ  float64
	imageURL string
}

var defaultMaps = []mapDefinition{
	{id: "AmbroseValley", label: "Ambrose Valley", scale: 900, originX: -370, originZ: -473, imageURL: "/minimaps/AmbroseValley_Minimap.png"},
	{id: "GrandRift", label: "Grand Rift", scale: 581, originX: -290, originZ: -290, imageURL: "/minimaps/GrandRift_Minimap.png"},
	{id: "Lockdown", label: "Lockdown", scale: 1000, originX: -500, originZ: -500, imageURL: "/minimaps/Lockdown_Minimap.jpg"},
}

type linearTransformer struct {
	def mapDefinition
}

func (t linearTransformer) MapID() string { return t.def.id }

func (t linearTransformer) Config() domain.MapInfo {
	return domain.MapInfo{
		ID:          t.def.id,
		Label:       t.def.label,
		Scale:       t.def.scale,
		OriginX:     t.def.originX,
		OriginZ:     t.def.originZ,
		ImageURL:    t.def.imageURL,
		ImageWidth:  imageSize,
		ImageHeight: imageSize,
	}
}

func (t linearTransformer) WorldToPixel(x, z float64) (float64, float64) {
	u := (x - t.def.originX) / t.def.scale
	v := (z - t.def.originZ) / t.def.scale
	pixelX := u * imageSize
	pixelY := (1 - v) * imageSize
	return pixelX, pixelY
}

type TransformerRegistry struct {
	mu           sync.RWMutex
	transformers map[string]domain.CoordinateTransformer
}

func NewTransformerRegistry() *TransformerRegistry {
	reg := &TransformerRegistry{transformers: make(map[string]domain.CoordinateTransformer)}
	for _, def := range defaultMaps {
		reg.transformers[def.id] = linearTransformer{def: def}
	}
	return reg
}

func (r *TransformerRegistry) Get(mapID string) (domain.CoordinateTransformer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.transformers[mapID]
	if !ok {
		return nil, fmt.Errorf("unknown map: %s", mapID)
	}
	return t, nil
}

func (r *TransformerRegistry) ListMaps() []domain.MapInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.MapInfo, 0, len(r.transformers))
	for _, def := range defaultMaps {
		if t, ok := r.transformers[def.id]; ok {
			out = append(out, t.Config())
		}
	}
	return out
}

func (r *TransformerRegistry) Register(t domain.CoordinateTransformer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transformers[t.MapID()] = t
}

func (r *TransformerRegistry) Transform(mapID string, x, z float64) (domain.Point, error) {
	t, err := r.Get(mapID)
	if err != nil {
		return domain.Point{}, err
	}
	px, py := t.WorldToPixel(x, z)
	return domain.Point{X: px, Y: py}, nil
}
