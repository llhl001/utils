package redis

import vmap "utils/container/map"

var (
	// Instance map
	instances = vmap.NewStrAnyMap(true)
)

// Instance returns an instance of redis client with specified group.
// The <name> param is unnecessary, if <name> is not passed,
// it returns a redis instance with default configuration group.
func Instance(name ...string) *Redis {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	v := instances.GetOrSetFuncLock(group, func() interface{} {
		if config, ok := GetConfig(group); ok {
			r := New(config)
			r.group = group
			return r
		}
		return nil
	})
	if v != nil {
		return v.(*Redis)
	}
	return nil
}
