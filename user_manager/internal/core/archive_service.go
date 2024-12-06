package core

type ArchiveService interface {
	StartArchiveLoop()
	Close() error
}
