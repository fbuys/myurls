package acceptance

// I relied heavily on this repo to create these helpers
// https://github.com/quii/go-graceful-shutdown/blob/main/acceptancetests/blackboxtestthings.go

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	baseBinName = "temp-testbinary"
	port        = "8080"
	Url         = "http://localhost:" + port
)

func LaunchTestProgram(path string) (cleanup func(), err error) {
	binName, err := buildBinary(path)
	if err != nil {
		return nil, err
	}

	kill, err := runServer(binName)

	cleanup = func() {
		if kill != nil {
			kill()
		}
		os.Remove(binName)
	}

	if err != nil {
		cleanup() // even though it's not listening correctly, the program could still be running
		return nil, err
	}

	return cleanup, nil
}

func buildBinary(path string) (string, error) {
	binName := randomString(10) + "-" + baseBinName

	build := exec.Command("go", "build", "-o", binName, path)

	if err := build.Run(); err != nil {
		return "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}
	return binName, nil
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func runServer(binName string) (kill func(), err error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	cmd := exec.Command(cmdPath)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cannot run temp converter: %s", err)
	}

	kill = func() {
		_ = cmd.Process.Kill()
	}

	err = WaitForServerListening()

	return kill, err
}

func WaitForServerListening() error {
	for i := 0; i < 30; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("nothing seems to be listening on localhost:%s", port)
}
