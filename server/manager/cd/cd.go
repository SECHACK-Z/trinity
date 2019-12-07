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
	db             *gorm.DB
	targetContexts []targetContext
}

type targetContext struct {
	repository string
	name       string
	ctx        context.Context
	canselFunc context.CancelFunc
}

func New(db *gorm.DB) *CDManager {
	db.AutoMigrate()
	cdManager := &CDManager{db: db}
	pubsub.GetWebookEvent.Sub(cdManager.onGetWebhook)
	return cdManager
}

func (m *CDManager) onGetWebhook(getWebhook pubsub.GetWebook) {
	repository := getWebhook.Repository
	for _, targetContext := range m.targetContexts {
		if targetContext.repository == repository {
			m.RegisterContinuousDelivery(&targetContext)
			break
		}
	}
	path := strings.Split(getWebhook.Repository, "/")
	newContext := targetContext{
		repository: getWebhook.Repository,
		name:       path[len(path)-1],
		ctx:        context.Background(),
		canselFunc: func() {},
	}
	m.targetContexts = append(m.targetContexts, newContext)
	m.RegisterContinuousDelivery(&m.targetContexts[len(m.targetContexts)-1])
}

func (m *CDManager) RegisterContinuousDelivery(target *targetContext) {
	pubsub.SystemEvent.Pub(pubsub.System{
		Time:    time.Now(),
		Type:    systemevent.CD_REGISTER,
		Message: target.repository,
	})

	target.canselFunc()
	ctx, cancelFunc := context.WithCancel(context.Background())
	target.ctx = ctx
	target.canselFunc = cancelFunc
	go m.run(target.repository, target)

}

func (m *CDManager) run(repository string, target *targetContext) {
	if err := os.Mkdir("tmp", 0777); err != nil {
		fmt.Println(err)
	}
	path := strings.Split(repository, "/")
	repositoryName := path[len(path)-1]

	directoryPath := "./tmp/" + repositoryName
	_, err := os.Stat(directoryPath)
	if err != nil {
		cmd := exec.Command("git", "clone", repository)
		cmd.Dir = "./tmp"
		if out, err := cmd.Output(); err != nil {
			fmt.Println("git clone failed", string(out))
			return
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = directoryPath
		if out, err := cmd.Output(); err != nil {
			fmt.Println("git pull failed")
			fmt.Println(string(out))
			return
		}
	}

	cmd := exec.Command("go", "build", "-o", "main")
	cmd.Dir = directoryPath
	if err := cmd.Run(); err != nil {
		fmt.Println("Error at build")
		return
	}

	cmd = exec.Command("./main")
	cmd.Dir = directoryPath
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to exec command: ", err)
	}
	fmt.Println("Successfully exec start.")
	for {
		select {
		case <-target.ctx.Done():
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
