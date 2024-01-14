package main

import (
  "fmt"
)

type Token struct {
  Data      string
  Recipient int
  TTL       int
}

type Node struct {
  ID       int
  Previous <-chan Token
  Next     chan<- Token
}

func (node *Node) processToken() {
  for {
    // С левого узла получаем токен
    token := <-node.Previous

    // Если сообщение этому узлу, то выходим из цикла
    if token.Recipient == node.ID {
      fmt.Printf("Узел с номером %d принял токен: Data:%s, TTL:%d\n", node.ID, token.Data, token.TTL)
      return
    }

    // Если сообщение не этому узлу, передаем токен дальше
    if token.TTL > 0 {
      token.TTL--
      node.Next <- token
    } else {
      fmt.Printf("Время жизни токена истекло на узле с номером %d", node.ID)
      return
    }
  }
}

func InputParameters(N *int, message *string, recipient *int, lifetime *int, startIndex *int) {
  fmt.Println("Число узлов:")
  fmt.Scanln(N)

  fmt.Println("Сообщение:")
  fmt.Scanln(message)

  fmt.Println("Узел-получатель:")
  fmt.Scanln(recipient)

  fmt.Println("TTL:")
  fmt.Scanln(lifetime)

  fmt.Println("Начальный узел:")
  fmt.Scanln(startIndex)
}

func main() {
  var N, recipient, lifetime, startIndex int
  var message string

  InputParameters(&N, &message, &recipient, &lifetime, &startIndex)

  // создание очередей
  channels := make([]chan Token, N)
  for i := range channels {
    channels[i] = make(chan Token)
  }

  // создание узлов
  nodes := make([]Node, N)
  for i := 0; i < N-1; i++ {
    nodes[i] = Node{ID: i, Previous: channels[i], Next: channels[i+1]}
    go nodes[i].processToken()
  }
  // связывание последнего узла с начальным
  nodes[N-1] = Node{ID: N - 1, Previous: channels[N-1], Next: channels[0]}
  go nodes[N-1].processToken()

  initialToken := Token{Data: message, Recipient: recipient, TTL: lifetime}
  channels[startIndex] <- initialToken

  // бесконечный цикл, чтобы программа не завершилась раньше времени
  for {
  }
}
