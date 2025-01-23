package service

import "github.com/nut-game/nano/pipeline"

type baseService struct {
	handlerHooks *pipeline.HandlerHooks
}

func (h *baseService) SetHandlerHooks(handlerHooks *pipeline.HandlerHooks) {
	h.handlerHooks = handlerHooks
}
