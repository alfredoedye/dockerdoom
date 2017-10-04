package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func runCmd(cmdstring string) {
	log.Printf("Executing the following command:  \"%v\"\n", cmdstring)
	parts := strings.Split(cmdstring, " ")
	cmd := exec.Command(parts[0], parts[1:len(parts)]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("The following command failed: \"%v\"\n", cmdstring)
	}
}

func outputCmd(cmdstring string) string {
	log.Printf("cmdString: \"%v\"\n", cmdstring)
	parts := strings.Split(cmdstring, " ")
	cmd := exec.Command(parts[0], parts[1:len(parts)]...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("The following command failed: \"%v\"\n", cmdstring)
	}
	log.Printf("Output: \"%v\"\n", output)
	return string(output)
}

func startCmd(cmdstring string) {
	log.Printf("STARTcmd:  \"%v\"\n", cmdstring)

	parts := strings.Split(cmdstring, " ")
	cmd := exec.Command(parts[0], parts[1:len(parts)]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		log.Fatalf("The following command failed: \"%v\"\n", cmdstring)
	}
}

func checkDockerImages(imageName, dockerBinary string) bool {
	log.Print("Checking Docker Images")
	output := outputCmd(fmt.Sprintf("%v images -q %v", dockerBinary, imageName))
	return len(output) > 0
}

func checkActiveDocker(dockerName, dockerBinary string) bool {
	log.Print("Checking Active Docker")
	return checkDocker(dockerName, dockerBinary, "-q")
}

func checkAllDocker(dockerName, dockerBinary string) bool {
	return checkDocker(dockerName, dockerBinary, "-aq")
}

func checkDocker(dockerName, dockerBinary, arg string) bool {
	output := outputCmd(fmt.Sprintf("%v ps %v", dockerBinary, arg))
	docker_ids := strings.Split(string(output), "\n")
	for _, docker_id := range docker_ids {
		if len(docker_id) == 0 {
			continue
		}
		output := outputCmd(fmt.Sprintf("%v inspect -f {{.Name}} %v", dockerBinary, docker_id))
		name := strings.TrimSpace(string(output))
		name = name[1:len(name)]
		if name == dockerName {
			return true
		}
	}
	return false
}

func socketLoop(listener net.Listener, dockerBinary, containerName string) {

	log.Print("socketLoop START")
	log.Printf("Listener = %v ", listener)
	log.Printf("dockerBinary = %v ", dockerBinary)
	log.Printf("containerName = %v", containerName)

	for true {
		log.Print("socketLoop For loop \n")
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		stop := false
		for !stop {
			log.Print("socketLoop stop\n")
			bytes := make([]byte, 40960)
			n, err := conn.Read(bytes)
			if err != nil {
				stop = true
			}
			bytes = bytes[0:n]
			strbytes := strings.TrimSpace(string(bytes))
			log.Printf(" strbytes: %v ", strbytes)

			if strbytes == "list" {
				log.Printf("\n List: \n")
				output := outputCmd(fmt.Sprintf("%v ps -q", dockerBinary))

				//cmd := exec.Command("/usr/local/bin/docker", "inspect", "-f", "{{.Name}}", "`docker", "ps", "-q`")
				outputstr := strings.TrimSpace(output)
				outputparts := strings.Split(outputstr, "\n")
				for _, part := range outputparts {
					output := outputCmd(fmt.Sprintf("%v inspect -f {{.Name}} %v", dockerBinary, part))
					name := strings.TrimSpace(output)
					name = name[1:len(name)]
					if name != containerName {
						_, err = conn.Write([]byte(name + "\n"))
						if err != nil {
							log.Fatal("Could not write to socker file")
						}
					}
				}
				conn.Close()
				stop = true
			} else if strings.HasPrefix(strbytes, "kill ") {
				log.Printf("socketloop KILL")
				parts := strings.Split(strbytes, " ")
				docker_id := strings.TrimSpace(parts[1])
				cmd := exec.Command(dockerBinary, "rm", "-f", docker_id)
				go cmd.Run()
				conn.Close()
				stop = true
			}
		}
	}
}

func main() {
	var socketFileFormat, containerName, vncPort, dockerBinary string
	var dockerWait int
	flag.StringVar(&socketFileFormat, "socketFileFormat", "/dockerdoom.socket", "Location and format of the socket file")
	flag.StringVar(&dockerBinary, "dockerBinary", "docker", "docker binary")
	flag.StringVar(&containerName, "containerName", "dockerdoom", "Name of the docker container running DOOM")
	flag.StringVar(&vncPort, "vncPort", "5900", "Port to open for VNC Viewer")
	flag.IntVar(&dockerWait, "dockerWait", 5, "Time to wait before checking if the container came up")

	flag.Parse()

	// Creacion del socket

	socketFile := fmt.Sprintf("dockerdoom.socket")

	listener, err := net.Listen("unix", socketFile)

	log.Print("socketFile: ", socketFile)

	if err != nil {
		log.Fatalf("Could not create socket file %v.\nYou could use \"rm -f %v\"", socketFile, socketFile)
	}

	socketLoop(listener, dockerBinary, containerName)
}
