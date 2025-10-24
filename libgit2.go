package main

/*
#cgo CFLAGS: -I./libgit2/include
#cgo LDFLAGS: -L./libgit2/build -lgit2 -lssl -lcrypto -lz

#include <git2.h>
#include <stdio.h>

extern void list_cb();
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
		return nil, to_error(error)
	}

	return repo, nil
}

// quick function to list all shit in a repository
func list_all(repo *C.git_repository) error {
	var obj *C.git_object
	error := C.git_revparse_single(&obj, repo, C.CString("HEAD^{tree}"))
	if error < 0 {
		return to_error(error)
	}

	tree := (*C.git_tree)(obj)

	error = C.git_tree_walk(tree, C.GIT_TREEWALK_PRE, (*[0]byte)(C.list_cb), nil)  // note, that we're getting our Go function *through* C
	if error < 0 {
		return to_error(error)
	}
	return nil
}

//export list_cb
func list_cb(root *C.char, entry *C.git_tree_entry, payload *C.void) {
	name := C.GoString(C.git_tree_entry_name(entry))
	fmt.Printf("entry: %s\n", name)
}


func to_error(error C.int) error {
	e := C.git_error_last()
	return fmt.Errorf("Error %d/%d: %s\n", error, e.klass, C.GoString(e.message))
}
