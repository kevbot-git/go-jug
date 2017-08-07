package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	addr = ":777"
)

func main() {

	firmataAdaptor := firmata.NewTCPAdaptor("192.168.0.20:3030")
	servo := gpio.NewServoDriver(firmataAdaptor, "0")
	led := gpio.NewLedDriver(firmataAdaptor, "1")

	work := func() {
		boiling := false
		server, err := socketio.NewServer(nil)
		if err != nil {
			log.Fatal(err)
		}

		server.On("connection", func(so socketio.Socket) {
			led.Off()
			so.Join("jug")

			if boiling {
				so.Emit("boiling")
			}

			so.On("disconnection", func() {
				fmt.Println("Seeya!")
				led.On()
			})

			so.On("boil", func() {
				if !boiling {
					boiling = true
					server.BroadcastTo("jug", "boiling", nil)
					fmt.Println("Turning on...")
					servo.Move(uint8(30))
					gobot.After(time.Second*5, func() {
						fmt.Println("Turning off...")
						servo.Move(uint8(120))
						server.BroadcastTo("jug", "boiled", nil)
						server.BroadcastTo("jug", "message", "JUG: I'm ready!")
						boiling = false
					})
				}
			})
		})
		server.On("error", func(so socketio.Socket, err error) {
			log.Println("error:", err)
		})

		fs := http.FileServer(http.Dir("www"))
		http.Handle("/", fs)
		http.Handle("/socket.io/", server)

		log.Println("Listening...")
		log.Fatal(http.ListenAndServe(addr, nil))
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)

	robot.Start()
}
