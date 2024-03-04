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

func (m *SharedResource) GetContainer() *SharedResourceContainer {
	return m.group
}

func (m *SharedResource) GetId() int {
	return m.id
}

func (m *SharedResource) Dispose() {
	// Go allows calling a function on nil.
	// It's a good thing because it allows the caller to avoid
	// checking if the resource is nil, which generally mean it's already disposed.
	//
	if m == nil {
		return
	}

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
	scriptContext JsContext
	isDisposed    bool

	nextResourceId int
	resourceMap    map[int]*SharedResource
	resourcesMutex sync.RWMutex

	parentContainer      *SharedResourceContainer
	nextContainer        *SharedResourceContainer
	previousContainer    *SharedResourceContainer
	childContainerHead   *SharedResourceContainer
	childContainersMutex sync.Mutex
}

func NewSharedResourceContainer(parent *SharedResourceContainer, ctx JsContext) *SharedResourceContainer {
	if (ctx == nil) && (parent != nil) {
		ctx = parent.scriptContext
	}

	m := &SharedResourceContainer{resourceMap: make(map[int]*SharedResource), scriptContext: ctx}

	if parent != nil {
		parent.saveChildContainer(m)
	}

	return m
}

func (m *SharedResourceContainer) CreateLock() (*JavascriptLock, *SharedResource) {
	jsLock := &JavascriptLock{}
	jsLock.Wait()
	return jsLock, m.NewSharedResource(jsLock, nil)
}

func (m *SharedResourceContainer) Dispose() {
	if m.isDisposed {
		return
	}
	m.isDisposed = true

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

	// Warning: resources are stored as a double in v8 side
	// doing that we can send a memory pointer, which can
	// exceed the size of a double. We don't use v8::external
	// the reason being than his memory isn't freed in the same
	// GC cycles doing that the memory can saturate in high load.

	m.resourcesMutex.Lock()
	id := m.nextResourceId

	if id > MaxResourceIdSize {
		m.resourcesMutex.Unlock()

		// This allows avoiding using too big integer
		// and going over double to int conversion capacity.
		//
		id = m.compactResourceId()
	}

	m.nextResourceId++
	m.resourceMap[id] = res
	m.resourcesMutex.Unlock()

	res.group = m
	res.id = id

	return res
}

const MaxResourceIdSize = 2147483647

var gCompactingMutex sync.Mutex

func (m *SharedResourceContainer) compactResourceId() int {
	gCompactingMutex.Lock()
	defer gCompactingMutex.Unlock()

	// Called by function while was in pause?
	if m.nextResourceId <= MaxResourceIdSize {
		m.resourcesMutex.Lock()
		return m.nextResourceId
	}

	// The pause isn't include in the caller lock,
	// which allows to free the current resources
	// while this pause is executing.
	//
	PauseMs(100)

	// We lock but we will not unlock
	// in order to let the caller use the lock
	// when exiting this function.
	//
	m.resourcesMutex.Lock()

	maxId := 0

	for key, _ := range m.resourceMap {
		if key > maxId {
			maxId = key
		}
	}

	maxId++

	println("Max id was ", m.nextResourceId, " and is now ", maxId)
	m.nextResourceId = maxId

	return maxId
}

func (m *SharedResourceContainer) GetScriptContext() JsContext {
	return m.scriptContext
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

type JavascriptLock struct {
	lock sync.Mutex
}

func (m *JavascriptLock) Wait() {
	m.lock.Lock()
}

func (m *JavascriptLock) Unlock() {
	m.lock.Unlock()
}

// OnReturnVoidAction allows to use the javascript 'progpReturnVoid'
// on a resource pointing to this object.
func (m *JavascriptLock) OnReturnVoidAction() error {
	m.Unlock()
	return nil
}
