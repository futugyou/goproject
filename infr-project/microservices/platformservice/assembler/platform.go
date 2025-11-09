package assembler

import (
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/viewmodel"
)

type PlatformAssembler struct{}

func (a *PlatformAssembler) ToPlatformDetailView(d *domain.Platform) *viewmodel.PlatformDetailView {
	mapper := ProjectAssembler{}
	projects := mapper.ToViewModels(func() []domain.PlatformProject {
		if d.Projects == nil {
			return []domain.PlatformProject{}
		}
		result := []domain.PlatformProject{}
		for _, p := range d.Projects {
			result = append(result, p)
		}
		return result
	}())
	return &viewmodel.PlatformDetailView{
		ID:       d.ID,
		Name:     d.Name,
		Url:      d.Url,
		Provider: d.Provider.String(),
		Properties: func() []viewmodel.Property {
			res := []viewmodel.Property{}
			for k, v := range d.Properties {
				res = append(res, viewmodel.Property{
					Key:   k,
					Value: v.Value,
				})
			}
			return res
		}(),
		Tags: d.Tags,
		Secrets: func() []viewmodel.Secret {
			res := []viewmodel.Secret{}
			for _, v := range d.Secrets {
				res = append(res, viewmodel.Secret{
					Key:       v.Key,
					VaultID:   v.Value,
					VaultKey:  v.VaultKey,
					MaskValue: v.VaultMaskValue,
				})
			}
			return res
		}(),
		Activate:  d.Activate,
		IsDeleted: d.IsDeleted,
		Projects:  projects,
	}
}

func (a *PlatformAssembler) ToPlatformViews(platforms []domain.Platform) []viewmodel.PlatformView {
	result := []viewmodel.PlatformView{}
	for _, platform := range platforms {
		result = append(result, viewmodel.PlatformView{
			ID:        platform.ID,
			Name:      platform.Name,
			Activate:  platform.Activate,
			Url:       platform.Url,
			Tags:      platform.Tags,
			IsDeleted: platform.IsDeleted,
			Provider:  platform.Provider.String(),
		})
	}

	return result
}
