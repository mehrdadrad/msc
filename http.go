package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type crawler struct {
	client *http.Client
	config *config
	ch     chan string
}

func start(ctx context.Context, config *config) {
	c := &crawler{
		config: config,
	}

	c.init()

	for i := 0; i < c.config.Workers; i++ {
		go c.worker(ctx)
	}

	c.loop(ctx)
}

func (c *crawler) init() {
	c.ch = make(chan string, 1)
	c.client = &http.Client{
		Timeout: time.Duration(c.config.Timeout) * time.Second,
	}

}

func (c *crawler) loop(ctx context.Context) {
	for {
		for _, site := range c.config.Sites {
			select {
			case c.ch <- site:
			case <-ctx.Done():
				return
			}
		}

		select {
		case <-time.After(time.Duration(c.config.Interval) * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func (c *crawler) worker(ctx context.Context) {
	for {
		site := <-c.ch
		req, err := http.NewRequest(c.config.Method, site, nil)
		if err != nil {
			log.Println("http.request", err)
			continue
		}
		resp, err := c.client.Do(req)
		if err != nil {
			log.Println("http.client", err)
			continue
		}

		resp.Body.Close()
	}
}
