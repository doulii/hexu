package watcher

import (
	"reflect"
	"strings"
)

func Update[T any](oldVal, newVal T) error {
	return update(reflect.ValueOf(oldVal), reflect.ValueOf(newVal))
}

func update(oldVal, newVal reflect.Value) error {
	if updateIfIsWatcher(oldVal, newVal) {
		return nil
	}

	switch oldVal.Kind() {
	case reflect.Struct:
		for i := 0; i < oldVal.NumField(); i++ {
			oldField := oldVal.Field(i)
			newField := newVal.Field(i)
			if err := update(oldField, newField); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if !oldVal.IsNil() && !newVal.IsNil() {
			return update(oldVal.Elem(), newVal.Elem())
		}

	case reflect.Map:
		// Watchers inside map values are matched by key.
		for _, key := range oldVal.MapKeys() {
			oldMapValue := oldVal.MapIndex(key)
			newMapValue := newVal.MapIndex(key)
			if oldMapValue.IsValid() && newMapValue.IsValid() {
				if err := update(oldMapValue, newMapValue); err != nil {
					return err
				}
			}
		}
	case reflect.Slice, reflect.Array:
		// Watchers inside slice/array cannot be updated safely.
		// Even when lengths are equal, matching by index is unsafe since the slice may have had elements removed and new ones inserted, breaking the correspondence.
	}
	return nil
}

func updateIfIsWatcher(oldVal, newVal reflect.Value) bool {
	if !(isWatcher(oldVal.Type()) && isWatcher(newVal.Type())) {
		return false
	}

	getNew := newVal.Addr().MethodByName("Get")
	setOld := oldVal.Addr().MethodByName("Set")
	setOld.Call(getNew.Call(nil))

	return true
}

func isWatcher(t reflect.Type) bool {
	return typeOrigin(t) == watcherTypeOrigin && t.PkgPath() == watcherPkgPath
}

func typeOrigin(t reflect.Type) string {
	return strings.Split(t.Name(), `[`)[0]
}

var (
	watcherTypeOrigin, watcherPkgPath = func() (string, string) {
		watcherType := reflect.TypeFor[Watcher[any]]()
		return typeOrigin(watcherType), watcherType.PkgPath()
	}()
)
