package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/arxdsilva/wpp"
)

func main() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		log.Fatal(fmt.Sprintf("error creating connection: %v\n", err))
	}
	err = login(wac)
	if err != nil {
		log.Fatal(fmt.Sprintf("error logging in: %v\n", err))
	}
	<-time.After(3 * time.Second)
	img, err := os.Open("image.jpg")
	if err != nil {
		log.Fatal(fmt.Sprintf("error reading file: %v", err))
	}
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "5521994445483@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: "Hello Gopher!",
		Content: img,
	}
	msgID, err := wac.Send(msg)
	if err != nil {
		log.Fatal(fmt.Sprintf("error sending message: %v", err))
	}
	fmt.Println("Message Sent -> ID : " + msgID)
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
