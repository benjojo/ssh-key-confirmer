package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strings"

	"golang.org/x/crypto/ssh"
)

var (
	quiet = flag.Bool("q", false, "quiet mode")
)

func main() {
	testKey := flag.String("i", "./id_ed25519.pub", "the key you want to confirm")
	portFlag := flag.Int("p", 22, "port of the SSH server")
	flag.Parse()

	target := flag.Arg(0)
	if target == "" {
		fmt.Printf("Usage:\nssh-key-confirmer -q (quiet) -i [path-to-public-key] -p [Port] root@192.0.0.1\n")
		os.Exit(8)
	}
	username, target := parseSSHTarget(target)

	testKeyBytes, err := ioutil.ReadFile(*testKey)
	if err != nil {
		fmt.Printf("failed to read public key, %v\n", err)
		os.Exit(5)
	}

	testPubKey, _, _, _, err = ssh.ParseAuthorizedKey(testKeyBytes)
	if err != nil {
		fmt.Printf("failed to parse public key, %v\n", err)
		os.Exit(5)
	}

	fS := fakeSigner{}
	bogoKey := makeBogoKey()
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(bogoKey, fS),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	_, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", target, *portFlag), config)
	if err != nil {
		if !strings.Contains(err.Error(), "no supported methods remain") {
			printAndExit(fmt.Sprintf("SSH conection error: %v\n", err), 1)
		}
		printAndExit(fmt.Sprintf("Key not found on user+server\n"), 2)
	}
}

func printAndExit(in string, code int) {
	if !*quiet {
		fmt.Print(in)
	}
	os.Exit(code)
}

var testPubKey ssh.PublicKey

type fakeSigner struct{}

func (f fakeSigner) PublicKey() ssh.PublicKey {
	return testPubKey
}

func (f fakeSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	printAndExit("Key is present on user+server\n", 0)
	return nil, nil
}

func parseSSHTarget(in string) (username string, target string) {
	bits := strings.Split(in, "@")
	if len(bits) != 2 {
		bits[0] = squarev6(bits[0])
		// get the user name
		userObj, err := user.Current()
		if err != nil {
			return "root", bits[0]
		}
		return userObj.Name, bits[0]
	}

	return bits[0], squarev6(bits[1])
}

func squarev6(in string) string {
	possibleAddr := net.ParseIP(in)
	if possibleAddr == nil {
		return in
	}

	// is it v6?
	if possibleAddr.To4() == nil {
		return fmt.Sprintf("[%v]", possibleAddr.String())
	}

	return in
}
