package tfconfig

func ModuleShouldHaveRemoteBackend(modules []*Module) Issues {
	issues := Issues{}
	for _, module := range modules {
		if len(module.Backends) == 0 {
			issues = append(issues, &Issue{
				Severity:   ERROR,
				ModulePath: module.Path,
				Message:    "No remote backend configured",
			})
		}
	}

	return issues
}

func ParseRules(modules []*Module) (issues Issues) {
	issues = append(issues, ModuleShouldHaveRemoteBackend(modules)...)

	return issues
}
