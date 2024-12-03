package samba

import (
	"net"
	"os"

	"github.com/hirochachacha/go-smb2"
	"github.com/rs/zerolog/log"
)

func UploadFile(filename string, fileBytes []byte) (err error) {
	conn, err := net.Dial("tcp", os.Getenv("SAMBA_ADDRESS"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     os.Getenv("SAMBA_USER"),
			Password: os.Getenv("SAMBA_PASSWORD"),
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	fs, err := s.Mount(os.Getenv("SAMBA_MOUNT"))
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	err = fs.WriteFile(filename, fileBytes, 0444)
	if err != nil {
		log.Error().Msg("Failed UploadFile File")
		return
	}
	log.Info().Msg("Successfully UploadFile File")
	return
}

func RemoveFile(filename string) (err error) {
	conn, err := net.Dial("tcp", os.Getenv("SAMBA_ADDRESS"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     os.Getenv("SAMBA_USER"),
			Password: os.Getenv("SAMBA_PASSWORD"),
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	fs, err := s.Mount(os.Getenv("SAMBA_MOUNT"))
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	err = fs.Remove(filename)
	if err != nil {
		log.Error().Msg("Failed remove file")
		return
	}
	log.Info().Msg("Successfully remove file")
	return
}
