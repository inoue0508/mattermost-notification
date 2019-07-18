package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Response レスポンスの形式
type Response struct {
	Status  int
	Message string
}

// Data data
type Data struct {
	ObjectKind       string         `json:"object_kind"`
	EventType        string         `json:"event_type"`
	User             UserInfo       `json:"user"`
	Project          ProjectInfo    `json:"project"`
	ObjectAttributes AttributeInfo  `json:"object_attributes"`
	Labels           []string       `json:"labels"`
	Changes          ChangeInfo     `json:"changes"`
	Repository       RepositoryInfo `json:"repository"`
	Assignees        []AssigneeInfo `json:"assignees"`
}

//UserInfo user information
type UserInfo struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// ProjectInfo projectの情報
type ProjectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CiConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

//AttributeInfo 属性情報
type AttributeInfo struct {
	AssigneeID        int            `json:"assignee_id"`
	AuthorID          int            `json:"author_id"`
	CreatedAt         string         `json:"created_at"`
	Description       string         `json:"description"`
	HeadPipelineID    int            `json:"head_pipeline_id"`
	ID                int            `json:"id"`
	Iid               int            `json:"iid"`
	LastEditedAt      string         `json:"last_edited_at"`
	LastEditedByID    string         `json:"last_edited_by_id"`
	MergeCommitSha    string         `json:"merge_commit_sha"`
	MergeError        string         `json:"merge_error"`
	MergeParams       Merge          `json:"merge_params"`
	MergeStatus       string         `json:"merge_status"`
	MergeUserID       string         `json:"merge_user_id"`
	MergeSucceeds     bool           `json:"merge_when_pipeline_succeeds"`
	MilestoneID       string         `json:"milestone_id"`
	SourceBranch      string         `json:"source_branch"`
	SourceProjectID   int            `json:"source_project_id"`
	State             string         `json:"state"`
	TargetBranch      string         `json:"target_branch"`
	TargetProjectID   int            `json:"target_project_id"`
	TimeEstimate      int            `json:"time_estimate"`
	Title             string         `json:"title"`
	UpdateAt          string         `json:"update_at"`
	UpdateByID        string         `json:"update_by_id"`
	URL               string         `json:"url"`
	Source            SourceInfo     `json:"source"`
	Target            TargetInfo     `json:"target"`
	LastCommit        LastCommitInfo `json:"last_commit"`
	WorkInProgress    bool           `json:"work_in_progress"`
	TotalTimeSpent    int            `json:"total_time_spent"`
	HumanTotal        int            `json:"humon_total_time_spent"`
	HumanTimeEstimate string         `json:"human_time_estimate"`
	AssigneeIDs       []int          `json:"assignee_ids"`
	Action            string         `json:"action"`
}

// Merge merge_params
type Merge struct {
	ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
}

// SourceInfo ソースブランチ情報
type SourceInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CiConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

//TargetInfo ターゲットブランチ情報
type TargetInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CiConfigPath      string `json:"ci_config_path"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
}

//LastCommitInfo 最終コミット情報
type LastCommitInfo struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	TimeStamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    AuthorInfo `json:"author"`
}

//AuthorInfo 著者情報
type AuthorInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

//ChangeInfo 変更点
type ChangeInfo struct {
	AuthorID        PastNowInt    `json:"author_id"`
	CreatedAt       PastNowString `json:"created_at"`
	Description     PastNowString `json:"description"`
	ID              PastNowInt    `json:"id"`
	Iid             PastNowInt    `json:"iid"`
	MergeParams     PastNowMerge  `json:"merge_params"`
	SourceBranch    PastNowString `json:"source_branch"`
	SourceProjectID PastNowInt    `json:"source_project_id"`
	TargetBranch    PastNowString `json:"target_branch"`
	TargetProjectID PastNowInt    `json:"target_project_id"`
	Title           PastNowString `json:"title"`
	UpdateAt        PastNowString `json:"updated_at"`
	TotalTimeSpent  PastNowInt    `json:"total_time_spent"`
}

//PastNowInt 変更点数値
type PastNowInt struct {
	Previous int `json:"previous"`
	Current  int `json:"current"`
}

//PastNowString 変更点文字
type PastNowString struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

//PastNowMerge 変更点マージ情報
type PastNowMerge struct {
	Previous Merge `json:"previous"`
	Current  Merge `json:"current"`
}

//RepositoryInfo レポジトリ情報
type RepositoryInfo struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

//AssigneeInfo アサイン者情報
type AssigneeInfo struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

//Replace 置換用構造体
type Replace struct {
	Title           string
	MergeRequestURL string
	Name            string
	TargetName      string
	ProjectName     string
	ProjectURL      string
	OriginName      string
	ToID            string
}

func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {

	url := "http://10.148.152.52:10080/hooks/93xop6iha3y9pbqbthsrpt8iue"
	returnMessage := createResponse(reqBody)
	if returnMessage == "" {
		fmt.Println("対象外")

	} else {

		fmt.Printf("Request Body: %v\n", string(reqBody))
		fmt.Printf("Response Body: %v\n", string(resBody))

		var jsonstr = []byte(returnMessage)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println()
		}
		defer resp.Body.Close()
	}

}

func createResponse(reqBody []byte) string {
	reqStr := string(reqBody)

	var data Data
	err := json.Unmarshal([]byte(reqStr), &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(data)
	fmt.Println()

	var jsonStr string
	if data.ObjectAttributes.Action == "open" {
		jsonStr = `{"text": "@{{.ToID}}\n#### 新たなマージリクエストをオープンしました。確認お願いします。:bow:\n
| タイトル | 内容                                   |
|:--------|:---------------------------------------|
| リクエスト名 | [{{.Title}}]({{.MergeRequestURL}}) |
| プロジェクト名 | [{{.ProjectName}}]({{.ProjectURL}}) |
| 依頼者 | {{.Name}}                                |
| 依頼先 | {{.TargetName}}"}                        |`
	} else if data.ObjectAttributes.Action == "merge" {
		jsonStr = `{"text": "#### マージリクエストを許可しました。:tada::tada:\n
| タイトル | 内容                                   |
|:--------|:---------------------------------------|
| リクエスト名 | [{{.Title}}]({{.MergeRequestURL}}) |
| プロジェクト名 | [{{.ProjectName}}]({{.ProjectURL}}) |
| 依頼者 | {{.OriginName}}                                |
| マージ者 | {{.TargetName}}"}                        |`
	} else {
		return ""
	}

	var resultMessage bytes.Buffer
	msg, err := template.New("myTemplate").Parse(jsonStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(99)
	}

	replace := Replace{
		Title:           data.ObjectAttributes.Title,
		MergeRequestURL: data.ObjectAttributes.URL,
		Name:            data.User.Name,
		TargetName:      data.Assignees[0].Name,
		ProjectName:     data.Project.Name,
		ProjectURL:      data.Project.WebURL,
		OriginName:      data.ObjectAttributes.LastCommit.Author.Name,
		ToID:            data.Assignees[0].Username,
	}

	err = msg.Execute(&resultMessage, replace)
	return resultMessage.String()

}

func main() {

	e := echo.New()

	e.Use(middleware.BodyDump(bodyDumpHandler))

	e.POST("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{
			Status:  http.StatusOK,
			Message: "aaa",
		})
	})

	// e.GET("/", func(c echo.Context) error {
	// 	return c.JSON(http.StatusOK, Response{
	// 		Status:  http.StatusOK,
	// 		Message: "aaa",
	// 	})
	// })

	// e.GET("/api", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, controller.GetAPI())
	// })

	e.Logger.Fatal(e.Start(":32333"))
}
