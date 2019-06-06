// Package xrt ...
//
// Example


xrt.NewJob("512m", os.TempDir()).WithInput("/foo.csv").RunMapReduceJob(
//
package xrt

import (
	"fmt"
	"io"
)

// JobWithoutInputOutput ...
type JobWithoutInputOutput interface {
	WithInput(string) JobWithInput
	WithOutput(string) JobWithOutput
	RunMapJob(int, func(Context) error) error
	RunMapReduceJob(
		int,
		func(Context, RecordWriter) error,
		int,
		func(Context, RecordReader) error,
	) error
}

type jobWithoutInputOutput struct {
	memory  string
	tempDir string
}

func (j jobWithoutInputOutput) WithInput(input string) JobWithInput {
	return jobWithInput{j.memory, j.tempDir, input}
}

func (j jobWithoutInputOutput) WithOutput(output string) JobWithOutput {
	return jobWithOutput{j.memory, j.tempDir, output}
}

func (j jobWithoutInputOutput) RunMapJob(
	mappers int,
	mapper func(Context) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		"",
		"",
		mappers,
		func(c *context) error {
			return mapper(c)
		},
		0,
		nil,
	)
}

func (j jobWithoutInputOutput) RunMapReduceJob(
	mappers int,
	mapper func(Context, RecordWriter) error,
	reducers int,
	reducer func(Context, RecordReader) error,
) error {
	return run2(
		j.memory,
		j.tempDir,
		"",
		"",
		mappers,
		func(c *context) error {
			return c.withRecordWriter(func(w RecordWriter) error {
				return mapper(c, w)
			})
		},
		reducers,
		func(c *context) error {
			return c.withRecordReader(func(r RecordReader) error {
				return reducer(c, r)
			})
		},
	)
}

// NewJob ...
func NewJob(memory, tempDir string) JobWithoutInputOutput {
	return jobWithoutInputOutput{memory, tempDir}
}

// JobWithInput ...
type JobWithInput interface {
	RunMapJob(int, func(Context, io.Reader) error) error
	RunMapReduceJob(
		int,
		func(Context, io.Reader, RecordWriter) error,
		int,
		func(Context, RecordReader) error,
	) error
	WithOutput(string) JobWithInputOutput
}

type jobWithInput struct {
	memory  string
	tempDir string
	input   string
}

func (j jobWithInput) WithOutput(output string) JobWithInputOutput {
	return jobWithInputOutput{j.memory, j.tempDir, j.input, output}
}

func (j jobWithInput) RunMapJob(
	mappers int,
	mapper func(Context, io.Reader) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		j.input,
		"",
		mappers,
		func(c *context) error {
			return c.withInputReader(func(r io.Reader) error {
				return mapper(c, r)
			})
		},
		0,
		nil,
	)
}

func (j jobWithInput) RunMapReduceJob(
	mappers int,
	mapper func(Context, io.Reader, RecordWriter) error,
	reducers int,
	reducer func(Context, RecordReader) error,
) error {
	return run2(
		j.memory,
		j.tempDir,
		j.input,
		"",
		mappers,
		func(c *context) error {
			return c.withInputReader(func(r io.Reader) error {
				return c.withRecordWriter(func(w RecordWriter) error {
					return mapper(c, r, w)
				})
			})
		},
		reducers,
		func(c *context) error {
			return c.withRecordReader(func(r RecordReader) error {
				return reducer(c, r)
			})
		},
	)
}

// JobWithOutput ...
type JobWithOutput interface {
	RunMapJob(int, func(Context, io.Writer) error) error
	RunMapReduceJob(
		int,
		func(Context, RecordWriter) error,
		int,
		func(Context, RecordReader, io.Writer) error,
	) error
	WithInput(string) JobWithInputOutput
}

type jobWithOutput struct {
	memory  string
	tempDir string
	output  string
}

func (j jobWithOutput) WithInput(input string) JobWithInputOutput {
	return jobWithInputOutput{j.memory, j.tempDir, input, j.output}
}

func (j jobWithOutput) RunMapJob(
	mappers int,
	mapper func(Context, io.Writer) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		"",
		j.output,
		mappers,
		func(c *context) error {
			return c.withOutputWriter(func(w io.Writer) error {
				return mapper(c, w)
			})
		},
		0,
		nil,
	)
}

func (j jobWithOutput) RunMapReduceJob(
	mappers int,
	mapper func(Context, RecordWriter) error,
	reducers int,
	reducer func(Context, RecordReader, io.Writer) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		"",
		j.output,
		mappers,
		func(c *context) error {
			return c.withRecordWriter(func(w RecordWriter) error {
				return mapper(c, w)
			})
		},
		reducers,
		func(c *context) error {
			return c.withRecordReader(func(r RecordReader) error {
				return c.withOutputWriter(func(w io.Writer) error {
					return reducer(c, r, w)
				})
			})
		},
	)
}

// JobWithInputOutput ...
type JobWithInputOutput interface {
	RunMapJob(int, func(Context, io.Reader, io.Writer) error) error
	RunMapReduceJob(
		int,
		func(Context, io.Reader, RecordWriter) error,
		int,
		func(Context, RecordReader, io.Writer) error,
	) error
}

type jobWithInputOutput struct {
	memory  string
	tempDir string
	input   string
	output  string
}

func (j jobWithInputOutput) RunMapJob(
	mappers int,
	mapper func(Context, io.Reader, io.Writer) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		j.input,
		j.output,
		mappers,
		func(c *context) error {
			return c.withInputReader(func(r io.Reader) error {
				return c.withOutputWriter(func(w io.Writer) error {
					return mapper(c, r, w)
				})
			})
		},
		0,
		nil,
	)
}

func (j jobWithInputOutput) RunMapReduceJob(
	mappers int,
	mapper func(Context, io.Reader, RecordWriter) error,
	reducers int,
	reducer func(Context, RecordReader, io.Writer) error,
) error {
	return run(
		j.memory,
		j.tempDir,
		j.input,
		j.output,
		mappers,
		func(c *context) error {
			return c.withInputReader(func(r io.Reader) error {
				return c.withRecordWriter(func(w RecordWriter) error {
					return mapper(c, r, w)
				})
			})
		},
		reducers,
		func(c *context) error {
			return c.withRecordReader(func(r RecordReader) error {
				return c.withOutputWriter(func(w io.Writer) error {
					return reducer(c, r, w)
				})
			})
		},
	)
}

func run2(
	memory string,
	tempDir string,
	input string,
	output string,
	mappers int,
	mapper func(*context) error,
	reducers int,
	reducer func(*context) error,
) error {
	if reducers <= 0 {
		return fmt.Errorf(
			"error: reducers must be set to 1 or more for MapReduce jobs",
		)
	}
	if reducer == nil {
		return fmt.Errorf(
			"error: reducer cannot be nil for MapReduce jobs",
		)
	}
	return run(
		memory,
		tempDir,
		input,
		output,
		mappers,
		mapper,
		reducers,
		reducer,
	)
}

func run(
	memory string,
	tempDir string,
	input string,
	output string,
	mappers int,
	mapper func(*context) error,
	reducers int,
	reducer func(*context) error,
) error {
	conf, err := newConfig(input, output, config)
	if err != nil {
		return err
	}
	if err := conf.setupFileSystem(); err != nil {
		return err
	}
	//	var once sync.Once
	//	done := make(chan struct{})
	//	startTime := time.Now()
	//	conf, err := newConfig(config, input, output)
	//	if err != nil {
	//		return err
	//	}
	//	log.Print("")
	//	log.Print("starting xrt job")
	//	log.Print("")
	//	log.Print("configuration:")
	//	log.Print("")
	//	if conf.hasInput() {
	//		log.Printf("  input:    %s", conf.input)
	//	}
	//	log.Printf("  mappers:  %d", conf.mappers)
	//	if conf.hasReducer() {
	//		log.Printf("  reducers: %d", conf.reducers)
	//	}
	//	if conf.hasOutput() {
	//		log.Printf("  output:   %s", conf.output)
	//	}
	//	log.Printf("  memory:   %s", conf.memoryString)
	//	log.Printf("  tempdir:  %s", conf.tempDir)
	//	if err := conf.prepareFileSystem(); err != nil {
	//		return err
	//	}
	//	log.Print("")
	//	log.Print("running mapper stage")
	//	log.Print("")
	//	var in input
	//	if conf.hasInput() {
	//		in = newInput(conf.input)
	//	}
	//	var out output
	//	if conf.hasOutput() {
	//		out = newOutput(conf)
	//	}
	//	var buf buffers
	//	if conf.hasReducer() {
	//		buf = newBuffer(conf)
	//	}
	//	startTimeMappers := time.Now()
	//	if err := runWorkers(
	//		conf.mappers,
	//		func(c *context) error {
	//			return mapper(c, in, out, buf)
	//		},
	//	); err != nil {
	//		rollback(once, done, err)
	//	}
	//	durationMappers := time.Since(startTimeMappers)
	//	log.Print("")
	//	var durationReducers time.Duration
	//	if conf.hasReducer() {
	//		log.Print("running reducer stage")
	//		log.Print("")
	//		startTimeReducers := time.Now()
	//		if err := runWorkers(
	//			conf.reducers,
	//			func(c *context) error {
	//				return reducer(c, buf)
	//			},
	//		); err != nil {
	//			rollback(once, done, err)
	//		}
	//		durationReducers = time.Since(startTimeReducers)
	//		log.Print("")
	//	}
	//	if hasOutput() {
	//		log.Print("committing")
	//		log.Print("")
	//		commit()
	//	}
	//	log.Printf("  mappers runtime: %s", durationMappers.String())
	//	if conf.hasReducer() {
	//		log.Printf("  reducers runtime: %s", durationReducers.String())
	//	}
	//	log.Printf("  total runtime: %s", time.Since(startTime).String())
	//	log.Print("")
	//	log.Print("success")
	//	if !hasOutput() {
	//		printOutput()
	//	}
	//	cleanup()
	return nil
}

//// startInterruptHandler launches a handler that will catch the first interrupt signal and attempt
//// a graceful termination (mainly to deal with ctrl-c)
//func startInterruptHandler(once sync.Once, exit ExitFunc) {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt)
//	go func() {
//		<-c
//		rollback(once, exit, errors.New("received interrupt - aborting job"))
//	}()
//}
//
////func runMany(workers int, worker func(context) error) error {
////	errc := make(chan error)
////	for i := 0; i < workers; i++ {
////		go func(wid int) { errc <- worker(context{wid, mappers, reducers}) }(i)
////	}
////	for i := 0; i < workers; i++ {
////		if err := <-errc; err != nil {
////			return err
////		}
////	}
////	return nil
////}
//func runMappers(conf *config, mapper MapperFunc) ([][]*buffer, error) {
//	in := newInput(conf)
//	out := new
//
//	var buffers [][]*buffer
//	buffers = make([][]*buffer, conf.mappers)
//	for i := range buffers {
//		buffers[i] = make([]*buffer, conf.reducers)
//	}
//
//	if conf.hasOutput() || conf.hasReducer() {
//		if tempDir, err = ioutil.TempDir(conf.tempDir, "xrt-"); err != nil {
//			return fmt.Errorf("xrt: failed initializing directory '%s' - %v", tempDir, err)
//		}
//	}
//	if conf.hasOutput() {
//		if err = os.Mkdir(tempOutput, 0700); err != nil {
//			return fmt.Errorf("xrt: failed initializing directory '%s' - %v", tempOutput, err)
//		}
//	}
//	if conf.hasReducer() {
//		if err = os.Mkdir(tempSpill, 0700); err != nil {
//			return fmt.Errorf("xrt: failed initializing directory '%s' - %v", tempSpill, err)
//		}
//	}
//
//	input := newInput(conf)
//	//	errc := make(chan error)
//	//	for i := 0; i < mappers; i++ {
//	//		go func(wid int) {
//	//			l := logger{wid}
//	//			errc <- worker(MapperContext{wid, mapper, mappers, reducers})
//	//		}(i)
//	//	}
//	//	for i := 0; i < workers; i++ {
//	//		if err := <-errc; err != nil {
//	//			return err
//	//		}
//	//	}
//	//	return nil
//}
//
//func mapWorker(c context) error {
//	c.log("mapper starting")
//	defer c.log("done")
//	if hasReducer() {
//		bufMem := memory / (mappers * reducers)
//		for i := range buffers[c.workerID] {
//			spillDir := path.Join(tempSpill, strconv.Itoa(c.workerID), strconv.Itoa(i))
//			buffers[c.workerID][i] = newBuffer(bufMem, spillDir)
//		}
//	}
//	if err := c.exec(mapper, mapStdinHandler, mapStdoutHandler, logStream); err != nil {
//		return err
//	}
//	if hasReducer() && reducers <= mappers {
//		c.log("sorting")
//		for i, b := range buffers[c.workerID] {
//			b.sort()
//			if b.needExternalSort() {
//				c.logf("merging %d spill files for reducer %d", b.spills, i)
//				if err := b.externalSort(); err != nil {
//					return err
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func mapStdinHandler(c context, w io.WriteCloser) error {
//	if hasInput() {
//		return inputStream(c, w, inputs)
//	}
//	return w.Close()
//}
//
//func mapStdoutHandler(c context, r io.ReadCloser) error {
//	if hasReducer() {
//		return intermediateMapStream(c, r, buffers[c.workerID])
//	}
//	return outputStream(c, r, tempOutput)
//}
//
//func reduceWorker(c context) error {
//	c.log("reducer starting")
//	defer c.log("done")
//	if hasReducer() && reducers > mappers {
//		c.log("sorting")
//		for i, b := range buffers[c.workerID] {
//			b.sort()
//			if b.needExternalSort() {
//				c.logf("merging %d spill files for reducer %d", b.spills, i)
//				if err := b.externalSort(); err != nil {
//					return err
//				}
//			}
//		}
//	}
//	records := 0
//	for i := range buffers {
//		records += buffers[i][c.workerID].records
//	}
//	c.logf("processing %d records", records)
//	return c.exec(reducer, reduceStdinHandler, reduceStdoutHandler, logStream)
//}
//
//func reduceStdinHandler(c context, w io.WriteCloser) error {
//	bufs := make([]*buffer, len(buffers))
//	for i := range buffers {
//		bufs[i] = buffers[i][c.workerID]
//	}
//	return intermediateReduceStream(c, w, bufs)
//}
//
//func reduceStdoutHandler(c context, r io.ReadCloser) error {
//	return outputStream(c, r, tempOutput)
//}
//
//// rollback ensure graceful termination of a failed job. It first executes the shutdownFunc
//// before cleaning up any temporary data stored on disk.
//func rollback(once sync.Once, exit ExitFunc, err error) {
//	once.Do(func() {
//		log.Print("error - attempting rollback")
//		log.Print("")
//		log.Print(err)
//		if exit != nil {
//			exit()
//		}
//		cleanup()
//		log.Print("failed")
//		os.Exit(1)
//	})
//}
//
//// commit ensures transactional termination of succesfull jobs. If the job is configured to create
//// output commit uses a directory move to transactioanlly "commit" the output from a temporary
//// folder to the final output folder.
//func commit() {
//	if err := os.Rename(tempOutput, output); err != nil {
//		log.Printf("  error moving output data from %s to %s - %v", tempOutput, output, err)
//		log.Printf("  temporary data directory %s was not removed", tempDir)
//		log.Print("failed")
//		os.Exit(1)
//	}
//}
//
//// printOutput reads the output from the temporary directory and copies it to stdout. This is used
//// by runs without output to display the output in the terminal instead of writting it to a file.
//func printOutput() {
//	files, err := ioutil.ReadDir(tempOutput)
//	if err != nil {
//		log.Printf("  error reading output data in %s - %v", tempOutput, err)
//		os.Exit(1)
//	}
//	for _, file := range files {
//		filename := path.Join(tempOutput, file.Name())
//		f, err := os.Open(filename)
//		if err != nil {
//			log.Printf("  error reading output data in %s - %v", filename, err)
//			os.Exit(1)
//		}
//		r := bufio.NewReader(f)
//		w := bufio.NewWriter(os.Stdout)
//		if _, err := io.Copy(w, r); err != nil {
//			log.Printf("  error copying output data in %s to stdout - %v", filename, err)
//		}
//		if err := w.Flush(); err != nil {
//			log.Printf("  error copying output data in %s to stdout - %v", filename, err)
//		}
//	}
//}
//
//// cleanup removes any remaining temporary files.
//func cleanup() {
//	// BUG this will break on windows since it does not allow removal of open files and by the
//	// time this is called it is possible fds in the tempdir are still open.
//	if err := os.RemoveAll(tempDir); err != nil {
//		log.Printf("  failed to remove temporary data directory %s - %v", tempDir, err)
//	}
//}
//
//func hasInput() bool {
//	return len(input) > 0
//}
//
//func hasMapper() bool {
//	return len(mapper) > 0
//}
//
//func hasReducer() bool {
//	return len(reducer) > 0
//}
//
//func hasOutput() bool {
//	return len(output) > 0
//}
