package server

import (
	"ac_bot/config"
	"ac_bot/internal/api/handler"
	"ac_bot/internal/domain"
	"ac_bot/internal/repository/postgres"
	"ac_bot/internal/service"
	"ac_bot/internal/service/tgs"

	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Server interface {
	Start(ctx context.Context)
}

type server struct {
	b   *bot.Bot
	env *Env // envorinments [postgres and etc...]
}

func NewServer(ctx context.Context, cfg *config.Config) Server {
	var err error

	s := server{env: NewEnv(ctx, cfg)}

	// storage
	storage := postgres.New(s.env.pool)

	// service
	svc, tsvc := service.New(ctx, storage, cfg)

	// handlers
	h := handler.New(cfg, svc)

	opts := []bot.Option{
		bot.WithMiddlewares(showMessageWithUserName),
		bot.WithDefaultHandler(h.Default),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, h.Callback),
	}

	s.b, err = bot.New(cfg.Token, opts...)
	if nil != err {
		log.Fatalf("create bot error: %s\n", err.Error())
	}

	// telegram service
	tsvc.Register(tgs.New(s.b))

	// --- Chat (Inline Mode)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.Default)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "/about", bot.MatchTypeExact, h.About)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "About", bot.MatchTypeExact, h.About)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "Post", bot.MatchTypeExact, h.PostStart)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "Post stop", bot.MatchTypeExact, h.PostStop)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, "/whoami", bot.MatchTypeExact, h.Whoami)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, domain.InPathOrgs, bot.MatchTypeExact, h.GetList)
	s.b.RegisterHandler(bot.HandlerTypeMessageText, domain.InPathRoles, bot.MatchTypeExact, h.GetList)
	// ---- Web App
	s.b.RegisterHandlerMatchFunc(webAppMatchFunc, h.WebApp)

	return &s
}

func (s *server) Start(ctx context.Context) {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.Print("[app] server is straring...")
	s.b.Start(ctx)

	defer s.env.Stop(ctx)
}

func webAppMatchFunc(update *models.Update) bool {
	log.Print("handleWebAppDataMatchFunc")
	return update.Message != nil && update.Message.WebAppData != nil && update.Message.WebAppData.Data != ""
}

func showMessageWithUserName(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Printf("[msg] %s say: %s", update.Message.From.FirstName, update.Message.Text)
		} else if update.CallbackQuery != nil {
			log.Printf("[clb] %s say: %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.Data)
		}
		next(ctx, b, update)
	}
}
