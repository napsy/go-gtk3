package gtk3

/*
#include <gtk/gtk.h>

static GtkEntry* to_GtkEntry(void* obj) {
	return GTK_ENTRY(obj);
}

*/
import "C"
import "unsafe"
import "github.com/norisatir/go-gtk3/gobject"

type Entry struct {
	object *C.GtkEntry
	*Container
}

//Create and return new entry Structure
func NewEntry() *Entry {
	e := &Entry{}

	o := C.gtk_entry_new()

	e.Container = NewContainer(unsafe.Pointer(o))

	e.object = C.to_GtkEntry(unsafe.Pointer(o))
	return e
}

// Conversion function for gobject registration map
func newEntryFromNative(obj unsafe.Pointer) interface{} {
	var entry Entry
	entry.object = C.to_GtkEntry(obj)
	entry.Container = NewContainer(unsafe.Pointer(obj))
	return &entry
}

func nativeFromEntry(e interface{}) *gobject.GValue {
	if entry, ok := e.(Entry); ok {
		gv := gobject.CreateCGValue(GtkType.ENTRY, entry.ToNative())
		return gv
	}
	return nil
}

func init() {
	// Register GtkEntry to gobject type system (in Go)
	gobject.RegisterCType(GtkType.ENTRY, newEntryFromNative)
	gobject.RegisterGoType(GtkType.ENTRY, nativeFromEntry)
}

// To be object-like
func (self Entry) ToNative() unsafe.Pointer {
	return unsafe.Pointer(self.object)
}

func (self Entry) Connect(s string, f interface{}, datas ...interface{}) {
	gobject.Connect(self, s, f, datas...)
}

func (self Entry) Set(properties map[string]interface{}) {
	gobject.Set(self, properties)
}

func (self Entry) Get(properties []string) map[string]interface{} {
	return gobject.Get(self, properties)
}

// Button interface
func (self *Entry) SetText(label string) {
	s := gobject.GString(label)
	//defer s.Free()
	C.gtk_entry_set_text(self.object, (*C.gchar)(s.GetPtr()))
}

func (self *Entry) GetText() string {
	l := C.gtk_entry_get_text(self.object)
	return gobject.GoString(unsafe.Pointer(l))
}

//TODO: gtk_button_get_event_window
