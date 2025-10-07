package models

type SyncResponse struct {
	Tasks                   []Task                   `json:"tasks,omitempty"`
	Tags                    []Tag                    `json:"tags,omitempty"`
	Spaces                  []Space                  `json:"spaces,omitempty"`
	RepetitiveTaskTemplates []RepetitiveTaskTemplate `json:"repetitiveTaskTemplates,omitempty"`
	LatestChangeID          int64                    `json:"latestChangeId"`
}
