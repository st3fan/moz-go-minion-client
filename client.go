package minion

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	CallbackEventScanState    = "scan-state"
	CallbackEventSessionState = "session-state"
)

type Client struct {
	Endpoint string
	ApiUser  string
	ApiKey   string
}

type Plan struct {
	Name        string `json:"name"`
	Created     int    `json:"created"`
	Description string `json:"description"`
}

type Site struct {
	URL     string   `json:"url"`
	Id      string   `json:"id"`
	Plans   []string `json:"plans"`
	Groups  []string `json:"groups"`
	Created int      `json:"created"`
}

type SitesResponse struct {
	Success bool   `json:"success"`
	Sites   []Site `json:"sites"`
}

// type ScanSummary struct {
// 	Id      string `json:"id"`
// 	Created int64  `json:"created"`
// 	State   string `json:"state"`
// }

type ScansResponse struct {
	Success bool   `json:"success"`
	Scans   []Scan `json:"scans"`
}

type ScanSessionIssue struct {
	Id       string `json:"id"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
}

type ScanSessionPlugin struct {
	Version string `json:"version"`
	Class   string `json:"class"`
	Weight  string `json:"weight"`
	Name    string `json:"name"`
}

type ScanSession struct {
	Id       string             `json:"id"`
	State    string             `json:"state"`
	Created  int64              `json:"created"`
	Started  int64              `json:"started"`
	Finished int64              `json:"finished"`
	Issues   []ScanSessionIssue `json:"issues"`
	Plugin   ScanSessionPlugin  `json:"plugin"`
}

type ScanConfiguration struct {
	Target string `json:"target"`
}

type Scan struct {
	Id            string            `json:"id"`
	State         string            `json:"state"`
	Created       int64             `json:"created"`
	Started       int64             `json:"started"`
	Finished      int64             `json:"finished"`
	Configuration ScanConfiguration `json:"configuration"`
	Sessions      []ScanSession     `json:"sessions"`
}

type ScanResponse struct {
	Success bool `json:"success"`
	Scan    Scan `json:"scan"`
}

type CreateSiteResponse struct {
	Success bool `json:"success"`
	Site    Site `json:"site"`
}

type GetPlansResponse struct {
	Success bool   `json:"success"`
	Plans   []Plan `json:"plans"`
}

func NewClient(endpoint, apiUser, apiKey string) (Client, error) {
	return Client{Endpoint: endpoint, ApiUser: apiUser, ApiKey: apiKey}, nil
}

func (c *Client) GetSites() ([]Site, error) {
	url := c.Endpoint + "/ws/sites"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sitesResponse := &SitesResponse{}
	if err = json.Unmarshal(body, sitesResponse); err != nil {
		return nil, err
	}

	return sitesResponse.Sites, nil
}

func (c *Client) GetScans(siteId string, plan string, limit int) ([]Scan, error) {
	url := c.Endpoint + "/ws/scans?site_id=" + siteId + "&plan_name=" + plan + "&limit=" + strconv.Itoa(limit)

	log.Printf("url is %+v", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	scansResponse := &ScansResponse{}
	if err = json.Unmarshal(body, scansResponse); err != nil {
		return nil, err
	}

	return scansResponse.Scans, nil
}

func (c *Client) GetScan(scanId string) (*Scan, error) {
	url := c.Endpoint + "/ws/scans/" + scanId

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	scanResponse := &ScanResponse{}
	if err = json.Unmarshal(body, scanResponse); err != nil {
		return nil, err
	}

	return &scanResponse.Scan, nil
}

func (c *Client) GetSitesByURL(siteURL string) ([]Site, error) {
	u := c.Endpoint + "/ws/sites?url=" + url.QueryEscape(siteURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sitesResponse := &SitesResponse{}
	if err = json.Unmarshal(body, sitesResponse); err != nil {
		return nil, err
	}

	return sitesResponse.Sites, nil
}

func (c *Client) CreateSite(siteTemplate Site) (Site, error) {
	encoded, err := json.Marshal(siteTemplate)
	if err != nil {
		return Site{}, err
	}

	u := c.Endpoint + "/ws/sites"

	client := &http.Client{}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(encoded))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return Site{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Site{}, err
	}

	createSiteResponse := &CreateSiteResponse{}
	if err = json.Unmarshal(body, createSiteResponse); err != nil {
		return Site{}, err
	}

	return createSiteResponse.Site, nil
}

func (c *Client) GetPlanByName(planName string) (Plan, error) {
	u := c.Endpoint + "/ws/plans?name=" + url.QueryEscape(planName)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return Plan{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Plan{}, err
	}

	getPlansResponse := &GetPlansResponse{}
	if err = json.Unmarshal(body, getPlansResponse); err != nil {
		return Plan{}, err
	}

	return getPlansResponse.Plans[0], nil
}

type CreateScanRequest struct {
	Plan        string `json:"planName"`
	SiteId      string `json:"siteId"`
	CallbackURL string `json:"callbackURL"`
}

type CreateScanResponse struct {
	Success bool `json:"success"`
	Scan    Scan `json:"scan"`
}

func (c *Client) CreateScan(siteId, planName, callbackURL string) (Scan, error) {
	createScanRequest := CreateScanRequest{
		Plan:        planName,
		SiteId:      siteId,
		CallbackURL: callbackURL,
	}

	encoded, err := json.Marshal(createScanRequest)
	if err != nil {
		return Scan{}, err
	}

	u := c.Endpoint + "/ws/scans/create"

	client := &http.Client{}
	req, err := http.NewRequest("PUT", u, bytes.NewBuffer(encoded))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Minion-Api-User", c.ApiUser)
	req.Header.Add("X-Minion-Api-Key", c.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return Scan{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Scan{}, err
	}

	createScanResponse := &CreateScanResponse{}
	if err = json.Unmarshal(body, createScanResponse); err != nil {
		return Scan{}, err
	}

	return createScanResponse.Scan, nil
}
