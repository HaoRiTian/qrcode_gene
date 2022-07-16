//go:generate go-winres make
package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "github.com/flopp/go-findfont"
    "os"
    qrcodeerr "qrcode_gene/error"
    "qrcode_gene/logic"
)

func setChinese() {
    if simhei, err := findfont.Find("simhei.ttf"); err != nil {
        return
    } else {
        os.Setenv("FYNE_FONT", simhei)
    }
}

func main() {
    setChinese()

    a := app.New()
    w := a.NewWindow("二维码生成")
    w.SetIcon(resourceHeadIconPng)
    iconPath := ""
    dl := dialog.NewFileOpen(func(ur fyne.URIReadCloser, err error) {
        if ur != nil && ur.URI() != nil {
            iconPath = ur.URI().Path()
        }
    }, w)
    input := widget.NewEntry()
    input.SetPlaceHolder("输入要生成的信息")

    selectIcon := widget.NewButtonWithIcon("选择Icon", theme.FileImageIcon(), func() {
        dl.Show()
    })
    inputBox := container.New(layout.NewVBoxLayout(), input, selectIcon)
    imgBox := container.New(layout.NewMaxLayout())

    b1 := widget.NewButtonWithIcon("生成", theme.DocumentCreateIcon(), func() {
        qrCodeImg, err := logic.GeneQrCode(input.Text, iconPath, nil)
        if err != nil {
            errDl := dialog.NewError(err, w)
            errDl.Show()
            return
        }
        img := canvas.NewImageFromImage(qrCodeImg)
        img.FillMode = canvas.ImageFillOriginal
        imgBox.RemoveAll()
        imgBox.Add(img)
    })
    b2 := widget.NewButtonWithIcon("保存", theme.DocumentSaveIcon(), func() {
        if imgBox.Objects != nil && imgBox.Objects[0] != nil {
            if err := logic.SaveImage(imgBox.Objects[0].(*canvas.Image).Image); err == nil {
                return
            }
        }
        errDl := dialog.NewError(qrcodeerr.ErrSaveQrCode, w)
        errDl.Show()
    })
    buttonBox := container.New(layout.NewHBoxLayout(), b1, b2)

    w.SetContent(container.New(layout.NewVBoxLayout(), inputBox, imgBox, buttonBox))
    w.Resize(fyne.NewSize(400, 600))
    w.ShowAndRun()

    os.Unsetenv("FYNE_FONT")
}
