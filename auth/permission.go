package auth

// A single permission entry, determines if you get read access, write access, or both.
type Permission struct {
	Read  bool
	Write bool
}

// Map of patterns and the permissions associated with them. Pattern should be
// a string that allows "*" wildcard. E.g. "foo.*" matches anything that starts
// with "foo.", where as "foo." matches only "foo.".
type Permissions map[string]*Permission
