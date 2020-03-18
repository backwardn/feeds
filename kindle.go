package feeds

import (
	"database/sql"
	"encoding/json"
	"github.com/jordan-wright/email"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/smtp"
)

const htmlDir = "articles"
const mobiDir = "output/mobi"

//const dbFilePath = "feeds.db"

type smtpCreds struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	From     string `json:"from"`
	User     string `json:"user"`
	Password string `json:"password"`
}
type destination struct {
	To string `json:"to"`
}

func DispatchToKindle(subject string, attachment string, c *sql.DB) error {
	targets := "SELECT targets.data, outputs.credentials FROM targets INNER JOIN outputs ON outputs.id = targets.output_id WHERE outputs.type = 'kindle'"
	r, err := c.Query(targets)
	if err != nil {
		return err
	}
	for r.Next() {
		var data []byte
		var credentials []byte

		r.Scan(&data, &credentials)

		var settings smtpCreds
		var target destination
		err = json.Unmarshal(credentials, &settings)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &target)
		if err != nil {
			return err
		}
		e := email.NewEmail()
		log.Printf("Emailing %s to %s", attachment, target.To)

		e.From = settings.From
		e.To = []string{target.To}
		e.Cc = []string{"marius.orcsik@gmail.com"}
		e.Subject = subject
		e.AttachFile(attachment)
		err = e.Send(settings.Server+":"+settings.Port, smtp.PlainAuth("", settings.User, settings.Password, settings.Server))
		if err != nil {
			return err
		}
	}
	return nil
}