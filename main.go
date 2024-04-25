// Present HTTP metrics, in prometheus format, from terraform cloud
//
// We fetch the stats only once an hour, and trigger that via requests
// to the health endpoint.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-tfe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
)

// opsQueued holds our metrics
var opsQueued = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: "terraform_cloud",
		Name:      "workspace_count",
		Help:      "Number of workspaces in the project.",
	},
	[]string{
		// The name of the project
		"project",
	},
)

// timestamp records the last time we fetched our stats
var timestamp time.Time

// health is our endpoint for health checking - if there has been more than
// an hour since we last fetched terraform cloud stats then do that.
func health(w http.ResponseWriter, req *http.Request) {

	now := time.Now()
	difference := now.Sub(timestamp)

	// More than an hour?  Refresh
	if difference.Hours() > 1 {
		fetchCloudStats()
		timestamp = now
	}

	fmt.Fprintf(w, "ok")
}

// main is our entrypoint
func main() {

	// register our guages
	prometheus.MustRegister(opsQueued)

	// bind the endpoints
	http.Handle("/metrics", promhttp.Handler())

	// load our initial stats
	err := fetchCloudStats()
	if err != nil {
		fmt.Printf("error fetching stats: %s\n", err)
		return
	}
	// record the time
	timestamp = time.Now()

	// bind our HTTP handlers
	http.HandleFunc("/health", health)
	http.HandleFunc("/", health)

	// serve
	fmt.Printf("listening on http://127.0.0.1:8090/\n")
	http.ListenAndServe(":8090", nil)
}

func fetchCloudStats() error {

	// Passing nil to tfe.NewClient method will also use the default configuration
	client, err := tfe.NewClient(tfe.DefaultConfig())
	if err != nil {
		return err
	}

	// Get the organizatons we're a member of
	orgs, err := client.Organizations.List(context.Background(), nil)
	if err != nil {
		return err
	}

	// No orgs?  Then we're done
	if len(orgs.Items) < 1 {
		fmt.Printf("Not a member of any organizations!\n")
		return nil
	}

	// Work on each organization
	for _, org := range orgs.Items {

		// Get the projects
		projects, err := client.Projects.List(context.Background(), org.Name, nil)
		if err != nil {
			return err
		}

		// For each project
		for _, proj := range projects.Items {

			// Get workspaces in this project
			workspaces, err := client.Workspaces.List(context.Background(), org.Name, &tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   10,
				},
				ProjectID: proj.ID,
			})
			if err != nil {
				return err
			}

			// bump the stats
			opsQueued.With(prometheus.Labels{"project": proj.Name}).Set(float64(workspaces.TotalCount))

			//			fmt.Printf("Project %s [%d/%d]\n", proj.Name, len(workspaces.Items), workspaces.TotalCount)
		}
	}

	return nil
}
