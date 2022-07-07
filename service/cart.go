package service

import (
	logging "github.com/sirupsen/logrus"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	"mall/serializer"
	"strconv"
)

//创建购物车
type CreateService struct {
	BossID uint `form:"boss_id" json:"boss_id"`
	Num    uint `form:"num" json:"num"`
}

func (service *CreateService) Create(id string, uid uint) serializer.Response {
	var product model.Product
	code := e.SUCCESS
	err := dao.DB.First(&product, id).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if product == (model.Product{}) {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	idInt, _ := strconv.Atoi(id)
	var cart model.Cart
	dao.DB.Where("user_id=? AND product_id=? AND boss_id=?", uid, id, product.BossID).Find(&cart)
	if cart == (model.Cart{}) {
		cart = model.Cart{
			UserID:    uid,
			ProductID: uint(idInt),
			BossID:    uint(product.BossID),
			Num:       1,
			MaxNum:    10,
			Check:     false,
		}
		err = dao.DB.Create(&cart).Error
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		return serializer.Response{
			Status: code,
			Data:   serializer.BuildCart(cart, product, service.BossID),
			Msg:    e.GetMsg(code),
		}
	} else if cart.Num < cart.MaxNum {
		cart.Num++
		err = dao.DB.Save(&cart).Error
		if err != nil {
			logging.Info(err)
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		return serializer.Response{
			Status: 201,
			Msg:    "商品已经在购物车了，数量+1",
			Data:   serializer.BuildCart(cart, product, service.BossID),
		}
	} else {
		return serializer.Response{
			Status: 202,
			Msg:    "超过最大上限",
		}
	}
}

//Show 订单
func (service *CreateService) Show(id string) serializer.Response {
	var carts []model.Cart
	code := e.SUCCESS
	err := dao.DB.Where("user_id=?", id).Find(&carts).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildCarts(carts),
		Msg:    e.GetMsg(code),
	}
}

//修改购物车信息
func (service *CreateService) Update(id string) serializer.Response {
	var cart model.Cart
	code := e.SUCCESS
	err := dao.DB.Where("id=?", id).Find(&cart).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	cart.Num = service.Num
	err = dao.DB.Save(&cart).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

//删除购物车
func (service *CreateService) Delete(pid string, uid uint) serializer.Response {
	var cart model.Cart
	code := e.SUCCESS
	err := dao.DB.Where("user_id=? AND product_id=?", uid, pid).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = dao.DB.Delete(&cart).Error
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
