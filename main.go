package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/spanner"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/kelseyhightower/envconfig"
	"github.com/sinmetal/srun/backend"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"
)

type EnvConfig struct {
	SpannerDatabase string `required:"true"`
	TracePrefix     string `default:"default"`
}

func main() {
	ctx := context.Background()

	var env EnvConfig
	if err := envconfig.Process("srun", &env); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("ENV_CONFIG %+v\n", env)

	project, err := metadatabox.ProjectID()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ProjectID:%s\n", project)
	{
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: project,
		})
		if err != nil {
			panic(err)
		}
		trace.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	sc, err := createClient(ctx, env.SpannerDatabase)
	if err != nil {
		panic(err)
	}
	// TODO Spanner Client Close

	ah := backend.AppHandlers{
		TweetStore: backend.NewTweetStore(sc),
	}

	http.HandleFunc("/tweetInsert", ah.TweetInsertHandler)
	http.HandleFunc("/tweetUpdate", ah.TweetUpdateHandler)
	http.HandleFunc("/tweetUpdateDML", ah.TweetUpdateDMLHandler)
	http.HandleFunc("/tweetUpdateBatchDML", ah.TweetUpdateBatchDMLHandler)
	http.HandleFunc("/tweetUpdateAndSelect", ah.TweetUpdateAndSelectHandler)
	http.HandleFunc("/tweetUpdateDMLAndSelect", ah.TweetUpdateDMLAndSelectHandler)
	http.HandleFunc("/", ah.HelloHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func createClient(ctx context.Context, db string, o ...option.ClientOption) (*spanner.Client, error) {
	config := spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened:           200,
			TrackSessionHandles: true,
		},
	}
	dataClient, err := spanner.NewClientWithConfig(ctx, db, config, o...)
	if err != nil {
		return nil, err
	}

	return dataClient, nil
}
