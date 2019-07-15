package watcher

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

// Watcher configures the watcher
type Watcher struct {
	MaxEvents  int
	Handler    func(e *watcher.Event) error
	Folder     string
	Interval   time.Duration
	Pattern    *regexp.Regexp
	isWatching bool
}

// Stop stops the watcher
func (c *Watcher) Stop() {
	c.isWatching = false
}

// Watch starts observing given folder at configured interval
func (c *Watcher) Watch(done chan bool) error {
	w := watcher.New()
	w.SetMaxEvents(c.MaxEvents)
	w.FilterOps(watcher.Write, watcher.Remove, watcher.Rename, watcher.Move)
	w.AddFilterHook(watcher.RegexFilterHook(c.Pattern, false))

	if err := w.AddRecursive(c.Folder); err != nil {
		return err
	}

	go func() {
		defer w.Close()
		log.Println(fmt.Sprintf("watching %s", c.Folder))
		c.isWatching = true

		for c.isWatching {
			select {
			case event := <-w.Event:
				err := c.Handler(&event)
				if err != nil {
					log.Println(fmt.Sprintf("watcher handler error: %s", err.Error()))
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			case <-time.After(time.Second):
				continue
			}
		}
		done <- true
	}()

	return w.Start(c.Interval)
}
