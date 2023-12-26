package handler

import "NewsChanel/internal/service"

type Handler struct {
	srv service.Service
}

func New(service service.Service) Handler {
	return Handler{
		srv: service,
	}
}
