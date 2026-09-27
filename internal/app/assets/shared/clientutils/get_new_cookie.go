package clientutils

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kartFr/Asset-Reuploader/internal/app/assets/shared/permissions"
	"github.com/kartFr/Asset-Reuploader/internal/app/config"
	"github.com/kartFr/Asset-Reuploader/internal/app/context"
	"github.com/kartFr/Asset-Reuploader/internal/app/request"
	"github.com/kartFr/Asset-Reuploader/internal/color"
	"github.com/kartFr/Asset-Reuploader/internal/console"
	"github.com/kartFr/Asset-Reuploader/internal/files"
)

var (
	cookieMu       sync.Mutex
	cookieInFlight bool
	cookieWaiters  []chan struct{}
)

func waitForCookieRefresh() bool {
	cookieMu.Lock()
	if !cookieInFlight {
		cookieInFlight = true
		cookieMu.Unlock()
		return true // caller is leader
	}
	ch := make(chan struct{})
	cookieWaiters = append(cookieWaiters, ch)
	cookieMu.Unlock()
	<-ch
	return false // follower: cookie already refreshed
}

func finishCookieRefresh() {
	cookieMu.Lock()
	cookieInFlight = false
	waiters := cookieWaiters
	cookieWaiters = nil
	cookieMu.Unlock()
	for _, ch := range waiters {
		close(ch)
	}
}

func GetNewCookie(ctx *context.Context, r *request.Request, m string) {
	if !waitForCookieRefresh() {
		ctx.PauseController.WaitIfPaused()
		return
	}
	defer finishCookieRefresh()

	pauseController := ctx.PauseController

	if !pauseController.Pause() {
		pauseController.WaitIfPaused()
		return
	}

	console.ClearScreen()

	client := ctx.Client
	inputErr := errors.New(m)
	for {
		fmt.Print(ctx.Logger.HistoryString())
		color.Error.Println(inputErr)

		i, err := console.LongInput("ROBLOSECURITY: ")
		console.ClearScreen()
		if err != nil {
			inputErr = err
			continue
		}

		fmt.Println("Authenticating cookie...")
		err = client.SetCookie(i)
		console.ClearScreen()
		if err != nil {
			inputErr = err
			continue
		}

		fmt.Println("Checking if account can edit universe...")
		err = permissions.CanEditUniverse(ctx, r)
		console.ClearScreen()
		if err != nil {
			inputErr = err
			continue
		}

		break
	}

	fmt.Print(ctx.Logger.HistoryString())

	if err := files.Write(config.Get("cookie_file"), client.Cookie); err != nil {
		ctx.Logger.Error("Failed to save cookie: ", err)
	}

	pauseController.Unpause()
}
