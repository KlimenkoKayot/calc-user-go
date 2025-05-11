package orderedmap

import (
	"sync"

	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

// OrderedMap — структура, которая сохраняет порядок элементов и позволяет обращаться по ключу
type OrderedMap struct {
	mu    *sync.RWMutex
	keys  []uint                             // Слайс для сохранения порядка ключей
	items map[uint]*models.RequestExpression // Мапа для хранения данных
	Len   int
}

// NewOrderedMap создает новый экземпляр OrderedMap
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		mu:    &sync.RWMutex{},
		keys:  make([]uint, 0),
		items: make(map[uint]*models.RequestExpression, 0),
		Len:   0,
	}
}

// Set добавляет или обновляет значение по ключу
func (om *OrderedMap) Set(key uint, value *models.RequestExpression) {
	om.mu.Lock()
	defer om.mu.Unlock()

	// Если ключ уже существует, обновляем значение
	if _, exists := om.items[key]; !exists {
		// Если ключ новый, добавляем его в слайс
		om.Len++
		om.keys = append(om.keys, key)
	}
	om.items[key] = value
}

// Возвращает значение по ключу
func (om *OrderedMap) Get(key uint) (*models.RequestExpression, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	value, exists := om.items[key]
	return value, exists
}

// Удаляет элемент по ключу
func (om *OrderedMap) Delete(key uint) {
	om.mu.Lock()
	defer om.mu.Unlock()

	// Удаляем ключ из мапы за O(1)
	delete(om.items, key)

	// Удаляем ключ из слайса за O(N)
	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			om.Len--
			break
		}
	}
}
