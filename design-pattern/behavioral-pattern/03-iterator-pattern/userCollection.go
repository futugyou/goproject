package main

type userCollection struct {
	users []*user
}

func (u *userCollection) createiterator() iterator {
	return &userIterator{users: u.users}
}
