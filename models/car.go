package models

import (
	"sync"
	"time"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"math/rand"
)

type Car struct {
	id          int
	parkingTime time.Duration
	image       *canvas.Image
	space       int
	exitImage   *canvas.Image
	x           float32 
	y           float32 
}

func NewCar(id int) *Car {
	randomNumber := rand.Intn(4) + 1

	imagePath := fmt.Sprintf("./assets/car%d.png", randomNumber)
	imageExitPath := fmt.Sprintf("./assets/car%d_exit.png", randomNumber)

	image := canvas.NewImageFromURI(storage.NewFileURI(imagePath))
	exitImage := canvas.NewImageFromURI(storage.NewFileURI(imageExitPath))
	return &Car{
		id:          id,
		parkingTime: time.Duration(rand.Intn(10)+10) * time.Second,
		image:       image,
		space:       0,
		exitImage:   exitImage,
		x:           325,
		y: 			 355,
	}
}

func (c *Car) Enter(p *Parking, carsContainer *fyne.Container) {
	p.GetSpaces() <- c.GetId()
	p.GetEntrance().Lock()

	spacesArray := p.GetSpacesArray()

	fmt.Printf("Auto %d ha entrado. Espacios ocupados: %d\n", c.GetId(), len(p.GetSpaces()))

	for i := 0; i < 6; i++ {
		c.image.Move(fyne.NewPos(c.image.Position().X, c.image.Position().Y-30))
		time.Sleep(time.Millisecond * 200)
	}

	p.GetEntrance().Unlock()

	for i := 0; i < len(spacesArray); i++ {
		if !spacesArray[i] {
			spacesArray[i] = true
			c.space = i
			c.y = 10
			c.x = 25 + float32(5+(i*30))
			c.image.Move(fyne.NewPos(c.x, c.y))
			break
		}
	}

	p.SetSpacesArray(spacesArray)
	carsContainer.Refresh()
}

func (c *Car) Leave(p *Parking, carsContainer *fyne.Container) {
	p.GetEntrance().Lock()
	<-p.GetSpaces()

	spacesArray := p.GetSpacesArray()
	spacesArray[c.space] = false
	p.SetSpacesArray(spacesArray)

	fmt.Printf("Auto %d ha salido. Espacios ocupados: %d\n", c.GetId(), len(p.GetSpaces()))
	p.GetEntrance().Unlock()

	for i := 0; i < 10; i++ {
		c.exitImage.Move(fyne.NewPos(c.exitImage.Position().X, c.exitImage.Position().Y+20))
		time.Sleep(time.Millisecond * 200)
	}

	carsContainer.Remove(c.exitImage)
	carsContainer.Refresh()
}

func (c *Car) Park(p *Parking, carsContainer *fyne.Container, wg *sync.WaitGroup) {
	for i := 0; i < 7; i++ {
		c.image.Move(fyne.NewPos(c.image.Position().X + 5, c.image.Position().Y + 3))
		time.Sleep(time.Millisecond * 200)
	}

	c.Enter(p, carsContainer)

	time.Sleep(c.parkingTime)

	carsContainer.Remove(c.image)
	c.exitImage.Resize(fyne.NewSize(30, 30))
	p.ExitQueue(carsContainer, c.exitImage, c.GetId())
	c.Leave(p, carsContainer)
	wg.Done()
}

func (c *Car) GetId() int {
	return c.id
}

func (c *Car) GetCarImage() *canvas.Image {
	return c.image
}

func (c *Car) MoveTo(x, y float32) {
	c.x = x
	c.y = y
	c.image.Move(fyne.NewPos(x, y))
}