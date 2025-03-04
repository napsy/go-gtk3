// gtype.go
// This file contains gobject type definition
package gobject
/*
#include <glib-object.h>

static inline const gchar* get_type_name(GValue v) { return G_VALUE_TYPE_NAME(&v); }

static inline const gchar* getTypeName(void * o) { return G_OBJECT_TYPE_NAME(o); }
static inline GType getTypeId(void * o) { return G_OBJECT_TYPE(o); }

static inline GObject* to_GObject(void* o) { return G_OBJECT(o); }

static inline GType getGValueType(GValue* v) {
	return G_VALUE_TYPE(v);
}

static inline const gchar* getGValueTypeName(GValue* v) {
	return G_VALUE_TYPE_NAME(v);
}

static inline gboolean pointer_in_gvalue(GValue* v) {
	return G_VALUE_HOLDS_POINTER(v);
}

static inline gboolean object_in_gvalue(GValue* v) {
	return G_VALUE_HOLDS_OBJECT(v);
}

*/
import "C"
import "unsafe"

//-----------------------------------------------------------------------
// GType
//-----------------------------------------------------------------------
type GType int

type GValueError struct {
	Err string
}

func (e GValueError) Error() string {
	return e.Err
}

// Function interface for converting C values to Go
type FuncToGo func(unsafe.Pointer) interface{}

// GoTypes registration map
// This map holds gobject types as GType
// and conversion function to Go type
// Every module has to register it's own type and conversion function
var gtypes map[GType]FuncToGo

// Get Type ID
func GetTypeID(obj unsafe.Pointer) GType {
	t := C.getTypeId(obj)
	return GType(t)
}

// Get Type Name
func GetTypeName(obj unsafe.Pointer) string {
	tn := C.getTypeName(obj)
	return GoString(unsafe.Pointer(tn))
}

func RegisterCType(typename GType, f FuncToGo) {
	gtypes[typename] = f
}

func ConvertToGo(obj unsafe.Pointer, typeid ...GType) (interface{}, error) {
	if obj == nil {
		return nil, GValueError{"Conversion Error"}
	}
	var id GType

	if len(typeid) > 0 {
		id = typeid[0]
	} else {
		id = GetTypeID(obj)
	}
	if f, ok := gtypes[id]; ok {
		res := f(obj)
		return res, nil
	}

	return nil, GValueError{"Unknown Type"}
}

// Function interface for converting Go values to C
type FuncToC func(interface{}) *GValue

var ctypes map[GType]FuncToC

func RegisterGoType(typeid GType, f FuncToC) {
	ctypes[typeid] = f
}

// GValue
type GValue struct {
	gtype GType
	value *C.GValue
}

func (self GValue) ToCGValue() unsafe.Pointer {
	return unsafe.Pointer(self.value)
}

func (self *GValue) ReInitializeType() {
	t := C.getGValueType(self.value)
	self.gtype = GType(t)
}

func (self GValue) GetTypeName() string {
	tn := C.getGValueTypeName(self.value)
	return GoString(unsafe.Pointer(tn))
}

func (self GValue) GetTypeID() GType {
	return self.gtype
}

func (self GValue) GetPtr() unsafe.Pointer {
	switch self.gtype {
	case G_TYPE_STRING:
		return unsafe.Pointer(C.g_value_get_string(self.value))
	case G_TYPE_BOOLEAN:
		b := C.g_value_get_boolean(self.value)
		return unsafe.Pointer(&b)
	case G_TYPE_CHAR:
		c := C.g_value_get_char(self.value)
		return unsafe.Pointer(&c)
	case G_TYPE_INT:
		i := C.g_value_get_int(self.value)
		return unsafe.Pointer(&i)
	case G_TYPE_LONG:
		l := C.g_value_get_long(self.value)
		return unsafe.Pointer(&l)
	case G_TYPE_INT64:
		i := C.g_value_get_int64(self.value)
		return unsafe.Pointer(&i)
	case G_TYPE_UCHAR:
		c := C.g_value_get_uchar(self.value)
		return unsafe.Pointer(&c)
	case G_TYPE_UINT:
		i := C.g_value_get_uint(self.value)
		return unsafe.Pointer(&i)
	case G_TYPE_ULONG:
		l := C.g_value_get_ulong(self.value)
		return unsafe.Pointer(&l)
	case G_TYPE_UINT64:
		i := C.g_value_get_uint64(self.value)
		return unsafe.Pointer(&i)
	case G_TYPE_FLOAT:
		f := C.g_value_get_float(self.value)
		return unsafe.Pointer(&f)
	case G_TYPE_DOUBLE:
		d := C.g_value_get_double(self.value)
		return unsafe.Pointer(&d)
	}

	// Now things get tricky
	// We know the type, but how is it stored?

	// Is it pointer?
	isIt := C.pointer_in_gvalue(self.value)
	if GoBool(unsafe.Pointer(&isIt)) {
		val := C.g_value_get_pointer(self.value)
		return unsafe.Pointer(val)
	}

	// Is it object?
	isIt = C.object_in_gvalue(self.value)
	if GoBool(unsafe.Pointer(&isIt)) {
		o := C.g_value_get_object(self.value)
		return unsafe.Pointer(o)
	}
	//Hmmm....
	ptr := C.g_value_peek_pointer(self.value)
	if ptr != nil {
		return unsafe.Pointer(ptr)
	}
	return nil
}

func (self GValue) Free() {
	C.g_value_unset(self.value)
}

func CreateCGValue(tn GType, object ...unsafe.Pointer) *GValue {
	var cv C.GValue
	C.g_value_init(&cv, C.GType(tn))

	// If no data, then return Gvalue initialized with default
	if len(object) == 0 {
		gv := GValue{tn, &cv}
		return &gv
	}
	obj := object[0]

	// Foundamental types are special
	// TODO: Handle more cases, like creating GValue from GdkEvents
	switch tn {
	case G_TYPE_STRING:
		C.g_value_take_string(&cv, (*C.gchar)(obj))
	case G_TYPE_BOOLEAN:
		C.g_value_set_boolean(&cv, *((*C.gboolean)(obj)))
	case G_TYPE_CHAR:
		C.g_value_set_char(&cv, *((*C.gchar)(obj)))
	case G_TYPE_INT:
		C.g_value_set_int(&cv, *((*C.gint)(obj)))
	case G_TYPE_LONG:
		C.g_value_set_long(&cv, *((*C.glong)(obj)))
	case G_TYPE_INT64:
		C.g_value_set_int64(&cv, *((*C.gint64)(obj)))
	case G_TYPE_UCHAR:
		C.g_value_set_uchar(&cv, *((*C.guchar)(obj)))
	case G_TYPE_UINT:
		C.g_value_set_uint(&cv, *((*C.guint)(obj)))
	case G_TYPE_ULONG:
		C.g_value_set_ulong(&cv, *((*C.gulong)(obj)))
	case G_TYPE_UINT64:
		C.g_value_set_uint64(&cv, *((*C.guint64)(obj)))
	case G_TYPE_FLOAT:
		C.g_value_set_float(&cv, *((*C.gfloat)(obj)))
	case G_TYPE_DOUBLE:
		C.g_value_set_double(&cv, *((*C.gdouble)(obj)))
	default:
		C.g_value_set_object(&cv, C.gpointer(obj))
	}

	gv := GValue{tn, &cv}
	return &gv
}

func ConvertToC(gotype interface{}) *GValue {
	// Test first for fundamental types
	switch gotype.(type) {
	case string:
		return ctypes[G_TYPE_STRING](gotype.(string))
	case bool:
		return ctypes[G_TYPE_BOOLEAN](gotype.(bool))
	case int8:
		return ctypes[G_TYPE_CHAR](gotype.(byte))
	case int:
		return ctypes[G_TYPE_INT](gotype.(int))
	case int32:
		return ctypes[G_TYPE_LONG](gotype.(int32))
	case int64:
		return ctypes[G_TYPE_INT64](gotype.(int64))
	case byte:
		return ctypes[G_TYPE_UCHAR](gotype.(uint8))
	case uint:
		return ctypes[G_TYPE_UINT](gotype.(uint))
	case uint32:
		return ctypes[G_TYPE_ULONG](gotype.(uint32))
	case uint64:
		return ctypes[G_TYPE_UINT64](gotype.(uint64))
	case float32:
		return ctypes[G_TYPE_FLOAT](gotype.(float32))
	case float64:
		return ctypes[G_TYPE_DOUBLE](gotype.(float64))
	case ObjectLike:
		o := gotype.(ObjectLike)
		t := GetTypeFromInstance(o.ToNative())
		if f, ok := ctypes[t]; ok {
			return f(o)
		}
	}
	return nil
}

func init() {
	C.g_type_init()
	gtypes = make(map[GType]FuncToGo)
	ctypes = make(map[GType]FuncToC)
}
