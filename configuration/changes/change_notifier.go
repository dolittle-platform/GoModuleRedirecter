package changes

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type ComponentName string

type ConfigurationChangeCallback func() error

type ConfigurationChangeNotifier interface {
	RegisterCallback(component ComponentName, callback ConfigurationChangeCallback) error
	TriggerChanged()
	Shutdown()
}

func NewConfigurationChangeNotifier(logger *zap.Logger) ConfigurationChangeNotifier {
	return &notifier{
		lock:              sync.Mutex{},
		channels:          make(map[ComponentName]chan<- struct{}),
		shutdownInitiated: false,
		shutdownCompleted: sync.WaitGroup{},
		logger:            logger,
	}
}

type notifier struct {
	lock              sync.Mutex
	channels          map[ComponentName]chan<- struct{}
	shutdownInitiated bool
	shutdownCompleted sync.WaitGroup
	logger            *zap.Logger
}

func (n *notifier) RegisterCallback(component ComponentName, callback ConfigurationChangeCallback) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, componentAlreadyRegistered := n.channels[component]; componentAlreadyRegistered {
		return errors.New("component already registered")
	}

	ch := make(chan struct{}, 1)
	n.channels[component] = ch
	go n.runConfigurationChangeCaller(component, ch, callback)

	return nil
}

func (n *notifier) TriggerChanged() {
	for _, ch := range n.channels {
		ch <- struct{}{}
	}
}

func (n *notifier) Shutdown() {
	n.shutdownInitiated = true
	n.TriggerChanged()
	n.shutdownCompleted.Wait()
}

func (n *notifier) runConfigurationChangeCaller(component ComponentName, notification <-chan struct{}, callback ConfigurationChangeCallback) {
	n.shutdownCompleted.Add(1)
	defer n.shutdownCompleted.Done()

	for {
		_ = <-notification

		if n.shutdownInitiated {
			return
		}

		n.logger.Info("Reloading configuration", zap.String("component", string(component)))

		if err := callback(); err != nil {
			n.logger.Warn("Failed to reload configuration", zap.String("component", string(component)), zap.Error(err))
		}
	}
}
