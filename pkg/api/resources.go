package api

import "time"

/*
 * Get
 */

// OnlineGetSafe represents the response of a GET /safe/UUID
type OnlineGetSafe struct {
	// _ref         string `json:"$ref"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	UUIDRef     string `json:"uuid_ref"`
}

// OnlineGetPlatform represents the reponse of a GET /platform/UUID
type OnlineGetPlatform struct {
	// _ref       string `json:"$ref"`
	Datacenter string `json:"datacenter"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type OnlineGetSSHKey struct {
	// _ref string `json:"$ref"`
	Desc        string `json:"description"`
	Fingerprint string `json:"fingerprint"`
	UUIDRef     string `json:"uuid_ref"`
}

type OnlineGetArchive struct {
	// _ref         string `json:"$ref"`
	CreationDate string         `json:"creation_date"`
	Description  string         `json:"description"`
	Name         string         `json:"name"`
	Parity       string         `json:"parity"`
	Status       string         `json:"status"`
	UUIDRef      string         `json:"uuid_ref"`
	Size         string         `json:"size"`
	Jobs         []OnlineGetJob `json:"current_jobs,omitempty"`
	Safe         OnlineGetSafe  `json:"safe"`
}

type OnlineGetArchives []OnlineGetArchive

func (o OnlineGetArchives) Len() int {
	return len(o)
}

func (o OnlineGetArchives) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o OnlineGetArchives) Less(i, j int) bool {
	date1, _ := time.Parse(time.RFC3339, o[i].CreationDate)
	date2, _ := time.Parse(time.RFC3339, o[j].CreationDate)
	return date2.Before(date1)
}

type OnlineBucketCredentials struct {
	Login    string            `json:"login"`
	Password string            `json:"password"`
	Protocol string            `json:"protocol"`
	SSHKeys  []OnlineGetSSHKey `json:"ssh_keys"`
	URI      string            `json:"uri"`
}

type OnlineGetBucket struct {
	// _ref         string `json:"$ref"`
	ArchivalDate string                    `json:"archival_date"`
	Credentials  []OnlineBucketCredentials `json:"credentials"`
	Status       string                    `json:"status"`
	UUIDRef      string                    `json:"uuid_ref"`
}

type OnlineGetLocation struct {
	// _ref         string `json:"$ref"`
	UUIDRef string `json:"uuid_ref"`
	Name    string `json:"name"`
}

type OnlineGetJob struct {
	// _ref     string `json:"$ref"`
	Progress int    `json:"progress"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	UUIDRef  string `json:"uuid_ref"`
}

/*
 * Post
 */

type OnlinePostSafe struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type OnlinePostArchive struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parity      string   `json:"parity,omitempty"`
	Protocols   []string `json:"protocols"`
	SSHKeys     []string `json:"ssh_keys"`
	Platforms   []string `json:"platforms"`
	Days        int      `json:"days"`
	LargeBucket bool
}

type OnlinePostUnArchive struct {
	Protocols  []string `json:"protocols"`
	SSHKeys    []string `json:"ssh_keys"`
	LocationID string   `json:"location_id"`
}

type OnlinePostResult struct {
	UUIDRef string `json:"uuid_ref"`
	Archive *struct {
		UUIDRef string `json:"uuid_ref"`
		Name    string `json:"name"`
		Status  string `json:"status"`
	} `json:"archive"`
}

/*
 * Patch
 */

type OnlinePatchArchive struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
