package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	image2 "github.com/aquasecurity/fanal/artifact/image"
	"github.com/aquasecurity/fanal/image"
	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/cache"

	// "github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/types"
	// "golang.org/x/xerrors"
)

func main() {
	imageFlag := flag.String("image", "", "image name")
	remoteFlag := flag.String("remote", "", "server url")
	outputTypeFlag := flag.String("output", "table", "output type such as table, json.")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	sc, cleanUp, err := initializeDockerScanner(ctx, *imageFlag, client.CustomHeaders{}, client.RemoteURL(*remoteFlag), time.Second*5000)
	if err != nil {
		fmt.Println("could not initialize scanner: %v", err)
	}

	defer cleanUp()

	results, err := sc.ScanArtifact(ctx, types.ScanOptions{
		VulnType:            []string{"os", "library"},
		ScanRemovedPackages: true,
		ListAllPackages:     true,
	})
	if err != nil {
		fmt.Println("could not scan image: %v", err)
	}

	if len(results) > 0 {
		fmt.Println("%d vulnerability/ies found", len(results[0].Vulnerabilities))

		if err = report.WriteResults(*outputTypeFlag, os.Stdout, []dbTypes.Severity{dbTypes.SeverityUnknown}, results, "", false); err != nil {
			fmt.Println("could not write results: s", err)
		}
	} else {
		fmt.Println("no vulnerabilities found for image %s", *imageFlag)
	}

}

func initializeDockerScanner(ctx context.Context, imageName string, customHeaders client.CustomHeaders, url client.RemoteURL, timeout time.Duration) (scanner.Scanner, func(), error) {
	scannerScanner := client.NewProtobufClient(url)
	clientScanner := client.NewScanner(customHeaders, scannerScanner)
	artifactCache := cache.NewRemoteCache(cache.RemoteURL(url), nil)
	dockerOption, err := types.GetDockerOption(timeout)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	imageImage, cleanup, err := image.NewDockerImage(ctx, imageName, dockerOption)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	artifact := image2.NewArtifact(imageImage, artifactCache, nil)
	scanner2 := scanner.NewScanner(clientScanner, artifact)
	return scanner2, func() {
		cleanup()
	}, nil
}
