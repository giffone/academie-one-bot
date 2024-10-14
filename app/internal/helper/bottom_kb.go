package helper

import (
	dm "ac_bot/internal/domain"
	"fmt"
	"log"
	"net/url"

	"github.com/go-telegram/bot/models"
)

func ParseBottomKb(keyboard dm.StKb) [][]models.KeyboardButton {
	log.Print(keyboard)
	tKb := make([][]models.KeyboardButton, len(keyboard))

	for i, kb := range keyboard {
		tKb[i] = make([]models.KeyboardButton, len(kb))
		for j, k := range kb {
			tKb[i][j] = models.KeyboardButton{
				Text: k.Text,
			}
			if k.WebAppURL != "" {
				tKb[i][j].WebApp = &models.WebAppInfo{
					URL: k.WebAppURL,
				}
			}
		}
	}

	return tKb
}

func StudentStKb(wa *dm.WebApp) dm.StKb {
	return dm.StKb{
		{
			dm.StKbButton{
				Text:      "QR Scan",
				WebAppURL: wa.Paths[dm.WebAppPathQrScan],
			},
		},
		{
			dm.StKbButton{
				Text:      "App",
				WebAppURL: wa.Addr,
			},
			dm.StKbButton{
				Text: "About",
			},
		},
	}
}

func GuestStKb(wa *dm.WebApp, orgs string) dm.StKb {
	return dm.StKb{
		{
			dm.StKbButton{
				Text: "About",
			},
		},
		{
			dm.StKbButton{
				Text: "Campus Invitation",
				WebAppURL: fmt.Sprintf("%s?orgs=%s",
					wa.Paths[dm.WebAppPathRegForm],
					url.QueryEscape(orgs),
				),
			},
		},
	}
}

func AdminStKbButton(wa *dm.WebApp, orgs, roles string) []dm.StKbButton {
	return []dm.StKbButton{
		{
			Text: "Admin",
			WebAppURL: fmt.Sprintf("%s?orgs=%s&roles=%s",
				wa.Paths[dm.WebAppPathAdmin],
				url.QueryEscape(orgs),
				url.QueryEscape(roles),
			),
		},
	}
}

func SecurityStKb(wa *dm.WebApp) dm.StKb {
	return dm.StKb{
		{
			dm.StKbButton{
				Text: "Post",
			},
			dm.StKbButton{
				Text: "Post stop",
			},
		},
	}
}
