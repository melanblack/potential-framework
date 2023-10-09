package cfghive

import "fmt"

func HiveSize(data map[string]HiveValue) uint {
	size := uint(0)
	for _, v := range data {
		if v.IsStoredType(HiveTypeSub) {
			sub, _ := v.Sub()
			size += HiveSize(sub)
		} else {
			size += 1
		}
	}
	return uint(size)
}

func HiveDump(data *map[string]HiveValue) {
	for k, v := range *data {
		if v.IsStoredType(HiveTypeSub) {
			sub, _ := v.Sub()
			HiveDump(&sub)
		} else {
			fmt.Printf(" k: %s, v: %s (%s)", k, v.TypeString(), v.Value())
		}
	}
}
