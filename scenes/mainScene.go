package scenes

import (
	"sync"
	"time"
	"estacionamiento_concurrente/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"fyne.io/fyne/v2/widget"
)

type MainScene struct {
	window fyne.Window
}

func NewMainScene(window fyne.Window) *MainScene {
	return &MainScene{
		window: window,
	}
}

var carsContainer = container.NewWithoutLayout()

func (s *MainScene) Show() {

	rectangle := canvas.NewRectangle(color.Transparent)
	rectangle.StrokeWidth = 2
	rectangle.StrokeColor = color.White
	rectangleW := 690
	rectangleH := 170
	rectangleX := 0
	rectangleY := 10
	rectangle.Resize(fyne.NewSize(float32(rectangleW), float32(rectangleH)))
	rectangle.Move(fyne.NewPos(float32(rectangleX), float32(rectangleY)))

	gate := canvas.NewRectangle(color.White)
	gateW := 150
	gateH := 10
	gateX := (rectangleW / 2) - (gateW / 2)
	gateY := rectangleH + 5
	gate.Resize(fyne.NewSize(float32(gateW), float32(gateH)))
	gate.Move(fyne.NewPos(float32(gateX), float32(gateY)))

	carsContainer.Add(rectangle)
	carsContainer.Add(gate)

	startButton := widget.NewButton("Iniciar Simulación", func() {
		s.Run()
	})

	
	content := container.NewVBox(
		startButton,
		carsContainer,
	)
	
	s.window.SetContent(content)
}

func (s *MainScene) Run() {
	p := models.NewParking(make(chan int, 20), &sync.Mutex{})
	poissonDist := models.NewPoissont()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			car := models.NewCar(id)
			carImage := car.GetCarImage()
			carImage.Resize(fyne.NewSize(30, 30))

			currentX := float32(325)
			currentY := float32(355 +(id%2 *60))

			car.MoveTo(currentX, currentY)

			carsContainer.Add(carImage)
			carsContainer.Refresh()

			car.Park(p, carsContainer, &wg)
		}(i)
		var randomPoissonNumber = poissonDist.CreateDist(float64(2))
		time.Sleep(time.Second * time.Duration(randomPoissonNumber))
	}

	wg.Wait()
	fmt.Println("Simulador terminado")
}
