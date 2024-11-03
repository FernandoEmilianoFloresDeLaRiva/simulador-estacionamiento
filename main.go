package main

import "estacionamiento_concurrente/views"

func main() {
	mainView := views.NewMainView()
	mainView.Run()
}
