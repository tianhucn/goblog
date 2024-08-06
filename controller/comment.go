package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"kandaoni.com/anqicms/config"
	"kandaoni.com/anqicms/model"
	"kandaoni.com/anqicms/provider"
	"kandaoni.com/anqicms/request"
	"kandaoni.com/anqicms/response"
	"strings"
)

func CommentPublish(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	if !strings.HasPrefix(ctx.RequestPath(false), currentSite.BaseURI) {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "Not Found",
		})
		return
	}
	// 支持返回为 json 或html， 默认 html
	returnType := ctx.PostValueTrim("return")
	if ok := SafeVerify(ctx, nil, returnType, "comment"); !ok {
		return
	}

	//登录状态的用户，发布不进审核，否则进审核
	status := uint(0)
	userId := ctx.Values().GetIntDefault("userId", 0)
	adminId := ctx.Values().GetIntDefault("adminId", 0)
	if adminId > 0 {
		status = 1
	} else {
		// 是否需要审核
		var contentVerify = true
		userGroup := ctx.Values().Get("userGroup")
		if userGroup != nil {
			group, ok := userGroup.(*model.UserGroup)
			if ok {
				contentVerify = !group.Setting.ContentNoVerify
			}
		}
		if contentVerify == false {
			// 不需要审核
			status = 1
		}
	}

	var req request.PluginComment
	// 采用post接收
	req.ArchiveId = uint(ctx.PostValueIntDefault("archive_id", 0))
	req.UserName = ctx.PostValueTrim("user_name")
	req.Ip = ctx.PostValueTrim("ip")
	req.Content = ctx.PostValueTrim("content")
	req.ParentId = uint(ctx.PostValueIntDefault("parent_id", 0))
	req.ToUid = uint(ctx.PostValueIntDefault("to_uid", 0))
	if userId > 0 {
		userInfo := ctx.Values().Get("userInfo")
		if userInfo != nil {
			user, ok := userInfo.(*model.User)
			if ok {
				req.UserName = user.UserName
			}
		}
	}

	req.Status = status
	req.UserId = uint(userId)
	if req.Ip == "" {
		req.Ip = ctx.RemoteAddr()
	}
	if req.ParentId > 0 {
		parent, err := currentSite.GetCommentById(req.ParentId)
		if err == nil {
			req.ToUid = parent.UserId
		}
	}

	comment, err := currentSite.SaveComment(&req)
	if err != nil {
		msg := currentSite.TplTr("SaveFailed")
		if returnType == "json" {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  msg,
			})
		} else {
			ShowMessage(ctx, msg, nil)
		}
		return
	}

	msg := currentSite.TplTr("PublishSuccessfully")
	if returnType == "json" {
		ctx.JSON(iris.Map{
			"code": config.StatusOK,
			"msg":  msg,
			"data": comment,
		})
	} else {
		var link string
		refer := ctx.GetReferrer()
		if refer.URL != "" {
			link = refer.URL
		}
		ShowMessage(ctx, currentSite.TplTr("PublishSuccessfully"), []Button{
			{Name: currentSite.TplTr("ClickToContinue"), Link: link},
		})
	}
}

func CommentPraise(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	if !strings.HasPrefix(ctx.RequestPath(false), currentSite.BaseURI) {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "Not Found",
		})
		return
	}
	var req request.PluginComment
	req.Id = uint(ctx.PostValueIntDefault("id", 0))

	comment, err := currentSite.GetCommentById(req.Id)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	comment.VoteCount += 1
	err = comment.Save(currentSite.DB)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	comment.Active = true

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  currentSite.TplTr("LikeSuccessfully"),
		"data": comment,
	})
}

func CommentList(ctx iris.Context) {
	currentSite := provider.CurrentSite(ctx)
	if !strings.HasPrefix(ctx.RequestPath(false), currentSite.BaseURI) {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "Not Found",
		})
		return
	}

	archiveId := uint(ctx.Params().GetIntDefault("id", 0))
	archive, err := currentSite.GetArchiveById(archiveId)
	if err != nil {
		ShowMessage(ctx, "Not Found", nil)
		return
	}
	archive.Link = currentSite.GetUrl("archive", archive, 0)
	ctx.ViewData("archive", archive)
	if webInfo, ok := ctx.Value("webInfo").(*response.WebInfo); ok {
		currentPage := ctx.URLParamIntDefault("page", 1)
		webInfo.Title = currentSite.TplTr("CommentShow", archive.Title)
		if currentPage > 1 {
			webInfo.Title += " - " + currentSite.TplTr("PageNum", currentPage)
		}
		webInfo.Keywords = archive.Keywords
		webInfo.Description = archive.Description
		webInfo.PageName = "comments"
		webInfo.CanonicalUrl = currentSite.GetUrl(fmt.Sprintf("/comment/%d(?page={page})", archive.Id), nil, currentPage)
		ctx.ViewData("webInfo", webInfo)
	}

	ctx.ViewData("archiveId", archiveId)
	tplName := "comment/list.html"
	if ViewExists(ctx, "comment_list.html") {
		tplName = "comment_list.html"
	}
	err = ctx.View(GetViewPath(ctx, tplName))
	if err != nil {
		ctx.Values().Set("message", err.Error())
	}
}
