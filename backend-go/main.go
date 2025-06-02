package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"graduate_work/algos"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

type HandlerFunc func(*websocket.Conn, any) error

type RequestType any

func handleWebSocket[T RequestType](w http.ResponseWriter, r *http.Request, constructor func(T) (algos.Algorithm, error)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка WebSocket:", err)
		return
	}

	var request T
	if err := conn.ReadJSON(&request); err != nil {
		fmt.Println("Ошибка чтения JSON:", err)
		return
	}

	algorithm, err := constructor(request)
	if err != nil {
		fmt.Println("Ошибка инициализации алгоритма:", err)
		return
	}

	go func() {
		defer func() {
			fmt.Println("Закрытие WebSocket соединения")
			conn.Close()
		}()

		disconnect := make(chan struct{})

		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Клиент отключился:", err)
					close(disconnect)
					return
				}
			}
		}()

		send := func(response algos.Response) error {
			select {
			case <-disconnect:
				return fmt.Errorf("соединение с клиентом закрыто")
			default:
			}

			// fmt.Println(time.Now())
			data, err := json.Marshal(response)
			if err != nil {
				fmt.Println("Ошибка формирования JSON:", err)
				return err
			}

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("Ошибка записи в WebSocket:", err)
				return err
			}

			return nil
		}

		algorithm.Run(send)
	}()

}

func main() {
	http.HandleFunc("/ws/AFSA", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, algos.NewAFSA)
	})
	http.HandleFunc("/ws/SFLA", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, algos.NewSFLA)
	})
	http.HandleFunc("/ws/firefly", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, algos.NewFA)
	})
	http.HandleFunc("/ws/ABC", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, algos.NewABC)
	})
	http.HandleFunc("/ws/GWO", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, algos.NewGWO)
	})

	fmt.Println("Сервер запущен на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
