package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

const (
	IP_ADDRESS = "localhost"
	PORT       = 8095
)

type Server struct {
	add net.UDPAddr
}

func NewServer() *Server {
	return &Server{
		add: net.UDPAddr{
			Port: PORT,
			IP:   net.ParseIP(IP_ADDRESS),
		},
	}
}

func (s *Server) Run() error {
	log.Print("Starting UDP server ", s.add.Port)
	conn, err := net.ListenUDP("udp", &s.add)
	if err != nil {
		return err
	}

	ListenLoop(conn)
	return nil
}

func ListenLoop(conn *net.UDPConn) {
	for {
		buff := make([]byte, 1024)
		_, clientAddr, _ := conn.ReadFromUDP(buff)

		msg := new(dns.Msg)
		err := msg.Unpack(buff)
		if err != nil {
			log.Println("Error decoding DNS response:", err)
			continue
		}

		route := msg.Question[0].Name
		log.Print(route)

		resp := router(msg, route)
		responseData, _ := resp.Pack()
		_, err = conn.WriteToUDP(responseData, clientAddr)
	}
}

func router(req *dns.Msg, route string) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetReply(req)

	var answer dns.RR

	parts := strings.Split(route, ".")

	if len(parts) < 3 || parts[0] != "movie" || parts[1] != "info" {
		answer, _ = dns.NewRR("movie.info. IN TXT \"invalid input\"")
		resp.Insert([]dns.RR{answer})
		return resp
	}

	var movie_name string
	for i := 2; i < len(parts); i += 1 {
		movie_name = movie_name + parts[i] + " "
	}

	movie, err := GetMovieDescription(API_KEY, movie_name)
	if err != nil {
		log.Print(err)
	}

	answer, _ = dns.NewRR(fmt.Sprintf("movie.info. IN TXT \"%s\"", movie.Description))

	resp.Insert([]dns.RR{answer})
	return resp
}

func main() {
	server := NewServer()
	err := server.Run()
	if err != nil {
		log.Println("Error ", err)
	}
}
