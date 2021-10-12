package main

type userIterator struct {
	index int
	users []*user
}

func (u *userIterator) hasnext() bool {
	return u.index < len(u.users)
}

func (u *userIterator) getnext() *user {
	if u.hasnext() {
		user := u.users[u.index]
		u.index++
		return user
	}
	return nil
}
