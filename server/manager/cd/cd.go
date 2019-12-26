package cd

import (
	"context"
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
	targetContexts map[string]targetContext
}

type targetContext struct {
	repository string
	name       string
	ctx        context.Context
	canselFunc context.CancelFunc
}

func New(db *gorm.DB) *CDManager {
	db.AutoMigrate()
	cdManager := &CDManager{db: db, targetContexts: map[string]targetContext{}}
	pubsub.GetWebhookEvent.Sub(cdManager.onGetWebhook)
	return cdManager
}

func (m *CDManager) onGetWebhook(getWebhook pubsub.GetWebhook) {
	repositoryName := getWebhook.Repository

	target, isExist := m.targetContexts[repositoryName]
	if isExist {
		target.canselFunc()
		delete(m.targetContexts, repositoryName)
	}

	path := strings.Split(getWebhook.Repository, "/")
	ctx, cancelFunc := context.WithCancel(context.Background())

	newContext := targetContext{
		repository: getWebhook.Repository,
		name:       path[len(path)-1],
		ctx:        ctx,
		canselFunc: cancelFunc,
	}

	m.targetContexts[repositoryName] = newContext
	go m.run(repositoryName, &newContext)
}

func (m *CDManager) run(repository string, target *targetContext) {
	path := strings.Split(repository, "/")
	repositoryName := path[len(path)-1]

	directoryPath := "./repository/" + repositoryName
	_, err := os.Stat(directoryPath)
	if err != nil {
		cmd := exec.Command("git", "clone", repository)
		cmd.Dir = "./repository"
		if _, err := cmd.Output(); err != nil {
			pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: "git clone failed"})
			return
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = directoryPath
		if _, err := cmd.Output(); err != nil {
			pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: "git pull failed"})
			return
		}
	}

	cmd := exec.Command("go", "build", "-o", "main")
	cmd.Dir = directoryPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.BUILD_FAILED, Message: string(out)})
		return
	}

	cmd = exec.Command("./main")
	cmd.Dir = directoryPath
	if err := cmd.Start(); err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.ERROR, Message: "Failed to exec command"})
		return
	}
	pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.APPLICATION_START})

	<-target.ctx.Done()
	pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.KILL_RECEIVED})
	if err := cmd.Process.Kill(); err != nil {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.KILL_FAILED})
	} else {
		pubsub.SystemEvent.Pub(pubsub.System{Time: time.Now(), Type: systemevent.KILL_SUCCESS})
	}
	return

}
