package commands

// type InsertSegmentCMD struct {
// 	Id   string
// 	Part string
// 	Body io.Reader
// }

// func NewInsertSegmentHandler(auth *kvstore.Client, repo domain.SegmentRepository) *InsertSegmentHandler {
// 	return &InsertSegmentHandler{
// 		repo: repo,
// 		auth: auth,
// 	}
// }

// type InsertSegmentHandler struct {
// 	repo domain.SegmentRepository
// 	auth *kvstore.Client
// }

// func (h *InsertSegmentHandler) InsertSegment(ctx context.Context, arg InsertSegmentCMD) error {
// 	if ok := h.auth.Exist(arg.Id); !ok {
// 		return fmt.Errorf("invalid stream key")
// 	}

// 	path := fmt.Sprintf("%s/%s", arg.Id, arg.Part)
// 	return h.repo.InsertVideoSegment(ctx, path, arg.Body)
// }
