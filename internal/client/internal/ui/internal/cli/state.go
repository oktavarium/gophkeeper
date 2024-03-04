package cli

type State int

const (
	MainState = iota
	LoginState
	RegisterState
	WorkState
	SettingsState
	StoreState
	LoginStoreState
)
