package appuser

var knownAuthProviders []string

func SetKnownAuthProviders(v []string) {
	if knownAuthProviders != nil {
		panic("knownAuthProviders already set")
	}
	knownAuthProviders = v
}
