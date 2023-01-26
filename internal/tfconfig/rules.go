package tfconfig

import "fmt"

func hasValidS3BackendConfig(attrs map[string]string) bool {
	// at least bucket and key should be defined
	bucket, bucketOk := attrs["bucket"]
	key, keyOk := attrs["key"]

	isValidBucket := bucketOk && bucket != ""
	isValidKey := keyOk && key != ""

	return isValidBucket && isValidKey
}

func ModuleShouldHaveRemoteBackend(modules []*Module) Issues {
	issues := Issues{}
	for _, module := range modules {

		if len(module.Backends) == 0 {
			issues = append(issues, &Issue{
				Severity:   ERROR,
				ModulePath: module.Path,
				Message:    "No remote backend configured",
			})

			continue
		}

		for _, backend := range module.Backends {
			// remote backends shouldn't be defined in overrides
			if isOverride(backend.Range.Filename) {
				continue
			}

			switch backend.Type {
			case "s3":
				if !hasValidS3BackendConfig(backend.Attributes) {
					issues = append(issues, &Issue{
						ModulePath: module.Path,
						Message:    "Invalid s3 backend configuration",
						Range:      &backend.Range,
						Severity:   ERROR,
					})
				}

			default:
				issues = append(issues, &Issue{
					Severity:   ERROR,
					ModulePath: module.Path,
					Message:    "No remote backend configured",
					Range:      &backend.Range,
				})
			}
		}
	}

	return issues
}

func hasValidLocalBackendConfig(attrs map[string]string) bool {
	// at least bucket and key should be defined
	path, pathOk := attrs["path"]

	return pathOk && path != ""
}

func ModuleShouldHaveLocalBackendOverride(modules []*Module) Issues {
	issues := Issues{}
	for _, module := range modules {

		var overrideBackend *Backend = nil
		for _, backend := range module.Backends {
			// local backends should only be defined in overrides
			if isOverride(backend.Range.Filename) && backend.Type == "local" {
				overrideBackend = backend
			}
		}

		if overrideBackend == nil {
			issues = append(issues, &Issue{
				Severity:   ERROR,
				ModulePath: module.Path,
				Message:    "No local backend override configured",
			})

			continue
		}

		if !hasValidLocalBackendConfig(overrideBackend.Attributes) {
			issues = append(issues, &Issue{
				ModulePath: module.Path,
				Message:    "Invalid local backend override configuration",
				Range:      &overrideBackend.Range,
				Severity:   ERROR,
			})
		}

		// if len(module.Backends) == 0 {
		// 	issues = append(issues, &Issue{
		// 		Severity:   ERROR,
		// 		ModulePath: module.Path,
		// 		Message:    "No local backend override configured",
		// 	})

		// 	continue
		// }

		// for _, backend := range module.Backends {
		// 	// local backends should only be defined in overrides
		// 	if !isOverride(backend.Range.Filename) {
		// 		continue
		// 	}

		// 	switch backend.Type {
		// 	case "local":
		// 		if !hasValidLocalBackendConfig(backend.Attributes) {
		// 			issues = append(issues, &Issue{
		// 				ModulePath: module.Path,
		// 				Message:    "Invalid local backend override configuration",
		// 				Range:      &backend.Range,
		// 				Severity:   ERROR,
		// 			})
		// 		}

		// 	default:
		// 		issues = append(issues, &Issue{
		// 			Severity:   ERROR,
		// 			ModulePath: module.Path,
		// 			Message:    "No local backend override configured",
		// 		})
		// 	}

		// }
	}

	return issues
}

func getValidBackendKey(backend *Backend) (string, bool) {
	// should be a combination of attributes that
	// make the particular backend unique
	switch backend.Type {
	case "s3":
		if hasValidS3BackendConfig(backend.Attributes) {
			bucket := backend.Attributes["bucket"]
			key := backend.Attributes["key"]

			return fmt.Sprintf("%s/%s", bucket, key), true
		}
	}

	return "", false
}

type ModuleBackend struct {
	ModulePath string
	Backend    *Backend
}

func ModuleShouldHaveUniqueBackend(modules []*Module) Issues {
	issues := Issues{}

	moduleBackends := map[string][]*ModuleBackend{}

	for _, module := range modules {
		for _, backend := range module.Backends {
			// remote backends shouldn't be defined in overrides
			if isOverride(backend.Range.Filename) {
				continue
			}

			if key, valid := getValidBackendKey(backend); valid {
				moduleBackends[key] = append(moduleBackends[key], &ModuleBackend{
					ModulePath: module.Path,
					Backend:    backend,
				})
			}
		}
	}

	// find backends that share remote configuration
	for _, moduleBackends := range moduleBackends {
		if len(moduleBackends) > 1 {
			for _, moduleBackend := range moduleBackends {
				issues = append(issues, &Issue{
					ModulePath: moduleBackend.ModulePath,
					Message:    "Duplicate remote backend configuration",
					Range:      &moduleBackend.Backend.Range,
					Severity:   ERROR,
				})
			}
		}
	}

	return issues
}

func ParseRules(modules []*Module, rules []func([]*Module) Issues) (issues Issues) {
	for _, check := range rules {
		issues = append(issues, check(modules)...)
	}

	return
}
