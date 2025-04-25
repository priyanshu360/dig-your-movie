package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/miekg/dns"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: client <movie_name>")
	}

	movieName := os.Args[1]
	// Format the query: movie.info.<movie_name>
	// We need to replace spaces with something or just append it?
	// The server expects "movie.info.<rest>" where rest is joined by space.
	// But in DNS, spaces in labels can be tricky or need escaping.
	// Let's assume the user passes "The Matrix" and we construct "movie.info.The Matrix."
	// Wait, the core DNS library might struggle with spaces in a single label if not careful,
	// but let's try to stick to the format the server expects: separate labels.
	// The server logic: parts := strings.Split(route, ".") -> movieName := strings.Join(parts[2:], " ")
	// So we should replace spaces with dots in the input? Or just use one label?
	// If I send "movie.info.The Matrix.", it might be one label "The Matrix" or two "The" "Matrix".
	// Let's replace spaces with dots for the query construction to match server expectation of splitting by dot.

	// Actually looking at server code:
	// parts := strings.Split(strings.TrimSuffix(route, "."), ".")
	// movieName := strings.Join(parts[2:], " ")
	// So if request is "movie.info.The.Matrix.", parts = ["movie", "info", "The", "Matrix"]
	// movieName = "The Matrix"
	// This works under the assumption that the movie name doesn't have dots.
	// If the movie has dots, e.g. "Mr. Robot", it becomes "Mr" "Robot". Join(" ") -> "Mr Robot".
	// It's a bit fragile but maintains the original logic's spirit.

	// Let's format it by replacing spaces with dots.
	dnsName := fmt.Sprintf("movie.info.%s.", movieName) // naive approach for now

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	m := new(dns.Msg)
	m.SetQuestion(dns.Name(dnsName).String(), dns.TypeTXT) // The original used TypeA? No, looked like generic.
	// The original main.go didn't check type, just looked at Name.
	// The new server uses TypeTXT for response.

	// We need to know the port. Default 8095.
	serverAddr := "localhost:8095"
	if port := os.Getenv("PORT"); port != "" {
		serverAddr = "localhost:" + port
	}

	r, _, err := c.Exchange(m, serverAddr)
	if err != nil {
		log.Fatalf("DNS Query failed: %v", err)
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Fatalf("DNS query failed with Rcode: %d", r.Rcode)
	}

	for _, ans := range r.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			for _, chunk := range txt.Txt {
				fmt.Println(chunk)
			}
		}
	}
}
