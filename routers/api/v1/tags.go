package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"gogin/pkg/app"
	"gogin/pkg/e"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"gogin/service/tag_service"
	"net/http"
)

// GetTags 获取多个文章标签
// @Summary 获取全部标签
// @Produce json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	//Name      string `form:"name" valid:"Required;MaxSize(100)"`
	//CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	//State     int    `form:"state" valid:"Range(0,1)"`
	Name      string `form:"name"`
	CreatedBy string `form:"created_by"`
	State     int    `form:"state"`
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	fmt.Printf("Name:%v,CreatedBy:%v", tagService.Name, tagService.CreatedBy)
	valid := validation.Validation{}
	valid.Required(tagService.Name, "name").Message("姓名不能为空")
	valid.Required(tagService.CreatedBy, "created_by").Message("创建人不能为空")
	valid.Range(tagService.State, 0, 1, "state").Message("状态只能为0或者1")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"required;Min(1)"`
	Name       string `from:"name" valid:"required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// EditTag 修改文章标签
// @Summary 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form EditTagForm
	)
	httpCode, errCode := app.BindAndValid(c, form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
	}
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists, err := tagService.ExistByID()
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if err = tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章标签

func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
