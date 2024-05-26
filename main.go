package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func startProcess(command string, args []string) (*os.Process, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	fmt.Println("Process started with PID:", cmd.Process.Pid)
	return cmd.Process, nil
}
func listProcesses() {
	// Trên Unix, bạn có thể sử dụng lệnh ps
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println(string(output))
}
func main() {
	var processes []*os.Process
	for {
		reader := bufio.NewReader(os.Stdin)
		var choice int
		fmt.Println("1. Start process")
		fmt.Println("2. Stop process")
		fmt.Println("3. List processes")
		fmt.Println("4. Exit")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			fmt.Println("Enter command to start process")
			command, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading command")
				continue
			}
			parts := strings.Fields(command)
			command = parts[0]
			args := parts[1:]
			process, err := startProcess(command, args)
			if err != nil {
				fmt.Println("Error starting process")
			} else {
				processes = append(processes, process)
			}
		case 2:
			if len(processes) == 0 {
				fmt.Println("No processes to stop")
				break
			}
			for i, proc := range processes {
				fmt.Printf("%d. PID: %d\n", i, proc.Pid)
			}
			fmt.Println("Enter index of process to stop")
			var index int
			fmt.Scan(&index)
			if index < 0 || index >= len(processes) {
				fmt.Println("Invalid index")
			} else {
				err := processes[index].Signal(syscall.SIGKILL)
				if err != nil {
					fmt.Println("Error stopping process")
				} else {
					fmt.Println("Process with PID", processes[index].Pid, "stopped")
					processes = append(processes[:index], processes[index+1:]...)
				}
			}
		case 3:
			listProcesses()
		case 4:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}
