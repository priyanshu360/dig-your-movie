package dns

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/priyanshu360/dig-your-movie/internal/config"
	"github.com/priyanshu360/dig-your-movie/internal/omdb"
)

type Server struct {
	addr       string
	omdbClient *omdb.Client
}

func NewServer(cfg *config.Config, omdbClient *omdb.Client) *Server {
	return &Server{
		addr:       fmt.Sprintf(":%d", cfg.Port),
		omdbClient: omdbClient,
	}
}

func (s *Server) Run() error {
	log.Printf("Starting UDP server on %s", s.addr)
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		buf := make([]byte, 512)
		_, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		go s.handleRequest(conn, clientAddr, buf)
	}
}

func (s *Server) handleRequest(conn *net.UDPConn, clientAddr *net.UDPAddr, buf []byte) {
	msg := new(dns.Msg)
	if err := msg.Unpack(buf); err != nil {
		log.Printf("Error decoding DNS message: %v", err)
		return
	}

	if len(msg.Question) == 0 {
		return
	}

	question := msg.Question[0]
	log.Printf("Received query: %s", question.Name)

	resp := s.router(msg, question.Name)
	respData, err := resp.Pack()
	if err != nil {
		log.Printf("Error packing response: %v", err)
		return
	}

	if _, err := conn.WriteToUDP(respData, clientAddr); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func (s *Server) router(req *dns.Msg, route string) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetReply(req)

	parts := strings.Split(strings.TrimSuffix(route, "."), ".")

	// Expecting query like: movie.info.<movie_name>
	// parts: ["movie", "info", "movie_name..."]
	if len(parts) < 3 || parts[0] != "movie" || parts[1] != "info" {
		txt, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN TXT \"invalid request format\"", req.Question[0].Name))
		resp.Answer = append(resp.Answer, txt)
		return resp
	}

	movieName := strings.Join(parts[2:], " ")
	movie, err := s.omdbClient.GetMovieDescription(context.Background(), movieName)

	var answerText string
	if err != nil {
		log.Printf("Error fetching movie: %v", err)
		answerText = fmt.Sprintf("Error: %v", err)
	} else {
		answerText = movie.Description
		// Truncate if too long for a single TXT record (255 chars limit usually, but max packet size matters)
		if len(answerText) > 250 {
			answerText = answerText[:247] + "..."
		}
	}

	txt, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN TXT \"%s\"", req.Question[0].Name, strings.ReplaceAll(answerText, "\"", "\\\"")))
	resp.Answer = append(resp.Answer, txt)

	return resp
}
