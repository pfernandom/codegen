package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/hashicorp/go-getter"
	safetemp "github.com/hashicorp/go-safetemp"
	"github.com/pfernandom/codegen/errors"
	"github.com/pfernandom/codegen/gen"
	"github.com/pfernandom/codegen/questions"
)

var Must = errors.MustNow

func main() {
	// pwd, err := os.Getwd()
	// Must(err)
	// Must(err)
	tmpDir, closer, err := safetemp.Dir("/tmp", "sample*dir")
	Must(err)
	defer closer.Close()

	outDir, closer2, err := safetemp.Dir("/tmp", "sample*dir")
	Must(err)
	defer closer2.Close()

	opts := []getter.ClientOption{}
	opts = append(opts, getter.WithProgress(NewProgressBar()))
	fmt.Println("tmpDir", tmpDir)

	client := getter.Client{
		Src: "git@github.com:pfernandom/codegen-example.git",
		Dst: tmpDir,
		// Pwd:     pwd,
		Mode:    getter.ClientModeDir,
		Options: opts,
	}
	errors.Must(client.Get()).OrFailWith("get error: %w", err)
	fmt.Println("tmpDir", tmpDir)
	for _, file := range errors.MustReturn(os.ReadDir(tmpDir)).OrFail() {
		fmt.Println(file.Name())
	}

	flags := gen.CmdFlags{
		DataFile:      filepath.Join(tmpDir, "data.json"),
		TemplateDir:   filepath.Join(tmpDir, "templates"),
		OutDir:        outDir,
		QuestionsFile: filepath.Join(tmpDir, "questions.json"),
	}
	questionsHandler := questions.NewQuestionsHandler()
	fmt.Printf("Creating project from %s and %s to %s\n", flags.DataFile, flags.TemplateDir, flags.OutDir)

	dg := gen.NewDataGen(flags, questionsHandler)

	Must(os.RemoveAll(flags.OutDir))
	Must(os.MkdirAll(flags.OutDir, 0755))
	Must(dg.Generate())
	fmt.Println("done")

	for _, file := range errors.MustReturn(os.ReadDir(outDir)).OrFail() {
		fmt.Println(file.Name())
	}

}

// ProgressBar wraps a github.com/cheggaaa/pb.Pool
// in order to display download progress for one or multiple
// downloads.
//
// If two different instance of ProgressBar try to
// display a progress only one will be displayed.
// It is therefore recommended to use DefaultProgressBar
type ProgressBar struct {
	// lock everything below
	lock sync.Mutex

	pool *pb.Pool

	pbs int
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		pool: pb.NewPool(),
	}
}

func ProgressBarConfig(bar *pb.ProgressBar, prefix string) {
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
}

// TrackProgress instantiates a new progress bar that will
// display the progress of stream until closed.
// total can be 0.
func (cpb *ProgressBar) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	cpb.lock.Lock()
	defer cpb.lock.Unlock()

	newPb := pb.New64(totalSize)
	newPb.SetCurrent(currentSize)
	ProgressBarConfig(newPb, filepath.Base(src))
	if cpb.pool == nil {
		cpb.pool = pb.NewPool()
		_ = cpb.pool.Start()
	}
	cpb.pool.Add(newPb)
	reader := newPb.NewProxyReader(stream)

	cpb.pbs++
	return &readCloser{
		Reader: reader,
		close: func() error {
			cpb.lock.Lock()
			defer cpb.lock.Unlock()

			newPb.Finish()
			cpb.pbs--
			if cpb.pbs <= 0 {
				_ = cpb.pool.Stop()
				cpb.pool = nil
			}
			return nil
		},
	}
}

type readCloser struct {
	io.Reader
	close func() error
}

func (c *readCloser) Close() error { return c.close() }
