package bot

import (
	"fmt"
	"strings"
)

var tags = [...]string{
	"language",
	"parody",
	"character",
	"group",
	"artist",
	"male",
	"female",
	"misc",
}

func noCategory(tag string) bool {
	return !strings.Contains(tag, ":")
}

func groupTags(tagPairs []string) TagMap {
	tagMap := make(TagMap)

	for _, tagPair := range tagPairs {
		tagSplit := strings.Split(tagPair, ":")
		if len(tagSplit) == 1 {
			tagMap["misc"] = append(tagMap["misc"], tagSplit[0])
		} else {
			tagMap[tagSplit[0]] = append(tagMap[tagSplit[0]], tagSplit[1])
		}
	}

	return tagMap
}

func formatTags(tagMap TagMap) (rawmsg string) {
	for _, tag := range tags {
		if tagVals, ok := tagMap[tag]; ok {
			rawmsg += fmt.Sprintf("`%s`:", tag)
			for i, tagVal := range tagVals {
				if i == len(tagVals)-1 {
					rawmsg += fmt.Sprintf("`%s`\n", tagVal)
				} else {
					rawmsg += fmt.Sprintf("`%s` ", tagVal)
				}
			}
		}
	}

	return
}
