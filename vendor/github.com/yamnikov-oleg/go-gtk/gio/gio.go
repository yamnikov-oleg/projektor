// +build !cgocheck

package gio

// #include "gio.go.h"
// #cgo pkg-config: gio-2.0
import "C"
import (
	"unsafe"

	"github.com/yamnikov-oleg/go-gtk/glib"
)

//-----------------------------------------------------------------------
// GIcon
//-----------------------------------------------------------------------

type Icon struct {
	GIcon *C.GIcon
}

func NewIconForString(str string) (*Icon, *glib.Error) {
	var err *C.GError
	ptr := C.CString(str)
	defer C.free(unsafe.Pointer(ptr))
	gicon := C.g_icon_new_for_string((*C.gchar)(ptr), &err)
	if err != nil {
		return nil, glib.ErrorFromNative(unsafe.Pointer(err))
	}
	return &Icon{gicon}, nil
}

func GetIconType() int {
	return int(C.g_icon_get_type())
}

func (icon *Icon) ToString() string {
	ptr := C.g_icon_to_string(icon.GIcon)
	defer C.free(unsafe.Pointer(ptr))
	return C.GoString((*C.char)(ptr))
}

//-----------------------------------------------------------------------
// GFile
//-----------------------------------------------------------------------

type File struct {
	GFile *C.GFile
}

func NewFileForPath(path string) *File {
	ptr := C.CString(path)
	defer C.free(unsafe.Pointer(ptr))
	return &File{C.g_file_new_for_path(ptr)}
}

func (f *File) QueryInfo(attributes string, flags FileQueryInfoFlags, cancellable unsafe.Pointer) (*FileInfo, *glib.Error) {
	var err *C.GError

	ptr := C.CString(attributes)
	defer C.free(unsafe.Pointer(ptr))

	fi := C.g_file_query_info(f.GFile, ptr, C.GFileQueryInfoFlags(flags), (*C.GCancellable)(cancellable), &err)
	if err != nil {
		return nil, glib.ErrorFromNative(unsafe.Pointer(err))
	}
	return &FileInfo{fi}, nil
}

type FileQueryInfoFlags int

const (
	FILE_QUERY_INFO_NONE              FileQueryInfoFlags = C.G_FILE_QUERY_INFO_NONE
	FILE_QUERY_INFO_NOFOLLOW_SYMLINKS                    = C.G_FILE_QUERY_INFO_NOFOLLOW_SYMLINKS
)

//-----------------------------------------------------------------------
// GFileInfo
//-----------------------------------------------------------------------

type FileInfo struct {
	GFileInfo *C.GFileInfo
}

func (fi *FileInfo) GetIcon() *Icon {
	return &Icon{C.g_file_info_get_icon(fi.GFileInfo)}
}
