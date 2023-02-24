package util

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

var validTags = []string{
	"Arrangement",
	"Verse", "Verse 1", "Verse 2", "Verse 3", "Verse 4", "Verse 5",
	"Chorus", "Chorus 1", "Chorus 2", "Chorus 3", "Chorus 4", "Chorus 5",
	"Pre-Chorus", "Bridge", "Tag", "Intro", "Outro", "Intermezzo",
}

func isValidTag(tag string) bool {
	for _, validTag := range validTags {
		if validTag == tag {
			return true
		}
	}

	return false
}

func ValidateChordSheet(chordsheet datatypes.JSON) error {
	csSchema := map[string]string{}

	err := json.Unmarshal([]byte(chordsheet.String()), &csSchema)
	if err != nil {
		return fmt.Errorf("could not parse chordsheet: %s", err.Error())
	}

	for key := range csSchema {
		if !isValidTag(key) {
			return fmt.Errorf("%s is not a valid tag", key)
		}
	}

	return nil
}
