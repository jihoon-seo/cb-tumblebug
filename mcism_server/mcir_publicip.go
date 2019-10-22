package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// https://github.com/cloud-barista/cb-spider/blob/master/cloud-control-manager/cloud-driver/interfaces/new-resources/PublicIPHandler.go
/*
type PublicIPReqInfo struct {
	Name         string
	KeyValueList []KeyValue
}

type PublicIPInfo struct {
	Name      string
	PublicIP  string
	OwnedVMID string
	Status    string

	KeyValueList []KeyValue
}
*/

type publicIpReq struct {
	//Id                string `json:"id"`
	ConnectionName string `json:"connectionName"`
	//CspPublicIpId     string `json:"cspPublicIpId"`
	CspPublicIpName string `json:"cspPublicIpName"`
	//PublicIp          string `json:"publicIp"`
	//OwnedVmId         string `json:"ownedVmId"`
	//ResourceGroupName string `json:"resourceGroupName"`
	Description  string     `json:"description"`
	KeyValueList []KeyValue `json:"keyValueList"`
}

type publicIpInfo struct {
	Id             string `json:"id"`
	ConnectionName string `json:"connectionName"`
	//CspPublicIpId   string `json:"cspPublicIpId"`
	CspPublicIpName string `json:"cspPublicIpName"`
	PublicIp        string `json:"publicIp"`
	OwnedVmId       string `json:"ownedVmId"`
	//ResourceGroupName string `json:"resourceGroupName"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	KeyValueList []KeyValue `json:"keyValueList"`
}

/* FYI
g.POST("/:nsId/resources/publicIp", restPostPublicIp)
g.GET("/:nsId/resources/publicIp/:publicIpId", restGetPublicIp)
g.GET("/:nsId/resources/publicIp", restGetAllPublicIp)
g.PUT("/:nsId/resources/publicIp/:publicIpId", restPutPublicIp)
g.DELETE("/:nsId/resources/publicIp/:publicIpId", restDelPublicIp)
g.DELETE("/:nsId/resources/publicIp", restDelAllPublicIp)
*/

// MCIS API Proxy: PublicIp
func restPostPublicIp(c echo.Context) error {

	nsId := c.Param("nsId")

	u := &publicIpReq{}
	if err := c.Bind(u); err != nil {
		return err
	}

	action := c.QueryParam("action")
	fmt.Println("[POST PublicIp requested action: " + action)
	if action == "create" {
		fmt.Println("[Creating PublicIp]")
		content, _ := createPublicIp(nsId, u)
		return c.JSON(http.StatusCreated, content)
		/*
			} else if action == "register" {
				fmt.Println("[Registering PublicIp]")
				content, _ := registerPublicIp(nsId, u)
				return c.JSON(http.StatusCreated, content)
		*/
	} else {
		mapA := map[string]string{"message": "You must specify: action=create"}
		return c.JSON(http.StatusFailedDependency, &mapA)
	}

}

func restGetPublicIp(c echo.Context) error {

	nsId := c.Param("nsId")

	id := c.Param("publicIpId")

	content := publicIpInfo{}

	fmt.Println("[Get publicIp for id]" + id)
	key := genResourceKey(nsId, "publicIp", id)
	fmt.Println(key)

	keyValue, _ := store.Get(key)
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===============================================")

	json.Unmarshal([]byte(keyValue.Value), &content)
	content.Id = id // Optional. Can be omitted.

	return c.JSON(http.StatusOK, &content)

}

func restGetAllPublicIp(c echo.Context) error {

	nsId := c.Param("nsId")

	var content struct {
		//Name string     `json:"name"`
		PublicIp []publicIpInfo `json:"publicIp"`
	}

	publicIpList := getPublicIpList(nsId)

	for _, v := range publicIpList {

		key := genResourceKey(nsId, "publicIp", v)
		fmt.Println(key)
		keyValue, _ := store.Get(key)
		fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
		publicIpTmp := publicIpInfo{}
		json.Unmarshal([]byte(keyValue.Value), &publicIpTmp)
		publicIpTmp.Id = v
		content.PublicIp = append(content.PublicIp, publicIpTmp)

	}
	fmt.Printf("content %+v\n", content)

	return c.JSON(http.StatusOK, &content)

}

func restPutPublicIp(c echo.Context) error {
	//nsId := c.Param("nsId")

	return nil
}

func restDelPublicIp(c echo.Context) error {

	nsId := c.Param("nsId")
	id := c.Param("publicIpId")

	err := delPublicIp(nsId, id)
	if err != nil {
		cblog.Error(err)
		mapA := map[string]string{"message": "Failed to delete the publicIp"}
		return c.JSON(http.StatusFailedDependency, &mapA)
	}

	mapA := map[string]string{"message": "The publicIp has been deleted"}
	return c.JSON(http.StatusOK, &mapA)
}

func restDelAllPublicIp(c echo.Context) error {

	nsId := c.Param("nsId")

	publicIpList := getPublicIpList(nsId)

	for _, v := range publicIpList {
		err := delPublicIp(nsId, v)
		if err != nil {
			cblog.Error(err)
			mapA := map[string]string{"message": "Failed to delete All publicIps"}
			return c.JSON(http.StatusFailedDependency, &mapA)
		}
	}

	mapA := map[string]string{"message": "All publicIps has been deleted"}
	return c.JSON(http.StatusOK, &mapA)

}

func createPublicIp(nsId string, u *publicIpReq) (publicIpInfo, error) {

	/* FYI
	type publicIpReq struct {
		//Id                string `json:"id"`
		ConnectionName string `json:"connectionName"`
		//CspPublicIpId     string `json:"cspPublicIpId"`
		CspPublicIpName string `json:"cspPublicIpName"`
		//PublicIp          string `json:"publicIp"`
		//OwnedVmId         string `json:"ownedVmId"`
		//ResourceGroupName string `json:"resourceGroupName"`
		Description string `json:"description"`
		KeyValueList []KeyValue `json:"keyValueList"`
	}
	*/

	url := SPIDER_URL + "/publicip?connection_name=" + u.ConnectionName

	method := "POST"

	//payload := strings.NewReader("{ \"Name\": \"" + u.CspPublicIpName + "\"}")
	type PublicIPReqInfo struct {
		Name         string
		KeyValueList []KeyValue
	}
	tempReq := PublicIPReqInfo{}
	tempReq.Name = u.CspPublicIpName
	tempReq.KeyValueList = u.KeyValueList
	payload, _ := json.MarshalIndent(tempReq, "", "  ")
	fmt.Println("payload: " + string(payload)) // for debug

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	fmt.Println("Called mockAPI.")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	// jhseo 191016
	//var s = new(imageInfo)
	//s := imageInfo{}
	type PublicIPInfo struct {
		Name      string
		PublicIP  string
		OwnedVMID string
		Status    string

		KeyValueList []KeyValue
	}
	temp := PublicIPInfo{}
	err2 := json.Unmarshal(body, &temp)
	if err2 != nil {
		fmt.Println("whoops:", err2)
	}

	content := publicIpInfo{}
	content.Id = genUuid()
	content.ConnectionName = u.ConnectionName
	content.CspPublicIpName = temp.Name // = u.CspPublicIpName
	content.PublicIp = temp.PublicIP
	content.OwnedVmId = temp.OwnedVMID
	content.Description = u.Description
	content.Status = temp.Status
	content.KeyValueList = temp.KeyValueList

	/* FYI
	type publicIpInfo struct {
		Id              string `json:"id"`
		ConnectionName  string `json:"connectionName"`
		//CspPublicIpId   string `json:"cspPublicIpId"`
		CspPublicIpName string `json:"cspPublicIpName"`
		PublicIp        string `json:"publicIp"`
		OwnedVmId       string `json:"ownedVmId"`
		//ResourceGroupName string `json:"resourceGroupName"`
		Description string `json:"description"`
		Status      string `json:"string"`
	}
	*/

	// cb-store
	fmt.Println("=========================== PUT createPublicIp")
	Key := genResourceKey(nsId, "publicIp", content.Id)
	/*
		mapA := map[string]string{
			"connectionName": content.ConnectionName,
			//"cspPublicIpId":     content.CspPublicIpId,
			"cspPublicIpName": content.CspPublicIpName,
			"publicIp":        content.PublicIp,
			"ownedVmId":       content.OwnedVmId,
			//"resourceGroupName": content.ResourceGroupName,
			"description": content.Description,
			"status":      content.Status}
		Val, _ := json.Marshal(mapA)
	*/
	Val, _ := json.Marshal(content)
	fmt.Println("Key: ", Key)
	fmt.Println("Val: ", Val)
	cbStorePutErr := store.Put(string(Key), string(Val))
	if cbStorePutErr != nil {
		cblog.Error(cbStorePutErr)
		return content, cbStorePutErr
	}
	keyValue, _ := store.Get(string(Key))
	fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
	fmt.Println("===========================")
	return content, nil
}

func getPublicIpList(nsId string) []string {

	fmt.Println("[Get publicIps")
	key := "/ns/" + nsId + "/resources/publicIp"
	fmt.Println(key)

	keyValue, _ := store.GetList(key, true)
	var publicIpList []string
	for _, v := range keyValue {
		//if !strings.Contains(v.Key, "vm") {
		publicIpList = append(publicIpList, strings.TrimPrefix(v.Key, "/ns/"+nsId+"/resources/publicIp/"))
		//}
	}
	for _, v := range publicIpList {
		fmt.Println("<" + v + "> \n")
	}
	fmt.Println("===============================================")
	return publicIpList

}

func delPublicIp(nsId string, Id string) error {

	fmt.Println("[Delete publicIp] " + Id)

	key := genResourceKey(nsId, "publicIp", Id)
	fmt.Println("key: " + key)

	keyValue, _ := store.Get(key)
	fmt.Println("keyValue: " + keyValue.Key + " / " + keyValue.Value)
	temp := publicIpInfo{}
	unmarshalErr := json.Unmarshal([]byte(keyValue.Value), &temp)
	if unmarshalErr != nil {
		fmt.Println("unmarshalErr:", unmarshalErr)
	}
	fmt.Println("temp.CspPublicIpName: " + temp.CspPublicIpName) // Identifier is subject to change.

	//url := "https://testapi.io/api/jihoon-seo/publicip/" + temp.CspPublicIpName + "?connection_name=" + temp.ConnectionName // for CB-Spider
	url := SPIDER_URL + "/publicip?connection_name=" + temp.ConnectionName // for testapi.io
	fmt.Println("url: " + url)

	method := "DELETE"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	fmt.Println("Called mockAPI.")
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	// delete publicIp info
	cbStoreDeleteErr := store.Delete(key)
	if cbStoreDeleteErr != nil {
		cblog.Error(cbStoreDeleteErr)
		return cbStoreDeleteErr
	}

	return nil
}
