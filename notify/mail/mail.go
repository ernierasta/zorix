package mail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/ernierasta/spock/shared"
)

// Send sends mail via smtp.
// Supports multiple recepients, TLS (port 465)/StartTLS(ports 25,587, any other).
// Mail should always valid (correctly encoded subject and body).
func Send(c shared.Check, n shared.Notification) {
	auth := smtp.PlainAuth("", n.User, n.Pass, n.Server)

	recipients := strings.Join(n.To, ", ")

	header := make(map[string]string)
	header["From"] = n.From
	header["To"] = recipients
	header["Date"] = c.Timestamp.Format(time.RFC1123Z)
	header["Subject"] = encodeRFC2047(n.Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(n.Text))

	SendMail(n.Server, n.Port, auth, false, n.From, n.To, []byte(message))
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	log.Println(addr.String())
	return strings.Trim(addr.String(), " <@>")
}

// SendMail connects to the server at addr, switches to TLS if
// possible, authenticates with the optional mechanism a if possible,
// and then sends an email from address from, to addresses to, with
// message msg.
// The addr must include a port, as in "mail.example.com:smtp".
//
// The addresses in the to parameter are the SMTP RCPT addresses.
//
// The msg parameter should be an RFC 822-style email with headers
// first, a blank line, and then the message body. The lines of msg
// should be CRLF terminated. The msg headers should usually include
// fields such as "From", "To", "Subject", and "Cc".  Sending "Bcc"
// messages is accomplished by including an email address in the to
// parameter but not including it in the msg headers.
//
// The SendMail function and the net/smtp package are low-level
// mechanisms and provide no support for DKIM signing, MIME
// attachments (see the mime/multipart package), or other mail
// functionality. Higher-level packages exist outside of the standard
// library.
func SendMail(host string, port int, a smtp.Auth, ignoreCert bool, from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	hostPort := fmt.Sprintf("%s:%d", host, port)

	c := &smtp.Client{}

	if port == 465 {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: ignoreCert,
			ServerName:         host,
		}
		conn, err := tls.Dial("tcp", hostPort, tlsconfig)
		if err != nil {
			return err
		}
		c, err = smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
	} else {
		var err error
		c, err = smtp.Dial(hostPort)
		if err != nil {
			return err
		}
	}
	defer c.Close()
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	if err = c.Hello(hostname); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName: host}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func check(err error) {
	if err != nil {
		log.Println("error sending mail, err: ", err)
	}
}

// validateLine checks to see if a line has CR or LF as per RFC 5321
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func tlsSend(auth smtp.Auth, server, srvAndPort, from, recipients, message string, isIgnoreCert bool) error {

	tlsconfig := &tls.Config{
		InsecureSkipVerify: isIgnoreCert,
		ServerName:         server,
	}
	conn, err := tls.Dial("tcp", srvAndPort, tlsconfig)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, server)
	err = client.Auth(auth)
	if err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	if err = client.Rcpt(recipients); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	client.Quit()

	return nil
}
