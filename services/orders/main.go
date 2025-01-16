package main

func main() {
	httpServer := NewHttpServer(":3000")
	go httpServer.Run()

	gRPCServer := NewGRPCServer(":8080")
	gRPCServer.Run()

}
