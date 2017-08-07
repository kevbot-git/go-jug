package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"github.com/googollee/go-socket.io"
)

const (
	addr = ":777"
)

func main() {
	
	go Serve()

	firmataAdaptor := firmata.NewTCPAdaptor("192.168.0.20:3030")
	servo := gpio.NewServoDriver(firmataAdaptor, "0")

	work := func() {
		gobot.Every(1*time.Second, func() {
			i := uint8(gobot.Rand(180))
			fmt.Println("Turning", i)
			servo.Move(i)
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}

func Serve() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	fs := http.FileServer(http.Dir("www"))
	http.Handle("/", fs)
	http.Handle("/socket.io/", server)

	log.Println("Listening...")
	http.ListenAndServe(addr, nil)
}