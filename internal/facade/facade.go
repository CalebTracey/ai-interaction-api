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
	"sync/atomic"
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
	apiRequest.Size = LargeImage
	g, ctx := errgroup.WithContext(ctx)

	// want x images total
	wantTotal := apiRequest.N
	//log.Infof("want total: %v", wantTotal)

	requestCount := 1
	if wantTotal/ChunkSize > 1 {
		requestCount = wantTotal / ChunkSize
	}

	//log.Infof("request count: %v", requestCount)
	responseChan := make(chan external.AIResponse, requestCount)
	workers := int32(requestCount)

	//log.Infof("chunking request...")
	// break up request according to ChunkSize and make requests concurrently
	for i := 1; i <= requestCount; i++ {
		// if last request, check if you don't need the full request size
		amount := i * ChunkSize

		if amount > wantTotal {
			// amount: 5.0, wantTotal: 3.0; want = 3.0
			amount = ChunkSize - (amount - wantTotal)
		}

		if amount < 0 {
			amount = wantTotal
		}
		//log.Infof("current request: %d; amount: %d", i, amount)
		apiRequest.N = amount

		g.Go(func() error {
			defer func() {
				// last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					//log.Infoln("last request, closing channel...")
					close(responseChan)
				}
			}()

			// check if the current request doesn't require the full image count
			// if not, set the request size to the remainder
			if _, err := s.requestChunk(ctx, apiRequest, responseChan); err != nil {
				return err
			}

			return nil
		})
	}

	g.Go(func() error {
		for response := range responseChan {
			apiResp.Result.Data = append(apiResp.Result.Data, response.Data...)
		}
		return nil
	})

	// wait for go routines to finish can catch any errors
	if err := g.Wait(); err != nil {
		return responseWithError(apiResp, err, http.StatusInternalServerError, "ImageRequest")
	}

	// TODO move mapping
	// TODO created seems broken at the source right now

	apiResp.Result.Created = len(apiResp.Result.Data)
	apiResp.Message.Count = len(apiResp.Result.Data)

	return apiResp
}

func (s Service) requestChunk(ctx context.Context, apiRequest *external.APIRequest, imageChan chan<- external.AIResponse) (resp external.AIResponse, err error) {
	if resp, err = s.DAO.ImageRequest(ctx, apiRequest); err == nil {
		select {
		case <-ctx.Done():
			log.Errorf("=== getImages: context cancelled")
			return resp, ctx.Err()
		case imageChan <- resp:
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
