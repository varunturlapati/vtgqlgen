package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/varunturlapati/vtgqlgen/datasource/db"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
)

func (r *RestRequests) GetRackFromNetbox(ctx context.Context, id int) (*entity.Rack, error) {
	var rack entity.Rack
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/dcim/racks/%v/", id))
	if err != nil {
		return &rack, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &rack, err
	}
	err = json.Unmarshal(body, &rack)
	if err != nil {
		return &rack, err
	}
	return &rack, nil
}

func (r *RestRequests) ListRacksFromNetbox(ctx context.Context) ([]*entity.Rack, error) {
	var res db.Result
	var racks []*entity.Rack
	//resp, err := http.Get("http://localhost:8000/api/dcim/racks/")
	nreq, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/api/dcim/racks/", nil)
	client := &http.Client{}
	nresp, err := client.Do(nreq)
	if err != nil {
		return nil, err
	}
	defer nresp.Body.Close()
	body, err := ioutil.ReadAll(nresp.Body)
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	for _, elem := range res.Results {
		tmpRack := &entity.Rack{
			Id:   elem.Id,
			Name: elem.Name,
			CustomFields: entity.CustomFields{
				RblxRackId:     elem.CustomFields.RblxRackId,
				DesignRevision: elem.CustomFields.DesignRevision,
				CageId:         elem.CustomFields.CageId,
			},
			Created: elem.Created,
		}
		racks = append(racks, tmpRack)
	}
	return racks, nil
}

func (r *RestRequests) GetServerFromNetbox(ctx context.Context, id int) (*entity.Server, error) {
	var Server entity.Server
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/dcim/devices/%v/", id))
	if err != nil {
		return &Server, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//log.Printf("Netbox server response body is %v and error is %v\n", body, err)
	if err != nil {
		return &Server, err
	}
	err = json.Unmarshal(body, &Server)
	log.Printf("Netbox server is %+v and json unmarshal error is %v\n", Server, err)
	if err != nil {
		return &Server, err
	}
	return &Server, nil
}

func (r *RestRequests) GetServerByAttrsFromNetbox(ctx context.Context, attrs *entity.ServerAttrs) (*entity.Server, error) {
	var Server entity.Server
	var res AllResults

	tmp := strings.Split(attrs.HostName, " ")
	qParam := strings.Join(tmp, "%20")
	clause := fmt.Sprintf("?name=%s", qParam)
	getSrvUrl := fmt.Sprintf("http://localhost:8000/api/dcim/devices/%v", clause)

	resp, err := http.Get(getSrvUrl)
	if err != nil {
		log.Printf("Err from rest for URL: %s is\n Err: %v", getSrvUrl, err)
		return &Server, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	// TODO - just taking the first one. Scope for doing something else
	if res.Results != nil && len(res.Results) > 0 {
		topRes := res.Results[0]
		Server.NetboxName = topRes["display_name"].(string)
		tmpRack := topRes["rack"]
		if tmpRack != nil {
			tmpRackMap := tmpRack.(map[string]interface{})
			Server.RackName = tmpRackMap["name"].(string)
		}
		log.Printf("Json unmarshal error with one server: %v\n", err)
	}
	return &Server, nil
}

func (r *RestRequests) ListServersFromNetbox(ctx context.Context) ([]*entity.Server, error) {
	var servers []*entity.Server
	var res AllResults
	//resp, err := http.Get("http://localhost:8000/api/dcim/Servers/")
	nreq, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/api/dcim/devices/", nil)
	client := &http.Client{}
	nresp, err := client.Do(nreq)
	if err != nil {
		return nil, err
	}
	defer nresp.Body.Close()
	body, err := ioutil.ReadAll(nresp.Body)
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	for _, elem := range res.Results {
		var tmpRackName, tmpDispName string
		tmpServer := elem//.(map[string]interface{})
		// log.Printf("TmpServer looks like %+v\n", tmpServer)
		tmpDispName = tmpServer["display_name"].(string)
		tmpId := tmpServer["id"].(float64)
		tmpRack := tmpServer["rack"]
		if tmpRack != nil {
			tmpRackMap := tmpRack.(map[string]interface{})
			tmpRackName = tmpRackMap["name"].(string)
		}
		retServer := &entity.Server{
			Id:         int(tmpId),
			RackName:   tmpRackName,
			NetboxName: tmpDispName,
		}
		servers = append(servers, retServer)
	}
	return servers, nil
}

type ListRacksByFruitIDsRow entity.Rack

func (r *RestRequests) ListRacksByFruitIDs(ctx context.Context, fruitIDs []int) ([]ListRacksByFruitIDsRow, error) {
	var retList []ListRacksByFruitIDsRow
	for _, fid := range fruitIDs {
		res, err := r.GetRackFromNetbox(ctx, fid)
		if err != nil {
			return nil, err
		}
		var tmp ListRacksByFruitIDsRow
		tmp.Id = res.Id
		tmp.Name = res.Name
		tmp.Created = res.Created
		tmp.CustomFields = res.CustomFields
		retList = append(retList, tmp)
	}
	return retList, nil
}
