package resolvers

import (
	"context"

	api "github.com/photoview/photoview/api/graphql"
	"github.com/photoview/photoview/api/graphql/models"
)

func (r *queryResolver) SiteInfo(ctx context.Context) (*models.SiteInfo, error) {
	return models.GetSiteInfo(r.DB(ctx))
}

type SiteInfoResolver struct {
	*Resolver
}

func (r *Resolver) SiteInfo() api.SiteInfoResolver {
	return &SiteInfoResolver{r}
}

func (SiteInfoResolver) FaceDetectionEnabled(ctx context.Context, obj *models.SiteInfo) (bool, error) {
	return false, nil
}
