package telegram_event

const msgHelp = `Welcome to the Real Estate Finder Bot! Use /search to find properties
based on filters like price, location, and type.
Set up /notify to get alerts for new listings matching your preferences.
Type /help for more information.`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
)
