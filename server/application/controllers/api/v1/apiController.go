package v1

import (
	"server2/application/infrastructure"
	"server2/application/presentation"
)

type V1APIHandlers struct {
	nodeRepo           infrastructure.NodeRepository
	responseRepo       infrastructure.ResponsesReporisory
	fastErrorResponses presentation.FastErrorResponses
}

func NewV1Handlers(NodeRepo infrastructure.NodeRepository, ResponseRepo infrastructure.ResponsesReporisory) *V1APIHandlers {
	return &V1APIHandlers{nodeRepo: NodeRepo, responseRepo: ResponseRepo, fastErrorResponses: presentation.NewFastErrorResponses()}
}
