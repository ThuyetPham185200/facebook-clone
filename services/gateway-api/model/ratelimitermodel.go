package model

type RateLimterModel struct {
	FeatureLimits   map[string]int
	IPLimit         int
	MaxRequestLimit int
}

func NewRateLimterModel(featureLimits map[string]int, ipLimit int, maxRequestLimit int) *RateLimterModel {
	// fmt.Printf("---------------------------------------\n")
	// for action, limit := range featureLimits {
	// 	fmt.Printf("Action: %s, Limit: %d\n", action, limit)
	// }
	// fmt.Printf("---------------------------------------\n")

	return &RateLimterModel{
		FeatureLimits:   featureLimits,
		IPLimit:         ipLimit,
		MaxRequestLimit: maxRequestLimit,
	}
}
