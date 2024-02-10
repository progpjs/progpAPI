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
	group     *sharedResourceGroup
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

//region sharedResourceGroup

type sharedResourceGroup struct {
	mutex  sync.RWMutex
	ptrMap map[int]*SharedResource
}

func newSharedResourceGroup() *sharedResourceGroup {
	return &sharedResourceGroup{ptrMap: make(map[int]*SharedResource)}
}

func (m *sharedResourceGroup) unSaveResource(res *SharedResource) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.ptrMap, res.id)
}

func (m *sharedResourceGroup) saveResource(id int, res *SharedResource) {
	res.group = m

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ptrMap[id] = res
}

//endregion

func GetSharedResource(id int) *SharedResource {
	idGroup := id % 10
	group := gSharedResourceGroups[idGroup]

	group.mutex.RLock()
	defer group.mutex.RUnlock()
	return group.ptrMap[id]
}

func NewSharedResource(value any, onDispose DisposeSharedResourceF) *SharedResource {
	gSharedResourceNextIdMutex.Lock()
	id := gSharedResourceNextId
	gSharedResourceNextId++
	gSharedResourceNextIdMutex.Unlock()

	m := &SharedResource{id: id, Value: value, onDispose: onDispose}
	runtime.SetFinalizer(m, (*SharedResource).finalizer)

	idGroup := id % 10
	group := gSharedResourceGroups[idGroup]
	group.saveResource(id, m)

	return m
}

type DisposeSharedResourceF func(value any)

var gSharedResourceGroups []*sharedResourceGroup
var gSharedResourceNextId = 1
var gSharedResourceNextIdMutex = sync.Mutex{}

func init() {
	gSharedResourceGroups = make([]*sharedResourceGroup, 10)

	for i := 0; i < 10; i++ {
		gSharedResourceGroups[i] = newSharedResourceGroup()
	}
}
