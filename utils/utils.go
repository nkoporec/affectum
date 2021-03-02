package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func GetDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't retrieve current user, error was: %s", err))
	}

	return filepath.Join(usr.HomeDir, "/affectum")
}

func CreateDir() {
	// Set up dir.
	affectumDir := GetDir()
	if _, err := os.Stat(affectumDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(affectumDir, 0755)
			os.Mkdir(filepath.Join(affectumDir, "files"), 0755)
			os.Create(filepath.Join(affectumDir, "affectum.log"))
		}
	}
}

func Logger(message string) {
	f, err := os.OpenFile(filepath.Join(GetDir(), "affectum.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(message)
}
