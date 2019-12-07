package cd

import (
	"context"
	"fmt"
	"main/pubsub"
	"main/pubsub/systemevent"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type CDManager struct {
	db         *gorm.DB
	ctx        context.Context
	canselFunc context.CancelFunc
}

func New(db *gorm.DB) *CDManager {
	db.AutoMigrate()
	cdManager := &CDManager{db: db}
	pubsub.UpdateConfigEvent.Sub(cdManager.onUpdateConfig)
	return cdManager
}

func (m *CDManager) onUpdateConfig(updateConfig pubsub.UpdateConfig) {

	for _, target := range config.Targets {
		if target.Repository != "" {
			m.RegisterContinuousDelivery(target.Repository)
		}
	}
}

func (m *CDManager) RegisterContinuousDelivery(repository string) {
	pubsub.SystemEvent.Pub(pubsub.System{
		Time:    time.Now(),
		Type:    systemevent.CD_REGISTER,
		Message: repository,
	})
	m.canselFunc()
	ctx, cancelFunc := context.WithCancel((context.Background()))
	m.ctx = ctx
	m.canselFunc = cancelFunc
	go m.run(repository)
}

func (m *CDManager) run(repository string) {
	if err := os.Mkdir("tmp", 0777); err != nil {
		fmt.Println(err)
	}
	path := strings.Split(repository, "/")
	repositoryName := path[len(path)-1]

	cmd := exec.Command("git", "clone", repository)
	cmd.Dir = "./tmp"
	if err := cmd.Run(); err != nil {
		fmt.Println("Error at git clone")
		return
	}

	cmd = exec.Command("go", "build", "-o", "main")
	cmd.Dir = "./tmp/" + repositoryName
	if err := cmd.Run(); err != nil {
		fmt.Println("Error at build")
		return
	}

	cmd = exec.Command("./main")
	cmd.Dir = "./tmp/" + repositoryName
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to exec command: ", err)
	}
	fmt.Println("Successfully exec start.")

	for {
		select {
		case <-m.ctx.Done():
			fmt.Println("Signal recieved.")
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("Failed to kill process:", err)
			} else {
				fmt.Println("Successfully killed process.")
			}
			return
		}
	}

}
