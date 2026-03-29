package injection

import (
	"reflect"
	"sort"
	"strings"
)

func findInjectable(injectables map[reflect.Type]reflect.Value, t reflect.Type) (reflect.Value, error) {
	found := injectables[t]
	if found.IsValid() {
		return found, nil
	}
	if t.Kind() != reflect.Interface {
		return reflect.Value{}, nil
	}
	matches := make([]reflect.Value, 0, 1)
	matchTypes := make([]string, 0, 1)
	for _, i := range injectables {
		if i.Type().Implements(t) {
			matches = append(matches, i)
			matchTypes = append(matchTypes, i.Type().String())
		}
	}
	switch len(matches) {
	case 0:
		return reflect.Value{}, nil
	case 1:
		return matches[0], nil
	default:
		sort.Strings(matchTypes)
		return reflect.Value{}, errAmbiguousInjectable(t, strings.Join(matchTypes, ", "))
	}
}

func transformArgs(args ...any) map[reflect.Type]reflect.Value {
	values := map[reflect.Type]reflect.Value{}
	for _, a := range args {
		if a == nil {
			continue
		}
		values[reflect.TypeOf(a)] = reflect.ValueOf(a)
	}
	return values
}
