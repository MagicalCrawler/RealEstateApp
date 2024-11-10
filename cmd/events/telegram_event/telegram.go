package telegram_event

import (
	"context"
	"errors"
	"log"

	"github.com/MagicalCrawler/RealEstateApp/cmd/clients/telegram"
	"github.com/MagicalCrawler/RealEstateApp/cmd/err"
	"github.com/MagicalCrawler/RealEstateApp/cmd/events"
)

type Processor struct {
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client) *Processor { //, storage storage.Storage)*Processor{
	return &Processor{
		tg:     client,
		offset: 0,
		// storage: storage,
	}
}
func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, er := p.tg.Updates(ctx, p.offset, limit)
	if er != nil {
		return nil, err.Wrap("can't get events", er)
	}

	log.Printf("Received %d updates", len(updates))

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return err.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, er := meta(event)
	if er != nil {
		return err.Wrap("can't process message", er)
	}

	if er := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); er != nil {
		return err.Wrap("can't process message", er)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, err.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}
func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
