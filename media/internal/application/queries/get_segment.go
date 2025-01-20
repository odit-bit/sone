package queries

// type VideoSegmentArgs struct {
// 	Key     string
// 	Segment string
// }

// type GetSegmentHandler struct {
// 	repo domain.SegmentFetcher
// }

// func NewGetSegmentHandler(repo domain.SegmentFetcher) *GetSegmentHandler {
// 	return &GetSegmentHandler{
// 		repo: repo,
// 	}
// }

// func (h *GetSegmentHandler) GetSegment(ctx context.Context, query VideoSegmentArgs) (*domain.Segment, error) {
// 	seg, err := h.repo.GetVideoSegment(ctx, query.Key, query.Segment)
// 	if err != nil {
// 		return nil, errors.Join(err, fmt.Errorf("get segment query"))
// 	}
// 	return seg, nil
// }
