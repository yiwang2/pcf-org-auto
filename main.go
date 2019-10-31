package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/rest"
	"os"
	"strconv"
	"strings"
	"time"
)

const org_name_prefix string = "HelloTestCreating800Org";

func main () {

	var cmd_create_800 string = "CREATE_800";
	var cmd_delete_800 string = "DELETE_800";
	var cmd_list string = "LIST";

	if os.Args[1] == cmd_create_800 {
		create800Orgs()
	} else if os.Args[1] == cmd_list {
		fmt.Println("List Orgs")
		list_organizations()
	} else if os.Args[1] == cmd_delete_800 {
		delete800Orgs()
	}
}

func create800Orgs () {
	const num_host int = 800;
	for i:=0; i < num_host; i++ {
		var org_name = org_name_prefix + strconv.Itoa(i)
		fmt.Println(org_name);
		create_organization(org_name);
		time.Sleep(500)
	}
}

func delete800Orgs () {
	var org_infos = get_organization_infos();
	for _, info :=range org_infos {
		if (strings.Index(info.Name, org_name_prefix) >=0) {
			delete_organization_with_guid(info.GUID);
			time.Sleep(500)
		}
	}
}

type Orgnaization struct {
	Name string `json:"name"`
	GUID string `json:"guid"`
}

type OrgnaizationResponse struct {
	Total_results int `json:"total_results"`
	Total_pages int `json:"total_pages"`
	Prev_url int `json:"prev_url"`
	Next_url int `json:"next_url"`
	Resources []*OrgnaizationResource `json:"resources"`
}

type OrgnaizationMetadata struct {
	GUID string `json:"guid"`
}

type OrgnaizationEntity struct {
	Name string `json:"name"`
}

type OrgnaizationResource struct {
	Metadata *OrgnaizationMetadata `json:"metadata"`
	Entity *OrgnaizationEntity `json:"entity"`
}


func delete_organization_with_guid(organization_guid string) {
	var end_point = "/v2/organizations/"+organization_guid
	var http_method = rest.Delete
	Headers := make(map[string]string)
	Headers["Authorization"] = getAuthoriziation()
	base_URL := getEndPoint() + end_point;

	request := rest.Request{
		Method:  http_method,
		BaseURL: base_URL,
		Headers:Headers,
	}

	response, err := rest.Send(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func create_organization (organization_name string) {
	var end_point = "/v2/organizations"
	var http_method = rest.Post
	Headers := make(map[string]string)
	Headers["Authorization"] = getAuthoriziation()
	base_URL := getEndPoint() + end_point;

	orgnaization_body := &Orgnaization{Name: organization_name}
	orgnaization_body_json, err := json.Marshal(orgnaization_body)
	var Body = []byte(string(orgnaization_body_json))
	request := rest.Request{
		Method:  http_method,
		BaseURL: base_URL,
		Headers:Headers,
		Body:Body,
	}

	response, err := rest.Send(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func list_organizations () {

	var orgnaizationResponse *OrgnaizationResponse = list_organization_result();
	org_resources := orgnaizationResponse.Resources
	for _, resource := range org_resources {
		fmt.Println(resource.Entity.Name)
	}
}

func list_organization_result () *OrgnaizationResponse{
	var end_point = "/v2/organizations"
	var http_method = rest.Get
	Headers := make(map[string]string)
	Headers["Authorization"] = getAuthoriziation()

	base_URL := getEndPoint() + end_point;
	request := rest.Request{
		Method:  http_method,
		BaseURL: base_URL,
		Headers:Headers,
	}

	response, err := rest.Send(request)
	if err != nil {
		fmt.Println(err)
		return nil;
	} else {
		organizations_response := []byte(response.Body)
		var organizations_response_struct OrgnaizationResponse
		json.Unmarshal(organizations_response, &organizations_response_struct)
		return &organizations_response_struct;

	}
}

func get_organization_infos() []*Orgnaization {

	var org_results []*Orgnaization

	var orgnaizationResponse *OrgnaizationResponse = list_organization_result();
	org_resources := orgnaizationResponse.Resources;
	for _, resource := range org_resources {
		org_result := &Orgnaization{};
		org_result.Name = resource.Entity.Name;
		org_result.GUID = resource.Metadata.GUID;
		org_results = append(org_results, org_result)
	}
	return org_results;
}

func getEndPoint () string{
	return os.Getenv("API_END_POINT")
}

func getAuthoriziation () string{
	return "bearer "+os.Getenv("API_KEY")
}