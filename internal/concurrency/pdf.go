package concurrency

import (
	"sync"

	pdfgen "github.com/Dimensionexpert/payslip/internal/genPDF"
)

type ConversionJob struct {
	Filepath string
	OutDir   string
}

type ConversionResult struct {
	Err      error
	Filepath string
}

func pdfWorker(
	workerID int,
	jobs <-chan ConversionJob,
	results chan<- ConversionResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		err := pdfgen.ConvertToPDF(
			job.Filepath,
			job.OutDir,
			workerID,
		)

		results <- ConversionResult{
			Filepath: job.Filepath,
			Err:      err,
		}
	}
}

// RunPDFConversion converts files using a fixed number of workers.
// onResult is called once for every completed conversion.
func RunPDFConversion(
	jobsList []ConversionJob,
	numWorkers int,
	onResult func(ConversionResult),
) []ConversionResult {

	if len(jobsList) == 0 {
		return nil
	}

	if numWorkers <= 0 {
		numWorkers = 1
	}

	if numWorkers > len(jobsList) {
		numWorkers = len(jobsList)
	}

	jobs := make(chan ConversionJob, len(jobsList))
	results := make(chan ConversionResult, len(jobsList))

	var wg sync.WaitGroup

	for workerID := 1; workerID <= numWorkers; workerID++ {
		wg.Add(1)

		go pdfWorker(
			workerID,
			jobs,
			results,
			&wg,
		)
	}

	for _, job := range jobsList {
		jobs <- job
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make([]ConversionResult, 0, len(jobsList))

	for result := range results {
		collected = append(collected, result)

		if onResult != nil {
			onResult(result)
		}
	}

	return collected
}
