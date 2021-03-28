package modules

import (
	"context"
	"html/template"
	"io"
	"path/filepath"
	"redirecter/configuration/changes"
	"redirecter/server/correlation"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type Writer interface {
	WriteGoGetResponse(writer io.Writer, url string, repository *Repository, ctx context.Context) error
	WriteUserResponse(writer io.Writer, url string, repository *Repository, ctx context.Context) error
	WriteNotFoundResponse(writer io.Writer, url string, ctx context.Context) error
}

func NewWriter(configuration Configuration, notifier changes.ConfigurationChangeNotifier, logger *zap.Logger) Writer {
	watcher, _ := fsnotify.NewWatcher()
	w := &writer{
		configuration: configuration,
		watcher:       watcher,
		logger:        logger,
	}
	go w.watcherLoop()
	w.updateWatcherFiles()
	notifier.RegisterCallback("writer", w.updateWatcherFiles)
	return w
}

type writer struct {
	configuration    Configuration
	watcher          *fsnotify.Watcher
	fileGoGet        string
	fileUser         string
	fileNotFound     string
	templateGoGet    *template.Template
	templateUser     *template.Template
	templateNotFound *template.Template
	logger           *zap.Logger
}

func (w *writer) updateWatcherFiles() error {
	pathGoGet := w.configuration.TemplateGoGetPath()
	w.logger.Info("Using go get template", zap.String("path", pathGoGet))
	w.updateWatchFile(w.fileGoGet, pathGoGet)
	w.fileGoGet = pathGoGet
	w.parseTemplateFile(pathGoGet)

	pathUser := w.configuration.TemplateUserPath()
	w.logger.Info("Using user template", zap.String("path", pathUser))
	w.updateWatchFile(w.fileUser, pathUser)
	w.fileUser = pathUser
	w.parseTemplateFile(pathUser)

	pathNotFound := w.configuration.TemplateNotFoundPath()
	w.logger.Info("Using not found template", zap.String("path", pathNotFound))
	w.updateWatchFile(w.fileNotFound, pathNotFound)
	w.fileNotFound = pathNotFound
	w.parseTemplateFile(pathNotFound)

	return nil
}

func (w *writer) updateWatchFile(previous, next string) {
	if previous != "" {
		if err := w.watcher.Remove(previous); err != nil {
			w.logger.Warn("Failed to stop watching template file", zap.Error(err), zap.String("path", previous))
		}
	}
	if err := w.watcher.Add(next); err != nil {
		w.logger.Warn("Failed to start watching template file", zap.Error(err), zap.String("path", next))
	}
}

func (w *writer) watcherLoop() {
LOOP:
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				break LOOP
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.parseTemplateFile(event.Name)
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				break LOOP
			}
			w.logger.Error("Error occured in template file watcher", zap.Error(err))
		}
	}
	w.logger.Error("File watcher loop stopped")
}

func (w *writer) parseTemplateFile(path string) {
	template, err := template.ParseFiles(path)
	if err != nil {
		w.logger.Error("Failed to parse template file", zap.Error(err), zap.String("path", path))
	}

	switch {
	case w.comparePaths(w.fileGoGet, path):
		w.logger.Info("Parsed go get template", zap.String("path", path), zap.Bool("ok", template != nil))
		w.templateGoGet = template
	case w.comparePaths(w.fileUser, path):
		w.logger.Info("Parsed user template", zap.String("path", path), zap.Bool("ok", template != nil))
		w.templateUser = template
	case w.comparePaths(w.fileNotFound, path):
		w.logger.Info("Parsed not found template", zap.String("path", path), zap.Bool("ok", template != nil))
		w.templateNotFound = template
	default:
		w.logger.Warn("Parsing template that is not used", zap.String("path", path))
	}
}

func (w *writer) comparePaths(left, right string) bool {
	leftPath, err := filepath.Abs(left)
	if err != nil {
		return false
	}
	rightPath, err := filepath.Abs(right)
	if err != nil {
		return false
	}
	return leftPath == rightPath
}

type foundTemplateData struct {
	Module        string
	Documentation string
	*Repository
}

func (w *writer) WriteGoGetResponse(writer io.Writer, url string, repository *Repository, ctx context.Context) error {
	w.logger.Info("Writing go get response", zap.String("url", url), zap.String("repository", repository.Source), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
	return w.templateGoGet.Execute(writer, foundTemplateData{
		Module:        url,
		Documentation: w.configuration.Documentation() + url,
		Repository:    repository,
	})
}

func (w *writer) WriteUserResponse(writer io.Writer, url string, repository *Repository, ctx context.Context) error {
	w.logger.Info("Writing user response", zap.String("url", url), zap.String("repository", repository.Source), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
	return w.templateUser.Execute(writer, foundTemplateData{
		Module:        url,
		Documentation: w.configuration.Documentation() + url,
		Repository:    repository,
	})
}

type notFoundTemplateData struct {
	Module string
}

func (w *writer) WriteNotFoundResponse(writer io.Writer, url string, ctx context.Context) error {
	w.logger.Info("Writing not found response", zap.String("url", url), zap.String("correlation", correlation.CorrelationFromContext(ctx)))
	return w.templateNotFound.Execute(writer, notFoundTemplateData{
		Module: url,
	})
}
