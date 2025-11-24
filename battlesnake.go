package main

type Config struct {
	Port       string
	MaxWorkers int8
}

func main() {
	c := Config{
		Port: "8000",
	}
	Run(&c)
}

func start(state GameState) {

}

func move() {

}
