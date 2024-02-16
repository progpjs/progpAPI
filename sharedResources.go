/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package progpAPI

import (
	"runtime"
	"sync"
)

//region SharedResource

type SharedResource struct {
	id        int
	Value     any
	group     *SharedResourceContainer
	onDispose DisposeSharedResourceF
}

func (m *SharedResource) finalizer() {
	m.Dispose()
}

func (m *SharedResource) GetId() int {
	return m.id
}

func (m *SharedResource) Dispose() {
	if m.group != nil {
		og := m.group
		m.group = nil
		og.unSaveResource(m)

		if m.onDispose != nil {
			m.onDispose(m.Value)
		}
	}
}

//endregion

//region SharedResourceContainer

type SharedResourceContainer struct {
	scriptIsolate ScriptIsolate

	nextResourceId int
	resourceMap    map[int]*SharedResource
	resourcesMutex sync.RWMutex

	parentContainer      *SharedResourceContainer
	nextContainer        *SharedResourceContainer
	previousContainer    *SharedResourceContainer
	childContainerHead   *SharedResourceContainer
	childContainersMutex sync.Mutex
}

func NewSharedResourceContainer(parent *SharedResourceContainer, iso ScriptIsolate) *SharedResourceContainer {
	m := &SharedResourceContainer{resourceMap: make(map[int]*SharedResource), scriptIsolate: iso}

	if parent != nil {
		parent.saveChildContainer(m)
	}

	return m
}

func (m *SharedResourceContainer) Dispose() {
	if m.parentContainer != nil {
		m.parentContainer.unSaveChildContainer(m)
	}

	for _, res := range m.resourceMap {
		res.Dispose()
	}
}

func (m *SharedResourceContainer) GetResource(resId int) *SharedResource {
	m.resourcesMutex.RLock()
	r := m.resourceMap[resId]
	m.resourcesMutex.RUnlock()
	return r
}

func (m *SharedResourceContainer) NewSharedResource(value any, onDispose DisposeSharedResourceF) *SharedResource {
	res := &SharedResource{Value: value, onDispose: onDispose}
	runtime.SetFinalizer(res, (*SharedResource).finalizer)

	m.resourcesMutex.Lock()
	id := m.nextResourceId
	m.nextResourceId++
	m.resourceMap[id] = res
	m.resourcesMutex.Unlock()

	res.group = m
	res.id = id

	return res
}

func (m *SharedResourceContainer) GetIsolate() ScriptIsolate {
	return m.scriptIsolate
}

func (m *SharedResourceContainer) saveChildContainer(child *SharedResourceContainer) {
	m.childContainersMutex.Lock()
	m.childContainersMutex.Unlock()

	child.parentContainer = m
	child.nextContainer = m.childContainerHead
	m.childContainerHead = child

	if child.nextContainer != nil {
		child.nextContainer.previousContainer = child
	}
}

func (m *SharedResourceContainer) unSaveChildContainer(child *SharedResourceContainer) {
	m.childContainersMutex.Lock()
	m.childContainersMutex.Unlock()

	if m.childContainerHead == child {
		m.childContainerHead = child.nextContainer
	}

	if child.nextContainer != nil {
		child.nextContainer.previousContainer = child.previousContainer
	}

	if child.previousContainer != nil {
		child.previousContainer.nextContainer = child.nextContainer
	}
}

func (m *SharedResourceContainer) unSaveResource(res *SharedResource) {
	m.resourcesMutex.Lock()
	delete(m.resourceMap, res.id)
	m.resourcesMutex.Unlock()
}

//endregion

func newSharedResource(value any, onDispose DisposeSharedResourceF) *SharedResource {
	m := &SharedResource{Value: value, onDispose: onDispose}

	// Make things way faster when dealing with self-managed resources.
	//
	if onDispose != nil {
		runtime.SetFinalizer(m, (*SharedResource).finalizer)
	}

	return m
}

type DisposeSharedResourceF func(value any)
