package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/varunturlapati/vtgqlgen/datasource/db"
	"github.com/varunturlapati/vtgqlgen/pkg/entity"
	"io/ioutil"
	"net/http"
)

func (r *RestRequests) GetRack(ctx context.Context, id int) (entity.ServerRack, error) {
	var rack entity.ServerRack
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/dcim/racks/%v/", id))
	if err != nil {
		return rack, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rack, err
	}
	err = json.Unmarshal(body, &rack)
	if err != nil {
		return rack, err
	}
	return rack, nil
}

func (r *RestRequests) ListRacks(ctx context.Context) ([]entity.ServerRack, error) {
	var res db.Result
	var racks []entity.ServerRack
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
		tmpRack := entity.ServerRack{
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
