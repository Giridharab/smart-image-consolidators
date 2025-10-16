package analyzer

import "smart-image-consolidator/configs"

func SuggestCanonicalBase(dockerfileContent string) string {
	for _, canonical := range configs.CanonicalBases {
		if contains(dockerfileContent, canonical.Original) {
			return canonical.Suggested
		}
	}
	return "No suggestion"
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}
