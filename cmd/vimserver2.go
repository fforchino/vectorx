package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/parser"
	"github.com/emiago/sipgo/sip"
	"github.com/emiago/sipgo/transport"
	"github.com/icholy/digest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	extIP, err := getMyExternalIP()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get external ip address!")
		return
	}
	extIP = extIP + ":5060"
	dst := flag.String("srv", "sip.linphone.org:5060", "Destination")
	username := flag.String("u", "", "SIP Username")
	password := flag.String("p", "", "Password")
	sipdebug := flag.Bool("sipdebug", false, "Turn on sipdebug")
	flag.Parse()

	// Make SIP Debugging available
	transport.SIPDebug = *sipdebug

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.StampMicro,
	}).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	// Setup UAC
	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent(*username),
		sipgo.WithUserAgentIP(extIP),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to setup user agent")
	}

	srv, err := sipgo.NewServer(ua)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to setup server handle")
	}

	client, err := sipgo.NewClient(ua)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to setup client handle")
	}

	ctx := context.TODO()
	go srv.ListenAndServe(ctx, "udp", extIP)

	// Wait that ouir server loads
	time.Sleep(1 * time.Second)

	// Create basic REGISTER request structure
	recipient := &sip.Uri{}
	parser.ParseUri(fmt.Sprintf("sip:%s@%s", *username, *dst), recipient)
	req := sip.NewRequest(sip.REGISTER, recipient, "SIP/2.0")
	req.AppendHeader(
		sip.NewHeader("Contact", fmt.Sprintf("<sip:%s@%s>", *username, extIP)),
	)

	// Send request and parse response
	// req.SetDestination(*dst)
	tx, err := client.TransactionRequest(req.Clone())
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to create transaction")
	}
	defer tx.Terminate()

	res, err := getResponse(tx)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to get response")
	}

	log.Info().Int("status", int(res.StatusCode())).Msg("Received status")
	if res.StatusCode() == 401 {
		// Get WwW-Authenticate
		wwwAuth := res.GetHeader("WWW-Authenticate")
		chal, err := digest.ParseChallenge(wwwAuth.Value())
		if err != nil {
			log.Fatal().Str("wwwauth", wwwAuth.Value()).Err(err).Msg("Fail to parse challenge")
		}

		// Reply with digest
		cred, _ := digest.Digest(chal, digest.Options{
			Method:   req.Method.String(),
			URI:      recipient.Host,
			Username: *username,
			Password: *password,
		})

		newReq := req.Clone()
		newReq.AppendHeader(sip.NewHeader("Authorization", cred.String()))

		tx, err = client.TransactionRequest(newReq)
		if err != nil {
			log.Fatal().Err(err).Msg("Fail to create transaction")
		}
		defer tx.Terminate()

		res, err = getResponse(tx)
		if err != nil {
			log.Fatal().Err(err).Msg("Fail to get response")
		}
	}

	if res.StatusCode() != 200 {
		log.Fatal().Msg("Fail to register")
	}

	log.Info().Msg("Client registered")
}

func getResponse(tx sip.ClientTransaction) (*sip.Response, error) {
	select {
	case <-tx.Done():
		return nil, fmt.Errorf("transaction died")
	case res := <-tx.Responses():
		return res, nil
	}
}

func getMyExternalIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return string(ip), err
}

/*package main

import (
	"fmt"
	"github.com/emiago/sipgo"
	"os"
	"os/signal"
	"syscall"
)

type SipClient struct {
	UserAgent *sipgo.UserAgent
}

func NewSipClient(sipAddress, password string) (*SipClient, error) {
	config := &sipgo.UserAgentOption{}

	ua, err := sipgo.NewUA(config)
	if err != nil {
		return nil, err
	}

	client := &SipClient{
		UserAgent: ua,
	}

	return client, nil
}

func (c *SipClient) sendTextMessage(to, message string) error {
	return c.UserAgent(to, message)
}

func (c *SipClient) onTextMessageReceived(f func(from, message string)) {
	c.UserAgent.OnMessage(f)
}

func main() {
	sipAddress := "sip:user@example.com"
	password := "your_password"

	_, err := NewSipClient(sipAddress, password)
	if err != nil {
		fmt.Println("Error initializing SIP client:", err)
		os.Exit(1)
	}

	// Send a text message
	err = client.sendTextMessage("sip:receiver@example.com", "Hello, this is a test message.")
	if err != nil {
		fmt.Println("Error sending text message:", err)
	}

	// Handle incoming text messages
	client.onTextMessageReceived(func(from, message string) {
		fmt.Printf("Received message from %s: %s\n", from, message)
	})

	// Wait for the application to be interrupted
	waitForInterrupt()

	// Close the SIP client
	client.UserAgent.Close()
}

func waitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
*/
