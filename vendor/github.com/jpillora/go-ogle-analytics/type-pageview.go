package ga

import "net/url"

//WARNING: This file was generated. Do not edit.

//Pageview Hit Type
type Pageview struct {
}

// NewPageview creates a new Pageview Hit Type.
func NewPageview() *Pageview {
	h := &Pageview{}
	return h
}

func (h *Pageview) addFields(v url.Values) error {
	return nil
}

func (h *Pageview) Copy() *Pageview {
	c := *h
	return &c
}
