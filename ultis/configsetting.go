package utils

var DefaultPreferences = map[string]interface{}{
	"appearance": map[string]interface{}{
		"theme":    "dark",
		"language": "en",
		"timezone": "Asia/Kolkata",
	},
	"notifications": map[string]interface{}{
		"email": map[string]interface{}{
			"enabled": true,
			"types":   []interface{}{"security", "updates", "newsletters"},
		},
		"push": map[string]interface{}{
			"enabled": false,
		},
	},
	"privacy": map[string]interface{}{
		"profileVisibility": "public",
		"showEmail":         false,
		"showPhoneNumber":   false,
	},
	"communication": map[string]interface{}{
		"emailFrequency":         "weekly",
		"preferredContactMethod": "email",
	},
	"accessibility": map[string]interface{}{
		"fontSize":     "medium",
		"highContrast": false,
		"reduceMotion": true,
	},
	"security": map[string]interface{}{
		"twoFactorEnabled":   true,
		"loginNotifications": true,
	},
}

func MergePreferences(defaults, input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range defaults {
		if inputVal, ok := input[key]; ok {
			switch valTyped := val.(type) {
			case map[string]interface{}:
				if inputMap, ok := inputVal.(map[string]interface{}); ok {
					result[key] = MergePreferences(valTyped, inputMap)
				} else {
					result[key] = val
				}
			default:
				result[key] = inputVal
			}
		} else {
			result[key] = val
		}
	}
	return result
}
