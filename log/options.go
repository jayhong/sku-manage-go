package log

type hookConfig struct {
	_server  string
	_logPath string
	_address string
	_host    string
}

type HookOption func(*hookConfig)

func ServiceOption(name string) HookOption {
	return func(cfg *hookConfig) {
		cfg._server = name
	}
}

func PathOption(logPath string) HookOption {
	return func(cfg *hookConfig) {
		cfg._logPath = logPath
	}
}

func AddressOption(address string) HookOption {
	return func(cfg *hookConfig) {
		cfg._address = address
	}
}

func HostOption(host string) HookOption {
	return func(cfg *hookConfig) {
		cfg._host = host
	}
}
