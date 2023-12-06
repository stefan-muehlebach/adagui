// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
	"bytes"
)

// Bool supports binding a bool value.
type Bool interface {
	DataItem
	Get() (bool)
	Set(bool)
}

type boundBool struct {
	base
	val *bool
}

// NewBool returns a bindable bool value that is managed internally.
func NewBool() Bool {
	var blank bool = false
	b := &boundBool{val: &blank}
        b.Init(b)
        return b
}

func (b *boundBool) Get() (bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return false
	}
	return *b.val
}

func (b *boundBool) Set(val bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalBool supports binding a bool value to an external value.
type ExternalBool interface {
	Bool
	Reload()
}

type boundExternalBool struct {
	boundBool
	old bool
}

// BindBool returns a new bindable value that controls the contents of the provided bool variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindBool(v *bool) ExternalBool {
	if v == nil {
		var blank bool = false
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalBool{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalBool) Set(val bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalBool) Reload() {
	b.Set(*b.val)
}

// Bytes supports binding a []byte value.
type Bytes interface {
	DataItem
	Get() ([]byte)
	Set([]byte)
}

type boundBytes struct {
	base
	val *[]byte
}

// NewBytes returns a bindable []byte value that is managed internally.
func NewBytes() Bytes {
	var blank []byte = nil
	b := &boundBytes{val: &blank}
        b.Init(b)
        return b
}

func (b *boundBytes) Get() ([]byte) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return nil
	}
	return *b.val
}

func (b *boundBytes) Set(val []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if bytes.Equal(*b.val, val) {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalBytes supports binding a []byte value to an external value.
type ExternalBytes interface {
	Bytes
	Reload()
}

type boundExternalBytes struct {
	boundBytes
	old []byte
}

// BindBytes returns a new bindable value that controls the contents of the provided []byte variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindBytes(v *[]byte) ExternalBytes {
	if v == nil {
		var blank []byte = nil
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalBytes{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalBytes) Set(val []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if bytes.Equal(b.old, val) {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalBytes) Reload() {
	b.Set(*b.val)
}

// Float supports binding a float64 value.
type Float interface {
	DataItem
	Get() (float64)
	Set(float64)
}

type boundFloat struct {
	base
	val *float64
}

// NewFloat returns a bindable float64 value that is managed internally.
func NewFloat() Float {
	var blank float64 = 0.0
	b := &boundFloat{val: &blank}
        b.Init(b)
        return b
}

func (b *boundFloat) Get() (float64) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return 0.0
	}
	return *b.val
}

func (b *boundFloat) Set(val float64) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalFloat supports binding a float64 value to an external value.
type ExternalFloat interface {
	Float
	Reload()
}

type boundExternalFloat struct {
	boundFloat
	old float64
}

// BindFloat returns a new bindable value that controls the contents of the provided float64 variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindFloat(v *float64) ExternalFloat {
	if v == nil {
		var blank float64 = 0.0
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalFloat{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalFloat) Set(val float64) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalFloat) Reload() {
	b.Set(*b.val)
}

// Int supports binding a int value.
type Int interface {
	DataItem
	Get() (int)
	Set(int)
}

type boundInt struct {
	base
	val *int
}

// NewInt returns a bindable int value that is managed internally.
func NewInt() Int {
	var blank int = 0
	b := &boundInt{val: &blank}
        b.Init(b)
        return b
}

func (b *boundInt) Get() (int) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return 0
	}
	return *b.val
}

func (b *boundInt) Set(val int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalInt supports binding a int value to an external value.
type ExternalInt interface {
	Int
	Reload()
}

type boundExternalInt struct {
	boundInt
	old int
}

// BindInt returns a new bindable value that controls the contents of the provided int variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindInt(v *int) ExternalInt {
	if v == nil {
		var blank int = 0
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalInt{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalInt) Set(val int) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalInt) Reload() {
	b.Set(*b.val)
}

// Rune supports binding a rune value.
type Rune interface {
	DataItem
	Get() (rune)
	Set(rune)
}

type boundRune struct {
	base
	val *rune
}

// NewRune returns a bindable rune value that is managed internally.
func NewRune() Rune {
	var blank rune = rune(0)
	b := &boundRune{val: &blank}
        b.Init(b)
        return b
}

func (b *boundRune) Get() (rune) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return rune(0)
	}
	return *b.val
}

func (b *boundRune) Set(val rune) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalRune supports binding a rune value to an external value.
type ExternalRune interface {
	Rune
	Reload()
}

type boundExternalRune struct {
	boundRune
	old rune
}

// BindRune returns a new bindable value that controls the contents of the provided rune variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindRune(v *rune) ExternalRune {
	if v == nil {
		var blank rune = rune(0)
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalRune{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalRune) Set(val rune) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalRune) Reload() {
	b.Set(*b.val)
}

// String supports binding a string value.
type String interface {
	DataItem
	Get() (string)
	Set(string)
}

type boundString struct {
	base
	val *string
}

// NewString returns a bindable string value that is managed internally.
func NewString() String {
	var blank string = ""
	b := &boundString{val: &blank}
        b.Init(b)
        return b
}

func (b *boundString) Get() (string) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return ""
	}
	return *b.val
}

func (b *boundString) Set(val string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalString supports binding a string value to an external value.
type ExternalString interface {
	String
	Reload()
}

type boundExternalString struct {
	boundString
	old string
}

// BindString returns a new bindable value that controls the contents of the provided string variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindString(v *string) ExternalString {
	if v == nil {
		var blank string = ""
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalString{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalString) Set(val string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalString) Reload() {
	b.Set(*b.val)
}

// Untyped supports binding a interface{} value.
type Untyped interface {
	DataItem
	Get() (interface{})
	Set(interface{})
}

type boundUntyped struct {
	base
	val *interface{}
}

// NewUntyped returns a bindable interface{} value that is managed internally.
func NewUntyped() Untyped {
	var blank interface{} = nil
	b := &boundUntyped{val: &blank}
        b.Init(b)
        return b
}

func (b *boundUntyped) Get() (interface{}) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if b.val == nil {
		return nil
	}
	return *b.val
}

func (b *boundUntyped) Set(val interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if *b.val == val {
		return
	}
	*b.val = val
	b.trigger()
}

// ExternalUntyped supports binding a interface{} value to an external value.
type ExternalUntyped interface {
	Untyped
	Reload()
}

type boundExternalUntyped struct {
	boundUntyped
	old interface{}
}

// BindUntyped returns a new bindable value that controls the contents of the provided interface{} variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func BindUntyped(v *interface{}) ExternalUntyped {
	if v == nil {
		var blank interface{} = nil
		v = &blank // never allow a nil value pointer
	}
	b := &boundExternalUntyped{}
	b.val = v
	b.old = *v
        b.Init(b)
	return b
}

func (b *boundExternalUntyped) Set(val interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.old == val {
		return
	}
	*b.val = val
	b.old = val
	b.trigger()
}

func (b *boundExternalUntyped) Reload() {
	b.Set(*b.val)
}
