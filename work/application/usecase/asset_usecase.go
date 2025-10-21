package usecase

import "idv/chris/MemoNest/domain/service"

type AssetUsecase struct {
	Img service.ImageProcessor
}

func (u *AssetUsecase) GetImageStoragePath(account, plain_text, name string) string {
	return u.Img.GetImageStoragePath(account, plain_text, name)
}

//-----------------------------------------------

func NewAssetUsecase(
	img service.ImageProcessor,
) *AssetUsecase {
	return &AssetUsecase{
		Img: img,
	}
}
