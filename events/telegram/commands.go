package telegram

import (
	"errors"
	"fmt"
	"github.com/sungatuly22/readAdvisorBot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command %s from %s", text, username)
	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatId int, pageURL string, username string) error {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("can't do command: %w", err)
	}
	if isExists {
		return p.tg.SendMessage(chatId, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return fmt.Errorf("can't save the page: %w", err)
	}
	if err := p.tg.SendMessage(chatId, msgSaved); err != nil {
		return fmt.Errorf("can't do command: %w", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatId int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return fmt.Errorf("can't do command: can't send random: %w", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return fmt.Errorf("can't do command: can't send random: %w", err)
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err != nil && u.Host != ""
}
