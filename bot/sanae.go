package bot

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/keyo-git/sanae-bot/api"

	// Postgres
	_ "github.com/lib/pq"
)

// Sanae ...
type Sanae struct {
	dg    *discordgo.Session
	exAPI *api.ExHentaiAPI
	db    *sql.DB
}

// NewSanae initalizes new Sanae object
func NewSanae(cfg *Config) (*Sanae, error) {
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	log.Println("Created Discord API handle")

	cookies := []*http.Cookie{}
	for _, cookie := range cfg.Cookies {
		cookies = append(cookies, &http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Domain:  ".exhentai.org",
			Path:    "/",
			Expires: time.Now().Add(356 * 24 * time.Hour),
		})
	}

	exAPI, err := api.NewExHentaiAPI(cookies)
	if err != nil {
		return nil, err
	}

	log.Println("Created ExHentai API handle")

	db, err := sql.Open("postgres", cfg.ConnStr)
	if err != nil {
		return nil, err
	}

	log.Println("Created Database handle")

	s := &Sanae{
		dg:    dg,
		exAPI: exAPI,
		db:    db,
	}

	s.dg.AddHandler(s.messageCreate)
	s.dg.AddHandler(s.messageReactionAdd)

	return s, err
}

// Open creates websocket connection to Discord
func (s *Sanae) Open() error { return s.dg.Open() }

// Close closes websocket connection to Discord
func (s *Sanae) Close() error { return s.dg.Close() }

// DbHandle returns handle to database
func (s *Sanae) DbHandle() *sql.DB { return s.db }

// ExAPI returns handle to ExHentai API
func (s *Sanae) ExAPI() *api.ExHentaiAPI { return s.exAPI }

// Sess returns handle to Discord API
func (s *Sanae) Sess() *discordgo.Session { return s.dg }

func (s *Sanae) messageCreate(sess *discordgo.Session, m CmdTrigger) {
	if m.Author.ID == sess.State.User.ID {
		return
	}

	m.Content = strings.Replace(m.Content, "\n", " ", -1)
	argv := strings.Split(m.Content, " ")

	for _, c := range GlobalRegistry.cmdRegistry {
		if c.Trigger == argv[0] {
			if err := c.Cmd(s, argv[1:], m); err != nil {
				log.Println(err)
			}
			break
		}
	}

	if err := ex(s, argv, m); err != nil {
		log.Println(err)
	}

}

func (s *Sanae) messageReactionAdd(sess *discordgo.Session, m ReactTrigger) {
	msg, err := sess.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println(err)
		return
	}

	if m.UserID == sess.State.User.ID || msg == nil ||
		msg.Author.ID != sess.State.User.ID {
		return
	}

	for _, c := range GlobalRegistry.reactionRegistry {
		if c.Emoji == m.Emoji.Name {
			if err = c.Cmd(s, m); err != nil {
				log.Println(err)
			}
			break
		}
	}
}
