package controller

import (
	"net/http"
	"strconv"
	"sync"

	coreController "github.com/CakeForKit/rsoi-lab1/internal/common/controller"
	"github.com/CakeForKit/rsoi-lab1/internal/common/custom_error"
	"github.com/CakeForKit/rsoi-lab1/internal/common/model"
	"github.com/CakeForKit/rsoi-lab1/internal/common/resolver"
	"github.com/CakeForKit/rsoi-lab1/internal/common/utils"
	"github.com/CakeForKit/rsoi-lab1/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var personCntrlr coreController.HttpController
var personCntrlrMutex sync.Mutex

type personController struct {
	service service.PersonService
}

func GetPersonController() (coreController.HttpController, error) {
	personCntrlrMutex.Lock()
	defer personCntrlrMutex.Unlock()

	if personCntrlr != nil {
		return personCntrlr, nil
	}

	personService, err := service.GetPersonService()
	if err != nil {
		return nil, err
	}

	personCntrlr = &personController{
		service: personService,
	}
	return personCntrlr, nil
}

func (controller *personController) RegisterHttpController(router *gin.Engine) {
	personRouter := router.Group("/api/v1/persons")

	personRouter.GET(":personId",
		controller.getById,
	)

	personRouter.GET("",
		controller.getAll,
	)

	personRouter.POST("",
		resolver.Resolver[model.Person],
		controller.create,
	)

	personRouter.PATCH("/:personId",
		resolver.Resolver[model.PersonUpdate],
		controller.update,
	)

	personRouter.DELETE("/:personId",
		controller.deleteById,
	)
}

func (controller *personController) getById(ctx *gin.Context) {
	id, ok := personID(ctx)
	if !ok {
		return
	}
	log.WithContext(ctx).Infof("UserController: GetById(id: %d): Start", id)
	if result, err := controller.service.GetById(ctx, id); err != nil {
		log.WithError(err).WithContext(ctx).Error("UserController: GetById(): Failed")
		utils.AbortContextWithError(ctx, err)
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, result)
	}
	log.WithContext(ctx).Info("UserController: GetById(): End")
}

func (controller *personController) getAll(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: GetAll(): Start")
	if result, err := controller.service.GetAll(ctx); err != nil {
		log.WithError(err).WithContext(ctx).Error("UserController: GetAll(): Failed")
		utils.AbortContextWithError(ctx, err)
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, result)
	}
	log.WithContext(ctx).Info("UserController: GetAll(): End")
}

func (controller *personController) create(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: Create(): Start")
	user := utils.MustGetRequestBody[model.Person](ctx)
	if user.Name == "" {
		utils.AbortContextWithError(ctx, custom_error.IllegalArgumentError("name is required"))
		return
	}
	if result, err := controller.service.Create(ctx, user); err != nil {
		log.WithError(err).WithContext(ctx).Error("UserController: Create(): Failed")
		utils.AbortContextWithError(ctx, err)
	} else {
		ctx.Header("Location", "/api/v1/persons/"+strconv.FormatUint(uint64(result.ID), 10))
		ctx.Status(http.StatusCreated)
	}
	log.WithContext(ctx).Info("UserController: Create(): End")
}

func (controller *personController) update(ctx *gin.Context) {
	log.WithContext(ctx).Info("UserController: Update(): Start")

	id, ok := personID(ctx)
	if !ok {
		return
	}
	user := utils.MustGetRequestBody[model.PersonUpdate](ctx)

	if user.Name == nil && user.Age == nil && user.Address == nil && user.Work == nil {
		utils.AbortContextWithError(ctx, custom_error.IllegalArgumentError("request body is empty"))
		return
	}

	if result, err := controller.service.Update(ctx, id, user); err != nil {
		log.WithError(err).WithContext(ctx).Error("UserController: Update(): Failed")
		utils.AbortContextWithError(ctx, err)
	} else {
		ctx.AbortWithStatusJSON(http.StatusOK, result)
	}
	log.WithContext(ctx).Info("UserController: Update(): End")
}

func (controller *personController) deleteById(ctx *gin.Context) {
	id, ok := personID(ctx)
	if !ok {
		return
	}
	log.WithContext(ctx).Infof("UserController: DeleteById(id: %d): Start", id)
	if err := controller.service.DeleteById(ctx, id); err != nil {
		log.WithError(err).WithContext(ctx).Error("UserController: DeleteById(): Failed")
		utils.AbortContextWithError(ctx, err)
	} else {
		ctx.Status(http.StatusNoContent)
	}
	log.WithContext(ctx).Info("UserController: DeleteById(): End")
}

func personID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("personId"), 10, 0)
	if err != nil || id == 0 {
		utils.AbortContextWithError(ctx, custom_error.IllegalArgumentError("personId must be a positive integer"))
		return 0, false
	}
	return uint(id), true
}
