// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_providers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	base "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers/__base"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

const POSITION_FILE = "_position.yaml"
const PROVIDER_COUNT = 52

var Factory = ModelProviderFactory{}

type ModelProviderFactory struct{}

type ModelProviderExtension struct {
	ProviderInstance *base.ModelProvider
	Name             string
	Position         int
}

func (f *ModelProviderFactory) GetProvidersFromDir() ([]*entities.ProviderEntity, error) {
	modelProviderExtensions, err := f.getMapProvidersExtensions()
	if err != nil {
		return nil, err
	}

	providerEntities, err := f.extensionConvertProviderEntity(modelProviderExtensions)

	if err != nil {
		return nil, err
	}

	return providerEntities, nil
}

func (f *ModelProviderFactory) extensionConvertProviderEntity(
	modelProviderExtensions map[string]*ModelProviderExtension,
) ([]*entities.ProviderEntity, error) {

	providers := make([]*entities.ProviderEntity, 0, PROVIDER_COUNT)

	for _, providerExtension := range modelProviderExtensions {
		modelProviderInstance := providerExtension.ProviderInstance
		if provider, err := modelProviderInstance.GetProviderSchema(); err != nil {
			return nil, err
		} else {
			provider.Position = providerExtension.Position
			providers = append(providers, provider)
		}
	}
	return providers, nil
}

// func (f *ModelProviderFactory) sortProviderEntityByPosition(
// 	providers []*entities.ProviderEntity,
// 	providerPositionMap map[string]int,
// ) {
// 	sort.Slice(providers, func(i, j int) bool {
// 		return providerPositionMap[providers[i].Provider] < providerPositionMap[providers[j].Provider]
// 	})
// }

func (f *ModelProviderFactory) GetPositionMap(fileDir string) (map[string]int, error) {
	positionInfo := make([]string, 0, PROVIDER_COUNT)
	positionFilePath := filepath.Join(fileDir, POSITION_FILE)
	positionFileContent, err := os.ReadFile(positionFilePath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	if err := yaml.Unmarshal(positionFileContent, &positionInfo); err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	positionIndexMap := make(map[string]int)

	for index, providerName := range positionInfo {
		positionIndexMap[strings.Trim(providerName, " ")] = index
	}

	return positionIndexMap, nil
}

func (f *ModelProviderFactory) resolveProviderExtensions(
	modelProviderResolvePaths []string,
	positionMap map[string]int,
) []*ModelProviderExtension {
	modelProviderExtensions := make([]*ModelProviderExtension, 0, PROVIDER_COUNT)
	for _, path := range modelProviderResolvePaths {
		modelProviderName := filepath.Base(path)
		modelProviderExtension := &ModelProviderExtension{
			Name:             modelProviderName,
			ProviderInstance: &base.ModelProvider{ModelConfPath: path},
			Position:         positionMap[modelProviderName],
		}

		modelProviderExtensions = append(modelProviderExtensions, modelProviderExtension)
	}

	return modelProviderExtensions
}

func (f *ModelProviderFactory) resolveProviderDirInfo() ([]fs.DirEntry, string, string, error) {
	_, fullFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return nil, "", "", errors.WithCode(code.ErrRunTimeCaller, "Fail to get runtime caller info")
	}

	fileDir := filepath.Dir(fullFilePath)

	dirEntries, err := os.ReadDir(fileDir)

	if err != nil {
		return nil, "", "", errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	return dirEntries, fullFilePath, fileDir, nil
}

func (f *ModelProviderFactory) resolveProviderDir(dirEntries []fs.DirEntry, fullFilePath string) ([]string, error) {
	modelProviderResolvePaths := make([]string, 0, PROVIDER_COUNT)

	providerDir := filepath.Dir(fullFilePath)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() && !strings.HasPrefix(dirEntry.Name(), "__") {
			modelProviderResolvePaths = append(
				modelProviderResolvePaths,
				fmt.Sprintf("%s/%s", providerDir, dirEntry.Name()),
			)
		}
	}
	return modelProviderResolvePaths, nil
}

func (f *ModelProviderFactory) getMapProvidersExtensions() (map[string]*ModelProviderExtension, error) {
	dirEntries, fullFilePath, fileDir, err := f.resolveProviderDirInfo()

	if err != nil {
		return nil, err
	}

	modelProviderResolvePaths, err := f.resolveProviderDir(dirEntries, fullFilePath)
	if err != nil {
		return nil, err
	}
	positionMap, err := f.GetPositionMap(fileDir)

	if err != nil {
		return nil, err
	}

	resolveProviderExtensions := f.resolveProviderExtensions(modelProviderResolvePaths, positionMap)

	return f.resolveMapProviderExtensions(resolveProviderExtensions), nil
}

func (f *ModelProviderFactory) resolveMapProviderExtensions(
	providerExtensions []*ModelProviderExtension,
) map[string]*ModelProviderExtension {

	providerMap := make(map[string]*ModelProviderExtension)

	for _, providerExtension := range providerExtensions {
		providerMap[providerExtension.Name] = providerExtension
	}
	return providerMap
}
