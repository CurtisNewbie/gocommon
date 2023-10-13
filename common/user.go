package common

import (
	"github.com/curtisnewbie/miso/miso"
)

type User struct {
	UserId   int
	UserNo   string
	Username string
	RoleNo   string
	IsNil    bool
}

const (
	UserIdTraceKey   = "id"
	UserNoTraceKey   = "userno"
	UsernameTraceKey = "username"
	RoleNoTraceKey   = "roleno"
)

var (
	nilUser                = User{IsNil: true}
	builtinPropagationKeys = []string{
		UserIdTraceKey,
		UserNoTraceKey,
		UsernameTraceKey,
		RoleNoTraceKey,
	}
)

func init() {
	LoadBuiltinPropagationKeys()
}

func LoadBuiltinPropagationKeys() {
	// load builtin propagation keys, so all dependents get the same behaviour
	for _, v := range builtinPropagationKeys {
		miso.AddPropagationKey(v)
	}
}

// Get a 'nil' User
func NilUser() User {
	return nilUser
}

// Get User from Rail (trace)
func GetUser(rail miso.Rail) User {
	idv := rail.CtxValInt(UserIdTraceKey)
	if idv <= 0 {
		return NilUser()
	}

	return User{
		UserId:   idv,
		Username: rail.CtxValStr(UsernameTraceKey),
		UserNo:   rail.CtxValStr(UserNoTraceKey),
		RoleNo:   rail.CtxValStr(RoleNoTraceKey),
		IsNil:    false,
	}
}
