package main

type button struct {
	command icommand
}

func (b *button) press() {
	b.command.exec()
}
