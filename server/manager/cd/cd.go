package cd

import (
	"context"
	"fmt"
	"main/pubsub"
	"os"
	"os/exec"
	"strings"

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
	for index, targetContext := range m.targetContexts {
		if targetContext.repository == repository {
			targetContext.canselFunc()
			// remove from slice
			m.targetContexts = append(m.targetContexts[:index], m.targetContexts[index+1:]...)
		}
	}
	path := strings.Split(getWebhook.Repository, "/")
	newContext := targetContext{
		repository: getWebhook.Repository,
		name:       path[len(path)-1],
		ctx:        context.Background(),
		canselFunc: func() {},
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	newContext.ctx = ctx
	newContext.canselFunc = cancelFunc

	m.targetContexts = append(m.targetContexts, newContext)
	go m.run(repository, &newContext)

}

func (m *CDManager) run(repository string, target *targetContext) {
	_ = os.Mkdir("repository", 0777)
	path := strings.Split(repository, "/")
	repositoryName := path[len(path)-1]

	directoryPath := "./repository/" + repositoryName
	_, err := os.Stat(directoryPath)
	if err != nil {
		cmd := exec.Command("git", "clone", repository)
		cmd.Dir = "./repository"
		if out, err := cmd.Output(); err != nil {
			fmt.Println("git clone failed", string(out))
			return
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = directoryPath
		if _, err := cmd.Output(); err != nil {
			fmt.Println("git pull failed")
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
