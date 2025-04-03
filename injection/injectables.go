package injection

import "reflect"

func findInjectable(injectables map[reflect.Type]reflect.Value, t reflect.Type) reflect.Value {
	found := injectables[t]
	if found.IsValid() {
		return found
	}
	if t.Kind() != reflect.Interface {
		return reflect.Value{}
	}
	for _, i := range injectables {
		if i.Type().Implements(t) {
			return i
		}
	}
	return reflect.Value{}
}

func transformArgs(args ...any) map[reflect.Type]reflect.Value {
	values := map[reflect.Type]reflect.Value{}
	for _, a := range args {
		values[reflect.TypeOf(a)] = reflect.ValueOf(a)
	}
	return values
}
