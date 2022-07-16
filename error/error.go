package qrcodeerr

import "errors"

var (
    ErrGeneQrCode        = errors.New("二维码生成失败")
    ErrInAvailData       = errors.New("无效信息")
    ErrInvalidIconFormat = errors.New("icon格式无效")
    ErrSaveQrCode        = errors.New("保存二维码失败")
    ErrCantGene          = errors.New("信息不能被生成")
)
