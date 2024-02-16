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
	nextResourceId     int
	ptrMap             map[int]*SharedResource
	childContainerHead *SharedResourceContainer

	parent        *SharedResourceContainer
	next          *SharedResourceContainer
	previous      *SharedResourceContainer
	scriptIsolate ScriptIsolate
}

func NewSharedResourceContainer(parent *SharedResourceContainer, iso ScriptIsolate) *SharedResourceContainer {
	m := &SharedResourceContainer{ptrMap: make(map[int]*SharedResource), scriptIsolate: iso}

	if parent != nil {
		parent.saveChildContainer(m)
	}

	return m
}

func (m *SharedResourceContainer) Dispose() {
	if m.parent != nil {
		m.parent.unSaveChildContainer(m)
	}

	for _, res := range m.ptrMap {
		res.Dispose()
	}
}

func (m *SharedResourceContainer) GetResource(resId int) *SharedResource {
	return m.ptrMap[resId]
}

func (m *SharedResourceContainer) NewSharedResource(value any, onDispose DisposeSharedResourceF) *SharedResource {
	id := m.nextResourceId
	m.nextResourceId++
	res := newSharedResource(id, value, onDispose)
	m.saveResource(id, res)
	return res
}

func (m *SharedResourceContainer) GetIsolate() ScriptIsolate {
	return m.scriptIsolate
}

func (m *SharedResourceContainer) saveChildContainer(child *SharedResourceContainer) {
	child.parent = m
	child.next = m.childContainerHead
	m.childContainerHead = child

	if child.next != nil {
		child.next.previous = child
	}
}

func (m *SharedResourceContainer) unSaveChildContainer(child *SharedResourceContainer) {
	if m.childContainerHead == child {
		m.childContainerHead = child.next
	}

	if child.next != nil {
		child.next.previous = child.previous
	}

	if child.previous != nil {
		child.previous.next = child.next
	}
}

func (m *SharedResourceContainer) unSaveResource(res *SharedResource) {
	delete(m.ptrMap, res.id)
}

func (m *SharedResourceContainer) saveResource(id int, res *SharedResource) {
	res.group = m
	m.ptrMap[id] = res
}

//endregion

func newSharedResource(id int, value any, onDispose DisposeSharedResourceF) *SharedResource {
	m := &SharedResource{id: id, Value: value, onDispose: onDispose}
	runtime.SetFinalizer(m, (*SharedResource).finalizer)
	return m
}

type DisposeSharedResourceF func(value any)
