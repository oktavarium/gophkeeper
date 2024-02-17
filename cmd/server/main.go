package main

func main() {
	if err := server.Run(); err != nil {
		panic(fmt.Errorf(err, "error running server"))
	}
}
