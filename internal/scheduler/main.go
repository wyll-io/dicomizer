package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
	"github.com/wyll-io/dicomizer/internal/check"
	"github.com/wyll-io/dicomizer/internal/storage"
)

func Create(
	ctx context.Context,
	cexpr string,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	pacs, aet, aec, aem, center string,
) (gocron.Scheduler, error) {
	location, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, err
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		return nil, err
	}

	if _, err := s.NewJob(
		gocron.CronJob(cexpr, false),
		gocron.NewTask(
			scheduleFunc,
			ctx, storageClient, dbClient, pacs, aet, aec, aem, center,
		),
	); err != nil {
		return nil, err
	}

	return s, nil
}

func scheduleFunc(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	pacs, aet, aec, aem, center string,
) {
	fmt.Println("Fetching patients info...")
	pInfos, err := dbClient.GetPatientsInfo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: to be parallelized
	for _, pInfo := range pInfos {
		fmt.Printf("Processing patient: \"%s\" (%s)\n", pInfo.Fullname, pInfo.PK)
		err := check.CheckPatientDCM(
			ctx,
			storageClient,
			dbClient,
			pacs,
			aet,
			aec,
			aem,
			center,
			pInfo,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Schedule done. Waiting for next run...")
}
