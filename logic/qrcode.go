package logic

import (
    "fmt"
    "github.com/nfnt/resize"
    "github.com/skip2/go-qrcode"
    "image"
    "image/draw"
    "image/jpeg"
    "image/png"
    "os"
    "path"
    "path/filepath"
    qrcodeerr "qrcode_gene/error"
    "qrcode_gene/util"
    "strings"
)

type Filter func(content string) error

var (
    bannedWords = []string{"赌博", "发牌", "色色", "毒", "约", "女优", "色情",
        "算命", "算卦", "神仙", "鬼怪", "魔鬼", "地狱",
        "回回", "鞑子", "高丽棒子", "老毛子", "血统", "杂种", "东亚病夫", "蛮夷", "男尊女卑", "大汉族主义"}
    defaultFilter Filter = func(content string) error {
        for _, str := range bannedWords {
            contains := strings.Contains(content, str)
            if contains {
                return qrcodeerr.ErrCantGene
            }
        }
        return nil
    }
)

func GeneQrCode(content string, iconPath string, filter Filter) (image.Image, error) {
    if content == "" {
        return nil, qrcodeerr.ErrInAvailData
    }

    var err error
    if filter != nil {
        if err = filter(content); err != nil {
            return nil, err
        }
    }

    if err = defaultFilter(content); err != nil {
        return nil, err
    }

    var img image.Image
    if iconPath == "" {
        img = getQr(content)
    } else {
        ext := filepath.Ext(iconPath)
        iconFile, _ := os.Open(iconPath)
        var iconImg image.Image
        switch ext {
        case ".jpg":
            iconImg, _ = jpeg.Decode(iconFile)
        case ".jpeg":
            iconImg, _ = jpeg.Decode(iconFile)
        case ".png":
            iconImg, _ = png.Decode(iconFile)
        default:
            return nil, qrcodeerr.ErrInvalidIconFormat
        }
        img = getQrWithIcon(content, iconImg)
    }

    if img == nil {
        return nil, qrcodeerr.ErrGeneQrCode
    } else {
        return img, nil
    }
}

func getQr(content string) image.Image {
    qrCode, err := qrcode.New(content, qrcode.Highest)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    qrCode.DisableBorder = true
    return qrCode.Image(256)
}

func getQrWithIcon(content string, iconImg image.Image) image.Image {
    var (
        bgImg  image.Image
        offset image.Point
    )
    qrCode, err := qrcode.New(content, qrcode.Highest)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    qrCode.DisableBorder = true
    bgImg = qrCode.Image(256)
    iconImg = resize.Resize(40, 40, iconImg, resize.Lanczos3)

    b := bgImg.Bounds()
    //居中设置icon到二维码图片
    offset = image.Pt((b.Max.X-iconImg.Bounds().Max.X)/2, (b.Max.Y-iconImg.Bounds().Max.Y)/2)
    m := image.NewRGBA(b)
    draw.Draw(m, b, bgImg, image.Point{X: 0, Y: 0}, draw.Src)
    draw.Draw(m, iconImg.Bounds().Add(offset), iconImg, image.Point{X: 0, Y: 0}, draw.Over)
    return m.SubImage(m.Bounds())
}

func SaveImage(img image.Image) error {
    fileName := geneAvailFileName("生成二维码.jpg")
    f, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer f.Close()
    err = jpeg.Encode(f, img, nil)
    if err != nil {
        return err
    }
    return nil
}

func geneAvailFileName(fileName string) string {
    fileSuffix := path.Ext(fileName)
    filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
    for i := 1; util.FileOrPathIsExists(filepath.Join(fileName)); i++ {
        fileName = filenameOnly + fmt.Sprintf("(%d)", i) + fileSuffix
    }
    return fileName
}
