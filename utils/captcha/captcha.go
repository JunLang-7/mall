package captcha

import (
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

// NewSlideCaptcha 创建一个新的滑动验证码生成器
func NewSlideCaptcha() slide.Captcha {
	// 创建一个滑动验证码生成器，并配置为生成一张图的验证码
	builder := slide.NewBuilder(slide.WithGenGraphNumber(1))

	// 从资源中获取验证码图片
	imgs, err := imagesv2.GetImages()
	if err != nil {
		panic(err)
	}

	// 从资源中获取验证码缺口图片
	graphs, err := tiles.GetTiles()
	if err != nil {
		panic(err)
	}

	// 将获取到的验证码缺口图片转换为 slide.GraphImage 类型，并设置到生成器中
	newGraphs := make([]*slide.GraphImage, 0, len(graphs))
	for _, g := range graphs {
		newGraphs = append(newGraphs, &slide.GraphImage{
			MaskImage:    g.MaskImage,
			OverlayImage: g.OverlayImage,
			ShadowImage:  g.ShadowImage,
		})
	}

	// 将验证码图片和缺口图片设置到生成器中
	builder.SetResources(slide.WithGraphImages(newGraphs), slide.WithBackgrounds(imgs))
	return builder.Make()
}
