package main

func main() {
	tv := &tv{}
	onCommand := &oncommand{
		device: tv,
	}

	offCommand := &offcommand{
		device: tv,
	}

	onButton := &button{
		command: onCommand,
	}

	onButton.press()

	offButton := &button{
		command: offCommand,
	}
	offButton.press()
}
