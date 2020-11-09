package seaweed

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chrislusf/seaweedfs/weed/storage/needle"

	//"github.com/qiaogw/pkg/store/driver/local"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/pb"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/security"
	"github.com/chrislusf/seaweedfs/weed/util"
	"github.com/chrislusf/seaweedfs/weed/wdclient"
)

var (
	//seaweedfs SeaOptions
	swGroup sync.WaitGroup
)

type SeaOptions struct {
	filerUrl          string
	filerAddress      string
	filerHost         string
	filerPort         uint64
	include           string
	replication       string
	collection        string
	ttl               string
	maxMB             int
	masterClient      *wdclient.MasterClient
	concurrenctFiles  int
	concurrenctChunks int
	grpcDialOption    grpc.DialOption
	masters           []string
	cipher            bool
	ttlSec            int32
}

type FileTask struct {
	sourceLocation     string
	file               *os.File
	destinationUrlPath string
	fileSize           int64
	fileMode           os.FileMode
	uid                uint32
	gid                uint32
}

//func init() {
//	seaweedfs.include = ""
//	seaweedfs.concurrenctFiles = 8
//	seaweedfs.concurrenctChunks = 8
//}
func NewSeaweed(filerUrls string) (*SeaOptions, error) {
	sw := &SeaOptions{
		include:           "",
		concurrenctFiles:  8,
		concurrenctChunks: 8,
	}

	util.LoadConfiguration("security", false)
	//fileOrDirs := []string{"./client.go"}
	//filerDestination := "http://192.168.0.140:8888/github/"
	filerUrl, err := url.Parse(filerUrls)
	if err != nil {
		fmt.Printf("The last argument should be a URL on filer: %v\n", err)
		return nil, err
	}
	urlPath := filerUrl.Path
	if !strings.HasSuffix(urlPath, "/") {
		fmt.Printf("The last argument should be a folder and end with \"/\": %v\n", err)
		return nil, err
	}

	if filerUrl.Port() == "" {
		fmt.Printf("The filer port should be specified.\n")
		return nil, err
	}

	filerPort, parseErr := strconv.ParseUint(filerUrl.Port(), 10, 64)
	if parseErr != nil {
		fmt.Printf("The filer port parse error: %v\n", parseErr)
		return nil, err
	}

	filerGrpcPort := filerPort + 10000
	filerGrpcAddress := fmt.Sprintf("%s:%d", filerUrl.Hostname(), filerGrpcPort)
	sw.grpcDialOption = security.LoadClientTLS(util.GetViper(), "grpc.client")

	masters, collection, replication, maxMB, cipher, err := readFilerConfig(sw.grpcDialOption, filerGrpcAddress)
	if err != nil {
		fmt.Printf("read from filer %s: %v\n", filerGrpcAddress, err)
		return nil, err
	}
	sw.filerUrl = filerGrpcAddress
	sw.filerAddress = filerGrpcAddress
	sw.filerHost = filerUrl.Host
	sw.filerPort = filerGrpcPort
	sw.collection = collection
	sw.replication = replication
	sw.maxMB = int(maxMB)
	sw.masters = masters
	sw.cipher = cipher

	ttl, err := needle.ReadTTL(sw.ttl)
	if err != nil {
		fmt.Printf("parsing ttl %s: %v\n", sw.ttl, err)
		return nil, err
	}
	sw.ttlSec = int32(ttl.Minutes()) * 60
	return sw, nil
}

func (s *SeaOptions) Put1(files *os.File, distPath string) error {

	fileCopyTaskChan := make(chan FileTask, s.concurrenctFiles)
	//urlPath := s.filerUrl + distPath
	go func() {
		defer close(fileCopyTaskChan)
		//for _, file := range files {
		if err := genFileTask(files, distPath, fileCopyTaskChan); err != nil {
			fmt.Fprintf(os.Stderr, "gen file list error: %v\n", err)
			//break
		}
		//}
	}()
	for i := 0; i < s.concurrenctFiles; i++ {
		swGroup.Add(1)
		go func() {
			defer swGroup.Done()
			worker := FileWorker{
				options:          s,
				filerHost:        s.filerHost,
				filerGrpcAddress: s.filerAddress,
			}
			if err := worker.copyFiles(fileCopyTaskChan); err != nil {
				fmt.Fprintf(os.Stderr, "seaweedfs file error: %v\n", err)
				return
			}
		}()
	}
	swGroup.Wait()
	return nil
}

func (s *SeaOptions) Put(files []*os.File, distPath string) error {
	fileCopyTaskChan := make(chan FileTask, s.concurrenctFiles)
	//urlPath := s.filerUrl + distPath
	go func() {
		defer close(fileCopyTaskChan)
		for _, file := range files {
			if err := genFileTask(file, distPath, fileCopyTaskChan); err != nil {
				fmt.Fprintf(os.Stderr, "gen file list error: %v\n", err)
				break
			}
		}
	}()
	for i := 0; i < s.concurrenctFiles; i++ {
		swGroup.Add(1)
		go func() {
			defer swGroup.Done()
			worker := FileWorker{
				options:          s,
				filerHost:        s.filerHost,
				filerGrpcAddress: s.filerAddress,
			}
			if err := worker.copyFiles(fileCopyTaskChan); err != nil {
				fmt.Fprintf(os.Stderr, "seaweedfs file error: %v\n", err)
				return
			}
		}()
	}
	swGroup.Wait()

	return nil
}

func readFilerConfig(grpcDialOption grpc.DialOption, filerGrpcAddress string) (masters []string, collection, replication string, maxMB uint32, cipher bool, err error) {
	err = pb.WithGrpcFilerClient(filerGrpcAddress, grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
		if err != nil {
			return fmt.Errorf("get filer %s configuration: %v", filerGrpcAddress, err)
		}
		masters, collection, replication, maxMB = resp.Masters, resp.Collection, resp.Replication, resp.MaxMb
		cipher = resp.Cipher
		return nil
	})
	return
}

func genFileTask(fileOrDir *os.File, destPath string, fileCopyTaskChan chan FileTask) error {

	fs, err := ioutil.ReadAll(fileOrDir)
	if err != nil {
		return err
	}
	s := len(fs)

	fileCopyTaskChan <- FileTask{
		file:               fileOrDir,
		sourceLocation:     "",
		destinationUrlPath: destPath,
		fileSize:           int64(s),
		//fileMode:           fi.Mode(),
		//uid: uid,
		//gid: gid,
	}

	return nil
}

type FileWorker struct {
	options          *SeaOptions
	filerHost        string
	filerGrpcAddress string
}

func (worker *FileWorker) copyFiles(fileCopyTaskChan chan FileTask) error {
	for task := range fileCopyTaskChan {
		if err := worker.doEachCopy(task); err != nil {
			return err
		}
	}
	return nil
}

func (worker *FileWorker) doEachCopy(task FileTask) error {

	if worker.options.include != "" {
		if ok, _ := filepath.Match(worker.options.include, filepath.Base(task.sourceLocation)); !ok {
			return nil
		}
	}

	// find the chunk count
	chunkSize := int64(worker.options.maxMB * 1024 * 1024)
	chunkCount := 1

	if chunkSize > 0 && task.fileSize > chunkSize {
		chunkCount = int(task.fileSize/chunkSize) + 1
	}

	if chunkCount == 1 {
		return worker.uploadFileAsOne(task, task.file)
	}

	return worker.uploadFileInChunks(task, task.file, chunkCount, chunkSize)
}

func (worker *FileWorker) uploadFileAsOne(task FileTask, f *os.File) error {

	// upload the file content
	fileName := filepath.Base(f.Name())
	//fmt.Println("fileName", fileName)
	mimeType := detectMime(f)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var chunks []*filer_pb.FileChunk
	var assignResult *filer_pb.AssignVolumeResponse
	var assignError error

	if task.fileSize > 0 {

		// assign a volume
		err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {

			request := &filer_pb.AssignVolumeRequest{
				Count:       1,
				Replication: worker.options.replication,
				Collection:  worker.options.collection,
				TtlSec:      worker.options.ttlSec,
				ParentPath:  task.destinationUrlPath,
			}

			assignResult, assignError = client.AssignVolume(context.Background(), request)
			if assignError != nil {
				return fmt.Errorf("assign volume failure %v: %v", request, assignError)
			}
			if assignResult.Error != "" {
				return fmt.Errorf("assign volume failure %v: %v", request, assignResult.Error)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failed to assign from %v: %v\n", worker.options.masters, err)
		}
		targetUrl := "http://" + assignResult.PublicUrl + "/" + assignResult.FileId

		uploadResult, err := operation.UploadData(targetUrl, fileName, worker.options.cipher, data, false, mimeType, nil, security.EncodedJwt(assignResult.Auth))
		if err != nil {
			return fmt.Errorf("upload data %v to %s: %v\n", fileName, targetUrl, err)
		}
		if uploadResult.Error != "" {
			return fmt.Errorf("upload %v to %s result: %v\n", fileName, targetUrl, uploadResult.Error)
		}
		fmt.Printf("uploaded %s to %s\n", fileName, targetUrl)

		chunks = append(chunks, uploadResult.ToPbFileChunk(assignResult.FileId, 0))

		fmt.Printf("copied %s => http://%s%s%s\n", fileName, worker.filerHost, task.destinationUrlPath, fileName)
	}

	if err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		request := &filer_pb.CreateEntryRequest{
			Directory: task.destinationUrlPath,
			Entry: &filer_pb.Entry{
				Name: fileName,
				Attributes: &filer_pb.FuseAttributes{
					Crtime:      time.Now().Unix(),
					Mtime:       time.Now().Unix(),
					Gid:         task.gid,
					Uid:         task.uid,
					FileSize:    uint64(task.fileSize),
					FileMode:    uint32(task.fileMode),
					Mime:        mimeType,
					Replication: worker.options.replication,
					Collection:  worker.options.collection,
					TtlSec:      worker.options.ttlSec,
				},
				Chunks: chunks,
			},
		}

		if err := filer_pb.CreateEntry(client, request); err != nil {
			return fmt.Errorf("update fh: %v", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("upload data %v to http://%s%s%s: %v\n", fileName, worker.filerHost, task.destinationUrlPath, fileName, err)
	}

	return nil
}

func (worker *FileWorker) uploadFileInChunks(task FileTask, f *os.File, chunkCount int, chunkSize int64) error {

	fileName := filepath.Base(f.Name())
	mimeType := detectMime(f)

	chunksChan := make(chan *filer_pb.FileChunk, chunkCount)

	concurrentChunks := make(chan struct{}, worker.options.concurrenctChunks)
	var wg sync.WaitGroup
	var uploadError error
	var collection, replication string

	fmt.Printf("uploading %s in %d chunks ...\n", fileName, chunkCount)
	for i := int64(0); i < int64(chunkCount) && uploadError == nil; i++ {
		wg.Add(1)
		concurrentChunks <- struct{}{}
		go func(i int64) {
			defer func() {
				wg.Done()
				<-concurrentChunks
			}()
			// assign a volume
			var assignResult *filer_pb.AssignVolumeResponse
			var assignError error
			err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
				request := &filer_pb.AssignVolumeRequest{
					Count:       1,
					Replication: worker.options.replication,
					Collection:  worker.options.collection,
					TtlSec:      worker.options.ttlSec,
					ParentPath:  task.destinationUrlPath,
				}

				assignResult, assignError = client.AssignVolume(context.Background(), request)
				if assignError != nil {
					return fmt.Errorf("assign volume failure %v: %v", request, assignError)
				}
				if assignResult.Error != "" {
					return fmt.Errorf("assign volume failure %v: %v", request, assignResult.Error)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Failed to assign from %v: %v\n", worker.options.masters, err)
			}
			if err != nil {
				fmt.Printf("Failed to assign from %v: %v\n", worker.options.masters, err)
			}

			targetUrl := "http://" + assignResult.PublicUrl + "/" + assignResult.FileId
			if collection == "" {
				collection = assignResult.Collection
			}
			if replication == "" {
				replication = assignResult.Replication
			}

			uploadResult, err, _ := operation.Upload(targetUrl, fileName+"-"+strconv.FormatInt(i+1, 10), worker.options.cipher, io.NewSectionReader(f, i*chunkSize, chunkSize), false, "", nil, security.EncodedJwt(assignResult.Auth))
			if err != nil {
				uploadError = fmt.Errorf("upload data %v to %s: %v\n", fileName, targetUrl, err)
				return
			}
			if uploadResult.Error != "" {
				uploadError = fmt.Errorf("upload %v to %s result: %v\n", fileName, targetUrl, uploadResult.Error)
				return
			}
			chunksChan <- uploadResult.ToPbFileChunk(assignResult.FileId, i*chunkSize)

			fmt.Printf("uploaded %s-%d to %s [%d,%d)\n", fileName, i+1, targetUrl, i*chunkSize, i*chunkSize+int64(uploadResult.Size))
		}(i)
	}
	wg.Wait()
	close(chunksChan)

	var chunks []*filer_pb.FileChunk
	for chunk := range chunksChan {
		chunks = append(chunks, chunk)
	}

	if uploadError != nil {
		var fileIds []string
		for _, chunk := range chunks {
			fileIds = append(fileIds, chunk.FileId)
		}
		operation.DeleteFiles(worker.options.masters[0], false, worker.options.grpcDialOption, fileIds)
		return uploadError
	}

	if err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		request := &filer_pb.CreateEntryRequest{
			Directory: task.destinationUrlPath,
			Entry: &filer_pb.Entry{
				Name: fileName,
				Attributes: &filer_pb.FuseAttributes{
					Crtime:      time.Now().Unix(),
					Mtime:       time.Now().Unix(),
					Gid:         task.gid,
					Uid:         task.uid,
					FileSize:    uint64(task.fileSize),
					FileMode:    uint32(task.fileMode),
					Mime:        mimeType,
					Replication: replication,
					Collection:  collection,
					TtlSec:      worker.options.ttlSec,
				},
				Chunks: chunks,
			},
		}

		if err := filer_pb.CreateEntry(client, request); err != nil {
			return fmt.Errorf("update fh: %v", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("upload data %v to http://%s%s%s: %v\n", fileName, worker.filerHost, task.destinationUrlPath, fileName, err)
	}

	fmt.Printf("copied %s => http://%s%s%s\n", fileName, worker.filerHost, task.destinationUrlPath, fileName)

	return nil
}

func detectMime(f *os.File) string {
	head := make([]byte, 512)
	f.Seek(0, io.SeekStart)
	n, err := f.Read(head)
	if err == io.EOF {
		return ""
	}
	if err != nil {
		fmt.Printf("read head of %v: %v\n", f.Name(), err)
		return ""
	}
	f.Seek(0, io.SeekStart)
	mimeType := http.DetectContentType(head[:n])
	if mimeType == "application/octet-stream" {
		return ""
	}
	return mimeType
}

func (worker *FileWorker) uploadFileDataAsOne(task FileTask, f *os.File) error {

	// upload the file content
	fileName := filepath.Base(f.Name())
	mimeType := detectMime(f)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var chunks []*filer_pb.FileChunk
	var assignResult *filer_pb.AssignVolumeResponse
	var assignError error

	if task.fileSize > 0 {

		// assign a volume
		err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {

			request := &filer_pb.AssignVolumeRequest{
				Count:       1,
				Replication: worker.options.replication,
				Collection:  worker.options.collection,
				TtlSec:      worker.options.ttlSec,
				ParentPath:  task.destinationUrlPath,
			}

			assignResult, assignError = client.AssignVolume(context.Background(), request)
			if assignError != nil {
				return fmt.Errorf("assign volume failure %v: %v", request, assignError)
			}
			if assignResult.Error != "" {
				return fmt.Errorf("assign volume failure %v: %v", request, assignResult.Error)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failed to assign from %v: %v\n", worker.options.masters, err)
		}
		targetUrl := "http://" + assignResult.PublicUrl + "/" + assignResult.FileId

		uploadResult, err := operation.UploadData(targetUrl, fileName, worker.options.cipher, data, false, mimeType, nil, security.EncodedJwt(assignResult.Auth))
		if err != nil {
			return fmt.Errorf("upload data %v to %s: %v\n", fileName, targetUrl, err)
		}
		if uploadResult.Error != "" {
			return fmt.Errorf("upload %v to %s result: %v\n", fileName, targetUrl, uploadResult.Error)
		}
		fmt.Printf("uploaded %s to %s\n", fileName, targetUrl)

		chunks = append(chunks, uploadResult.ToPbFileChunk(assignResult.FileId, 0))

		fmt.Printf("copied %s => http://%s%s%s\n", fileName, worker.filerHost, task.destinationUrlPath, fileName)
	}

	if err := pb.WithGrpcFilerClient(worker.filerGrpcAddress, worker.options.grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		request := &filer_pb.CreateEntryRequest{
			Directory: task.destinationUrlPath,
			Entry: &filer_pb.Entry{
				Name: fileName,
				Attributes: &filer_pb.FuseAttributes{
					Crtime:      time.Now().Unix(),
					Mtime:       time.Now().Unix(),
					Gid:         task.gid,
					Uid:         task.uid,
					FileSize:    uint64(task.fileSize),
					FileMode:    uint32(task.fileMode),
					Mime:        mimeType,
					Replication: worker.options.replication,
					Collection:  worker.options.collection,
					TtlSec:      worker.options.ttlSec,
				},
				Chunks: chunks,
			},
		}

		if err := filer_pb.CreateEntry(client, request); err != nil {
			return fmt.Errorf("update fh: %v", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("upload data %v to http://%s%s%s: %v\n", fileName, worker.filerHost, task.destinationUrlPath, fileName, err)
	}

	return nil
}

func Get(grpcDialOption grpc.DialOption, filerGrpcAddress string) (masters []string, collection, replication string, maxMB uint32, cipher bool, err error) {
	err = pb.WithGrpcFilerClient(filerGrpcAddress, grpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
		if err != nil {
			return fmt.Errorf("get filer %s configuration: %v", filerGrpcAddress, err)
		}
		masters, collection, replication, maxMB = resp.Masters, resp.Collection, resp.Replication, resp.MaxMb
		cipher = resp.Cipher
		return nil
	})
	return
}
