package facade

import (
	"context"
	"fmt"
	"github.com/calebtracey/ai-interaction-api/external"
	"github.com/calebtracey/ai-interaction-api/internal/dao/openai"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"math"
	"net/http"
	"strconv"
)

const (
	LargeImage  = "1024x1024"
	MediumImage = "512x512"
	SmallImage  = "256x256"

	ChunkSize = 5
)

type ServiceI interface {
	GenerateImage(ctx context.Context, apiRequest *external.APIRequest) (apiResp *external.APIResponse)
}

type Service struct {
	DAO openai.DAOI
}

func (s Service) GenerateImage(ctx context.Context, apiRequest *external.APIRequest) *external.APIResponse {
	apiResp := new(external.APIResponse)
	//mu := new(sync.Mutex)

	//mu.Lock()
	apiRequest.Size = LargeImage
	//mu.Unlock()

	amount := apiRequest.N

	log.Infof("apiRequest amount: %d", amount)
	g, ctx := errgroup.WithContext(ctx)

	imageChan := make(chan external.GenImages, amount)
	defer close(imageChan)

	// break up request according to ChunkSize and make requests concurrently
	log.Infof("Chunking Request...")
	for i := 0; i < amount; i += ChunkSize {
		amount = requestAmount(i, amount)
		log.Infof("Current amount: %d", amount)
		g.Go(func() error {

			//defer func() {
			//	// last one out closes shop
			//	if atomic.AddInt32(&workers, -1) == 0 {
			//		log.Infoln("last request, closing channel...")
			//		close(resultChan)
			//	}
			//}()

			// check if the current request doesn't require the full image count
			// if not, set the request size to the remainder
			log.Infoln("==^ getting images...")
			if _, err := s.getImages(ctx, apiRequest, imageChan); err != nil {
				return err
			}

			log.Infoln("==^ goroutine returning...")
			return nil
		})
	}

	g.Go(func() error {
		defer close(imageChan)
		idx := 0
		for images := range imageChan {
			log.Infof("reading image channel #%d\n", idx)
			apiResp.Result.Data = append(apiResp.Result.Data, images...)
			//atomic.AddInt32(&created, int32(res.Created))
			//imageResponses[idx] = &res.Data
			idx++
		}
		return nil
	})

	// wait for go routines to finish can catch any errors
	if err := g.Wait(); err != nil {
		return responseWithError(apiResp, err, http.StatusInternalServerError, "ImageRequest")
	}

	// TODO move mapping
	// TODO created seems broken at the source right now
	//totalCreated := int(created)
	//apiResp.Result.Created = totalCreated
	//apiResp.Message.Count = totalCreated
	log.Infof("Final OpenAI result: %v", apiResp.Result)
	log.Infoln("returning final APIResponse...")

	return apiResp
}

func requestAmount(i, amount int) int {
	end := i + ChunkSize
	if end > amount {
		//todo this is all not right
		return end - amount
	} else {
		return ChunkSize
	}
}

func (s Service) getImages(ctx context.Context, apiRequest *external.APIRequest, imageChan chan<- external.GenImages) (resp external.AIResponse, err error) {
	if resp, err = s.DAO.ImageRequest(ctx, apiRequest); err == nil {
		select {
		case <-ctx.Done():
			log.Errorf("=== getImages: context cancelled")
			return resp, ctx.Err()
		case imageChan <- resp.Data:
			log.Infoln("===^ resultChan received result")
			//	log.Infoln("=== default select")
		}
		return resp, nil

	} else {
		log.Error(err)
		return resp, err
	}
}

func remainder(totalRequests, requestAmount int) int {
	r := math.Remainder(float64(totalRequests), float64(requestAmount))
	rCeil := math.Ceil(r)
	return int(rCeil)
}

func iterator(size int) []struct{} {
	return make([]struct{}, size)
}

// responseWithError adds an error log and returns the response
func responseWithError(resp *external.APIResponse, err error, code int, trace string) *external.APIResponse {
	resp.Message.ErrorLog = external.ErrorLogs{{
		ExceptionType: http.StatusText(code),
		StatusCode:    strconv.Itoa(code),
		Trace:         fmt.Sprintf("%s: error: %v", trace, err),
		RootCause:     err.Error(),
	}}
	return resp
}
