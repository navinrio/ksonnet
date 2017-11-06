package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/ksonnet/ksonnet/metadata/app"
	"github.com/ksonnet/ksonnet/metadata/parts"
	"github.com/ksonnet/ksonnet/metadata/registry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// AddRegistry adds a registry with `name`, `protocol`, and `uri` to
// the current ksonnet application.
func (m *manager) AddRegistry(name, protocol, uri, version string) (*registry.Spec, error) {
	appSpec, err := m.AppSpec()
	if err != nil {
		return nil, err
	}

	// Add registry reference to app spec.
	registryManager, err := makeGitHubRegistryManager(&app.RegistryRefSpec{
		Name:     name,
		Protocol: protocol,
		URI:      uri,
	})
	if err != nil {
		return nil, err
	}

	err = appSpec.AddRegistryRef(registryManager.RegistryRefSpec)
	if err != nil {
		return nil, err
	}

	// Retrieve the contents of registry.
	registrySpec, err := m.getOrCacheRegistry(registryManager)
	if err != nil {
		return nil, err
	}

	// Write registry specification back out to app specification.
	specBytes, err := appSpec.Marshal()
	if err != nil {
		return nil, err
	}

	err = afero.WriteFile(m.appFS, string(m.appYAMLPath), specBytes, defaultFilePermissions)
	if err != nil {
		return nil, err
	}

	return registrySpec, nil
}

func (m *manager) GetRegistry(name string) (*registry.Spec, string, error) {
	registryManager, protocol, err := m.getRegistryManager(name)
	if err != nil {
		return nil, "", err
	}

	regSpec, exists, err := m.registrySpecFromFile(m.registryPath(registryManager))
	if !exists {
		return nil, "", fmt.Errorf("Registry '%s' does not exist", name)
	} else if err != nil {
		return nil, "", err
	}

	return regSpec, protocol, nil
}

func (m *manager) CacheDependency(registryName, libID, libName, libVersion string) (*parts.Spec, error) {
	// Retrieve application specification.
	appSpec, err := m.AppSpec()
	if err != nil {
		return nil, err
	}

	if _, ok := appSpec.Libraries[libName]; ok {
		return nil, fmt.Errorf("Library '%s' already exists", libName)
	}

	// Retrieve registry manager for this specific registry.
	regRefSpec, exists := appSpec.GetRegistryRef(registryName)
	if !exists {
		return nil, fmt.Errorf("Registry '%s' does not exist", registryName)
	}

	registryManager, _, err := m.getRegistryManagerFor(regRefSpec)
	if err != nil {
		return nil, err
	}

	// Get all directories and files first, then write to disk. This
	// protects us from failing with a half-cached dependency because of
	// a network failure.
	directories := []AbsPath{}
	files := map[AbsPath][]byte{}
	parts, libRef, err := registryManager.ResolveLibrary(
		libID,
		libName,
		libVersion,
		func(relPath string, contents []byte) error {
			files[appendToAbsPath(m.vendorPath, relPath)] = contents
			return nil
		},
		func(relPath string) error {
			directories = append(directories, appendToAbsPath(m.vendorPath, relPath))
			return nil
		})
	if err != nil {
		return nil, err
	}

	// Add library to app specification, but wait to write it out until
	// the end, in case one of the network calls fails.
	appSpec.Libraries[libName] = libRef
	appSpecData, err := appSpec.Marshal()
	if err != nil {
		return nil, err
	}

	log.Infof("Retrieved %d files", len(files))

	for _, dir := range directories {
		if err := m.appFS.MkdirAll(string(dir), defaultFolderPermissions); err != nil {
			return nil, err
		}
	}

	for path, content := range files {
		if err := afero.WriteFile(m.appFS, string(path), content, defaultFilePermissions); err != nil {
			return nil, err
		}
	}

	err = afero.WriteFile(m.appFS, string(m.appYAMLPath), appSpecData, defaultFilePermissions)
	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (m *manager) registryDir(regManager registry.Manager) AbsPath {
	return appendToAbsPath(m.registriesPath, regManager.RegistrySpecDir())
}

func (m *manager) registryPath(regManager registry.Manager) AbsPath {
	return appendToAbsPath(m.registriesPath, regManager.RegistrySpecFilePath())
}

func (m *manager) getRegistryManager(registryName string) (registry.Manager, string, error) {
	appSpec, err := m.AppSpec()
	if err != nil {
		return nil, "", err
	}

	regRefSpec, exists := appSpec.GetRegistryRef(registryName)
	if !exists {
		return nil, "", fmt.Errorf("Registry '%s' does not exist", registryName)
	}

	return m.getRegistryManagerFor(regRefSpec)
}

func (m *manager) getRegistryManagerFor(registryRefSpec *app.RegistryRefSpec) (registry.Manager, string, error) {
	var err error
	var manager registry.Manager
	var protocol string

	switch registryRefSpec.Protocol {
	case "github":
		manager, err = makeGitHubRegistryManager(registryRefSpec)
		protocol = "github"
	default:
		return nil, "", fmt.Errorf("Invalid protocol '%s'", registryRefSpec.Protocol)
	}

	if err != nil {
		return nil, "", err
	}

	return manager, protocol, nil
}

func (m *manager) registrySpecFromFile(path AbsPath) (*registry.Spec, bool, error) {
	registrySpecFile := string(path)
	exists, err := afero.Exists(m.appFS, registrySpecFile)
	if err != nil {
		return nil, false, err
	}

	isDir, err := afero.IsDir(m.appFS, registrySpecFile)
	if err != nil {
		return nil, false, err
	}

	// NOTE: case where directory of the same name exists should be
	// fine, most filesystems allow you to have a directory and file of
	// the same name.
	if exists && !isDir {
		registrySpecBytes, err := afero.ReadFile(m.appFS, registrySpecFile)
		if err != nil {
			return nil, false, err
		}

		registrySpec := registry.Spec{}
		err = json.Unmarshal(registrySpecBytes, &registrySpec)
		if err != nil {
			return nil, false, err
		}
		return &registrySpec, true, nil
	}

	return nil, false, nil
}

func (m *manager) getOrCacheRegistry(gh registry.Manager) (*registry.Spec, error) {
	// Check local disk cache.
	registrySpecFile := m.registryPath(gh)
	registrySpec, exists, err := m.registrySpecFromFile(registrySpecFile)
	if !exists {
		return nil, fmt.Errorf("Registry '%s' does not exist", gh.MakeRegistryRefSpec().Name)
	} else if err != nil {
		return nil, err
	}

	// If failed, use the protocol to try to retrieve app specification.
	registrySpec, err = gh.FetchRegistrySpec()
	if err != nil {
		return nil, err
	}

	registrySpecBytes, err := registrySpec.Marshal()
	if err != nil {
		return nil, err
	}

	// NOTE: We call mkdir after getting the registry spec, since a
	// network call might fail and leave this half-initialized empty
	// directory.
	registrySpecDir := appendToAbsPath(m.registriesPath, gh.RegistrySpecDir())
	err = m.appFS.MkdirAll(string(registrySpecDir), defaultFolderPermissions)
	if err != nil {
		return nil, err
	}

	err = afero.WriteFile(m.appFS, string(registrySpecFile), registrySpecBytes, defaultFilePermissions)
	if err != nil {
		return nil, err
	}

	return registrySpec, nil
}