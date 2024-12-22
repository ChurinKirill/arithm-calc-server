package main

func main() {
	err := InitializeLogger("log.txt")
	if err != nil {
		return
	}
	StartServer()
	ShutdownLogger()
}
