package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"workflow-service/model"
)

func InitWorkflowRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/workflow/page", WorkflowPageListHandler)
	r.Get(common.BASE_CONTEXT+"/workflow", WorkflowListHandler)
	r.Post(common.BASE_CONTEXT+"/workflow", UpsertWorkflowHandler)
	r.Delete(common.BASE_CONTEXT+"/workflow/{id}", DeleteWorkflowHandler)
	r.Post(common.BASE_CONTEXT+"/workflow/batch-delete", batchDeleteWorkflowHandler)
	r.Post(common.BASE_CONTEXT+"/workflow/batch-upsert", batchUpsertWorkflowHandler)
	r.Get(common.BASE_CONTEXT+"/workflow/groupby", WorkflowGroupbyHandler)
}

// @Summary GroupBy
// @Description GroupBy, for example,  _select=level, then return  {level_val1:sum1,level_val2:sum2}, _where can input status=0
// @Tags Workflow
// @Param _select query string true "_select"
// @Param _where query string false "_where"
// @Produce  json
// @Success 200 {object} common.Response{data=[]map[string]any} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /workflow/groupby [get]
func WorkflowGroupbyHandler(w http.ResponseWriter, r *http.Request) {

	common.CommonGroupby(w, r, common.GetDaprClient(), "o_workflow")
}

// @Summary batch update
// @Description batch update
// @Tags Workflow
// @Accept  json
// @Param entities body []map[string]any true "objects array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /workflow/batch-upsert [post]
func batchUpsertWorkflowHandler(w http.ResponseWriter, r *http.Request) {

	var entities []map[string]any
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam)
		return
	}

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.WorkflowTableInfo.Name, model.Workflow_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary page query
// @Description page query, _page(from 1 begin), _page_size, _order, and others fields, status=1, name=$like.%CAM%
// @Tags Workflow
// @Param _page query int true "current page"
// @Param _page_size query int true "page size"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param name query string false "name"
// @Param status query string false "status"
// @Param content query string false "content"
// @Param description query string false "description"
// @Param type query string false "type"
// @Param workflow_id query string false "workflow_id"
// @Param running_id query string false "running_id"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Workflow}} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /workflow/page [get]
func WorkflowPageListHandler(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam)
		return
	}
	common.CommonPageQuery[model.Workflow](w, r, common.GetDaprClient(), "o_workflow", "id")

}

// @Summary query objects
// @Description query objects
// @Tags Workflow
// @Param _select query string false "_select"
// @Param _order query string false "order"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param name query string false "name"
// @Param status query string false "status"
// @Param content query string false "content"
// @Param description query string false "description"
// @Param type query string false "type"
// @Param workflow_id query string false "workflow_id"
// @Param running_id query string false "running_id"
// @Produce  json
// @Success 200 {object} common.Response{data=[]model.Workflow} "objects array"
// @Failure 500 {object} common.Response ""
// @Router /workflow [get]
func WorkflowListHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonQuery[model.Workflow](w, r, common.GetDaprClient(), "o_workflow", "id")
}

// @Summary save
// @Description save
// @Tags Workflow
// @Accept       json
// @Param item body model.Workflow true "object"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Workflow} "object"
// @Failure 500 {object} common.Response ""
// @Router /workflow [post]
func UpsertWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	var val model.Workflow
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	beforeHook, exists := common.GetUpsertBeforeHook("Workflow")
	if exists {
		v, err1 := beforeHook(r, val)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
		val = v.(model.Workflow)
	}

	err = common.DbUpsert[model.Workflow](r.Context(), common.GetDaprClient(), val, model.WorkflowTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary delete
// @Description delete
// @Tags Workflow
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Workflow} "object"
// @Failure 500 {object} common.Response ""
// @Router /workflow/{id} [delete]
func DeleteWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Workflow")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "o_workflow", "id", "id")
}

// @Summary batch delete
// @Description batch delete
// @Tags Workflow
// @Accept  json
// @Param ids body []string true "id array"
// @Produce  json
// @Success 200 {object} common.Response ""
// @Failure 500 {object} common.Response ""
// @Router /workflow/batch-delete [post]
func batchDeleteWorkflowHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam)
		return
	}
	beforeHook, exists := common.GetBatchDeleteBeforeHook("Workflow")
	if exists {
		_, err1 := beforeHook(r, ids)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "o_workflow", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
