package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gogin/pkg/app"
	"gogin/pkg/e"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service/article_service"
	"gogin/service/tag_service"
	"net/http"
)

// GetArticle 获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)

		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// 获取多个文章

func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}
	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		ID:       tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	article, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
	}
	data := make(map[string]interface{})
	data["list"] = article
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 新增文章
type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"required;Min(1)"`
	Title         string `form:"title" valid:"required;MaxSize(100)"`
	Desc          string `form:"title" valid:"required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

func AddArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)
	httpCode, errCode := app.BindAndValid(c, form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	tagService := tag_service.Tag{
		ID: form.TagID,
	}
	exists, err := tagService.ExistByID()
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}
	if err = articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章
type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

func EditArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form EditArticleForm
	)
	httpCode, errCode := app.BindAndValid(c, form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exist, err := articleService.ExistByID()
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
	}
	tagService := tag_service.Tag{ID: form.TagID}
	exist, err = tagService.ExistByID()
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_TAG, nil)
	}

	if err = articleService.Edit(); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_EDIT_ARTICLE_FAIL, nil)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章

func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if err := articleService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
