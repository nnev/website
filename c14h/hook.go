package main

import (
	"log"
	"os/exec"
)

func RunHook() error {
	cmd := exec.Command("/bin/sh", "-c", *hook)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Could not run hook:", err)
		log.Println("=== BEGIN OUTPUT ===")
		log.Print(string(output))
		log.Println("=== END OUTPUT ===")
		return err
	}

	return nil
}
