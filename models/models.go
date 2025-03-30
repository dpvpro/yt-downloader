package models

import (
	"sync"
	"time"
)

type DownloadStatus string

const (
	StatusPending   DownloadStatus = "pending"
	StatusDownloading  DownloadStatus = "downloading"
	StatusConverting   DownloadStatus = "converting"
	StatusCompleted DownloadStatus = "completed"
	StatusFailed    DownloadStatus = "failed"
)

type VideoInfo struct {
	URL       string
	Title     string
	Status    DownloadStatus
	Error     string
	FileName  string
	CreatedAt time.Time
}

type DownloadRequest struct {
	ID        string
	URLs      []string
	Videos    []VideoInfo
	UseProxy  bool
	ProxyURL  string
	CreatedAt time.Time
}

type DownloadStore struct {
	mu       sync.RWMutex
	requests map[string]*DownloadRequest
}

func NewDownloadStore() *DownloadStore {
	return &DownloadStore{
		requests: make(map[string]*DownloadRequest),
	}
}

func (s *DownloadStore) Add(request *DownloadRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests[request.ID] = request
}

func (s *DownloadStore) Get(id string) (*DownloadRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	req, ok := s.requests[id]
	return req, ok
}

func (s *DownloadStore) UpdateVideoStatus(requestID string, url string, status DownloadStatus, errorMsg string, fileName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	req, ok := s.requests[requestID]
	if !ok {
		return
	}
	
	for i, video := range req.Videos {
		if video.URL == url {
			req.Videos[i].Status = status
			req.Videos[i].Error = errorMsg
			if fileName != "" {
				req.Videos[i].FileName = fileName
			}
			break
		}
	}
}