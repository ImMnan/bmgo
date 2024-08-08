/*
Copyright © 2024 Manan Patel - github.com/immnan
*/
package find

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type inforesponseA struct {
	Result result
	Error  errorResult `json:"error"`
}
type result struct {
	Name               string `json:"name"`
	Owner              owner
	MembersCount       int `json:"membersCount"`
	WorkspacesCount    int `json:"workspacesCount"`
	Plan               plan
	CloudProviders     []string
	Credits            int             `json:"credits"`
	Expiration         int             `json:"expiration"`
	Features           map[string]bool `json:"features"`
	HasPrivateLocation bool            `json:"hasPrivateLocations"`
	IsPayingAccount    bool            `json:"isPayingAccount"`
}

type owner struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type plan struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	ReportRetention    int    `json:"reportRetention"`
	ThreadsPerEngine   int    `json:"threadsPerEngine"`
	TotalCredits       int    `json:"totalCredits"`
	Concurrency        int    `json:"concurrency"`
	Engines            int    `json:"engines"`
	PaymentServiceType string `json:"paymentServiceType"`
	MaxWorkspaces      int    `json:"maxWorkspaces"`
	MaxParrallelTests  int    `json:"maxParallelTests"`
}

func getAccountId(accountId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()
	accountIdStr := strconv.Itoa(accountId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/accounts/"+accountIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(apiId, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseObject inforesponseA
		json.Unmarshal(bodyText, &responseObject)
		if responseObject.Error.Code == 0 {
			accountName := responseObject.Result.Name
			ownerEmail := responseObject.Result.Owner.Email
			workspaceCount := responseObject.Result.WorkspacesCount
			memberCount := responseObject.Result.MembersCount
			concurrency := responseObject.Result.Plan.Concurrency
			engines := responseObject.Result.Plan.Engines
			paymentType := responseObject.Result.Plan.PaymentServiceType
			maxParrellelTests := responseObject.Result.Plan.MaxParrallelTests
			maxWorkspaces := responseObject.Result.Plan.MaxWorkspaces

			accountPlanId := responseObject.Result.Plan.Id
			//	accountPlanName := responseObject.Result.Plan.Name
			accountReportRet := responseObject.Result.Plan.ReportRetention
			accountThreadsPE := responseObject.Result.Plan.ThreadsPerEngine
			accountCredits := responseObject.Result.Credits
			accountExpiration := int64(responseObject.Result.Expiration)
			mytimeExpiration := time.Unix(accountExpiration, 0)
			expirationTimeStr := fmt.Sprint(mytimeExpiration)

			tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			// Print headers
			fmt.Fprintln(tabWriter, "NAME\tOWNER\tWORKSPACES\tUSERS")
			//fmt.Printf("\n%-25s %-35s %-15s %-10s\n", "NAME", "OWNER", "WORKSPACES", "USERS")
			//fmt.Printf("%-25s %-35s %-15d %-10d", accountName, ownerEmail, workspaceCount, memberCount)
			fmt.Fprintf(tabWriter, "%s\t%s\t%d\t%d\n", accountName, ownerEmail, workspaceCount, memberCount)
			tabWriter.Flush()
			fmt.Printf("\nPlanId: %s\nTotal Credits: %d\nReport Retention: %d\nThreads Per Engine: %d\nExpiration: %s\nMax Workspaces: %d\nMax Parallel Tests: %d\nConcurrency: %d\nEngines: %d\nPayment Service Type: %s\n",
				accountPlanId, accountCredits, accountReportRet, accountThreadsPE, expirationTimeStr[0:16],
				maxWorkspaces, maxParrellelTests, concurrency, engines, paymentType)

			cloudProviders := []string{}
			for i := 0; i < len(responseObject.Result.CloudProviders); i++ {
				cloudProlist := responseObject.Result.CloudProviders[i]
				cloudProviders = append(cloudProviders, cloudProlist)
			}
			tabWriter.Flush()
			fmt.Printf("\nSupported cloud providers: %v", cloudProviders)
			fmt.Println("\n[!] Navigate to Account settings: https://a.blazemeter.com/app/#/settings/admin/accounts/" + accountIdStr)
			fmt.Println("\n[!] Features enabled for the account:")
			for k, v := range responseObject.Result.Features {
				if v {
					fmt.Printf("%-35s %-10v\n", k, v)
				}
			}
		} else {
			errorCode := responseObject.Error.Code
			errorMessage := responseObject.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}

type inforesponseBodyWS struct {
	Result resultWS
	Error  errorResult `json:"error"`
}

type resultWS struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	MembersCount int    `json:"membersCount"`
	AccountId    int    `json:"accountId"`
}

type teamInfo struct {
	Name      string      `json:"name"`
	CreatedAt string      `json:"created_at"`
	UserCount int         `json:"user_count"`
	Buckets   int         `json:"bucket_count"`
	Owned_by  owned_by    `json:"owned_by"`
	Plan      teamPlan    `json:"plan"`
	Error     errorResult `json:"error"`
}

type owned_by struct {
	OwnerEmail string `json:"email"`
	OwnerUUID  string `json:"uuid"`
}
type teamPlan struct {
	PlanUUID          string `json:"uuid"`
	Name              string `json:"name"`
	Max_requests      int    `json:"max_requests"`
	Max_collaborators int    `json:"max_collaborators"`
	Max_buckets       int    `json:"max_buckets"`
}

func getWorkspace(workspaceId int, rawOutput bool) {
	apiId, apiSecret := Getapikeys()

	workspaceIdStr := strconv.Itoa(workspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://a.blazemeter.com/api/v4/workspaces/"+workspaceIdStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(apiId, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseObjectWS inforesponseBodyWS
		json.Unmarshal(bodyText, &responseObjectWS)
		if responseObjectWS.Error.Code == 0 {
			workspaceName := responseObjectWS.Result.Name
			members := responseObjectWS.Result.MembersCount
			accountId := responseObjectWS.Result.AccountId
			enabled := responseObjectWS.Result.Enabled

			fmt.Printf("\n%-30s %-10s %-10s %-10s\n", "NAME", "ACCOUNT", "MEMBERS", "ENABLED")
			fmt.Printf("%-30s %-10d %-10d %-10t", workspaceName, accountId, members, enabled)
			fmt.Printf("\n-")
			fmt.Println("\n[!] Navigate to Workspace settings: https://a.blazemeter.com/app/#/settings/admin/workspaces/"+workspaceIdStr, "\n-")
		} else {
			errorCode := responseObjectWS.Error.Code
			errorMessage := responseObjectWS.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}

}

func getTeamInfo(teamId string, rawOutput bool) {
	Bearer := fmt.Sprintf("Bearer %v", GetPersonalAccessToken())
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.runscope.com/teams/"+teamId, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", Bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if rawOutput {
		fmt.Printf("%s\n", bodyText)
	} else {
		var responseBodyTeamInfo teamInfo
		json.Unmarshal(bodyText, &responseBodyTeamInfo)
		if responseBodyTeamInfo.Error.Status == 0 {
			fmt.Printf("\n%-40s %-8s %-8s %-16s", "NAME", "USERS", "BUCKETS", "CREATED")
			teamName := responseBodyTeamInfo.Name
			teamCreated := responseBodyTeamInfo.CreatedAt
			teamUsers := responseBodyTeamInfo.UserCount
			teamBuckets := responseBodyTeamInfo.Buckets
			fmt.Printf("\n%-40s %-8d %-8d %-16s", teamName, teamUsers, teamBuckets, teamCreated[0:14])
			fmt.Println("\n-")
			teamOwnerEmail := responseBodyTeamInfo.Owned_by.OwnerEmail
			teamOwnerUUID := responseBodyTeamInfo.Owned_by.OwnerUUID
			teamPlanUUID := responseBodyTeamInfo.Plan.PlanUUID
			teamPlanName := responseBodyTeamInfo.Plan.Name
			teamMaxReq := responseBodyTeamInfo.Plan.Max_requests
			teamMaxbuckets := responseBodyTeamInfo.Plan.Max_buckets
			teamMaxCollaborators := responseBodyTeamInfo.Plan.Max_collaborators
			fmt.Printf("TEAM OWNER UUID: %v\nTEAM OWNER EMAIL: %v\n-", teamOwnerUUID, teamOwnerEmail)
			fmt.Printf("\nTEAM PLAN ID: %v\nTEAM PLAN NAME: %v\nTEAM MAX REQUESTS: %v\nTEAM MAX BUCKETS: %v\nTEAM MAX COLLABORATORS: %v",
				teamPlanUUID, teamPlanName, teamMaxReq, teamMaxbuckets, teamMaxCollaborators)
			fmt.Println("\n-")
		} else {
			errorCode := responseBodyTeamInfo.Error.Status
			errorMessage := responseBodyTeamInfo.Error.Message
			fmt.Printf("\nError code: %v\nError Message: %v\n\n", errorCode, errorMessage)
		}
	}
}
