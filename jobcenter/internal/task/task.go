package task

import (
	"github.com/go-co-op/gocron"
	"jobcenter/internal/logic"
	"jobcenter/internal/svc"
	"time"
)


type Task struct {
	s   *gocron.Scheduler
	ctx *svc.ServiceContext
}

func NewTask(ctx *svc.ServiceContext) *Task {
	return &Task{
		s:   gocron.NewScheduler(time.UTC),
		ctx: ctx,
	}
}


func (t *Task) Run() {



	t.s.Every(1).Minute().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("1m")
	})
	t.s.Every(3).Minute().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("3m")
	})
	t.s.Every(5).Minute().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("5m")
	})
	t.s.Every(15).Minute().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("15m")
	})
	t.s.Every(30).Minute().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("30m")
	})
	t.s.Every(1).Hour().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("1H")
	})
	t.s.Every(2).Hour().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("2H")
	})
	t.s.Every(4).Hour().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("4H")
	})
	t.s.Every(1).Day().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("1D")
	})
	t.s.Every(1).Week().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("1W")
	})
	t.s.Every(1).Month().Do(func() {
		logic.NewKline(t.ctx.Config.Okx, t.ctx.MongoClient, t.ctx.KafkaClient, t.ctx.Cache).Do("1M")
	})

	
	t.s.Every(1).Minute().Do(func() {
		logic.NewRate(t.ctx.Config.Okx, t.ctx.Cache).Do()
	})

}


func (t *Task) StartBlocking() {
	t.s.StartBlocking()
}

func (t *Task) Stop() {
	t.s.Stop()
}
