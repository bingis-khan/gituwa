package main

/*
#cgo CFLAGS: -I./libgit2/include
#cgo LDFLAGS: -L./libgit2/build -lgit2 -lssl -lcrypto -lz

#include <git2.h>
*/
import "C"
import "unsafe"
import "fmt"

func libgit_init() {
	C.git_libgit2_init()
}

func repository_open(path string) (*C.git_repository, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))  // not sure if this resource is used after repo shit. maybe I should make a datatype which stores all the resources?

	var repo *C.git_repository
	error := C.git_repository_open(&repo, cpath)
	if error < 0 {
		e := C.git_error_last()
		return nil, fmt.Errorf("Error %d/%d: %s\n", error, e.klass, C.GoString(e.message))
	}

	return repo, nil
}
