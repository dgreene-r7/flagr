package handler

import (
	"github.com/checkr/flagr/pkg/config"
	"github.com/checkr/flagr/pkg/entity"
	"github.com/checkr/flagr/swagger_gen/models"
	"github.com/checkr/flagr/swagger_gen/restapi/operations"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/constraint"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/distribution"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/evaluation"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/export"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/flag"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/health"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/segment"
	"github.com/checkr/flagr/swagger_gen/restapi/operations/variant"
	"github.com/go-openapi/runtime/middleware"
)

var getDB = entity.GetDB

// Setup initialize all the handler functions
func Setup(api *operations.FlagrAPI) {
	if config.Config.EvalOnlyMode {
		setupCRUD(api, true)
		setupHealth(api)
		setupEvaluation(api)
		return
	}

	setupHealth(api)
	setupEvaluation(api)
	setupCRUD(api, false)
	setupExport(api)
}

func setupCRUD(api *operations.FlagrAPI, readOnly bool) {
	c := NewCRUD()

	// flags
	api.FlagFindFlagsHandler = flag.FindFlagsHandlerFunc(c.FindFlags)
	api.FlagGetFlagHandler = flag.GetFlagHandlerFunc(c.GetFlag)
	api.FlagGetFlagSnapshotsHandler = flag.GetFlagSnapshotsHandlerFunc(c.GetFlagSnapshots)
	api.FlagGetFlagEntityTypesHandler = flag.GetFlagEntityTypesHandlerFunc(c.GetFlagEntityTypes)

	// segments
	api.SegmentFindSegmentsHandler = segment.FindSegmentsHandlerFunc(c.FindSegments)

	// constraints
	api.ConstraintFindConstraintsHandler = constraint.FindConstraintsHandlerFunc(c.FindConstraints)

	// distributions
	api.DistributionFindDistributionsHandler = distribution.FindDistributionsHandlerFunc(c.FindDistributions)

	// variants
	api.VariantFindVariantsHandler = variant.FindVariantsHandlerFunc(c.FindVariants)

	// Return early if we only get read-only CRUD operations
	if readOnly {
		return
	}

	// flags
	api.FlagCreateFlagHandler = flag.CreateFlagHandlerFunc(c.CreateFlag)
	api.FlagPutFlagHandler = flag.PutFlagHandlerFunc(c.PutFlag)
	api.FlagDeleteFlagHandler = flag.DeleteFlagHandlerFunc(c.DeleteFlag)
	api.FlagSetFlagEnabledHandler = flag.SetFlagEnabledHandlerFunc(c.SetFlagEnabledState)

	// segments
	api.SegmentCreateSegmentHandler = segment.CreateSegmentHandlerFunc(c.CreateSegment)
	api.SegmentPutSegmentHandler = segment.PutSegmentHandlerFunc(c.PutSegment)
	api.SegmentDeleteSegmentHandler = segment.DeleteSegmentHandlerFunc(c.DeleteSegment)
	api.SegmentPutSegmentsReorderHandler = segment.PutSegmentsReorderHandlerFunc(c.PutSegmentsReorder)

	// constraints
	api.ConstraintCreateConstraintHandler = constraint.CreateConstraintHandlerFunc(c.CreateConstraint)
	api.ConstraintPutConstraintHandler = constraint.PutConstraintHandlerFunc(c.PutConstraint)
	api.ConstraintDeleteConstraintHandler = constraint.DeleteConstraintHandlerFunc(c.DeleteConstraint)

	// distributions
	api.DistributionPutDistributionsHandler = distribution.PutDistributionsHandlerFunc(c.PutDistributions)

	// variants
	api.VariantCreateVariantHandler = variant.CreateVariantHandlerFunc(c.CreateVariant)
	api.VariantPutVariantHandler = variant.PutVariantHandlerFunc(c.PutVariant)
	api.VariantDeleteVariantHandler = variant.DeleteVariantHandlerFunc(c.DeleteVariant)
}

func setupEvaluation(api *operations.FlagrAPI) {
	ec := GetEvalCache()
	ec.Start()

	e := NewEval()
	api.EvaluationPostEvaluationHandler = evaluation.PostEvaluationHandlerFunc(e.PostEvaluation)
	api.EvaluationPostEvaluationBatchHandler = evaluation.PostEvaluationBatchHandlerFunc(e.PostEvaluationBatch)

	if config.Config.RecorderEnabled {
		// Try GetDataRecorder to catch fatal errors before we start the evaluation api
		GetDataRecorder()
	}
}

func setupHealth(api *operations.FlagrAPI) {
	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(
		func(health.GetHealthParams) middleware.Responder {
			return health.NewGetHealthOK().WithPayload(&models.Health{Status: "OK"})
		},
	)
}

func setupExport(api *operations.FlagrAPI) {
	api.ExportGetExportSqliteHandler = export.GetExportSqliteHandlerFunc(exportSQLiteHandler)
	api.ExportGetExportEvalCacheJSONHandler = export.GetExportEvalCacheJSONHandlerFunc(exportEvalCacheJSONHandler)
}
