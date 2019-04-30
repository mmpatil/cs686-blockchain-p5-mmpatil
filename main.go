package main

import (
	"log"
	"net/http"
	"os"

	"./p3"
)

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
// 	})

// 	log.Fatal(http.ListenAndServe(":8080", nil))

// }

func main() {
	router := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}
}

// func main() {
// 	peers := data.NewPeerList(5, 4)
// 	peers.Add("1111", 1)
// 	peers.Add("4444", 4)
// 	peers.Add("-1-1", -1)
// 	peers.Add("0000", 0)
// 	peers.Add("2121", 21)
// 	peers.Rebalance()
// 	expected := data.NewPeerList(5, 4)
// 	expected.Add("1111", 1)
// 	expected.Add("4444", 4)
// 	expected.Add("2121", 21)
// 	expected.Add("-1-1", -1)
// 	fmt.Println(reflect.DeepEqual(peers, expected))

// 	peers = data.NewPeerList(5, 2)
// 	peers.Add("1111", 1)
// 	peers.Add("4444", 4)
// 	peers.Add("-1-1", -1)
// 	peers.Add("0000", 0)
// 	peers.Add("2121", 21)
// 	peers.Rebalance()
// 	expected = data.NewPeerList(5, 2)
// 	expected.Add("4444", 4)
// 	expected.Add("2121", 21)
// 	fmt.Println(reflect.DeepEqual(peers, expected))

// 	peers = data.NewPeerList(5, 4)
// 	peers.Add("1111", 1)
// 	peers.Add("7777", 7)
// 	peers.Add("9999", 9)
// 	peers.Add("11111111", 11)
// 	peers.Add("2020", 20)
// 	peers.Rebalance()
// 	expected = data.NewPeerList(5, 4)
// 	expected.Add("1111", 1)
// 	expected.Add("7777", 7)
// 	expected.Add("9999", 9)
// 	expected.Add("2020", 20)
// 	fmt.Println(reflect.DeepEqual(peers, expected))
// }
