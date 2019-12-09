package healthcheck

import (
	"context"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"main/pubsub"
	"main/pubsub/systemevent"
	"net/http"
	"time"
)

type HealthCheckManager struct {
	db *gorm.DB
	ctx context.Context
	cancelFunc context.CancelFunc
}

func New(db *gorm.DB) *HealthCheckManager {
	db.AutoMigrate()
	ctx, cancelFunc := context.WithCancel(context.Background())
	healthCheckManager := &HealthCheckManager{db: db, ctx: ctx, cancelFunc: cancelFunc }
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
			m.AddHealthCheck("http://" + target.Proxy)
		}
	}
}

func (m *HealthCheckManager) AddHealthCheck(target string) {
	pubsub.SystemEvent.Pub(pubsub.System{
		Time:    time.Now(),
		Type:    systemevent.HEALTH_CHECK_REGISTER,
		Message: target,
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	m.ctx = ctx
	m.cancelFunc = cancelFunc
	go m.run(target,m.ctx)
}

func (m *HealthCheckManager) run(target string, ctx context.Context) {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			req, err := http.NewRequest("GET", target, nil)
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
					Message: "",
				})
			}
			pubsub.HealthCheckEvent.Pub(pubsub.HealthCheck{
				Target:  target,
				Status:  res.StatusCode,
				Message: string(body),
			})
		case <-ctx.Done():
			break
		}
	}
}
