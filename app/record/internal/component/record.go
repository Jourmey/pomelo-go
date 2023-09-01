package component

import (
	"fmt"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/component"
)

type record struct {
	svcCtx *svc.ServiceContext
}

func (r *record) Init() {
	fmt.Println("record.Init")
}

func (r *record) AfterInit() {
	fmt.Println("record.AfterInit")
}

func (r *record) BeforeShutdown() {
	fmt.Println("record.BeforeShutdown")
}

func (r *record) Shutdown() {
	fmt.Println("record.Shutdown")
}

func NewRecord(svcCtx *svc.ServiceContext) component.Component {
	return &record{
		svcCtx: svcCtx,
	}
}
