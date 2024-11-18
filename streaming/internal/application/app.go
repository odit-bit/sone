package application

import (
	"context"

	"github.com/odit-bit/sone/streaming/internal/application/commands"
	"github.com/odit-bit/sone/streaming/internal/application/queries"
	"github.com/odit-bit/sone/streaming/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		CreateStreamKey(ctx context.Context, cmd commands.CreateStreamKeyCMD) (domain.Key, error)
	}

	Queries interface {
		GetSegment(ctx context.Context, query queries.GetSegment) (*domain.Segment, error)
		GetEntries(ctx context.Context, query queries.GetEntry) (<-chan domain.Entry, error)
	}

	commandHandler struct {
		*commands.CreateStreamKeyHandler
	}

	queriesHandler struct {
		*queries.GetSegmentHandler
		*queries.GetEntryHandler
	}
)

var _ App = (*Application)(nil)

type Application struct {
	commandHandler
	queriesHandler
}

func New(segmentrepo domain.SegmentRepository, entryRepo domain.EntryRepository, streamkeyRepo domain.StreamKeyRepository) Application {
	return Application{
		commandHandler: commandHandler{
			CreateStreamKeyHandler: commands.NewCreateStreamKeyHandler(streamkeyRepo),
		},
		queriesHandler: queriesHandler{
			GetSegmentHandler: queries.NewGetSegmentHandler(segmentrepo),
			GetEntryHandler:   queries.NewGetEntryHandler(entryRepo),
		},
	}
}
