/*
Copyright © 2026 mikhi32
*/
package cmd

// Imports
import (
	"context"
	"fmt"
	"slices"
	"strings"

	"NewHakuneko/rippers"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Don't ask why I made a function here... arrays in go are way too complicated lol.
func getGigaReaders() []string {
	gigaReaders := [...]string{
		"comic-days.com",
		"ichicomi.com",
		"shonenjumpplus.com",
		"tonarinoyj.jp",
		"andsofa.com",
		"morningtwo.com",
		"getsumagakichi.com",
		"bibliosirius.com",
		"kuragebunch.com",
		"comicbunch-kai.com",
		"viewer.heros-web.com",
		"comicborder.com",
		"comic-gardo.com",
		"comic-zenon.com",
		"magcomi.com",
		"comic-action.com",
		"comic-trail.com",
		"feelweb.jp",
		"www.sunday-webry.com",
		"comic-ogyaaa.com",
		"comic-earthstar.com",
		"ourfeel.jp",
		"comic-seasons.com",
		"comic-y-ours.com",
	}
	return gigaReaders[:]
}

func getComiciReaders() []string {
	comiciReaders := [...]string{
		"rimacomiplus.jp",
		"younganimal.com",
		"takecomic.jp",
		"kimicomi.com",
		"youngchampion.jp",
		"mangalt.jp",
		"mangaspa.nikkan-spa.jp",
		"piacomic.jp",
		"comic-room-base.com",
		"manga-zegra.com",
		"comic-growl.com",
		"asacomi.jp",
		"comicpash.jp",
		"comic.j-nbooks.jp",
		"hayacomic.jp",
		"championcross.jp",
		"kansai.mag-garden.co.jp",
		"comicride.jp",
		"carula.jp",
		"comic-medu.com",
		"comics.manga-bang.com",
		"bigcomics.jp",
		"studio.booklista.co.jp",
		"hanayume.com",
	}
	return comiciReaders[:]
}

func getPixiv() []string {
	comiciReaders := [...]string{
		"www.pixiv.net",
	}
	return comiciReaders[:]
}

// Main function to organize everything
func RipMain(url string, cookie string, ctx context.Context, folder string) {
	if url == "" {
		runtime.EventsEmit(ctx, "error-emit", "You need to send something for us to rip!")
	} else {
		// Check what reader we are using.
		if !strings.HasPrefix(url, "https://") {
			runtime.EventsEmit(ctx, "error-emit", "Please submit a url!")
		} else {
			fmt.Println(strings.Split(url, "/")[2])
			gigaReaders := getGigaReaders()
			comiciReaders := getComiciReaders()
			pixiv := getPixiv()
			isGiga := slices.Contains(gigaReaders, strings.Split(url, "/")[2])
			isComici := slices.Contains(comiciReaders, strings.Split(url, "/")[2])
			isPixiv := slices.Contains(pixiv, strings.Split(url, "/")[2])
			// Rip accordingly.
			if isGiga {
				rippers.RipGiga(url, cookie, ctx, folder)
			} else if isComici {
				rippers.RipComici(url, cookie, ctx, folder)
				} else if isPixiv {
				rippers.RipPixiv(url, cookie, ctx, folder)
			} else {
				runtime.EventsEmit(ctx, "error-emit", "We currently do not support this url. Apologies!")
			}
		}
	}
}
