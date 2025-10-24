package main

import "fmt"
import "os"

func main() {
	fmt.Println("asd asd cock");
	libgit_init()

	repo, err := repository_open(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)  // maybe copy the exit status of the libgit error?
	}

	err = list_all(repo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
