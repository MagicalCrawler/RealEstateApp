package telegram_event

import (
	"context"
	"log"
	"net/url"
	"strings"	
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (er error) {
	// defer func() { er = e.WrapIfErr("can't do command: save page", er) }()

	// page := &storage.Page{
	// 	URL:      pageURL,
	// 	UserName: username,
	// }

	// isExists, er := p.storage.IsExists(ctx, page)
	// if er != nil {
	// 	return er
	// }
	// if isExists {
	// 	return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
	// }

	// if er := p.storage.Save(ctx, page); er != nil {
	// 	return er
	// }

	if er := p.tg.SendMessage(ctx, chatID, msgSaved); er != nil {
		return er
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (er error) {
	// defer func() { er = err.WrapIfErr("can't do command: can't send random", er) }()

	// page, er := p.storage.PickRandom(ctx, username)
	// if er != nil && !errors.Is(er, storage.ErrNoSavedPages) {
	// 	return er
	// }
	// if errors.Is(er, storage.ErrNoSavedPages) {
	// 	return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	// }

	// if er := p.tg.SendMessage(ctx, chatID, page.URL); er != nil {
	// 	return er
	// }

	// return p.storage.Remove(ctx, page)
	return
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, er := url.Parse(text)

	return er == nil && u.Host != ""
}
