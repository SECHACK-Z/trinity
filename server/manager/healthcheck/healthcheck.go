package healthcheck

import (
	"context"
	"io/ioutil"
	"main/config"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type HealthCheckManager struct {
	db         *gorm.DB
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New(db *gorm.DB) *HealthCheckManager {
	db.AutoMigrate()
	ctx, cancelFunc := context.WithCancel(context.Background())
	healthCheckManager := &HealthCheckManager{db: db, ctx: ctx, cancelFunc: cancelFunc}
	pubsub.UpdateConfigEvent.Sub(healthCheckManager.onUpdateConfig)
	return healthCheckManager
}

func (m *HealthCheckManager) onUpdateConfig(updateConfig pubsub.UpdateConfig) {
	config := updateConfig.Config
	m.cancelFunc()
	ctx, cancelFunc := context.WithCancel(context.Background())
	m.ctx = ctx
	m.cancelFunc = cancelFunc
	for _, target := range config.Targets {
		if target.HealthCheck {
			m.AddHealthCheck(target)
		}
	}
}

func (m *HealthCheckManager) AddHealthCheck(target config.Target) {
	pubsub.SystemEvent.Pub(pubsub.System{
		Time:    time.Now(),
		Type:    systemevent.HEALTH_CHECK_REGISTER,
		Message: target.Proxy,
	})

	go m.run(target, m.ctx)
}

func (m *HealthCheckManager) run(target config.Target, ctx context.Context) {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			req, err := http.NewRequest("GET", "http://"+target.Proxy, nil)
			if err != nil {
				// TODO
				continue
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				pubsub.HealthCheckEvent.Pub(pubsub.HealthCheck{
					Target:  target,
					Status:  502,
					Message: err.Error(),
				})
				continue
			}
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				pubsub.HealthCheckEvent.Pub(pubsub.HealthCheck{
					Target:  target,
					Status:  res.StatusCode,
					Message: err.Error(),
				})
				continue
			}
			pubsub.HealthCheckEvent.Pub(pubsub.HealthCheck{
				Target:  target,
				Status:  res.StatusCode,
				Message: string(body),
			})

			// check favicon
			res, err = http.Get("http://" + target.Proxy + "/favicon.ico")
			if err == nil {
				if res.Header.Get("Content-type") != "image/x-icon" {
					pubsub.NewsEvent.Pub(pubsub.News{
						Type:   "favicon is missing",
						Target: target,
					})
				}
			}
		case <-ctx.Done():
			break
		}
	}
}
