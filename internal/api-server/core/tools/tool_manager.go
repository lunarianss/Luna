package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"gopkg.in/yaml.v3"
)

type ToolManager struct {
	providerRuntimes []*biz_entity.ToolProviderRuntime
}

func NewToolManager() *ToolManager {
	return &ToolManager{}
}

func (tm *ToolManager) ListBuiltInProviders() ([]*biz_entity.ToolProviderRuntime, error) {

	if err := tm.resolveRuntimePath(); err != nil {
		return nil, err
	}

	if err := tm.unMarshalProvider(); err != nil {
		return nil, err
	}

	if err := tm.unMarshalTools(); err != nil {
		return nil, err
	}

	return tm.providerRuntimes, nil
}

func (tm *ToolManager) unMarshalTools() error {
	var tools []*biz_entity.ToolRuntimeConfiguration
	for _, providerRuntime := range tm.providerRuntimes {
		tools = make([]*biz_entity.ToolRuntimeConfiguration, 0, 3)
		toolPath := fmt.Sprintf("%s/%s", providerRuntime.ConfPath, "tools")

		dirs, err := os.ReadDir(toolPath)

		if err != nil {
			return errors.WithSCode(code.ErrRunTimeCaller, err.Error())
		}

		for _, dir := range dirs {
			if dir.IsDir() || !strings.HasSuffix(dir.Name(), ".yaml") {
				continue
			}

			toolYamlPath := fmt.Sprintf("%s/%s", toolPath, dir.Name())
			toolYamlByte, err := os.ReadFile(toolYamlPath)

			if err != nil {
				return errors.WithSCode(code.ErrRunTimeCaller, err.Error())
			}

			toolRuntime := &biz_entity.ToolRuntimeConfiguration{
				ToolStaticConfiguration: &biz_entity.ToolStaticConfiguration{},
			}

			if err := yaml.Unmarshal(toolYamlByte, toolRuntime.ToolStaticConfiguration); err != nil {
				return errors.WithSCode(code.ErrDecodingYaml, err.Error())
			}

			tools = append(tools, toolRuntime)

			providerRuntime.Tools = tools
		}
	}

	return nil
}

func (tm *ToolManager) unMarshalProvider() error {
	for _, providerRuntime := range tm.providerRuntimes {
		confPath := providerRuntime.ConfPath
		toolProviderPath := fmt.Sprintf("%s/%s.yaml", confPath, providerRuntime.ToolProviderName)

		toolBytes, err := os.ReadFile(toolProviderPath)

		if err != nil {
			return errors.WithSCode(code.ErrRunTimeCaller, err.Error())
		}

		if err := yaml.Unmarshal(toolBytes, providerRuntime.ToolProviderStatic); err != nil {
			return errors.WithSCode(code.ErrDecodingYaml, err.Error())
		}
	}
	return nil
}

func (tm *ToolManager) ResolveProviderPath(provider string) (string, error) {
	_, fullFilePath, _, ok := runtime.Caller(0)

	providerRuntime := &biz_entity.ToolProviderRuntime{
		ToolProviderStatic: &biz_entity.ToolProviderStatic{},
	}

	if !ok {
		return "", errors.WithSCode(code.ErrRunTimeCaller, "Fail to get runtime caller info")
	}

	toolsDir := filepath.Join(filepath.Dir(fullFilePath), "provider", "builtin", provider)
	toolProviderPath := fmt.Sprintf("%s/%s.yaml", toolsDir, provider)

	toolBytes, err := os.ReadFile(toolProviderPath)

	if err != nil {
		return "", errors.WithSCode(code.ErrRunTimeCaller, err.Error())
	}

	if err := yaml.Unmarshal(toolBytes, providerRuntime.ToolProviderStatic); err != nil {
		return "", errors.WithSCode(code.ErrDecodingYaml, err.Error())
	}

	return fmt.Sprintf("%s/_assets/%s", toolsDir, providerRuntime.Identity.Icon), nil
}

func (tm *ToolManager) resolveRuntimePath() error {
	var providerRuntimes []*biz_entity.ToolProviderRuntime

	_, fullFilePath, _, ok := runtime.Caller(0)

	if !ok {
		return errors.WithSCode(code.ErrRunTimeCaller, "Fail to get runtime caller info")
	}

	toolRootDir := filepath.Dir(fullFilePath)

	toolsDir := filepath.Join(toolRootDir, "provider", "builtin")

	dirEntries, err := os.ReadDir(toolsDir)

	if err != nil {
		return errors.WithSCode(code.ErrRunTimeCaller, err.Error())
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			providerRuntimes = append(providerRuntimes, &biz_entity.ToolProviderRuntime{
				ConfPath:           fmt.Sprintf("%s/%s", toolsDir, dirEntry.Name()),
				ToolProviderStatic: &biz_entity.ToolProviderStatic{},
				ToolProviderName:   dirEntry.Name(),
			})
		}
	}

	tm.providerRuntimes = providerRuntimes
	return nil
}
