package model

type Column struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	BaseType    string      `json:"base_type"`
	FieldRef    []string    `json:"field_ref"`
	SpecialType interface{} `json:"special_type"`
	Fingerprint struct {
		Global struct {
			DistinctCount int     `json:"distinct-count"`
			Nil           float64 `json:"nil%"`
		} `json:"global"`
		Type struct {
			TypeNumber struct {
				Min float64 `json:"min"`
				Q1  float64 `json:"q1"`
				Q3  float64 `json:"q3"`
				Max float64 `json:"max"`
				Sd  float64 `json:"sd"`
				Avg float64 `json:"avg"`
			} `json:"type/Number"`
		} `json:"type"`
	} `json:"fingerprint"`
}
