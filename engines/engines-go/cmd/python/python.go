package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {

	log.Println("executing command")
	cmd := exec.Command("/usr/bin/python3", "-c", "import time; print(time.time());")
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	log.Println("running cmd.Output()")
	out, err := cmd.Output()

	if err != nil {
		log.Fatalf("error while Output command: %v", err)
		return
	}

	log.Println("getting command output")
	log.Printf("output: '%s'", string(out))

	log.Println("returning")

}
