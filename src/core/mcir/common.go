package mcir

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	//uuid "github.com/google/uuid"
	"github.com/cloud-barista/cb-spider/interface/api"
	"github.com/cloud-barista/cb-tumblebug/src/core/common"
	"github.com/go-resty/resty/v2"
	"github.com/xwb1989/sqlparser"

	// CB-Store
	cbstore_utils "github.com/cloud-barista/cb-store/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// CB-Store
//var cblog *logrus.Logger
//var store icbs.Store

//var SPIDER_REST_URL string

func init() {
	//cblog = config.Cblogger
	//store = cbstore.GetStore()
	//SPIDER_REST_URL = os.Getenv("SPIDER_REST_URL")
}

// DelAllResources deletes all TB MCIR object of given resourceType
func DelAllResources(nsId string, resourceType string, forceFlag string) error {

	nsId = common.ToLower(nsId)

	resourceIdList := ListResourceId(nsId, resourceType)

	if len(resourceIdList) == 0 {
		return nil
	}

	for _, v := range resourceIdList {
		err := DelResource(nsId, resourceType, v, forceFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

// DelResource deletes the TB MCIR object
func DelResource(nsId string, resourceType string, resourceId string, forceFlag string) error {

	//fmt.Println("[Delete " + resourceType + "] " + resourceId)
	fmt.Printf("DelResource() called; %s %s %s \n", nsId, resourceType, resourceId) // for debug

	//check, lowerizedResourceId, err := LowerizeAndCheckResource(nsId, resourceType, resourceId)
	//resourceId = lowerizedResourceId
	nsId = common.ToLower(nsId)
	resourceId = common.ToLower(resourceId)
	check, err := CheckResource(nsId, resourceType, resourceId)

	if !check {
		errString := "The " + resourceType + " " + resourceId + " does not exist."
		//mapA := map[string]string{"message": errString}
		//mapB, _ := json.Marshal(mapA)
		err := fmt.Errorf(errString)
		//return http.StatusNotFound, mapB, err
		return err
	}

	if err != nil {
		common.CBLog.Error(err)
		return err
	}

	key := common.GenResourceKey(nsId, resourceType, resourceId)
	fmt.Println("key: " + key)

	keyValue, _ := common.CBStore.Get(key)
	/*
		if keyValue == nil {
			mapA := map[string]string{"message": "Failed to find the resource with given ID."}
			mapB, _ := json.Marshal(mapA)
			err := fmt.Errorf("Failed to find the resource with given ID.")
			return http.StatusNotFound, mapB, err
		}
	*/
	//fmt.Println("keyValue: " + keyValue.Key + " / " + keyValue.Value)

	//cspType := common.GetResourcesCspType(nsId, resourceType, resourceId)

	if os.Getenv("SPIDER_CALL_METHOD") == "REST" {

		var url string

		// Create Req body
		type JsonTemplate struct {
			ConnectionName string
		}
		tempReq := JsonTemplate{}

		switch resourceType {
		case common.StrImage:
			// delete image info
			err := common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				//return http.StatusInternalServerError, nil, err
				return err
			}

			sql := "DELETE FROM `image` WHERE `id` = '" + resourceId + "';"
			fmt.Println("sql: " + sql)
			// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
			_, err = sqlparser.Parse(sql)
			if err != nil {
				common.CBLog.Error(err)
			}

			stmt, err := common.MYDB.Prepare(sql)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = stmt.Exec()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Data deleted successfully..")
			}

			//return http.StatusOK, nil, nil
			return nil
		case common.StrSpec:
			// delete spec info

			//get related recommend spec
			//keyValue, err := common.CBStore.Get(key)
			content := TbSpecInfo{}
			err := json.Unmarshal([]byte(keyValue.Value), &content)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			err = common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			//delete related recommend spec
			err = DelRecommendSpec(nsId, resourceId, content.Num_vCPU, content.Mem_GiB, content.Storage_GiB)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			sql := "DELETE FROM `spec` WHERE `id` = '" + resourceId + "';"
			fmt.Println("sql: " + sql)
			// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
			_, err = sqlparser.Parse(sql)
			if err != nil {
				common.CBLog.Error(err)
			}

			stmt, err := common.MYDB.Prepare(sql)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = stmt.Exec()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Data deleted successfully..")
			}

			//return http.StatusOK, nil, nil
			return nil
		case common.StrSSHKey:
			temp := TbSshKeyInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}
			tempReq.ConnectionName = temp.ConnectionName
			url = common.SPIDER_REST_URL + "/keypair/" + temp.Name //+ "?connection_name=" + temp.ConnectionName
		case common.StrVNet:
			temp := TbVNetInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}
			tempReq.ConnectionName = temp.ConnectionName
			url = common.SPIDER_REST_URL + "/vpc/" + temp.Name //+ "?connection_name=" + temp.ConnectionName
		case common.StrSecurityGroup:
			temp := TbSecurityGroupInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}
			tempReq.ConnectionName = temp.ConnectionName
			url = common.SPIDER_REST_URL + "/securitygroup/" + temp.Name //+ "?connection_name=" + temp.ConnectionName
		/*
			case "subnet":
				temp := subnetInfo{}
				json.Unmarshal([]byte(keyValue.Value), &content)
				return content.CspSubnetId
			case "publicIp":
				temp := publicIpInfo{}
				json.Unmarshal([]byte(keyValue.Value), &temp)
				tempReq.ConnectionName = temp.ConnectionName
				url = common.SPIDER_REST_URL + "/publicip/" + temp.CspPublicIpName //+ "?connection_name=" + temp.ConnectionName
			case "vNic":
				temp := vNicInfo{}
				json.Unmarshal([]byte(keyValue.Value), &temp)
				tempReq.ConnectionName = temp.ConnectionName
				url = common.SPIDER_REST_URL + "/vnic/" + temp.CspVNicName //+ "?connection_name=" + temp.ConnectionName
		*/
		default:
			err := fmt.Errorf("invalid resourceType")
			//return http.StatusBadRequest, nil, err
			return err
		}

		fmt.Println("url: " + url)

		client := resty.New()

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(tempReq).
			//SetResult(&SpiderSpecInfo{}). // or SetResult(AuthSuccess{}).
			//SetError(&AuthError{}).       // or SetError(AuthError{}).
			Delete(url)

		if err != nil {
			common.CBLog.Error(err)
			err := fmt.Errorf("an error occurred while requesting to CB-Spider")
			return err
		}

		/*
			if res.StatusCode == 400 || res.StatusCode == 401 {
				fmt.Println("HTTP Status code 400 Bad Request or 401 Unauthorized.")
				err := fmt.Errorf("HTTP Status code 400 Bad Request or 401 Unauthorized")
				common.CBLog.Error(err)
				return res, err
			}

			// delete vNet info
			err := common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				return res, err
			}

			return res, nil
		*/

		fmt.Println("HTTP Status code " + strconv.Itoa(resp.StatusCode()))
		switch {
		case forceFlag == "true":
			url += "?force=true"
			fmt.Println("forceFlag == true; url: " + url)

			_, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(tempReq).
				//SetResult(&SpiderSpecInfo{}). // or SetResult(AuthSuccess{}).
				//SetError(&AuthError{}).       // or SetError(AuthError{}).
				Delete(url)

			if err != nil {
				common.CBLog.Error(err)
				err := fmt.Errorf("an error occurred while requesting to CB-Spider")
				return err
			}

			err = common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				//return res.StatusCode, body, err
				return err
			}
			//return res.StatusCode, body, nil
			return nil
		case resp.StatusCode() >= 400 || resp.StatusCode() < 200:
			err := fmt.Errorf(string(resp.Body()))
			common.CBLog.Error(err)
			//return res.StatusCode, body, err
			return err
		default:
			err := common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				//return res.StatusCode, body, err
				return err
			}
			//return res.StatusCode, body, nil
			return nil
		}

	} else {

		// CCM API 설정
		ccm := api.NewCloudResourceHandler()
		err := ccm.SetConfigPath(os.Getenv("CBTUMBLEBUG_ROOT") + "/conf/grpc_conf.yaml")
		if err != nil {
			common.CBLog.Error("ccm failed to set config : ", err)
			return err
		}
		err = ccm.Open()
		if err != nil {
			common.CBLog.Error("ccm api open failed : ", err)
			return err
		}
		defer ccm.Close()

		switch resourceType {
		case common.StrImage:
			// delete image info
			err := common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				//return http.StatusInternalServerError, nil, err
				return err
			}

			sql := "DELETE FROM `image` WHERE `id` = '" + resourceId + "';"
			fmt.Println("sql: " + sql)
			// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
			_, err = sqlparser.Parse(sql)
			if err != nil {
				common.CBLog.Error(err)
			}

			stmt, err := common.MYDB.Prepare(sql)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = stmt.Exec()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Data deleted successfully..")
			}

			//return http.StatusOK, nil, nil
			return nil
		case common.StrSpec:
			// delete spec info

			//get related recommend spec
			content := TbSpecInfo{}
			err := json.Unmarshal([]byte(keyValue.Value), &content)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			err = common.CBStore.Delete(key)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			//delete related recommend spec
			err = DelRecommendSpec(nsId, resourceId, content.Num_vCPU, content.Mem_GiB, content.Storage_GiB)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			sql := "DELETE FROM `spec` WHERE `id` = '" + resourceId + "';"
			fmt.Println("sql: " + sql)
			// https://stackoverflow.com/questions/42486032/golang-sql-query-syntax-validator
			_, err = sqlparser.Parse(sql)
			if err != nil {
				common.CBLog.Error(err)
			}

			stmt, err := common.MYDB.Prepare(sql)
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = stmt.Exec()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Data deleted successfully..")
			}
			return nil

		case common.StrSSHKey:
			temp := TbSshKeyInfo{}
			err := json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			_, err = ccm.DeleteKeyByParam(temp.ConnectionName, temp.Name, forceFlag)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

		case common.StrVNet:
			temp := TbVNetInfo{}
			err := json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			_, err = ccm.DeleteVPCByParam(temp.ConnectionName, temp.Name, forceFlag)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

		case common.StrSecurityGroup:
			temp := TbSecurityGroupInfo{}
			err := json.Unmarshal([]byte(keyValue.Value), &temp)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

			_, err = ccm.DeleteSecurityByParam(temp.ConnectionName, temp.Name, forceFlag)
			if err != nil {
				common.CBLog.Error(err)
				return err
			}

		default:
			err := fmt.Errorf("invalid resourceType")
			return err
		}

		err = common.CBStore.Delete(key)
		if err != nil {
			common.CBLog.Error(err)
			return err
		}
		return nil

	}
}

type SpiderNameIdSystemId struct {
	NameId   string
	SystemId string
}

type SpiderAllListWrapper struct {
	AllList SpiderAllList
}

type SpiderAllList struct {
	MappedList     []SpiderNameIdSystemId
	OnlySpiderList []SpiderNameIdSystemId
	OnlyCSPList    []SpiderNameIdSystemId
}

// Response struct for InspectResources
type TbInspectResourcesResponse struct {
	// ResourcesOnCsp       interface{} `json:"resourcesOnCsp"`
	// ResourcesOnSpider    interface{} `json:"resourcesOnSpider"`
	// ResourcesOnTumblebug interface{} `json:"resourcesOnTumblebug"`
	ResourcesOnCsp       []resourceOnCspOrSpider `json:"resourcesOnCsp"`
	ResourcesOnSpider    []resourceOnCspOrSpider `json:"resourcesOnSpider"`
	ResourcesOnTumblebug []resourceOnTumblebug   `json:"resourcesOnTumblebug"`
}

type resourceOnCspOrSpider struct {
	Id          string `json:"id"`
	CspNativeId string `json:"cspNativeId"`
}

type resourceOnTumblebug struct {
	Id          string `json:"id"`
	CspNativeId string `json:"cspNativeId"`
	NsId        string `json:"nsId"`
	McisId      string `json:"mcisId"`
	Type        string `json:"type"`
	ObjectKey   string `json:"objectKey"`
}

// ListResourceStatus returns the state list of TB MCIR objects of given resourceType
func ListResourceStatus(connConfig string, resourceType string) (interface{}, error) {

	nsList := common.ListNsId()
	// var TbResourceList []string
	var TbResourceList []resourceOnTumblebug
	for _, ns := range nsList {
		/*
			resourceListInNs := ListResourceId(ns, resourceType)
			for i, _ := range resourceListInNs {
				resourceListInNs[i] = ns + "/" + resourceListInNs[i]
			}
			TbResourceList = append(TbResourceList, resourceListInNs...)
		*/

		resourceListInNs, err := ListResource(ns, resourceType)
		if err != nil {
			common.CBLog.Error(err)
			err := fmt.Errorf("an error occurred while getting resource list")
			return nil, err
		}
		if resourceListInNs == nil {
			continue
		}

		switch resourceType {
		case common.StrVNet:
			resourcesInNs := resourceListInNs.([]TbVNetInfo) // type assertion
			for _, resource := range resourcesInNs {
				temp := resourceOnTumblebug{}
				temp.Id = resource.Id
				temp.CspNativeId = resource.CspVNetId
				temp.NsId = ns
				//temp.McisId = ""
				temp.Type = resourceType
				temp.ObjectKey = common.GenResourceKey(ns, resourceType, resource.Id)

				TbResourceList = append(TbResourceList, temp)
			}
		case common.StrSecurityGroup:
			resourcesInNs := resourceListInNs.([]TbSecurityGroupInfo) // type assertion
			for _, resource := range resourcesInNs {
				temp := resourceOnTumblebug{}
				temp.Id = resource.Id
				temp.CspNativeId = resource.CspSecurityGroupId
				temp.NsId = ns
				//temp.McisId = ""
				temp.Type = resourceType
				temp.ObjectKey = common.GenResourceKey(ns, resourceType, resource.Id)

				TbResourceList = append(TbResourceList, temp)
			}
		case common.StrSSHKey:
			resourcesInNs := resourceListInNs.([]TbSshKeyInfo) // type assertion
			for _, resource := range resourcesInNs {
				temp := resourceOnTumblebug{}
				temp.Id = resource.Id
				temp.CspNativeId = resource.CspSshKeyName
				temp.NsId = ns
				//temp.McisId = ""
				temp.Type = resourceType
				temp.ObjectKey = common.GenResourceKey(ns, resourceType, resource.Id)

				TbResourceList = append(TbResourceList, temp)
			}
		}

	}

	client := resty.New()
	client.SetAllowGetMethodPayload(true)

	// Create Req body
	type JsonTemplate struct {
		ConnectionName string
	}
	tempReq := JsonTemplate{}
	tempReq.ConnectionName = connConfig

	var spiderRequestURL string
	switch resourceType {
	case common.StrVNet:
		spiderRequestURL = common.SPIDER_REST_URL + "/allvpc"
	case common.StrSecurityGroup:
		spiderRequestURL = common.SPIDER_REST_URL + "/allsecuritygroup"
	case common.StrSSHKey:
		spiderRequestURL = common.SPIDER_REST_URL + "/allkeypair"
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(tempReq).
		SetResult(&SpiderAllListWrapper{}). // or SetResult(AuthSuccess{}).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).
		Get(spiderRequestURL)

	if err != nil {
		common.CBLog.Error(err)
		err := fmt.Errorf("an error occurred while requesting to CB-Spider")
		return nil, err
	}

	fmt.Println("HTTP Status code " + strconv.Itoa(resp.StatusCode()))
	switch {
	case resp.StatusCode() >= 400 || resp.StatusCode() < 200:
		err := fmt.Errorf(string(resp.Body()))
		common.CBLog.Error(err)
		//return res.StatusCode, body, err
		return nil, err
	default:
	}

	temp, _ := resp.Result().(*SpiderAllListWrapper) // type assertion

	result := TbInspectResourcesResponse{}
	result.ResourcesOnTumblebug = TbResourceList
	// result.ResourcesOnCsp = append((*temp).AllList.MappedList, (*temp).AllList.OnlyCSPList...)
	// result.ResourcesOnSpider = append((*temp).AllList.MappedList, (*temp).AllList.OnlySpiderList...)

	for _, v := range (*temp).AllList.MappedList {
		tmpObj := resourceOnCspOrSpider{}
		tmpObj.Id = v.NameId
		tmpObj.CspNativeId = v.SystemId

		result.ResourcesOnCsp = append(result.ResourcesOnCsp, tmpObj)
		result.ResourcesOnSpider = append(result.ResourcesOnSpider, tmpObj)
	}

	for _, v := range (*temp).AllList.OnlySpiderList {
		tmpObj := resourceOnCspOrSpider{}
		tmpObj.Id = v.NameId
		tmpObj.CspNativeId = v.SystemId

		result.ResourcesOnSpider = append(result.ResourcesOnSpider, tmpObj)
	}

	for _, v := range (*temp).AllList.OnlyCSPList {
		tmpObj := resourceOnCspOrSpider{}
		tmpObj.Id = v.NameId
		tmpObj.CspNativeId = v.SystemId

		result.ResourcesOnCsp = append(result.ResourcesOnCsp, tmpObj)
	}

	return result, nil
}

// ListResourceId returns the list of TB MCIR object IDs of given resourceType
func ListResourceId(nsId string, resourceType string) []string {

	nsId = common.ToLower(nsId)

	if resourceType == common.StrImage ||
		resourceType == common.StrSSHKey ||
		resourceType == common.StrSpec ||
		resourceType == common.StrVNet ||
		//resourceType == "subnet" ||
		//resourceType == "publicIp" ||
		//resourceType == "vNic" ||
		resourceType == common.StrSecurityGroup {
		// continue
	} else {
		return []string{"invalid resource type"}
	}

	fmt.Println("[Get " + resourceType + " list")
	key := "/ns/" + nsId + "/resources/" + resourceType
	fmt.Println(key)

	keyValue, _ := common.CBStore.GetList(key, true)

	var resourceList []string
	for _, v := range keyValue {
		//if !strings.Contains(v.Key, "vm") {
		resourceList = append(resourceList, strings.TrimPrefix(v.Key, "/ns/"+nsId+"/resources/"+resourceType+"/"))
		//}
	}
	for _, v := range resourceList {
		fmt.Println("<" + v + "> \n")
	}
	fmt.Println("===============================================")
	return resourceList

}

// ListResource returns the list of TB MCIR objects of given resourceType
func ListResource(nsId string, resourceType string) (interface{}, error) {

	nsId = common.ToLower(nsId)

	if resourceType == common.StrImage ||
		resourceType == common.StrSSHKey ||
		resourceType == common.StrSpec ||
		resourceType == common.StrVNet ||
		//resourceType == "subnet" ||
		//resourceType == "publicIp" ||
		//resourceType == "vNic" ||
		resourceType == common.StrSecurityGroup {
		// continue
	} else {
		errString := "Cannot list " + resourceType + "s."
		err := fmt.Errorf(errString)
		return nil, err
	}

	fmt.Println("[Get " + resourceType + " list")
	key := "/ns/" + nsId + "/resources/" + resourceType
	fmt.Println(key)

	keyValue, err := common.CBStore.GetList(key, true)
	keyValue = cbstore_utils.GetChildList(keyValue, key)

	if err != nil {
		common.CBLog.Error(err)
		/*
			fmt.Println("func ListResource; common.CBStore.GetList gave error")
			var resourceList []string
			for _, v := range keyValue {
				resourceList = append(resourceList, strings.TrimPrefix(v.Key, "/ns/"+nsId+"/resources/"+resourceType+"/"))
			}
			for _, v := range resourceList {
				fmt.Println("<" + v + "> \n")
			}
			fmt.Println("===============================================")
		*/
		return nil, err
	}
	if keyValue != nil {
		switch resourceType {
		case common.StrImage:
			res := []TbImageInfo{}
			for _, v := range keyValue {
				tempObj := TbImageInfo{}
				err = json.Unmarshal([]byte(v.Value), &tempObj)
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
				res = append(res, tempObj)
			}
			return res, nil
		case common.StrSecurityGroup:
			res := []TbSecurityGroupInfo{}
			for _, v := range keyValue {
				tempObj := TbSecurityGroupInfo{}
				err = json.Unmarshal([]byte(v.Value), &tempObj)
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
				res = append(res, tempObj)
			}
			return res, nil
		case common.StrSpec:
			res := []TbSpecInfo{}
			for _, v := range keyValue {
				tempObj := TbSpecInfo{}
				err = json.Unmarshal([]byte(v.Value), &tempObj)
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
				res = append(res, tempObj)
			}
			return res, nil
		case common.StrSSHKey:
			res := []TbSshKeyInfo{}
			for _, v := range keyValue {
				tempObj := TbSshKeyInfo{}
				err = json.Unmarshal([]byte(v.Value), &tempObj)
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
				res = append(res, tempObj)
			}
			return res, nil
		case common.StrVNet:
			res := []TbVNetInfo{}
			for _, v := range keyValue {
				tempObj := TbVNetInfo{}
				err = json.Unmarshal([]byte(v.Value), &tempObj)
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
				res = append(res, tempObj)
			}
			return res, nil
		}

		//return true, nil
	}

	return nil, nil // When err == nil && keyValue == nil
}

// GetAssociatedObjectCount returns the number of MCIR's associated Tumblebug objects
func GetAssociatedObjectCount(nsId string, resourceType string, resourceId string) (int, error) {
	nsId = common.ToLower(nsId)
	resourceId = common.ToLower(resourceId)
	check, err := CheckResource(nsId, resourceType, resourceId)

	if !check {
		errString := "The " + resourceType + " " + resourceId + " does not exist."
		//mapA := map[string]string{"message": errString}
		//mapB, _ := json.Marshal(mapA)
		err := fmt.Errorf(errString)
		return -1, err
	}

	if err != nil {
		common.CBLog.Error(err)
		return -1, err
	}
	fmt.Println("[Get count] " + resourceType + ", " + resourceId)

	key := common.GenResourceKey(nsId, resourceType, resourceId)
	//fmt.Println(key)

	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		common.CBLog.Error(err)
		return -1, err
	}
	if keyValue != nil {
		inUseCount := int(gjson.Get(keyValue.Value, "associatedObjectList.#").Int())
		return inUseCount, nil
	}
	errString := "Cannot get " + resourceType + " " + resourceId + "."
	err = fmt.Errorf(errString)
	return -1, err
}

// GetAssociatedObjectList returns the list of MCIR's associated Tumblebug objects
func GetAssociatedObjectList(nsId string, resourceType string, resourceId string) ([]string, error) {

	var result []string

	//check, lowerizedResourceId, err := LowerizeAndCheckResource(nsId, resourceType, resourceId)
	//resourceId = lowerizedResourceId
	nsId = common.ToLower(nsId)
	resourceId = common.ToLower(resourceId)
	check, err := CheckResource(nsId, resourceType, resourceId)

	if !check {
		errString := "The " + resourceType + " " + resourceId + " does not exist."
		//mapA := map[string]string{"message": errString}
		//mapB, _ := json.Marshal(mapA)
		err := fmt.Errorf(errString)
		return nil, err
	}

	if err != nil {
		common.CBLog.Error(err)
		return nil, err
	}
	fmt.Println("[Get count] " + resourceType + ", " + resourceId)

	key := common.GenResourceKey(nsId, resourceType, resourceId)
	//fmt.Println(key)

	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		common.CBLog.Error(err)
		return nil, err
	}
	if keyValue != nil {
		/*
			objList := gjson.Get(keyValue.Value, "associatedObjectList")
			objList.ForEach(func(key, value gjson.Result) bool {
				result = append(result, value.String())
				return true
			})
		*/

		/*
			switch resourceType {
			case common.StrImage:
				res := TbImageInfo{}
				json.Unmarshal([]byte(keyValue.Value), &res)
				//result = res.
			case common.StrSecurityGroup:
				res := TbSecurityGroupInfo{}
				json.Unmarshal([]byte(keyValue.Value), &res)

			case common.StrSpec:
				res := TbSpecInfo{}
				json.Unmarshal([]byte(keyValue.Value), &res)

			case common.StrSSHKey:
				res := TbSshKeyInfo{}
				json.Unmarshal([]byte(keyValue.Value), &res)
				result = res.AssociatedObjectList
			case common.StrVNet:
				res := TbVNetInfo{}
				json.Unmarshal([]byte(keyValue.Value), &res)

			}
		*/

		type stringList struct {
			AssociatedObjectList []string `json:"associatedObjectList"`
		}
		res := stringList{}
		err = json.Unmarshal([]byte(keyValue.Value), &res)
		if err != nil {
			common.CBLog.Error(err)
			return nil, err
		}
		result = res.AssociatedObjectList

		return result, nil
	}
	errString := "Cannot get " + resourceType + " " + resourceId + "."
	err = fmt.Errorf(errString)
	return nil, err
}

// UpdateAssociatedObjectList adds or deletes the objectKey (currently, vmKey) to/from TB object's associatedObjectList
func UpdateAssociatedObjectList(nsId string, resourceType string, resourceId string, cmd string, objectKey string) ([]string, error) {

	nsId = common.ToLower(nsId)
	resourceId = common.ToLower(resourceId)
	/*
		check, err := CheckResource(nsId, resourceType, resourceId)

		if !check {
			errString := "The " + resourceType + " " + resourceId + " does not exist."
			//mapA := map[string]string{"message": errString}
			//mapB, _ := json.Marshal(mapA)
			err := fmt.Errorf(errString)
			return -1, err
		}

		if err != nil {
			common.CBLog.Error(err)
			return -1, err
		}
	*/
	fmt.Println("[Set count] " + resourceType + ", " + resourceId)

	key := common.GenResourceKey(nsId, resourceType, resourceId)
	//fmt.Println(key)

	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		common.CBLog.Error(err)
		return nil, err
	}

	if keyValue != nil {
		objList, _ := GetAssociatedObjectList(nsId, resourceType, resourceId)
		switch cmd {
		case common.StrAdd:
			for _, v := range objList {
				if v == objectKey {
					errString := objectKey + " is already associated with " + resourceType + " " + resourceId + "."
					err = fmt.Errorf(errString)
					return nil, err
				}
			}
			// fmt.Println("len(objList): " + strconv.Itoa(len(objList))) // for debug
			// fmt.Print("objList: ")                                     // for debug
			// fmt.Println(objList)                                       // for debug

			var anyJson map[string]interface{}
			json.Unmarshal([]byte(keyValue.Value), &anyJson)
			if anyJson["associatedObjectList"] == nil {
				array_to_be := []string{objectKey}
				// fmt.Println("array_to_be: ", array_to_be) // for debug

				anyJson["associatedObjectList"] = array_to_be
			} else { // anyJson["associatedObjectList"] != nil
				array_as_is := anyJson["associatedObjectList"].([]interface{})
				// fmt.Println("array_as_is: ", array_as_is) // for debug

				array_to_be := append(array_as_is, objectKey)
				// fmt.Println("array_to_be: ", array_to_be) // for debug

				anyJson["associatedObjectList"] = array_to_be
			}
			updatedJson, _ := json.Marshal(anyJson)
			// fmt.Println(string(updatedJson)) // for debug

			keyValue.Value = string(updatedJson)
		case common.StrDelete:
			var foundKey int
			var foundVal string
			for k, v := range objList {
				if v == objectKey {
					foundKey = k
					foundVal = v
					break
				}
			}
			if foundVal == "" {
				errString := "Cannot find the associated object " + objectKey + "."
				err = fmt.Errorf(errString)
				return nil, err
			} else {
				keyValue.Value, err = sjson.Delete(keyValue.Value, "associatedObjectList."+strconv.Itoa(foundKey))
				if err != nil {
					common.CBLog.Error(err)
					return nil, err
				}
			}
		}

		if err != nil {
			common.CBLog.Error(err)
			//return content, res.StatusCode, body, err
			return nil, err
		}
		err = common.CBStore.Put(key, keyValue.Value)
		if err != nil {
			common.CBLog.Error(err)
			//return content, res.StatusCode, body, err
			return nil, err
		}
		/*
			keyValue, _ := common.CBStore.Get(key)
			//fmt.Println("<" + keyValue.Key + "> \n" + keyValue.Value)
			fmt.Println("===========================")
			to_be = int8(gjson.Get(keyValue.Value, "inUseCount").Uint())
			return to_be, nil
		*/

		result, _ := GetAssociatedObjectList(nsId, resourceType, resourceId)
		return result, nil
	}
	errString := "Cannot get " + resourceType + " " + resourceId + "."
	err = fmt.Errorf(errString)
	return nil, err
}

// GetResource returns the requested TB MCIR object
func GetResource(nsId string, resourceType string, resourceId string) (interface{}, error) {

	//check, lowerizedResourceId, err := LowerizeAndCheckResource(nsId, resourceType, resourceId)
	//resourceId = lowerizedResourceId
	nsId = common.ToLower(nsId)
	resourceId = common.ToLower(resourceId)
	check, err := CheckResource(nsId, resourceType, resourceId)

	if !check {
		errString := "The " + resourceType + " " + resourceId + " does not exist."
		//mapA := map[string]string{"message": errString}
		//mapB, _ := json.Marshal(mapA)
		err := fmt.Errorf(errString)
		return nil, err
	}

	if err != nil {
		common.CBLog.Error(err)
		return nil, err
	}
	fmt.Println("[Get resource] " + resourceType + ", " + resourceId)

	key := common.GenResourceKey(nsId, resourceType, resourceId)
	//fmt.Println(key)

	keyValue, err := common.CBStore.Get(key)
	if err != nil {
		common.CBLog.Error(err)
		return nil, err
	}
	if keyValue != nil {
		switch resourceType {
		case common.StrImage:
			res := TbImageInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &res)
			if err != nil {
				common.CBLog.Error(err)
				return nil, err
			}
			return res, nil
		case common.StrSecurityGroup:
			res := TbSecurityGroupInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &res)
			if err != nil {
				common.CBLog.Error(err)
				return nil, err
			}
			return res, nil
		case common.StrSpec:
			res := TbSpecInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &res)
			if err != nil {
				common.CBLog.Error(err)
				return nil, err
			}
			return res, nil
		case common.StrSSHKey:
			res := TbSshKeyInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &res)
			if err != nil {
				common.CBLog.Error(err)
				return nil, err
			}
			return res, nil
		case common.StrVNet:
			res := TbVNetInfo{}
			err = json.Unmarshal([]byte(keyValue.Value), &res)
			if err != nil {
				common.CBLog.Error(err)
				return nil, err
			}
			return res, nil
		}

		//return true, nil
	}
	errString := "Cannot get " + resourceType + " " + resourceId + "."
	err = fmt.Errorf(errString)
	return nil, err
}

// CheckResource returns the existence of the TB MCIR resource in bool form.
func CheckResource(nsId string, resourceType string, resourceId string) (bool, error) {

	// Check parameters' emptiness
	if nsId == "" {
		err := fmt.Errorf("CheckResource failed; nsId given is null.")
		return false, err
	} else if resourceType == "" {
		err := fmt.Errorf("CheckResource failed; resourceType given is null.")
		return false, err
	} else if resourceId == "" {
		err := fmt.Errorf("CheckResource failed; resourceId given is null.")
		return false, err
	}

	// Check resourceType's validity
	if resourceType == common.StrImage ||
		resourceType == common.StrSSHKey ||
		resourceType == common.StrSpec ||
		resourceType == common.StrVNet ||
		resourceType == common.StrSecurityGroup {
		//resourceType == "subnet" ||
		//resourceType == "publicIp" ||
		//resourceType == "vNic" {
		// continue
	} else {
		err := fmt.Errorf("invalid resource type")
		return false, err
	}

	lowerizedNsId := common.ToLower(nsId)
	lowerizedResourceId := common.ToLower(resourceId)

	fmt.Println("[Check resource] " + resourceType + ", " + lowerizedResourceId)

	key := common.GenResourceKey(lowerizedNsId, resourceType, lowerizedResourceId)
	//fmt.Println(key)

	keyValue, _ := common.CBStore.Get(key)
	if keyValue != nil {
		return true, nil
	}
	return false, nil

}

/*
func convertSpiderResourceToTumblebugResource(resourceType string, i interface{}) (interface{}, error) {
	if resourceType == "" {
		err := fmt.Errorf("CheckResource failed; resourceType given is null.")
		return nil, err
	}

	// Check resourceType's validity
	if resourceType == common.StrImage ||
		resourceType == common.StrSSHKey ||
		resourceType == common.StrSpec ||
		resourceType == common.StrVNet ||
		resourceType == common.StrSecurityGroup {
		//resourceType == "subnet" ||
		//resourceType == "publicIp" ||
		//resourceType == "vNic" {
		// continue
	} else {
		err := fmt.Errorf("invalid resource type")
		return nil, err
	}

}
*/

// https://stackoverflow.com/questions/45139954/dynamic-struct-as-parameter-golang

type ReturnValue struct {
	CustomStruct interface{}
}

type NameOnly struct {
	Name string
}

// GetNameFromStruct accepts any struct for argument, and returns
func GetNameFromStruct(u interface{}) string {
	var result = ReturnValue{CustomStruct: u}

	//fmt.Println(result)

	msg, ok := result.CustomStruct.(NameOnly)
	if ok {
		//fmt.Printf("Message1 is %s\n", msg.Name)
		return msg.Name
	} else {
		return ""
	}
}

//func createResource(nsId string, resourceType string, u interface{}) (interface{}, int, []byte, error) {
